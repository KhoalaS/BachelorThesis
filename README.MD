# Selbstüberwachende Approximationsalgorithmen | Implementierung und Evaluierung eines 3-HS Algorithmus

## Reduction Rules
The ususal signature of a reduction rule looks as follows:
```go
func NameRule(g HyperGraph, c map[int32]bool) (HyperGraph, map[int32]bool)
```
We take both a HyperGraph g and a *Set* c as arguments and return both a HyperGraph and a *Set*.

We did not implement these rules as member functions to the HyperGraph struct, since we want to preserve the old graph instances. We do this because these rules will be called in a Search-Tree algorithm. This could impact the performance heavily. A better solution would be compute the current graph based on the partial solution c. It is sufficient for small toy examples as of now.



## Datastructures

### Vertex

```go
type Vertex struct {
	id int32
	data any
}
```
The $\texttt{Vertex}$ datatype has two fields. The field $\texttt{id}$ is an arbitrary identifier and $\texttt{data}$ serves as a placeholder for actual data associated with the vertex.


### Edge

```go
type Edge struct {
	v map[int32]bool
}
```
The $\texttt{Edge}$ datatype has one field. The field $\texttt{v}$ is a map with keys of type $\texttt{int32}$ and values of type $\texttt{bool}$. 

When working with the endpoints of an edge, we are usually not interested in the associated values, since we never mutate the edges.

This simulates a $\texttt{Set}$ datatype while allowing faster access times than simple arrays/slices.

## Hypergraph
```go
type HyperGraph struct {
	Vertices map[int32]Vertex
	Edges map[int32]Edge
}
```
The $\texttt{HyperGraph}$ datatype has two field. Both fields $\texttt{Vertices}$ and $\texttt{Edges}$ are maps with keys of type $\texttt{int32}$ and values of type $\texttt{Vertex}$ and $\texttt{Edge}$ respectively.

We chose this $\texttt{Set}$-like datastructure over lists again because of faster access times, but also operations that remove edges/vertices are easier.
