# Selbstüberwachende Approximationsalgorithmen | Implementierung und Evaluierung eines 3-HS Algorithmus

## Reduction Rules
The ususal signature of a reduction rule looks as follows:
```go
func NameRule(g HyperGraph, c map[int32]bool)
```
We take both a HyperGraph g and a *Set* c as arguments and mutate them.

We could also implement these rules as member functions of the HyperGraph struct.

### Edge Domination Rule
```go
func EdgeDominationRule(g HyperGraph, c map[int32]bool)
```
> A hyperedge $e$ is dominated by another hyperedge
$f$ if $f \subset e$. In that case, delete $e$.

The naive approach would be to compare every edge $e$ with eachother and test if $e$ is a subset of the other edge $f$. This results in a quadratic time complexity of $O(|E|^2\cdot d)$.

There are little optimizations possible, like checking the number vertices of the edges, but obviously the upper bound remains the same.

A better optimization can be achieved by partitioning the set of edges into two disjoint subsets $sub$ and $dom$. While the set $sub$ only contains edges that are possibly subsets of other edges, the set $dom$ only contains edges that can be dominated. Since we assume that edges of degree 1 are already removed from our graph, we know that $dom=\left\{ e| e\in E \text{ with } |e| =3 \right\}$.

This results in a time complexity of $r_s\cdot r_d \cdot|E|^2 \cdot d$, where $r_s$ denotes the ratio $\frac{|sub|}{|E|}$ and $r_d$ denotes $\frac{|dom|}{|E|}$. This is slightly better but still close to the original complexity if $|r_s - r_d|$ is big.

We therefore employ paralellism to solve the problem in somewhat feasable time. The $sub$ set is partitioned into $\texttt{nCPU}$ many subsets of (near) equal size, where $\texttt{nCPU}$ denotes the number of physical CPU cores of the system. We then solve the problem in parallel for all of these subsets.

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
	edgeCounter int32
	Degree int32
}
```
The $\texttt{HyperGraph}$ datatype has four field. Both fields $\texttt{Vertices}$ and $\texttt{Edges}$ are maps with keys of type $\texttt{int32}$ and values of type $\texttt{Vertex}$ and $\texttt{Edge}$ respectively.

We chose this $\texttt{Set}$-like datastructure over lists again because of faster access times, but also operations that remove edges/vertices are easier.

The field $\texttt{edgeCounter}$ is an internal counter used to assign ids to added edges. 

The field $\texttt{Degree}$ specifies the maximum degree of the graph. 
