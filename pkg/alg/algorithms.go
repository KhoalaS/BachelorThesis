package alg

import (
	"bufio"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/KhoalaS/BachelorThesis/pkg/hypergraph"
)

var Labels = []string{"kTiny", "kVertDom", "kEdgeDom", "kSmall", "kTri", "kExtTri", "kApVertDom", "kApDoubleVertDom", "kSmallEdgeDegTwo", "kFallback"}

var Ratios = map[string]IntTuple{
	"kTiny":            {A: 1, B: 1},
	"kSmall":           {A: 2, B: 1},
	"kTri":             {A: 3, B: 2},
	"kExtTri":          {A: 4, B: 2},
	"kApVertDom":       {A: 2, B: 1},
	"kApDoubleVertDom": {A: 2, B: 1},
	"kSmallEdgeDegTwo": {A: 4, B: 2},
	"kFallback":        {A: 3, B: 1},
}

func LoggingThreeHS_F3ApprPoly(g *hypergraph.HyperGraph, c map[int32]bool, graphtype string, masterfilename string, iteration int, outdir string) map[string]int {

	vSize := len(g.Vertices)
	eSize := len(g.Edges)

	header := "Ratio;"
	header += strings.Join(Labels, ";") + "\n"

	os.Mkdir(outdir, 0700)

	logfilename := fmt.Sprintf("%s/%s_%.2f_%d.csv", outdir, graphtype, float64(len(g.Edges))/float64(len(g.Vertices)), iteration)
	logfile, err := os.Create(logfilename)
	if err != nil {
		log.Fatalf("Could not create file %s", logfilename)
	}
	logWriter := bufio.NewWriter(logfile)

	fMasterFilename := fmt.Sprintf("%s/%s", outdir, masterfilename)
	masterfile, err := os.OpenFile(fMasterFilename, os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			masterfile, _ = os.Create(fMasterFilename)
			masterfile.WriteString(header)
		} else {
			log.Fatalf("Could not open file %s", fMasterFilename)
		}
	}

	defer logfile.Close()
	defer masterfile.Close()

	logWriter.WriteString(header)

	execs := MakeExecs()
	msg := ""

	for len(g.Edges) > 0 {
		execs = ApplyRules(g, c, execs, 0)
		execs["kFallback"] += hypergraph.F3TargetLowDegree(g, c)

		msg = fmt.Sprintf("%f;", GetRatio(execs))
		for _, v := range Labels {
			msg += fmt.Sprintf("%d;", execs[v])
		}
		msg = msg[:len(msg)-1] + "\n"
		logWriter.WriteString(msg)
	}
	masterfile.WriteString(fmt.Sprintf("%s;%d;%d;%d\n", msg[0:len(msg)-1], vSize, eSize, len(c)))
	logWriter.Flush()
	return execs
}

