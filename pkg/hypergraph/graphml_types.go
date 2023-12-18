package hypergraph

import "encoding/xml"

type GraphMl struct {
	XMLName xml.Name `xml:"graphml"`
	Graph MLGraph `xml:"graph"`
}

type MLGraph struct {
	XMLName xml.Name `xml:"graph"`
	Nodes []MLNode `xml:"node"`
	Edges []MLHyperEdge `xml:"hyperedge"`
}

type MLNode struct {
	XMLName xml.Name `xml:"node"`
	Id int32 `xml:"id,attr"`
	Data MLData `xml:"data"`
}

type MLData struct {
	XMLName xml.Name `xml:"data"`
	Key string `xml:"key,attr"`
	Value any `xml:",chardata"` 
}

type MLHyperEdge struct {
	XMLName xml.Name `xml:"hyperedge"`
	Endpoints []MLEndpoint `xml:"endpoint"`
}

type MLEndpoint struct {
	XMLName xml.Name `xml:"endpoint"`
	Node int32 `xml:"node,attr"`
}

