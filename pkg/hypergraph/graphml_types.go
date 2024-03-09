package hypergraph

import "encoding/xml"

type GraphMl struct {
	XMLName xml.Name `xml:"graphml"`
	Graph   MLGraph  `xml:"graph"`
}

type StdGraphMl struct {
	XMLName xml.Name   `xml:"graphml"`
	Graph   MLStdGraph `xml:"graph"`
}

type MLGraph struct {
	XMLName    xml.Name      `xml:"graph"`
	Nodes      []MLNode      `xml:"node"`
	HyperEdges []MLHyperEdge `xml:"hyperedge"`
}

type MLStdGraph struct {
	XMLName xml.Name    `xml:"graph"`
	Nodes   []MLStdNode `xml:"node"`
	Edges   []MLEdge    `xml:"edge"`
}

type MLNode struct {
	XMLName xml.Name `xml:"node"`
	Id      int32    `xml:"id,attr"`
	Data    MLData   `xml:"data"`
}

type MLStdNode struct {
	XMLName xml.Name `xml:"node"`
	Id      string   `xml:"id,attr"`
}

type MLData struct {
	XMLName xml.Name `xml:"data"`
	Key     string   `xml:"key,attr"`
	Value   int      `xml:",chardata"`
}

type MLHyperEdge struct {
	XMLName   xml.Name     `xml:"hyperedge"`
	Endpoints []MLEndpoint `xml:"endpoint"`
}

type MLEdge struct {
	XMLName xml.Name `xml:"edge"`
	Source  string   `xml:"source,attr"`
	Target  string   `xml:"target,attr"`
}

type MLEndpoint struct {
	XMLName xml.Name `xml:"endpoint"`
	Node    int32    `xml:"node,attr"`
}