func LoggingThreeHS_F3ApprPolyFrontier(g *hypergraph.HyperGraph, c map[int32]bool, graphtype string, masterfilename string, iteration int, outdir string) map[string]int {

	vSize := len(g.Vertices)
	eSize := len(g.Edges)

	header := "Ratio;"
	header += strings.Join(Labels, ";") + "\n"

	os.Mkdir(outdir, 0700)

	logfilename := fmt.Sprintf("%s/%s_%.2f_%d.csv", outdir, graphtype, float64(len(g.Edges))/float64(len(g.Vertices)), iteration)
	logfile, err := os.Create(logfilename)
	if err != nil {
		log.Fatalf("Could not create file %s", logfilename)
	}
	logWriter := bufio.NewWriter(logfile)

	fMasterFilename := fmt.Sprintf("%s/%s", outdir, masterfilename)
	masterfile, err := os.OpenFile(fMasterFilename, os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			masterfile, _ = os.Create(fMasterFilename)
			masterfile.WriteString(header[0 : len(header)-1])
			masterfile.WriteString(";Vertices;Edges;HittingSet\n")
		} else {
			log.Fatalf("Could not open file %s", fMasterFilename)
		}
	}

	defer logfile.Close()
	defer masterfile.Close()

	logWriter.WriteString(header)

	msg := ""

	execs := MakeExecs()
	ApplyRules(g, c, execs, 0)
	expDepth := 2

	e := hypergraph.F3TargetLowDegreeDetect(g)
	if e != -1 {
		execs["kFallback"] += 1
		for v := range g.Edges[e].V {
			c[v] = true
		}
	}

	msg = fmt.Sprintf("%f;", GetRatio(execs))
	for _, v := range Labels {
		msg += fmt.Sprintf("%d;", execs[v])
	}
	msg = msg[:len(msg)-1] + "\n"
	logWriter.WriteString(msg)

	if len(g.Edges) == 0 {
		logWriter.Flush()
		masterfile.WriteString(msg)
		return execs
	}

	gf := hypergraph.F3_ExpandFrontier(g, e, expDepth)
	fmt.Println(len(gf.Edges))

	for len(g.Edges) > 0 {
		expand := make(map[int32]bool)
		ApplyRulesFrontier(gf, g, c, execs, expand)
		if len(expand) > 0 {
			fmt.Println("Expand")
			fmt.Println(execs)
			gf = hypergraph.ExpandFrontier(g, expDepth, expand)
			msg = fmt.Sprintf("%f;", GetRatio(execs))
			for _, v := range Labels {
				msg += fmt.Sprintf("%d;", execs[v])
			}
			msg = msg[:len(msg)-1] + "\n"
			logWriter.WriteString(msg)
			continue
		}

		isSmall := false
		e := hypergraph.F3TargetLowDegreeDetect(g)
		if e == -1 {
			fmt.Println("Could not find size 3 edge")
			e = hypergraph.F2Detect(g)
			isSmall = true
		}

		for v := range g.Edges[e].V {
			c[v] = true
		}

		gf = hypergraph.F3_ExpandFrontier(g, e, expDepth)
		if isSmall {
			execs["kSmall"] += 1
		} else {
			execs["kFallback"] += 1
		}

		msg = fmt.Sprintf("%f;", GetRatio(execs))
		for _, v := range Labels {
			msg += fmt.Sprintf("%d;", execs[v])
		}
		msg = msg[:len(msg)-1] + "\n"
		logWriter.WriteString(msg)
	}
	masterfile.WriteString(fmt.Sprintf("%s;%d;%d;%d\n", msg[0:len(msg)-1], vSize, eSize, len(c)))
	logWriter.Flush()
	return execs
}

func ThreeHS_F3ApprPoly(g *hypergraph.HyperGraph, c map[int32]bool, prio int) map[string]int {
	execs := MakeExecs()
	f3 := 0

	for len(g.Edges) > 0 {
		execs = ApplyRules(g, c, execs, prio)
		prio = 0
		kFallback := hypergraph.F3TargetLowDegree(g, c)
		execs["kFallback"] += kFallback
		f3++
		//prio = nextPrio
	}
	return execs
}

func ApplyRules(g *hypergraph.HyperGraph, c map[int32]bool, execs map[string]int, prio int) map[string]int {

	switch prio {
	case 2:
		exec := hypergraph.SmallTriangleRule(g, c)
		execs["kTri"] += exec
	}

	for {
		kTiny := hypergraph.RemoveEdgeRule(g, c, hypergraph.TINY)
		kVertDom := hypergraph.VertexDominationRule(g, c)
		kTiny += hypergraph.RemoveEdgeRule(g, c, hypergraph.TINY)
		//kEdgeDom := hypergraph.EdgeDominationRule(g)
		kEdgeDom := 0
		kApVertDom := hypergraph.ApproxVertexDominationRule(g, c)
		kApDoubleVertDom := hypergraph.ApproxDoubleVertexDominationRule(g, c)
		kSmallEdgeDegTwo := hypergraph.SmallEdgeDegreeTwoRule(g, c)
		kTri := hypergraph.SmallTriangleRule(g, c)
		kExtTri := hypergraph.ExtendedTriangleRule(g, c)
		//kSmall := hypergraph.RemoveEdgeRule(g, c, hypergraph.SMALL)
		kSmall := 0

		execs["kTiny"] += kTiny
		execs["kVertDom"] += kVertDom
		execs["kEdgeDom"] += kEdgeDom
		execs["kTri"] += kTri
		execs["kExtTri"] += kExtTri
		execs["kSmall"] += kSmall
		execs["kApVertDom"] += kApVertDom
		execs["kApDoubleVertDom"] += kApDoubleVertDom
		execs["kSmallEdgeDegTwo"] += kSmallEdgeDegTwo

		if kTiny+kTri+kSmall+kApVertDom+kApDoubleVertDom+kEdgeDom+kVertDom+kExtTri+kSmallEdgeDegTwo == 0 {
			break
		}
	}

	fmt.Println(execs)

	return execs
}

