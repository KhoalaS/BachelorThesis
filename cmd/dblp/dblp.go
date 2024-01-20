package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"

	"golang.org/x/text/encoding/ianaindex"
)

type Publication struct {
	Authors []Author `xml:"author"`
}

type Author struct {
	XMLName xml.Name `xml:"author"`
	Aux string `xml:"aux,attr"`
	Bibtex string `xml:"bibtex,attr"`
	Orcid string `xml:"orcid,attr"`
	Label string `xml:"label,attr"`
	Type string `xml:"type,attr"`
	Name string `xml:",chardata"`
}

func main(){

	in := flag.String("i", "./dblp.xml", "filepath to dblp.xml") 
	flag.Parse()

	f, err := os.Open(*in)
	if err != nil{
		log.Fatal("Could not open dblp.xml")
	}
	decoder := xml.NewDecoder(f)
	decoder.CharsetReader = func(charset string, reader io.Reader) (io.Reader, error) {
		enc, err := ianaindex.IANA.Encoding(charset)
		if err != nil {
			return nil, fmt.Errorf("charset %s: %s", charset, err.Error())
		}
		return enc.NewDecoder().Reader(reader), nil
	}
	var entities map[string]string
	entitiyFile, err := os.Open("../entities.json")
	if err != nil{
		log.Fatal("Could not open entities.json")
	}
	entitiyBytes, _ := io.ReadAll(entitiyFile)
	json.Unmarshal(entitiyBytes, &entities)
	entitiyFile.Close()

	var authorIdCounter int32 =0
	authorId := make(map[string]int32)
	coauthor := make(map[int32]map[int32]bool)

	pubtypes := make(map[string]int)

	decoder.Entity = entities
	r := regexp.MustCompile("article|inproceedings|proceedings|book|incollection|phdthesis|mastersthesis")
	count := 0

	for {
		t, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF){
				log.Default().Println("Finished reading file")
				break
			}
			log.Fatal(err)
		}

		switch tt := t.(type){
		case xml.StartElement:
			if tt.Name.Local == "dblp" {
				continue
			}
			if r.MatchString(tt.Name.Local) {
				pubtypes[tt.Name.Local]++
				var p Publication
				err := decoder.DecodeElement(&p, &tt)
				if err != nil {
					log.Default().Println(err)
				}

				for _, auth := range p.Authors {
					if _, ex := authorId[auth.Name]; !ex {
						authorId[auth.Name] = authorIdCounter
						authorIdCounter++
					}
				}

				if len(p.Authors) > 1 {
					a := p.Authors[0]
					for i:=1; i<len(p.Authors); i++ {
						if _ ,ex := coauthor[authorId[a.Name]]; !ex {
							coauthor[authorId[a.Name]] = make(map[int32]bool)
						}
						coauthor[authorId[a.Name]][authorId[p.Authors[i].Name]] = true
					}
				}
				
				count++
				fmt.Printf("Processed %d publications\r", count)
			}
		}
	}
	f.Close()

	fmt.Println()
	pubSum := 0
	for _, val := range pubtypes {
		pubSum += val
	}
	fmt.Printf("Processed %d publications\n", pubSum)
	fmt.Printf("Found %d authors\n", len(authorId))

	for t, v := range pubtypes {
		fmt.Println(t, v)
	}

	fmt.Println("Writing coauthor graph to file")

	out, err := os.Create("./coauthor.txt")
	if err != nil {
		log.Fatal("Could not open file coauthor.txt")
	}
	defer out.Close()

	count = 0
	bw := bufio.NewWriterSize(out, 8192)
	
	for k, val := range coauthor {
		for v := range val {
			line := strconv.Itoa(int(k)) + " " + strconv.Itoa(int(v)) + "\n"
			bw.Write([]byte(line))
		}
		fmt.Printf("Processed %d authors\r", count)
		count++
	}
	err = bw.Flush()
	if err != nil {
		log.Fatal("Could not flush buffer for coauthor.txt")
	}
	fmt.Println()
	fmt.Println("Processed all edges")
}