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

### Edge Hashing

```go
func getHash(arr []int32) uint32
```
Time Complexity: $n + n\cdot \log(n)$, where $n$ denotes the size of $\texttt{arr}$.

We start by sorting $\texttt{arr}$ with a $\textnormal{Quick-Sort}$-Algorithm. We then join the elements of $\texttt{arr}$ with the delimiter "|", producing a string of the form $|id_0|id_1|\dots|id_n|$. We then hash this string using the pure Go implementation of the [xxHash](https://github.com/Cyan4973/xxHash) hashing algorithm. It is probably the fastest non-cryptographic hash function to date. This will yield a $\texttt{uint32}$ hashsum that we return.

Whenever we refer to $\textit{the hash of an edge}$ we refer to the output of this function, using the endpoints of the edge as the $\texttt{arr}$ argument.

### Compute Subsets of Size $s$

```go
func getSubsetsRec(arr *[]int32, i int, n int, s int, data *[]int32, index int, subsets *list.List)
```
Time Complexity: $\binom{n}{s}$

Let us explain all the arguments first. 

- $\texttt{arr}$ is a pointer to an array. This array is basically the input set whose subsets we want to compute. 
- $\texttt{i}$ is an index over $\texttt{arr}$
- $\texttt{n}$ is the size of $\texttt{arr}$. 
- $\texttt{s}$ is the size of the computed subsets
- $\texttt{data}$ is a temporary array
- $\texttt{index}$ is an index over $\texttt{data}$
- $\texttt{subsets}$ is a list that will store all the subsets of size $\texttt{s}$ we find

The implementation looks as follows.

```go
func getSubsetsRec(arr *[]int32, i int, n int, s int, data *[]int32, index int, subsets *list.List){
    if index ==  s{
        subset := make([]int32, s)    
        for j:= 0; j < index; j++ {
            subset[j] = (*data)[j]
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

Given an array of integers and an integer target $t$, return indices of the two numbers such that they add up to $t$.

```go
func twoSum(items map[int32]int32, t int32) ([]int32, bool)
```
Time Complexity: $n$, where $n$ denotes the size of $\texttt{items}$.

Let us again explain the arguments:

- $\texttt{items}$ is a map that associates an integer with another integer
- $\texttt{t}$ is the target value
- $\texttt{return}$ values are a $\texttt{int32}$ array containing two keys $\texttt{a,b}$ of $\texttt{items}$, s.t. the sum of $\texttt{items[a]+items[b]}$ is equal to the target $\texttt{t}$ and a bool that indicates if such a pair was found.

The implementation looks as follows:
```go
func twoSum(items map[int32]int32, t int32) ([]int32, bool) { 
    lookup := make(map[int32]int32)

    for key, val := range items {
        if _, ex := lookup[t - val]; ex {
            return []int32{key, lookup[t - val]}, true
        } else {
            lookup[val] = key
        }
    }
    return nil, false
}
```
We start by creating a map called $\texttt{lookup}$. We then iterate other our entries $\texttt{items}$, each entry destructuring to a $\texttt{key,val}$ pair. We check if the entry $\texttt{lookup[t-val]}$ exists. 

- If the entry exist, we return an array containing $\texttt{val}$ and $\texttt{lookup[t-val]}$ 
- If the entry does not exist, we add a new entry to the lookup map with $\texttt{lookup[val] = key}$.

## Reduction Rules
The ususal signature of a reduction rule looks as follows:
```go
func NameRule(g HyperGraph, c map[int32]bool) int32
```
We take both a HyperGraph g and a *Set* c as arguments and mutate them. We then return the number of rule executions.

We prioritize time complexity over memory complexity when implementing rules, which does not equate to ignoring memory complexity completely.

#### Notation
We usually use the $\texttt{teletype font}$ to refer to variable names or code. We will often refer to the $\texttt{id}$ field of a vertex struct as $\texttt{v}$ and will use the notion of a vertex $v$ and the $\texttt{id}$ of a vertex interchangebly. We use $d$ to refer to the degree of the graph.

#### Executions

Reduction rules are usually meant to be applied exhaustively. We actually do not want to do that, because the algorithms for the rules rely on building up auxialliary data strucutures. Since rebuilding these structures takes a not negligeble amount of time, we want to reuse them as much as possible. 

Trying to apply a rule exhaustively in a single execution will introduce some difficulties:

- Some rules will mutate the graph, especially the set of edges. This can potentially create or eliminate a structure in the graph that a rule is targeting.
- Can the auxialliary datastructure be mutated to reflect these changes in linear time?

Luckily for us, most of these problems can be solved using Go. We usually have two possible ways of implementing a rule:

1. __Outer Loop:__ We wrap our main algorithm in an outer loop. We then apply the rule one time per loop iteration. The outer loop will break when the rule can't be applied anymore. This way we can reuse our auxialliary data structures.
2. __Pseudo Single Loop:__ We find all structures a rule is targeting in a single loop. We do this by iterating over the data structures and simultaniously deleting the targets, mutating the structure we are iterating over. This can possibly eliminate targets which we did not reach yet, or create new targets which we already passed. This requires the use of an outer loop. Since we can't decouple the deletion of the targets from the main algorithm, we inevitably introduce a quadratic time complexity in a worst case scenario. 

### Tiny/Small Edge Rule
- tiny edges: Delete all hyperedges of size one and place the corresponding vertices into the hitting set
- small edges: If $e$ is a hyperedge of size two, i.e., $e = \{ x,y \}$, then put both $x$
and $y$ into the hitting set.
```go
func RemoveEdgeRule(g HyperGraph, c map[int32]bool, t int) 
```

#### $O(|E|)$ Algortihm
Iterate over all edges of the graph and mark all edges of size $t$ in a set $\texttt{remEdges}$. We then iterate over all edges $\texttt{remEdges}$. We put the endpoints of the current edge $e$ in our partial solution. We then delete all edges that are edge adjacent to $e$ from $\texttt{remEdges}$ and our graph.

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

Specifically applied to $d=3$, this results in a time complexity of:

$$
\begin{align*}
T &= |E| \cdot 3\cdot \log (3) + (|E| \cdot (11 + 24 \cdot \log(3))) \\
&=|E| \cdot(3\cdot\log(3) + (11 + 24 \cdot \log(3)))
\end{align*}
$$

__Lemma.__ _We can either call this algorithm exhaustively until there are no more edge domination situations, or find all domination situations in a single execution of the algorithm._

_Proof._ Assume that our algorithm finds an edge domination situation. Then simply removing the dominated edge will not create nor eliminate an edge domination situation. It is therefore safe to remove all dominated edges in a single execution of the rule.

This also allows us to parallelize the main part of the algorithm, where we check each edge in our $dom$ set. We can achieve a speedup of $\approx 2$ on a six-core CPU an a pseudo-random graph with one million vertices and two million edges.

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

- $\texttt{vDeg map[int32]int32}:$ this map associates a vertex $v$ with $deg(v)$.
- $\texttt{vSubCount make(map[int32]map[int32]int32)}:$ this map associates a vertex $v$ with all other vertices that are vertex-adjacent to $v$. 

Both maps can be computed in time $|E|\cdot d^2$.

Example:

Let $E=\{ \{ 1,2,3 \}, \{ 1,2,4 \} \}$. Then $\texttt{vSub}$ and $\texttt{vSubCount}$ will look as follows,

$$
\begin{align*}
\texttt{vDeg} = \{\quad
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

We then iterate other $\texttt{vSubCount}$, we will refer to the current value in the iteration as $\texttt{vSubCount[v]}$. We then use a $\textnormal{Two-Sum}$-Algorithm to compute and return the first pair in $\texttt{vSubCount[v]}$, s.t. for the pair $\texttt{[a,b]}$ holds, $\texttt{vSubCount[v][a]} + \texttt{vSubCount[v][b]} = \texttt{vDeg[v]}+1$. If such a pair exists, then we conclude that for every edge $f$ such that $v \in f$, it holds for $f$ that, either $\texttt{a}\in f$ or $\texttt{b}\in f$.

**Lemma.** _The outlined procedure above is correct, under the assumption that the underlying graph does not contain any duplicate edges._

_Proof._ Let $G$ be a hypergraph. We first remove all edges of size one with the _Tiny Edge Rule_ and eliminate any duplicate edges. We then construct our two maps $\texttt{vSub}$ and $\texttt{vSubCount}$. Let $\texttt{v}$ be an entry in $\texttt{vSubCount}$ and $n=\texttt{vSub[v]}$. Now let $\texttt{sol=[a,b]}$ be the result of calling our $\textnormal{Two-Sum}$ implementation on $\texttt{vSubCount[v]}$ with a target sum of $n=\texttt{vDeg[v]}+1$. 

_Proposition._ If $\texttt{sol}$ is non-empty, then the edge $\{ \texttt{v,a,b} \}$ exists. 

Let $\texttt{sol=[a,b]}$ be the solution obtained by calling our $\textnormal{Two-Sum}$ algorithm on $\texttt{vSubCount[v]}$ with a target sum of $n=\texttt{vDeg[v]}+1$. For the sake of contradiction let us assume that the edge $\{ \texttt{v,a,b} \}$ does not exist. Since our graph does not contain duplicate edges and does not contain $\{ \texttt{v,a,b} \}$, there exist $\texttt{vDeg[v]}+1$ many edges that contain either $\texttt{a,v}$ or $\texttt{b,v}$. This however contradicts that there only exist $\texttt{vDeg[v]}$ many edges containing $\texttt{v}$. Therefore it must be, that the assumption that $\{ \texttt{v,a,b} \}$ does not exist, is false.

Since $\{ \texttt{v,a,b} \}$ exists, $\texttt{a}$ and $\texttt{b}$ can only occur $n-2 = \texttt{vDeg[v]} - 1$ times in other edges containing $\texttt{v}$. Since duplicate edges of  $\{ \texttt{v,a,b} \}$ can not exist, we know that every other edge containing $\texttt{v}$ also contains $\texttt{a}$ or $\texttt{b}$, but not both simultaniously. $\square$

We then add the two vertices in the solution to our partial solution $c$.

__Idea__: The initial idea for this algorithm involved the usage of an incidence matrix, where edges are identified by the rows and the vertices are identified by the columns. To check the Domination condition for a vertex $v$, the algortihm would select all edges/columns that contain $v$ and then add up the columns. Now let $n$ be the amount of edges containing $v$. If there exist two entries in the resulting column that have a combined value of $n+1$, then the rule applies for $v$ under the assumption that there are no duplicate edges. This would result in an algorithm with a time complexity of $|E|+|V|+|V|^2\cdot|E|$. 


### Approximative Double Vertex Domination Rule

- approximative double vertex domination: Assume there is a hyperedge $e =
\{x, y, a\}$ and another vertex $b$ such that, whenever $x$ or $y$ belong to some
hyperedge $h$, then $a$ or $b$ also belong to $h$. Then, we put $a$ and $b$ together
into the hitting set that we produce.

#### $\mathcal{O}(|E| + |V|\cdot |E|^2)$ Algorithm

We start by iterating over the edges of our graph and building an incidence list.

We then iterate over the edges again. If the current edge $e$ is not of size 3 we continue with the iteration. If $|e| = 3$ we iterate over the endpoints of $e$. We assign the variable $\texttt{a}$ to $v$ and create a map $\texttt{vCount map[int32]int32}$. For each vertex $w \in e\setminus v$ we look at the the edges incident to $w$, we then for each vertex $u$ of these edges, increment $\texttt{vCount[u]}$. We ignore edges that are incident to $\texttt{a}$ and do not care about vertices that are in $e$. We keep track of the amount of edges we considered with the variable $\texttt{xyCount}$.

If we find an entry $\texttt{vCount[b]}$, s.t. $\texttt{vCount[b] = xyCount}$, we found a double vertex domination situation in $e$.

### Small Triangle Rule
- small triangle situation: Assume there are three small hyperedges $e = \{y, z\}$, $f = \{x, y\}$, $g = \{x, z\}$. This describes a triangle situation $(e, f, g)$. Then, we put $\{x, y, z\}$ together into the hitting set, and we can even choose another hyperedge of size three to worsen the ratio.
```go
func SmallTriangleRule(g HyperGraph, c map[int32]bool)
```

#### $O(|E|+|V|^2)$ Algorithm
Again the quadratic exponent in $|V|$ looks scarier than it is. This happens beacause we have to remove the triangle.

We start by constructing an edge-adjacency list $\texttt{adjList}$ for all edges of size two. We then iterate over the entries of the list. For the current entry $\texttt{adjList[v]}$ we compute all subsets of size two of the entry. If both vertices of the subset are vertex-adjacent to each other, then we found a small triangle situation. If we find a triangle situation we put the corresponding vertices in our partial solution and alter the adjacency list to reflect these changes. We do this by iterating over all vertices that are vertex-adjacent to the triangle. For every vertex $w$ of these vertices we delete all vertices of the triangle from the entry $\texttt{adjList\lbrack}w\texttt{\rbrack}$.

This last step will introduce the quadratic complexity, since in the worst case, for a vertex $v$ in a triangle, there could exist $|V|$ many size two edges that contain $v$. This worst case occurs very rarely, which justifies using this quadratic algorithm. We could alternatively move the last step of the algorithm outside of the loop, and wrap both procedures with an outer loop which breaks if we dont find any more triangles. This simulates calling the rule exhaustively, while achieving a linear time complexity.

## Self-Monitoring
Each of the reduction rule functions returns a $\texttt{int32}$ value, which indicates the number of rule executions. We can then use these values to calculate the local ratio.

## Testing

Every reduction rule is tested for their correctness with unit tests. We create small graphs in these tests, that contain structures, which the rules are targeting. We then test for the elements in the partial solution and the amount of edges left in the graph.

## Possible Optimizations

Right now we sometimes use the Go $\texttt{apppend}$ function, to add elements to a slice. Calls to $\texttt{apppend}$ are more expensive than an ordinary assignment to a fixed size slice. 

Using fixed size slices could improve performance slightly in cases where we know the maximum amount of elements we will add to a slice.