func ThreeHS_F3ApprPolyFrontier(g *hypergraph.HyperGraph, c map[int32]bool) map[string]int {
	execs := MakeExecs()
	ApplyRules(g, c, execs, 0)
	expDepth := 2

	if len(g.Edges) == 0 {
		return execs
	}

	e := hypergraph.F3TargetLowDegreeDetect(g)
	if e != -1 {
		execs["kFallback"] += 1
		for v := range g.Edges[e].V {
			c[v] = true
		}
	}

	gf := hypergraph.F3_ExpandFrontier(g, e, expDepth)

	for len(g.Edges) > 0 {
		expand := make(map[int32]bool)
		ApplyRulesFrontier(gf, g, c, execs, expand)
		if len(expand) > 0 {
			gf = hypergraph.ExpandFrontier(g, expDepth, expand)
			continue
		}

		isSmall := false
		e := hypergraph.F3TargetLowDegreeDetect(g)
		if e == -1 {
			fmt.Println("Could not find size 3 edge")
			e = hypergraph.F2Detect(g)
			isSmall = true
		}

		for v := range g.Edges[e].V {
			c[v] = true
		}

		gf = hypergraph.F3_ExpandFrontier(g, e, expDepth)
		if isSmall {
			execs["kSmall"] += 1
		} else {
			execs["kFallback"] += 1
		}
	}
	return execs
}

func ApplyRulesFrontier(gf *hypergraph.HyperGraph, g *hypergraph.HyperGraph, c map[int32]bool, execs map[string]int, expand map[int32]bool) {
	for {
		kTiny := hypergraph.S_RemoveEdgeRule(gf, g, c, hypergraph.TINY, expand)
		kVertDom := hypergraph.S_VertexDominationRule(gf, g, c, expand)
		kTiny += hypergraph.S_RemoveEdgeRule(gf, g, c, hypergraph.TINY, expand)
		//kEdgeDom := hypergraph.S_EdgeDominationRule(gf, g, expand)
		kEdgeDom := 0
		kApVertDom := hypergraph.S_ApproxVertexDominationRule(gf, g, c, expand)
		kApDoubleVertDom := hypergraph.S_ApproxDoubleVertexDominationRule2(gf, g, c, expand)
		kSmallEdgeDegTwo := hypergraph.S_SmallEdgeDegreeTwoRule(gf, g, c, expand)
		kTri := hypergraph.S_SmallTriangleRule(gf, g, c, expand)
		kExtTri := hypergraph.S_ExtendedTriangleRule(gf, g, c, expand)
		//kSmall := hypergraph.S_RemoveEdgeRule(gf, g, c, hypergraph.SMALL, expand)
		kSmall := 0

		execs["kTiny"] += kTiny
		execs["kVertDom"] += kVertDom
		execs["kEdgeDom"] += kEdgeDom
		execs["kTri"] += kTri
		execs["kExtTri"] += kExtTri
		execs["kSmall"] += kSmall
		execs["kApVertDom"] += kApVertDom
		execs["kApDoubleVertDom"] += kApDoubleVertDom
		execs["kSmallEdgeDegTwo"] += kSmallEdgeDegTwo
		if kTiny+kTri+kSmall+kEdgeDom+kVertDom+kExtTri+kApVertDom+kApDoubleVertDom+kSmallEdgeDegTwo == 0 {
			break
		}
	}
}

