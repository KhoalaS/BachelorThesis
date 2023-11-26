# Self-monitoring Approximation Algorithms | Implementation of a 3-HS Algorithm

## Programming Language
We chose the [Go](https://go.dev/) language. It is a statically typed, compiled language, with a C-like syntax. It is using a garbage collector to handle memory management, which at first seems off-putting for an application like this. Since Go's GC is very efficient, we are not worried about that fact. If the situation arises in which we do need to handle memory manually, we can utilize Go's $\texttt{unsafe}$ package in conjunction with C-interop.

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

### Hypergraph
```go
type HyperGraph struct {
	Vertices map[int32]Vertex
	Edges map[int32]Edge
	edgeCounter int32
	Degree int32
}
```
The $\texttt{HyperGraph}$ datatype has four fields. Both fields $\texttt{Vertices}$ and $\texttt{Edges}$ are maps with keys of type $\texttt{int32}$ and values of type $\texttt{Vertex}$ and $\texttt{Edge}$ respectively.

We chose this $\texttt{Set}$-like datastructure over lists again because of faster access times, but also operations that remove edges/vertices are built-in to the map type.

The field $\texttt{edgeCounter}$ is an internal counter used to assign ids to added edges. 

The field $\texttt{Degree}$ specifies the maximum degree of the graph. 

## Misc. Algorithms/Utilities

### IdValueHolder Struct
```go
type IdValueHolder struct {
	Id    int32
	Value int32
}
```
The struct has two fields. The field $\texttt{Id}$ is an id and the field $\texttt{Value}$ holds a $\texttt{int32}$ value. We use this struct in out custom $\textnormal{Two-Sum}$ implementation, we will elaborate further in the $\textnormal{Two-Sum}$ section.

### getHash

```go
func getHash(arr []int32) uint32
```
Time Complexity: $n + n\cdot \log(n)$, where $n$ denotes the size of $\texttt{arr}$.

We start by sorting $\texttt{arr}$ with a $\textnormal{Quick-Sort}$-Algorithm. We then join the elements of $\texttt{arr}$ with the delimiter "|", producing a string of the form $|id_0|id_1|\dots|id_n|$. We then hash this string using the pure Go implementation of the [xxHash](https://github.com/Cyan4973/xxHash) hashing algorithm. It is probably the fastest non-cryptographic hash function to date. This will yield a $\texttt{uint32}$ hashsum that we return.

Whenever we refer to $\textit{the hash of an edge}$ we refer to the output of this function, using the endpoints of the edge as the $\texttt{arr}$ argument.

### getSubsetsRec
```go
func getSubsetsRec(arr *[]int32, i int, n int, s int, data *[]int32, index int, subsets *list.List)
```
Time Complexity: $\binom{n}{r}$

Let us explain all the arguments first. 
- $\texttt{arr}$ is a pointer to an array. This array is basically the input set whose subsets we want to compute. 
- $\texttt{i}$ is used to keep track of the recursion depth
- $\texttt{n}$ is the size of $\texttt{arr}$. 
- $\texttt{s}$ is the size of the computed subsets
- $\texttt{data}$ 
- $\texttt{index}$ is the current position, on which we want to insert a value into $\texttt{data}$
- $\texttt{subsets}$ is a list that will store all the subsets of size $\texttt{s}$ we find

The implementation looks as follows.

```go
func getSubsetsRec(arr *[]int32, i int, n int, s int, data *[]int32, index int, subsets *list.List){
    if index ==  s{
        subset := make([]int32, s)    
        for j:= 0; j < index; j++ {
            subset[j] = *data[j]
        }
        subsets.PushBack(subset)
        return
    }

    if i >= n {
        return
    }

    (*data)[index] = (*arr)[i]
    
    getSubsetsRec(arr, i+1, n, s, data, index+1, subsets)
    getSubsetsRec(arr, i+1, n, s, data, index, subsets)
}
```
The algorithm is pretty simple. 

1. We either put the $i^{th}$ element of $\texttt{arr}$ into our temporary array $\texttt{data}$ and increment $\texttt{index}$ and $\texttt{i}$.
2. Or we only increment $\texttt{i}$.

We have two base conditions.

1. If $\texttt{index}$ is equal to $\texttt{s}$, then we copy the contents of $\texttt{data[0:s]}$ into an array and push it onto $\texttt{subsets}$.
2. If $\texttt{i}$ is greater than or equal $\texttt{n}$. 

Lists in Go are not very memory efficient, but since we exclusively call this function with $\texttt{arr}$ representing the vertices in an edge, the value $\texttt{n}$ is usually fixed at 3. The raised memory problems occur at values of $\texttt{n}\geq 10000$, thus justifying the continued usage of lists.

### Two-Sum
#### $\mathcal{O}(n)$ Algorithm

## Reduction Rules
The ususal signature of a reduction rule looks as follows:
```go
func NameRule(g HyperGraph, c map[int32]bool)
```
We take both a HyperGraph g and a *Set* c as arguments and mutate them.

We prioritize time complexity over memory complexity when implementing rules, which does not equate to ignoring memory complexity completely.

### Notation
We usually use the $\texttt{terminal typeface}$ to refer to variable names or code. We will often refer to the $\texttt{id}$ field of a vertex struct as $\texttt{v}$ and will use the notion of a vertex $v$ and the $\texttt{id}$ of a vertex interchangebly. We use $d$ to refer to the degree of the graph.

### Tiny/Small Edge Rule
- tiny edges: Delete all hyperedges of size one and place the corresponding vertices into the hitting set
- small edges: If $e$ is a hyperedge of size two, i.e., $e = \{ x,y \}$, then put both $x$
and $y$ into the hitting set.
```go
func RemoveEdgeRule(g HyperGraph, c map[int32]bool, t int) 
```

#### $O(|E|)$ Algortihm
Iterate over all edges of the graph and mark all edges of size $t$. Then put all endpoints of these edges into the partial solution and remove the marked edges and all edges that are edge-adjacent to marked edges. 

### Edge Domination Rule
- (hyper)edge domination: A hyperedge $e$ is _dominated_ by another hyperedge $f$ if $f\subset e$. In that case, delete $e$.

```go
func EdgeDominationRule(g HyperGraph, c map[int32]bool)
```

#### $O(|E|)$ Algorithm

We partition our set of edges into two disjoint sets $sub$ and $dom$. The set $dom$ will contain edges that could possibly be dominated. The set $sub$ will contain hashes of edges $e$ that could dominate another edge. We compute these hashes as follows.

First we store the ids of the endpoints of $e$ in an array. We then use the built-in $\texttt{sort}$ package to sort the array with a ${\textit{Quick-Sort}}$ algorithm. We then compute a string $\texttt{in}$ by joining the ids with a non-numerical delimiter like "|". We then use a [hash function](https://github.com/OneOfOne/xxhash) to obtain a hashsum of type $\texttt{uint32}$.

We then iterate over the set $dom$ and compute every strict subset of the current edge $f$. For each of these subsets, we test if the "$\textit{hash of the subset}$" is present in our set $sub$. If it is then $f$ is dominated by another edge. 

The exact time complexity is as follows:

$$
\begin{align*}
T &= |E| \cdot d\cdot \log (d) + (|E| \cdot (d + 2^d + (2^d \cdot d \cdot \log(d)))) \\
\end{align*}
$$

where $d$ denotes the degree of the graph.

##### Remarks
We technically dont have to hash the string $\texttt{in}$ for our purposes. But we can possibly improve the lookup speed of the set $sub$ if we do hash the keys first.

### Approximative Vertex Domination Rule
- approximative vertex domination: Assume there is a hyperedge $e = \{ x,y,z \}$ such that, whenever $x$ belongs to some hyperedge $h$, then $y$ or $z$ also belong to $h$. Then, we put $y$ and $z$ together into the hitting set that we produce.
```go
func ApproxVertexDominationRule(g HyperGraph, c map[int32]bool)
```

#### $O(|E| + |V|^2)$ Algorithm
The quadratic exponent in $|V|$ looks scary at first, but will only occur in the case if there exists a vertex $v$ s.t. for every other vertex $w$ there exists an edge $e$ with $\{ v,w \} \subset e$.

We first construct two maps:

- $\texttt{vSub map[int32]int32}:$ this map associates a vertex $v$ with the amount of edges containing the $v$.
- $\texttt{vSubCount make(map[int32]map[int32]int32)}:$ this map associates a vertex $v$ with all other vertices that are vertex-adjacent to $v$. 

Both maps can be computed in time $|E|\cdot d^2$.

Example:

Let $E=\{ \{ 1,2,3 \}, \{ 1,2,4 \} \}$. Then $\texttt{vSub}$ and $\texttt{vSubCount}$ will look as follows,

$$
\begin{align*}
\texttt{vSub} = \{\quad
 &1: 2,\\
 &2: 2, \\
 &3: 1, \\	
 &4: 1 \quad \}
\end{align*}
\newline

\begin{align*}
\texttt{vSubCount} = \{\quad
 &1: \{ 2: 2, \quad 3:1, \quad 4:1\},\\
 &2: \{ 1: 2, \quad 3:1, \quad 4:1\}, \\
 &3: \{ 1: 1, \quad 2:1\}, \\
 &4: \{ 1: 1, \quad 2:1\} \quad \}	
\end{align*}
$$

We then iterate other $\texttt{vSubCount}$, we will refer to the current value in the iteration as $\texttt{vSubCount[v]}$. We then use a $\textnormal{Two-Sum}$-Algorithm to compute and return the first pair in $\texttt{vSubCount[v]}$, s.t. for the pair $\texttt{[a,b]}$ holds, $\texttt{a} + \texttt{b} = \texttt{vSub[v]}+1$. If such a pair exists, then we conclude that for every edge $f$ such that $v \in f$, it holds for $f$ that, either $\texttt{a}\in f$ or $\texttt{b}\in f$.

**Lemma.** The outlined procedure above is correct, under the assumption that the underlying graph does not contain any duplicate edges.

_Proof._ Let $G$ be a hypergraph. We first remove all edges of size one with the _Tiny Edge Rule_ and eliminate any duplicate edges. We then construct our two maps $\texttt{vSub}$ and $\texttt{vSubCount}$. Let $\texttt{v}$ be an entry in $\texttt{vSubCount}$ and $n=\texttt{vSub[v]}$. Now let $\texttt{sol=[a,b]}$ be the result of calling our $\textnormal{Two-Sum}$ implementation on $\texttt{vSubCount[v]}$ with a target sum of $n=\texttt{vSub[v]}+1$. 

_Proposition._ If $\texttt{sol}$ is non-empty, then the edge $\{ \texttt{v,a,b} \}$ exists. 

Let $\texttt{sol=[a,b]}$ be the solution obtained by calling our $\textnormal{Two-Sum}$ algorithm on $\texttt{vSubCount[v]}$ with a target sum of $n=\texttt{vSub[v]}+1$. For the sake of contradiction let us assume that the edge $\{ \texttt{v,a,b} \}$ does not exist. Since our graph does not contain duplicate edges and does not contain $\{ \texttt{v,a,b} \}$, there exist $\texttt{vSub[v]}+1$ many edges that contain either $\texttt{a,v}$ or $\texttt{b,v}$. This however contradicts that there only exist $\texttt{vSub[v]}$ many edges containing $\texttt{v}$. Therefore it must be, that the assumption that $\{ \texttt{v,a,b} \}$ does not exist, is false.

Since $\{ \texttt{v,a,b} \}$ exists, $\texttt{a}$ and $\texttt{b}$ can only occur $n-2 = \texttt{vSubCount[v]} - 1$ times in other edges containing $\texttt{v}$. Since duplicate edges of  $\{ \texttt{v,a,b} \}$ can not exist, we know that every other edge containing $\texttt{v}$ also contains $\texttt{a}$ or $\texttt{b}$, but not both simultaniously. $\square$

We then add the two vertices in the solution to our partial solution $c$.

[__MULTI PASS LOOP__]

__Idea__: The initial idea for this algorithm involved the usage of an incidence matrix, where edges are identified by the rows and the vertices are identified by the columns. To check the Domination condition for a vertex $v$, the algortihm would select all edges/columns that contain $v$ and then add up the columns. Now let $n$ be the amount of edges containing $v$. If there exist two entries in the resulting column that have a combined value of $n+1$, then the rule applies for $v$ under the assumption that there are no duplicate edges. This would result in an algorithm with a time complexity of $|E|+|V|+|V|^2\cdot|E|$. 


### Small Triangle Rule
- small triangle situation: Assume there are three small hyperedges $e = \{y, z\}$, $f = \{x, y\}$, $g = \{x, z\}$. This describes a triangle situation $(e, f, g)$. Then, we put $\{x, y, z\}$ together into the hitting set, and we can even choose another hyperedge of size three to worsen the ratio.
```go
func SmallTriangleRule(g HyperGraph, c map[int32]bool)
```

#### $O(|E|+|V|^2)$ Algorithm
Again the quadratic exponent in $|V|$ looks scarier than it is. This happens beacause we have to remove the triangle.

We start by constructing an adjacency list $\texttt{adjList}$ for all edges of size two. We then iterate over the entries of the list. For the current entry $\texttt{adjList[v]}$ we compute all subsets of size two of the entry. If both vertices of the subset are vertex-adjacent to each other, then we found a small triangle situation. If we find a triangle situation we put the corresponding vertices in our partial solution and alter the adjacency list to reflect these changes. We do this by iterating over all vertices that are vertex-adjacent to the triangle. For every vertex $w$ of these vertices we delete all vertices of the triangle from the entry $\texttt{adjList\lbrack}w\texttt{\rbrack}$.

This last step will introduce the quadratic complexity, since in the worst case, for a vertex $v$ in a triangle, there could exist $|V|$ many size two edges that contain $v$. In randomly generated graphs, a triangle situation occurs very rarely, which justifies using this quadratic algorithm. We could alternatively move the last step of the algorithm outside of the loop, and wrap both procedures with an outer loop which breaks if we dont find any more triangles. This simulates calling the rule exhaustively, while achieving a linear time complexity.

## Possible Optimizations

Right now we sometimes use the Go $\texttt{apppend}$ function, to add elements to a slice. Calls to $\texttt{apppend}$ are more expensive than an ordinary assignment to a fixed size slice. 

Using fixed size slices could improve performance slightly in cases where we know the maximum amount of elements we will add to a slice.
