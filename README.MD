# Selbstüberwachende Approximationsalgorithmen | Implementierung und Evaluierung eines 3-HS Algorithmus

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

#### $O(|E|\cdot d^3)$ Algorithm

We partition our set of edges into two disjoint sets $sub$ and $dom$. The set $dom$ will contain edges that could possibly be dominated. The set $sub$ will contain hashes of edges $e$ that could dominate another edge. We compute these hashes as follows.

First we store the ids of the endpoints of $e$ in an array. We then use the built-in $\texttt{sort}$ package to sort the array with a ${\textit{Quick-Sort}}$ algorithm. We then compute a string $\texttt{in}$ by joining the ids with a non-numerical delimiter like "|". We then use a [hash function](https://github.com/OneOfOne/xxhash) to obtain a hashsum of type $\texttt{uint32}$.

We then iterate over the set $dom$ and compute every strict subset of the current edge $f$. For each of these subsets, we test if the "$\textit{hash of the subset}$" is present in our set $sub$. If it is then $f$ is dominated by another edge. 

The exact time complexity is as follows:
$$
\begin{align*}
T &= |E| + (|E| \cdot \left(((d-1) \cdot (d\cdot(d-1))) + \sum_{i=1}^{d-1}{\binom{d}{i}}\right)) \\
&< |E| + |E| \cdot \left(d^3 + \sum_{i=1}^{d-1}{\binom{d}{i}}\right)
\end{align*}
$$

##### Remarks
We technically dont have to hash the string $\texttt{in}$ for our purposes. But we can possibly improve the lookup speed of the set $sub$ if we do hash the keys first.

### Approximative Vertex Domination Rule
```go
func ApproxVertexDominationRule(g HyperGraph, c map[int32]bool)
```

#### $O(|E|+ |V|^2)$
The quadratic exponent $|V|$ looks scary at first, but will only occur in the case if there exists a vertex $v$ s.t. for every other vertex $w$ there exists an edge $e$ with $\{ v,w \} \subset e$.


First we want to associate every endpoint of an edge $e$ with the other endpoints in $e$. We do this by iterating over all edges in our graph $G$ and effectively computing all subsets of $e$ with size $|e|-1$. For each of these subsets $e_{sub}$ we associate the vertex $w \in e\setminus e_{sub}$ with the hash of $e_{sub}$. We also keep track of the amount of times other vertices are adjacent to $w$. 

By following this procedure we obtain two maps, the first one called $\texttt{vSub}$ associates a vertex $v$ with the hashes of subsets of edges that result in removing $v$ from edges containing $v$. The second map $\texttt{vSubCount}$ associates a vertex $v$ with all other vertices that are adjacent to $v$. Note that $\texttt{vSubCount}$ also keeps track of the number of times a vertex is adjacent to $v$.

Example:

Let $E=\{ \{ 1,2,3 \}, \{ 1,2,4 \} \}$. Then $\texttt{vSub}$ and $\texttt{vSubCount}$ will look as follows,

$$
\begin{align*}
\texttt{vSub} = \{\quad
 &1: h(2,3), h(2,4)\\
 &2: h(1,3), h(1,4) \\
 &3: h(1,2) \\	
 &4: h(1,2) \quad \}
\end{align*}
\newline

\begin{align*}
\texttt{vSubCount} = \{\quad
 &1: \{ 2: 2, \quad 3:1, \quad 4:1\}\\
 &2: \{ 1: 2, \quad 3:1, \quad 4:1\} \\
 &3: \{ 1: 1, \quad 2:1\} \\
 &4: \{ 1: 1, \quad 2:1\} \quad \}	

\end{align*}
$$

where $h(x,y)$ denotes the hash of an edge with endpoints $x,y.$

We then iterate other $\texttt{vSubCount}$, we will refer to the current value in the iteration as $\texttt{vSubCount[v]}$. Using the $\texttt{IdValueHolder}$ struct (a simple struct that holds an id and a value field) we store all entries of our current value $\texttt{vSubCount[v]}$, in an array $\texttt{arr}$. We then use a $\textnormal{Two-Sum}$-Algorithm to compute all pairs in $\texttt{arr}$, s.t. for each pair $(\texttt{a},\texttt{b})$ holds, $\texttt{a.value} + \texttt{b.value} = \texttt{len(vSub[v])}+1$. We then check for each of these pairs if $\texttt{vSub[v]}\texttt{\lbrack}h(\texttt{a.id},\texttt{b.id})\texttt{\rbrack}$ exists. If such a pair exists, then we conclude that for every edge $f$ such that $v \in f$, it holds for $f$ that, either $\texttt{a.id}\in f$ or $\texttt{b.id}\in f$.

**Lemma.** The outlined procedure above is correct, under the assumption that the underlying graph does not contain any duplicate edges.

_Proof._ Let $G$ be a hypergraph. We first remove all edges of size one with the _Tiny Edge Rule_ and eliminate any duplicate edges. We then construct our two maps $\texttt{vSub}$ and $\texttt{vSubCount}$. Let $\texttt{v}$ be an entry in $\texttt{vSubCount}$ and $n$ be the length of $\texttt{vSub[v]}$. Now let $\texttt{solutions}$ be the result of calling our $\textnormal{Two-Sum}$ implementation on the array version of $\texttt{vSubCount[v]}$. For every solution $sol \in \texttt{solutions}, sol = (\texttt{a},\texttt{b})$ we can consider two possible cases:

Let $hash$ be $h(\texttt{a}, \texttt{b})$.

__Case 1:__ $\texttt{vSub[v]\lbrack}hash\texttt{\rbrack}$ exists. Hence the edge $\{ v,a,b \}$ does exist. Let $occur$ be $\texttt{|vSub[v]|}+1$, the amount of edges which contain $\texttt{v}$ plus one. Since we used this $occur$ value as our target sum in the $\textnormal{Two-Sum}$ call, the following holds for $\texttt{a}$ and $\texttt{b}$, 

$$\texttt{a.value} + \texttt{b.value} = occur$$

Since we know that the edge $\{ v,a,b \}$ exists, we can conclude that $\texttt{a}$ or $\texttt{b}$ are occuring $occur-2 = |\texttt{vSub[v]}|-1$ times in other edges containing $\texttt{v}$. Under the assumption that we dont have duplicate edges, $\texttt{a}$ or $\texttt{b}$ can't occur in the same edge. Therfore all edges with endpoint $\texttt{v}$ either contain $\texttt{a}$ or $\texttt{b}$.

__Case 2:__ $\texttt{vSub[v]\lbrack}hash\texttt{\rbrack}$ does not exist. The edge $\{ v,a,b \}$ does not exist.

$\square$

We are only interested in one valid solution per entry in $\texttt{vSub[v]}$. We then add the two vertices to our partial solution $c$.

The only problem right now is that we have to call this rule exhaustively, since putting vertices in our partial solution mutates the graphs edges. 
## Possible Optimizations

Right now we sometimes use the Go $\texttt{apppend}$ function, to add elements to a slice. Calls to $\texttt{apppend}$ are more expensive than an ordinary assignment to a fixed size slice. 

Using fixed size slices could improve performance slightly in cases where we know the maximum amount of elements we will add to a slice.
