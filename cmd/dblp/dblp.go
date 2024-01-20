package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"

	"github.com/KhoalaS/BachelorThesis/pkg/hypergraph"
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

	in := flag.String("-i", "./dblp.xml", "filepath to dblp.xml") 
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
	entitiyFile, _ := os.Open("../entities.json")
	entitiyBytes, _ := io.ReadAll(entitiyFile)
	json.Unmarshal(entitiyBytes, &entities)

	var authorIdCounter int32 =0
	authorId := make(map[string]int32)
	coauthor := make(map[int32]map[int32]bool)

	decoder.Entity = entities
	r := regexp.MustCompile("article|inproceedings|proceedings|book|incollection|phdthesis|mastersthesis|www")
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

				hypergraph.GetSubsetsRec3[Author](p.Authors, 2, func (sub []Author)  {
					if _ ,ex := coauthor[authorId[sub[0].Name]]; !ex {
						coauthor[authorId[sub[0].Name]] = make(map[int32]bool)
					}
					coauthor[authorId[sub[0].Name]][authorId[sub[1].Name]] = true
				})
				count++
				fmt.Printf("Processed %d publications\r", count)
			}
		}
	}
}