package alg

import (
	_ "embed"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/KhoalaS/BachelorThesis/pkg/hypergraph"
)

var Labels = []string{"kTiny", "kVertDom", "kEdgeDom", "kSmall", "kTri", "kExtTri", "kApVertDom", "kApDoubleVertDom", "kSmallEdgeDegTwo", "kSmallEdgeDegTwo2", "kFallback"}

var Ratios = map[string]IntTuple{
	"kTiny":             {A: 1, B: 1},
	"kSmall":            {A: 2, B: 1},
	"kTri":              {A: 3, B: 2},
	"kExtTri":           {A: 4, B: 2},
	"kApVertDom":        {A: 2, B: 1},
	"kApDoubleVertDom":  {A: 2, B: 1},
	"kSmallEdgeDegTwo":  {A: 4, B: 2},
	"kSmallEdgeDegTwo2": {A: 3, B: 2},
	"kFallback":         {A: 3, B: 1},
}

func LoggingThreeHS_F3ApprPoly(g *hypergraph.HyperGraph, c map[int32]bool, graphtype string, masterfilename string, iteration int, outdir string) map[string]int {

	vSize := len(g.Vertices)
	eSize := len(g.Edges)

	header := "Ratio;"
	header += strings.Join(Labels, ";")

	os.Mkdir(outdir, 0700)

	fMasterFilename := fmt.Sprintf("%s/%s", outdir, masterfilename)
	masterfile, err := os.OpenFile(fMasterFilename, os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			masterfile, _ = os.Create(fMasterFilename)
			masterfile.WriteString(header)
			masterfile.WriteString(";Vertices;Edges;HittingSet;Opt;Time\n")
		} else {
			log.Fatalf("Could not open file %s", fMasterFilename)
		}
	}

	defer masterfile.Close()

	execs := MakeExecs()
	msg := ""
	start := time.Now()

	for len(g.Edges) > 0 {
		execs = ApplyRules(g, c, execs, 0)
		execs["kFallback"] += hypergraph.F3TargetLowDegree(g, c)
	}

	stop := time.Since(start).Seconds()

	msg = fmt.Sprintf("%f;", GetRatio(execs))
	for _, v := range Labels {
		msg += fmt.Sprintf("%d;", execs[v])
	}
	msg = msg[:len(msg)-1]
	masterfile.WriteString(fmt.Sprintf("%s;%d;%d;%d;%d;%.2f\n", msg, vSize, eSize, len(c), GetEstOpt(execs), RoundUp(stop, 2)))
	return execs
}