func ApplyRulesRand(g *hypergraph.HyperGraph, c map[int32]bool, execs map[string]int, prio int) map[string]int {

	switch prio {
	case 2:
		exec := hypergraph.SmallTriangleRule(g, c)
		execs["kTri"] += exec
	}

	for {

		kApVertDom := 0
		kApDoubleVertDom := 0
		kSmallEdgeDegTwo := 0
		kTri := 0
		kExtTri := 0
		kSmall := 0
		kVertDom := 0
		kTiny := 0
		kEdgeDom := 0

		perm := make([]int, 9)
		for i := range perm {
			perm[i] = i
		}

		Shuffle[int](perm)

		for i := 0; i < 9; i++ {
			switch perm[i] {
			case 0:
				kApVertDom = hypergraph.ApproxVertexDominationRule(g, c)
			case 1:
				kApDoubleVertDom = hypergraph.ApproxDoubleVertexDominationRule2(g, c)
			case 2:
				kSmallEdgeDegTwo = hypergraph.SmallEdgeDegreeTwoRule(g, c)
			case 3:
				kTri = hypergraph.SmallTriangleRule(g, c)
			case 4:
				kExtTri = hypergraph.ExtendedTriangleRule(g, c)
			case 5:
				kSmall = hypergraph.RemoveEdgeRule(g, c, hypergraph.SMALL)
			case 6:
				kVertDom = hypergraph.VertexDominationRule(g, c)
			case 7:
				kTiny = hypergraph.RemoveEdgeRule(g, c, hypergraph.TINY)
			case 8:
				kEdgeDom = hypergraph.EdgeDominationRule(g)
			}
		}

		execs["kTiny"] += kTiny
		execs["kVertDom"] += kVertDom
		execs["kEdgeDom"] += kEdgeDom
		execs["kTri"] += kTri
		execs["kExtTri"] += kExtTri
		execs["kSmall"] += kSmall
		execs["kApVertDom"] += kApVertDom
		execs["kApDoubleVertDom"] += kApDoubleVertDom
		execs["kSmallEdgeDegTwo"] += kSmallEdgeDegTwo

		if kTiny+kEdgeDom+kVertDom+kTri+kSmall+kApVertDom+kApDoubleVertDom+kSmallEdgeDegTwo+kExtTri == 0 {
			break
		}
	}

	return execs
}

func ApplyRulesFrontierRand(gf *hypergraph.HyperGraph, g *hypergraph.HyperGraph, c map[int32]bool, execs map[string]int, expand map[int32]bool) map[string]int {

	for {

		kApVertDom := 0
		kApDoubleVertDom := 0
		kSmallEdgeDegTwo := 0
		kTri := 0
		kExtTri := 0
		kSmall := 0
		kVertDom := 0
		kTiny := 0
		kEdgeDom := 0

		perm := make([]int, 9)
		for i := range perm {
			perm[i] = i
		}

		Shuffle[int](perm)

		for i := 0; i < 9; i++ {
			switch perm[i] {
			case 0:
				kApVertDom = hypergraph.S_ApproxVertexDominationRule(gf, g, c, expand)
			case 1:
				kApDoubleVertDom = hypergraph.S_ApproxDoubleVertexDominationRule2(gf, g, c, expand)
			case 2:
				kSmallEdgeDegTwo = hypergraph.S_SmallEdgeDegreeTwoRule(gf, g, c, expand)
			case 3:
				kTri = hypergraph.S_SmallTriangleRule(gf, g, c, expand)
			case 4:
				kExtTri = hypergraph.S_ExtendedTriangleRule(gf, g, c, expand)
			case 5:
				kSmall = hypergraph.S_RemoveEdgeRule(gf, g, c, hypergraph.SMALL, expand)
			case 6:
				kVertDom = hypergraph.S_VertexDominationRule(gf, g, c, expand)
			case 7:
				kTiny = hypergraph.S_RemoveEdgeRule(gf, g, c, hypergraph.TINY, expand)
			case 8:
				kEdgeDom = hypergraph.S_EdgeDominationRule(gf, g, expand)

			}
		}

		execs["kTiny"] += kTiny
		execs["kVertDom"] += kVertDom
		execs["kEdgeDom"] += kEdgeDom
		execs["kTri"] += kTri
		execs["kExtTri"] += kExtTri
		execs["kSmall"] += kSmall
		execs["kApVertDom"] += kApVertDom
		execs["kApDoubleVertDom"] += kApDoubleVertDom
		execs["kSmallEdgeDegTwo"] += kSmallEdgeDegTwo

		if kTiny+kEdgeDom+kVertDom+kTri+kSmall+kApVertDom+kApDoubleVertDom+kSmallEdgeDegTwo+kExtTri == 0 {
			break
		}
	}

	return execs
}

func GreedyHighDeg(g *hypergraph.HyperGraph, c map[int32]bool) {

	for len(g.Edges) > 0 {
		max := 0
		var remVertex int32 = -1

		for v := range g.Vertices {
			d := g.Deg(v)
			if d > max {
				max = d
				remVertex = v
			}
		}

		for e := range g.IncMap[remVertex] {
			g.RemoveEdge(e)
		}
		c[remVertex] = true
		fmt.Println(len(g.Edges))
	}
}