func LoggingThreeHS_F3ApprPolyFrontier(g *hypergraph.HyperGraph, c map[int32]bool, graphtype string, masterfilename string, iteration int, outdir string) map[string]int {

	vSize := len(g.Vertices)
	eSize := len(g.Edges)

	header := "Ratio;"
	header += strings.Join(Labels, ";")

	os.Mkdir(outdir, 0700)

	fMasterFilename := fmt.Sprintf("%s/%s", outdir, masterfilename)
	masterfile, err := os.OpenFile(fMasterFilename, os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			masterfile, _ = os.Create(fMasterFilename)
			masterfile.WriteString(header)
			masterfile.WriteString(";Vertices;Edges;HittingSet;Opt;Time\n")
		} else {
			log.Fatalf("Could not open file %s", fMasterFilename)
		}
	}

	defer masterfile.Close()

	msg := ""

	start := time.Now()

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
	msg = msg[:len(msg)-1]

	if len(g.Edges) == 0 {
		masterfile.WriteString(msg)
		return execs
	}

	gf := hypergraph.F3_ExpandFrontier(g, e, expDepth)

	for len(g.Edges) > 0 {
		expand := make(map[int32]bool)
		ApplyRulesFrontier(gf, g, c, execs, expand)
		if len(expand) > 0 {
			gf = hypergraph.ExpandFrontier(g, expDepth, expand)
			continue
		}

		e := hypergraph.F3TargetLowDegreeDetect(g)
		if e == -1 {
			continue
		}

		for v := range g.Edges[e].V {
			c[v] = true
		}

		gf = hypergraph.F3_ExpandFrontier(g, e, expDepth)
		execs["kFallback"] += 1

	}

	stop := time.Since(start).Seconds()

	msg = fmt.Sprintf("%f;", GetRatio(execs))
	for _, v := range Labels {
		msg += fmt.Sprintf("%d;", execs[v])
	}
	msg = msg[:len(msg)-1]
	masterfile.WriteString(fmt.Sprintf("%s;%d;%d;%d;%d;%.2f\n", msg, vSize, eSize, len(c), GetEstOpt(execs), RoundUp(stop, 2)))
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
		kSmallEdgeDegTwo, kSmallEdgeDegTwo2 := hypergraph.SmallEdgeDegreeTwoRule(g, c)
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
		execs["kSmallEdgeDegTwo2"] += kSmallEdgeDegTwo2

		if kTiny+kTri+kSmall+kApVertDom+kApDoubleVertDom+kEdgeDom+kVertDom+kExtTri+kSmallEdgeDegTwo+kSmallEdgeDegTwo2 == 0 {
			break
		}
	}
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

		e := hypergraph.F3TargetLowDegreeDetect(g)
		if e == -1 {
			continue
		}

		for v := range g.Edges[e].V {
			c[v] = true
		}

		gf = hypergraph.F3_ExpandFrontier(g, e, expDepth)
		execs["kFallback"] += 1
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
		kSmallEdgeDegTwo, kSmallEdgeDegTwo2 := hypergraph.S_SmallEdgeDegreeTwoRule(gf, g, c, expand)
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
		execs["kSmallEdgeDegTwo2"] += kSmallEdgeDegTwo2

		if kTiny+kTri+kSmall+kEdgeDom+kVertDom+kExtTri+kApVertDom+kApDoubleVertDom+kSmallEdgeDegTwo+kSmallEdgeDegTwo2 == 0 {
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
		kSmallEdgeDegTwo2 := 0
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
				kApDoubleVertDom = hypergraph.ApproxDoubleVertexDominationRule(g, c)
			case 2:
				kSmallEdgeDegTwo, kSmallEdgeDegTwo2 = hypergraph.SmallEdgeDegreeTwoRule(g, c)
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
		execs["kSmallEdgeDegTwo2"] += kSmallEdgeDegTwo2

		if kTiny+kEdgeDom+kVertDom+kTri+kSmall+kApVertDom+kApDoubleVertDom+kSmallEdgeDegTwo+kExtTri+kSmallEdgeDegTwo2 == 0 {
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
		kSmallEdgeDegTwo2 := 0
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
				kSmallEdgeDegTwo, kSmallEdgeDegTwo2 = hypergraph.S_SmallEdgeDegreeTwoRule(gf, g, c, expand)
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
		execs["kSmallEdgeDegTwo2"] += kSmallEdgeDegTwo2

		if kTiny+kEdgeDom+kVertDom+kTri+kSmall+kApVertDom+kApDoubleVertDom+kSmallEdgeDegTwo+kExtTri+kSmallEdgeDegTwo2 == 0 {
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
	}
}

func ApplyRulesSingle(gf *hypergraph.HyperGraph, g *hypergraph.HyperGraph, c map[int32]bool, execs map[string]int, expand map[int32]bool) {
	kTiny := 0
	kVertDom := 0
	kEdgeDom := 0

	for {
		old := kTiny + kVertDom + kEdgeDom

		kTiny += hypergraph.S_RemoveEdgeRule(gf, g, c, hypergraph.TINY, expand)
		kVertDom += hypergraph.S_VertexDominationRule(gf, g, c, expand)
		kEdgeDom += hypergraph.S_EdgeDominationRule(gf, g, expand)

		new := kTiny + kVertDom + kEdgeDom
		if old == new {
			break
		}
	}

	execs["kTiny"] += kTiny
	execs["kVertDom"] += kVertDom
	execs["kEdgeDom"] += kEdgeDom

	kApVertDom := hypergraph.FS_ApproxVertexDominationRule(gf, g, c, expand)
	if kApVertDom > 0 {
		execs["kApVertDom"] += kApVertDom
		return
	}
	kApDoubleVertDom := hypergraph.FS_ApproxDoubleVertexDominationRule(gf, g, c, expand)
	if kApDoubleVertDom > 0 {
		execs["kApDoubleVertDom"] += kApDoubleVertDom
		return
	}
	kSmallEdgeDegTwo, kSmallEdgeDegTwo2 := hypergraph.FS_SmallEdgeDegreeTwoRule(gf, g, c, expand)
	if kSmallEdgeDegTwo+kSmallEdgeDegTwo2 > 0 {
		execs["kSmallEdgeDegTwo"] += kSmallEdgeDegTwo
		execs["kSmallEdgeDegTwo2"] += kSmallEdgeDegTwo2
		return
	}
	kTri := hypergraph.FS_SmallTriangleRule(gf, g, c, expand)
	if kTri > 0 {
		execs["kTri"] += kTri
		return
	}
	kExtTri := hypergraph.FS_ExtendedTriangleRule(gf, g, c, expand)
	if kExtTri > 0 {
		execs["kExtTri"] += kExtTri
		return
	}
	kSmall := hypergraph.FS_RemoveEdgeRule(gf, g, c, hypergraph.SMALL, expand)
	if kSmall > 0 {
		execs["kSmall"] += kSmall
		return
	}
}

func findEntry(gf *hypergraph.HyperGraph, g *hypergraph.HyperGraph, c map[int32]bool, execs map[string]int, expand map[int32]bool) bool {
	kApVertDom := hypergraph.FS_ApproxVertexDominationRule(gf, g, c, expand)
	if kApVertDom > 0 {
		execs["kApVertDom"] += kApVertDom
		return true
	}
	kApDoubleVertDom := hypergraph.FS_ApproxDoubleVertexDominationRule(gf, g, c, expand)
	if kApDoubleVertDom > 0 {
		execs["kApDoubleVertDom"] += kApDoubleVertDom
		return true
	}
	kSmallEdgeDegTwo, kSmallEdgeDegTwo2 := hypergraph.FS_SmallEdgeDegreeTwoRule(gf, g, c, expand)
	if kSmallEdgeDegTwo+kSmallEdgeDegTwo2 > 0 {
		execs["kSmallEdgeDegTwo"] += kSmallEdgeDegTwo
		execs["kSmallEdgeDegTwo2"] += kSmallEdgeDegTwo2
		return true
	}
	kTri := hypergraph.FS_SmallTriangleRule(gf, g, c, expand)
	if kTri > 0 {
		execs["kTri"] += kTri
		return true
	}
	kExtTri := hypergraph.FS_ExtendedTriangleRule(gf, g, c, expand)
	if kExtTri > 0 {
		execs["kExtTri"] += kExtTri
		return true
	}
	kSmall := hypergraph.FS_RemoveEdgeRule(gf, g, c, hypergraph.SMALL, expand)
	if kSmall > 0 {
		execs["kSmall"] += kSmall
		return true
	}
	return false
}

func PreProcessOnly(g *hypergraph.HyperGraph, c map[int32]bool, execs map[string]int, expand map[int32]bool) {
	kTiny := 0
	kVertDom := 0
	kEdgeDom := 0

	for {
		old := kTiny + kVertDom + kEdgeDom

		kTiny += hypergraph.FS_TinyEdgeRule(g, c, expand)
		kVertDom += hypergraph.FS_VertexDominationRule(g, expand)
		kEdgeDom += hypergraph.FS_EdgeDominationRule(g, expand)

		new := kTiny + kVertDom + kEdgeDom
		if old == new {
			break
		}
	}
	execs["kTiny"] += kTiny
	execs["kVertDom"] += kVertDom
	execs["kEdgeDom"] += kEdgeDom
}

func ThreeHS_F3ApprPolyFrontierSingle(g *hypergraph.HyperGraph, c map[int32]bool, logging bool) map[string]int {
	var masterfile *os.File
	var err error
	if logging {
		header := "Ratio;"
		header += strings.Join(Labels, ";")

		outdir := "./data/dblp_str3"
		os.Mkdir(outdir, 0700)

		masterfilename := fmt.Sprintf("master_CUSTOM_%d", time.Now().Unix())
		fMasterFilename := fmt.Sprintf("%s/%s", outdir, masterfilename)
		masterfile, err = os.OpenFile(fMasterFilename, os.O_APPEND|os.O_WRONLY, 0755)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				masterfile, _ = os.Create(fMasterFilename)
				masterfile.WriteString(header)
				masterfile.WriteString(";Vertices;Edges;HittingSet;Opt;Time\n")
			} else {
				log.Fatalf("Could not open file %s", fMasterFilename)
			}
		}
	}

	execs := MakeExecs()
	expand := make(map[int32]bool)

	start := time.Now()

	PreProcessOnly(g, c, execs, expand)

	if len(g.Edges) == 0 {
		return execs
	}

	expDepth := 1

	gf := hypergraph.ExpandFrontier(g, expDepth, expand)

	for len(g.Edges) > 0 {
		fmt.Println(execs, len(g.Edges))
		expand := make(map[int32]bool)
		ApplyRulesSingle(gf, g, c, execs, expand)

		if len(expand) > 0 {
			gf = hypergraph.ExpandFrontier(g, expDepth, expand)
			continue
		} else {
			entry := findEntry(g, g, c, execs, expand)
			gf = hypergraph.ExpandFrontier(g, expDepth, expand)
			if entry {
				continue
			}
		}

		e := hypergraph.F3TargetLowDegreeDetect(g)
		if e == -1 {
			fmt.Println("No size 3 edge")
			continue
		}

		for v := range g.Edges[e].V {
			c[v] = true
		}

		gf = hypergraph.F3_ExpandFrontier(g, e, expDepth)
		execs["kFallback"] += 1
	}

	if logging {
		stop := time.Since(start).Seconds()

		msg := fmt.Sprintf("%f;", GetRatio(execs))
		for _, v := range Labels {
			msg += fmt.Sprintf("%d;", execs[v])
		}
		msg = msg[:len(msg)-1]
		masterfile.WriteString(fmt.Sprintf("%s;%d;%d;%d;%d;%.2f\n", msg, vSize, eSize, len(c), GetEstOpt(execs), RoundUp(stop, 2)))
	}

	return execs
}
