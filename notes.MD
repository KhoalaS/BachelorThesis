---
    title: Implementation and evaluation of a self-monitoring approximation algorithm for 3-Hitting-Set
    link-citations: true

...
\RestyleAlgo{ruled}
\DontPrintSemicolon
\SetAlgoVlined
\LinesNumbered

# Programming Language

We chose the [Go](https://go.dev/) language. It is a statically typed, compiled language, with a C-like syntax. It is using a garbage collector to handle memory management, which at first seems off-putting for an application like this. Since Go's GC is very efficient, we are not worried about that fact. If the situation arises in which we do need to handle memory manually, we can utilize Go's $\texttt{unsafe}$ package in conjunction with C-interop.

# Datastructures

## Vertex

```go
type Vertex struct {
	id int32
	data any
}
```

The $\texttt{Vertex}$ datatype has two fields. The field $\texttt{id}$ is an arbitrary identifier and $\texttt{data}$ serves as a placeholder for actual data associated with the vertex.

## Edge

```go
type Edge struct {
	v map[int32]bool
}
```

The $\texttt{Edge}$ datatype has one field. The field $\texttt{v}$ is a map with keys of type $\texttt{int32}$ and values of type $\texttt{bool}$.

When working with the endpoints of an edge, we are usually not interested in the associated values, since we never mutate the edges.

This simulates a _Set_ datatype while allowing faster access times than simple arrays/slices.

## Hypergraph

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

# Misc. Algorithms/Utilities

## Edge Hashing

```go
func getHash(arr []int32) string
```

Time Complexity: $n + n\cdot \log(n)$, where $n$ denotes the size of $\texttt{arr}$.

We start by sorting $\texttt{arr}$ with a $\textsc{Quick-Sort}$-Algorithm. We then join the elements of $\texttt{arr}$ with the delimiter $"|"$, returning a string of the form $"|id_0|id_1|\dots|id_n|"$.

Whenever we refer to $\textit{the hash of an edge}$ we refer to the output of this function, using the endpoints of the edge as the $\texttt{arr}$ argument.

## Compute Subsets of Size $s$

\todo{rewrite with new algorithm}

```go
func getSubsetsRecMain(arr *[]int32, i int, n int, s int, data *[]int32, index int, subsets *list.List)
```

Time Complexity: $2^n$

Let us explain all the arguments first.

-   $\texttt{arr}$ is a pointer to an array. This array is basically the input set whose subsets we want to compute.
-   $\texttt{i}$ is an index over $\texttt{arr}$
-   $\texttt{n}$ is the size of $\texttt{arr}$.
-   $\texttt{s}$ is the size of the computed subsets
-   $\texttt{data}$ is a temporary array
-   $\texttt{index}$ is an index over $\texttt{data}$
-   $\texttt{subsets}$ is a list that will store all the subsets of size $\texttt{s}$ we find

The implementation looks as follows.

```go
func getSubsetsRecMain(arr []int32, i int, n int, s int, data []int32, index int, subsets *list.List) {
	if index == s {
		subset := make([]int32, s)
		for j := 0; j < index; j++ {
			subset[j] = data[j]
		}
		subsets.PushBack(subset)
		return
	}

	if i >= n {
		return
	}

	data[index] = arr[i]

	getSubsetsRecMain(arr, i+1, n, s, data, index+1, subsets)
	getSubsetsRecMain(arr, i+1, n, s, data, index, subsets)
}
```

The algorithm is pretty simple. For every element $x$ in $\texttt{arr}$, we either put $x$ into $\texttt{data}$ or we do not. We have two base conditions.

1. If $\texttt{index} = \texttt{s}$, then we copy the contents of $\texttt{data}$ into an array and add it to $\texttt{subsets}$.
2. If $\texttt{i} \geq \texttt{n}$, then we have considered all elements for the current path of the search tree.

Lists in Go are not very memory efficient, but since we exclusively call this function with $\texttt{arr}$ representing the vertices in an edge, the value $\texttt{n}$ is usually fixed at 3. The raised memory problems occur at values of $\texttt{n}\geq 10000$, thus justifying the continued usage of lists.

For the case where we have to compute a lot of subsets, we provide a slightly different version of this function. Instead of passing in the $\texttt{subsets}$ list, we pass in a callback function that shall be called whenever we find a subset, using the found subset as an argument.

## Two-Sum

Given an array of integers and an integer target $t$, return indices of the two numbers such that they add up to $t$.

Time Complexity: $n$, where $n$ denotes the size of $\texttt{items}$.

$$
\begin{algorithm}[H]
\nllabel{label}

\KwIn{An array of integers $arr$, a target value $t$}
\KwOut{Two indices $a,b$, such that $arr[a] + arr[b] = t$, a boolean indicating if a solution was found}
\BlankLine

$lookup \gets \texttt{map}[\mathbb{N}]\mathbb{N}$\;

\For{$i\gets 0$ \KwTo $len(arr)$}{
    \lIf{$lookup[t-arr[i]]$ exists}{
        \KwRet $(i, lookup[t-arr[i]]), true$\;
    }
    \Else{
        $lookup[arr[i]] \gets i$\;
    } 
}
\KwRet $nil, false$\;

\caption{An algorithm for the \textssc{Two-Sum} problem}
\end{algorithm}
$$

We start by creating a map called $lookup$. We then iterate other the elements of $arr$, checking if the entry $lookup[t-arr[i]]$ exists.

-   If the entry exists, we return a pair $(i, lookup[t-arr[i]])$ and the boolean value  $true$ since we found a solution.
-   If the entry does not exist, we add a new entry to the $lookup$ map using $arr[i]$ as key and $i$ as value.

If no solution was found, we return $nil$ and the boolean value $false$.   


This algorithm is an ingredient for the implementation of one of the reduction rules, specifically the Approximative Vertex Domination Rule. The actual implementation is accepting a map instead of an array as its first parameter.

# Hypergraph Models

## First Testing Model

```go
func GenerateTestGraph(n int32, m int32, tinyEdges bool) *HyperGraph
```

Let us explain the arguments first:

-   $\texttt{n}$ are the number of vertices the graph will have
-   $\texttt{m}$ is the amount of edges the graph will at most have
-   $\texttt{tinyEdges}$ when $\texttt{false}$ indicates that we do not want to generate edges of size 1.

We use a very naive approach for generating (pseudo-)random graphs.

We first create an empty Hypergraph struct and add $\texttt{n}$ many vertices to that graph.

We then compute a random $\texttt{float32}$ value $\texttt{r}$ in the half-open interval $[ 0.0, 1.0 )$. This value will be used to determine the size of an edge $e$. The edges are distributed based on their size as follows:

$$
size(\texttt{r})=
\begin{cases}
    1 & \quad \texttt{r} < 0.01\\
    2 & \quad 0.01 \leq \texttt{r} < 0.60\\
    3 & \quad \text{else}
\end{cases}
$$

We then store the result of $size(\texttt{r})$ in a variable $\texttt{d}$. We then randomly pick vertices in the half-open interval $[ 0, \texttt{n} )$, until we have picked $\texttt{d}$ many distinct vertices. We then check if an edge containing these vertices already exists. If it does not exist we add it to our graph.

That results in the graph having at most $\texttt{m}$ edges and not exactly $\texttt{m}$, since we did not want to artificially saturate the graph with edges.

Using an established model like the [Erdős–Rényi Model](https://www.semanticscholar.org/paper/On-random-graphs.-I.-Erdos-R%C3%A9nyi/e77e986f717b353329e2202c384c5499191a0646) would be more favourable. But since we just needed something tangible to start testing, this "model" will suffice for the time being.

One could also look into generating random bipartite graphs that translate back to a hypergraph with the desired vertex and edge numbers.

## Preferential Attachment Hypergraph Model

In the Preferential Attachment Model, one will add edges to an existing graph, with a probability proportional to the degree of the endpoints of that edge. This edge will either contain a newly added vertex, or will be comprised of vertices already part of the graph.

We will use an implementation by Antelmi et al. as reference [@SimpleHypergraphs.js], which is part of their work on _SimpleHypergraphs.jl_ [@DBLP:journals/im/SpagnuoloCSPSKA20], a hypergraph software library written in the Julia language. The implementation is based on a preferential attachment model proposed by Avin et al. in [@DBLP:journals/corr/AvinLP15].

```go
func GeneratePrefAttachmentGraph(n int, p float64, maxEdgesize int32)
```

-   $\texttt{n}$ is the amount of vertices the graph will have
-   $\texttt{p}$ is the probability of adding a new vertex to the graph
-   $\texttt{maxEdgesize}$ is the maximum size of a generated edge

## Preferential Attachment Hypergraph Model with high Modularity (Giroire et al.)

Looking at the first preferential attachment model, one can see that the resulting graph will be one big community.

We therefore also considered a model with high modularity proposed by Giroire et al.[@Giroire_Nisse_Trolliet_Sulkowska_2022].

# Reduction Rules

The ususal signature of a reduction rule looks as follows:

```go
func NameRule(g HyperGraph, c map[int32]bool) int32
```

We take both a \texttt{HyperGraph} struct \texttt{g} and a _Set_ \texttt{c} as arguments and mutate them. We then return the number of rule executions.

We prioritize time complexity over memory complexity when implementing rules, which does not equate to ignoring memory complexity completely.

## Executions

The proposed rules are meant to be applied exhaustively. A one to one implementation of a rule will only find one of the structures the rule is targeting. Calling such a rule implementation exhaustively will take polynomial time, but will be very inefficient in regards to memory writes and execution time. It is therefor necessary to design the algorithms for the rules with the aspect of exhaustive application in mind. 

The general outline of an algorithm will look as follows,

1. Construct auxilliary data structures that are used to find parts of the graph, for which the rule can be applied. 
2. As long as we can apply rules do:
   
    Iterate over the auxilliary data structure, or the vertices/edges themselves.
      - Identify targets of the rule.
      - Mutate the graph according to the rule.
      - Mutate the auxilliary data structure according to the rule.

This way we minimize memory writes by reusing existing data structures we built in step 1, but also save on execution time, since we mutate and iterate the auxilliary data structure at the same time. 

## Algorithms

The presented algorithms are agnostic in regards to the underlying hypergraph implementation. That means that the datatype used to represent the hypergraph, does not need to have special characteristics. Although it could be beneficial to incorporate some of the used data structures into the hypergraph itself. 

Further we introduce a $map$ data structure in our pseudocode with syntax $\texttt{map}[A]B$, which describes a mapping from $A$ into $B$. We use square brackets to indicate access and mutation of the mapping, e.g. $\gamma[0] \gets 1$. The map type also exposes a primitive function with the signature $delete(\gamma, x)$, which simply means that we want to delete the entry with key $x$ from our map. When iterating over a map in a \texttt{for}-loop we destructure the entries into a $(key, value)$-pair, e.g. $\texttt{for}\; (\_,v) \in map\; \texttt{do}$. Note that unused values are omitted with the underscore symbol. 

### Tiny/Small Edge Rule

-   tiny edges: Delete all hyperedges of size one and place the corresponding vertices into the hitting set.
-   small edges: If $e$ is a hyperedge of size two, i.e., $e = \{ x,y \}$, then put both $x$
    and $y$ into the hitting set.

$O(|E|^2)$ **Algortihm**. Iterate over all edges of the graph and mark all edges of size $t$ in a set $rem$. We then iterate over all edges in $rem$. We put the endpoints of the current edge $e$ in our partial solution. Finally we delete all edges that are edge-adjacent to $e$ from $rem$ and our graph.

$$
\begin{algorithm}[H]
\nllabel{TinySmall}

\KwIn{A hypergraph $G=(V,E)$, a set $C$, an integer $t$ denoting the size of the edges to be removed}
\KwOut{An integer denoting the number of rule applications.}

\BlankLine
$rem \gets \emptyset$\;
$inc \gets \texttt{map}[V]2^E$\;
$exec \gets 0$\;

\For{$e \in E$}{
    \If{$|e| = t$}{
        $rem \gets rem \cup \{e\}$
    }
    \For{$v \in V$}{
        $inc[v] \gets inc[v] \cup e$
    }
}

\For{$e \in rem$}{
    $exec \gets exec+1$\;
    \For{$v \in e$}{
        $C \gets C \cup \{v\}$\;
        $V \gets V\setminus \{v\}$\;
        \For{$f\in inc[v]$}{
            $rem \gets rem\setminus \{f\}$\;
            $E \gets E\setminus \{f\}$\;
        }
    }
}
\KwRet $exec$\;

\caption{Algorithm for exhaustive application of Tiny/Small Edge Rule}
\end{algorithm}
$$

### Edge Domination Rule

-   (hyper)edge domination: A hyperedge $e$ is _dominated_ by another hyperedge $f$ if $f\subset e$. In that case, delete $e$.

$O(|E|)$ **Algorithm**. We partition our set of edges into two disjoint sets $sub$ and $dom$. The set $dom$ will contain edges that could possibly be dominated. The set $sub$ will contain hashes of edges $e$ that could dominate another edge.

We then iterate over the set $dom$ and compute every strict subset of the current edge $f$. For each of these subsets, we test if the hash of the subset is present in our set $sub$. If it is then $f$ is dominated by another edge.

$$
\begin{algorithm}[H]
\nllabel{edom}
\SetKwFunction{subsets}{getSubsetsRec}
\KwIn{A hypergraph $G=(V,E)$ without size one edges, a set $C$}
\KwOut{An integer denoting the number of rule applications.}
\BlankLine

$sub \gets \emptyset$\;
$dom \gets \emptyset$\;
$exec \gets 0$\;

\For{$e \in E$}{
    \eIf{$|e| = 2$}{
        $sub \gets sub\cup \{ hash(e) \}$\;
    }{
        $dom \gets dom\cup \{ e \}$\;
    }
}

\If{$|sub| = 0$}{
    \KwRet $exec$\;
}
\BlankLine

\For{$e \in dom$}{
    $subsets \gets \subsets{e, 2}$\;
    \For{$f \in subsets$}{
        \If{$hash(f) \in sub$}{
            $E \gets E \setminus \{ e \}$\;
            $exec \gets exec+1$\;
            \textbf{break}\;
        }
    }
}

\KwRet $exec$\;

\caption{Algorithm for exhaustive application of Edge Domination Rule}
\end{algorithm}
$$

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

\begin{lemma}
\hyperref[edom]{Algorithm 3} finds all edges of G that are dominated, iff G has no size one edges.
\end{lemma}

_Proof._ Let $e$ be a dominated edge. Since there are no size one edges, edges with size two cannot be dominated. Thus $e$ has to be of size three. Then simply removing $e$ will not create or eliminate an edge domination situation. It is therefore sufficient to only check size three edges for the domination condition. 

This also allows us to parallelize the main part of the algorithm, where we check each edge in our $dom$ set. We can achieve a speedup of $\approx 2$ on a six-core CPU an a pseudo-random graph with one million vertices and two million edges.

### Vertex Domination Rule

-   A vertex $x$ is dominated by a vertex $y$ if, whenever $x$
    belongs to some hyperedge $e$, then $y$ also belongs to $e$. Then, we can simply
    delete $x$ from the vertex set and from all edges it belongs to.

$O(|E| + |V|^2)$ **Algorithm**.
We first construct two maps:

-   $vDeg$: this map associates a vertex $v$ with $deg(v)$.
-   $inc$: this map associates a vertex $v$ with all other edges that are incident to $v$.

We then iterate over the vertices of the graph. For the current vertex $v$, we compute a map $vCount$. This map keeps track of the amount of times a vertex $w$ is adjacent to $v$. The vertex $v$ is dominated, if one of the entries in $vCount$ is equal to $deg(v)$. In that case we remove $v$ from all edges and our vertex set.

$$
\begin{algorithm}[H]
\nllabel{vdom}

\KwIn{A hypergraph $G=(V,E)$ }
\KwOut{An integer denoting the number of rule applications.}
\BlankLine

$vDeg \gets \texttt{map}[V]\mathbb{N}$\;
$inc \gets \texttt{map}[V]2^E$\;

\For{$e \in E$}{
    \For{$v \in e$}{
        $vDeg[v] \gets vDeg[v]+1$\;
        $inc[v] \gets inc[v] \cup \{ e \}$\;
    }
}

$outer \gets true$\;
\While{$outer$}{
    $outer \gets false$\;
    \For{$v \in V$}{
        $vCount \gets \texttt{map}[V]\mathbb{N}$\;
        \For{$e \in inc[v]$}{
            \For{$w \in e$}{
                $vCount[w] \gets vCount[w]+1$\;
            }
        }
        $delete(vCount, v)$\;
        \BlankLine
        $dom \gets false$\;
        \For{$(\_, val) \in vCount$}{
            \If{$val = vDeg[v]$}{
                $dom = true$\;
                $\textbf{break}$\;
            }
        }
        \If{$dom$}{
            $outer = true$\;
            \For{$e \in inc[v]$}{
                $e \gets e\setminus \{ v \}$;
            }
            $V \gets V\setminus \{ v \}$\;
            $delete(inc, v)$\;
            $exec \gets exec+1$\;
        }
    }
}
\KwRet $exec$\;

\caption{Algorithm for exhaustive application of Vertex Domination Rule}
\end{algorithm}
$$

### Approximative Vertex Domination Rule

-   approximative vertex domination: Assume there is a hyperedge $e = \{ x,y,z \}$ such that, whenever $x$ belongs to some hyperedge $h$, then $y$ or $z$ also belong to $h$. Then, we put $y$ and $z$ together into the hitting set that we produce.

$O(|E| + |V|\cdot |E|)$ __Algorithm__. The additional factor $|E|$ looks scary at first, but will only occur in the worst case, if there exists a vertex $v$, that is incident to all edges in $G$.

We first construct two maps:

-   $vDeg$: this map associates a vertex $v$ with $deg(v)$.
-   $adjCount$: this map associates a vertex $v$ with the amount of times another vertex is adjacent to $v$.

Example:

Let $E=\{ \{ 1,2,3 \}, \{ 1,2,4 \} \}$. Then $vDeg$ and $adjCount$ will look as follows,

$$
\begin{align*}
vDeg = \{\quad&1: 2,\\
 &2: 2, \\
 &3: 1, \\
 &4: 1 \quad \}
\end{align*}
$$

$$
\begin{align*}
adjCount = \{\quad&1: \{ 2: 2, \quad 3:1, \quad 4:1\},\\
 &2: \{ 1: 2, \quad 3:1, \quad 4:1\}, \\
 &3: \{ 1: 1, \quad 2:1\}, \\
 &4: \{ 1: 1, \quad 2:1\} \quad \}
\end{align*}
$$

We then iterate other $adjCount$, refering to the current value in the iteration as $adjCount[v]$. We then use the $\textsc{Two-Sum}$-Algorithm to compute and return the first pair in $adjCount[v]$, s.t. for the pair $(a,b)$ holds, $adjCount[v][a] + adjCount[v][b] = vDeg[v]+1$. If such a pair exists, then we conclude that for every edge $f$ such that $v \in f$, it holds for $f$ that, either $a\in f$ or $b\in f$.

\begin{lemma} The outlined procedure above is correct, under the assumption that the underlying graph does not contain any duplicate edges.
\end{lemma}

_Proof._ Let $G$ be a hypergraph. We first remove all edges of size one with the _Tiny Edge Rule_. We then construct our two maps $vDeg$ and $adjCount$. Let $v$ be an entry in $adjCount$ and $sol=(a,b)$ be the result of calling our $\textsc{Two-Sum}$ implementation on $adjCount[v]$ with a target sum of $n=vDeg[v]+1$.

_Proposition._ If $sol$ is non-empty, then the edge $\{ v,a,b \}$ exists.

Let $sol=(a,b)$ be the solution obtained by calling our $\textsc{Two-Sum}$ algorithm on $adjCount[v]$ with a target sum of $n=vDeg[v]+1$. For the sake of contradiction let us assume that the edge $\{ v,a,b \}$ does not exist. Since our graph does not contain duplicate edges and does not contain $\{ v,a,b \}$, there exist $vDeg[v]+1$ many edges that contain either $\{a,v\}$ or $\{b,v\}$. This however contradicts that there only exist $vDeg[v]$ many edges containing $v$. Therefore it must be, that the assumption that $\{ v,a,b \}$ does not exist, is false.

Since $\{ v,a,b \}$ exists, $a$ and $b$ can only occur $n-2 = vDeg[v] - 1$ times in other edges containing $v$. Since duplicate edges of $\{ v,a,b \}$ can not exist, we know that every other edge containing $v$ also contains $a$ or $b$, but not both simultaniously. $\square$

We then add the two vertices in the solution $sol$ to our partial solution $c$.

$$
\begin{algorithm}[H]
\nllabel{apvdom}

\KwIn{A hypergraph $G=(V,E)$, a set $C$}
\KwOut{An integer denoting the number of rule applications.}
\BlankLine

$vDeg \gets \texttt{map}[V]\mathbb{N}$\;
$adjCount \gets \texttt{map}[V]\texttt{map}[V]\mathbb{N}$\;
$inc \gets \texttt{map}[V]2^E$\;
$exec \gets 0$\;

\For{$e \in E$}{
    \For{$v \in e$}{
        $vDeg[v] \gets vDeg[v]+1$\;
        $inc[v] \gets inc[v] \cup \{ e \}$\;
        \For{$w \in e$}{
            \If{$w \neq v$}{
                $adjCount[v][w] \gets adjCount[v][w]+1$\;
            }
        }
    }
}

$outer \gets true$\;
\While{$outer$}{
    $outer \gets false$\;
    \For{$(v, count) \in adjCount$}{
        $sol, ex \gets \textsc{Two-Sum}(count, vDeg[v]+1)$\;
        \If{\textbf{not} $ex$}{
            \textbf{continue}\;
        }
        $outer \gets true$\;
        $exec \gets exec+1$\;

        \For{$w \in sol$}{
            $C\gets C\cup \{ w \}$\;
            \For{$e \in inc[w]$}{
                \For{$x \in e$}{
                    \If{$x = w$}{
                        \textbf{continue}\;
                    }
                    $vDeg[x] \gets vDeg[x]-1$\;
                    \For{$y \in e\setminus x$}{
                        $adjCount[x][y] \gets adjCount[x][y]-1$\;
                    }
                }
                $E \gets E\setminus \{ e \}$\;
            }
            $V \gets V\setminus \{ w \}$\;
            $delete(inc, w)$\;
			$delete(vDeg, w)$\;
			$delete(adjCount, w)$\;
        }
    }
}

\KwRet $exec$\;

\caption{Algorithm for exhaustive application of Approximative Vertex Domination Rule}
\end{algorithm}
$$

**Idea**: The initial idea for this algorithm involved the usage of an incidence matrix, where edges are identified by the rows and the vertices are identified by the columns. To check the _Domination Condition_ for a vertex $v$, the algortihm would select all edges/columns that contain $v$ and then add up the columns. Now let $n$ be the amount of edges containing $v$. If there exist two entries in the resulting column that have a combined value of $n+1$, then the rule applies for $v$ under the assumption that there are no duplicate edges. This would result in an algorithm with a worse time complexity of $|V|+|V|^2\cdot|E|$, and a high memory complexity of $|V|\cdot|E|$.

### Approximative Double Vertex Domination Rule

-   approximative double vertex domination: Assume there is a hyperedge $e =
\{x, y, a\}$ and another vertex $b$ such that, whenever $x$ or $y$ belong to some
    hyperedge $h$, then $a$ or $b$ also belong to $h$. Then, we put $a$ and $b$ together
    into the hitting set that we produce.

$\mathcal{O}(|E| + |E|^2)$ __Algorithm__. We start by iterating over the edges of our graph and building an incidence list. We then iterate over the edges again. If the current edge $e$ is not of size 3 we continue with the iteration. If $|e| = 3$ we iterate over the endpoints of $e$. We assign the variable $\texttt{a}$ to $v$ and create a map $\texttt{vCount map[int32]int32}$. For each vertex $w \in e\setminus v$ we look at the the edges incident to $w$, we then for each vertex $u$ of these edges, increment $\texttt{vCount[u]}$. We ignore edges that are incident to $\texttt{a}$ and do not care about vertices that are in $e$. We keep track of the amount of edges we considered with the variable $\texttt{xyCount}$.

If we find an entry $\texttt{vCount[b]}$, s.t. $\texttt{vCount[b] = xyCount}$, we found a double vertex domination situation in $e$.

### Small Triangle Rule

-   small triangle situation: Assume there are three small hyperedges $e = \{y, z\}$, $f = \{x, y\}$, $g = \{x, z\}$. This describes a triangle situation $(e, f, g)$. Then, we put $\{x, y, z\}$ together into the hitting set, and we can even choose another hyperedge of size three to worsen the ratio.

```go
func SmallTriangleRule(g HyperGraph, c map[int32]bool)
```

$O(|E|+|V|^2)$ __Algorithm__. Again the quadratic exponent in $|V|$ looks scarier than it is. This happens beacause we have to remove the triangle. We start by constructing an edge-adjacency list $\texttt{adjList}$ for all edges of size two. We then iterate over the entries of the list. For the current entry $\texttt{adjList[v]}$ we compute all subsets of size two of the entry. If both vertices of the subset are vertex-adjacent to each other, then we found a small triangle situation. If we find a triangle situation we put the corresponding vertices in our partial solution and alter the adjacency list to reflect these changes. We do this by iterating over all vertices that are vertex-adjacent to the triangle. For every vertex $w$ of these vertices we delete all vertices of the triangle from the entry $\texttt{adjList\lbrack}w\texttt{\rbrack}$.

This last step will introduce the quadratic complexity, since in the worst case, for a vertex $v$ in a triangle, there could exist $|V|$ many size two edges that contain $v$. This worst case occurs very rarely, which justifies using this quadratic algorithm. We could alternatively move the last step of the algorithm outside of the loop, and wrap both procedures with an outer loop which breaks if we dont find any more triangles. This simulates calling the rule exhaustively, while achieving a linear time complexity.

### Small Edge Degree 2 Rule

-   small edge degree 2: Let v be a vertex of degree 2, and let the two hyperedges
    containing $v$ be $e = \{x, v\}$ and $f = \{v, y, z\}$. Then we can select a hyperedge $g$
    that contains one of the neighbors of $v$ in $f$ but not $x$, for example $g = \{u, w, z\}
$(when $y = w$ is possible as a special case) or $g = \{u, z\}$. We put $x, u$ and $z$
    and $w$ (when existing) into the hitting set.

#### Outline

-   find deg 2 vertex
-   Check for at least one adjacent size two edge
-   iterate over the size 3 edge neighbors $w$ of $v$, or if both are size 2 try both
-   select a edge adjacent to $w$ that does not contain $x$

data structures:
degMap
incMap

## Self-Monitoring

Each of the reduction rule functions returns a $\texttt{int32}$ value, which indicates the number of rule executions.

We store the ratios for each rule in a map of the form:

```go
var ratios = map[string]pkg.IntTuple{
	"kRulename": {A:1, B:1},
}
```

Where $\texttt{A}$ denotes the amount of vertices put into the partial solution by a single rule execution. And $\texttt{B}$ denotes the amount of vertices present in an optimal solution.

We can then use these values to calculate the estimated approximation factor as follows:

```go
execs := ApplyRules(g, c)

var nom float64 = 0
var denom float64 = 0

for key, val := range execs {
	nom += float64(ratios[key].A * val)
	denom += float64(ratios[key].B * val)
}
```

We conducted some preliminary testing on graphs with 10, 100, 1000 and 10000 vertices. We calculateted the estimated approximation factor for these graphs. We did this for ratios $r = \frac{|E|}{|V|}$ of 1 to 20.

![Estimated Approximation Factor for Graphs with 10, 100, 1K, 10K Vertices for $r\in[1,20]$](https://raw.githubusercontent.com/KhoalaS/BachelorThesis/main/img/line_fac.png)

We also looked at the number of rule executions for $r=1$ and $r=10$.

![Number of Rule Executions, $r=1$](https://raw.githubusercontent.com/KhoalaS/BachelorThesis/main/img/bar_ex_1.png)

![Number of Rule Executions, $r=10$](https://raw.githubusercontent.com/KhoalaS/BachelorThesis/main/img/bar_ex_10.png)

Note that the reduction rules alone are sufficient to compute a Hitting-Set for our graphs. The est. ratio is quite low, since we do not need to put whole size 3 edges into our hitting set. This is due to the way our random-graph model works. This will not work for all graph classes, namely 3-uniform graphs.

We can observe a large performance hit, once our graph is near 3-uniform. We present some ideas that could potentially speedup the execution for these graphs.

**Factor-3 Rule Function**. Instead of only executing the Factor-3 Rule exactly one time, we could consider using a function instead. This function would compute a value based on the current iteration. Experiments with linear,exponential and square root functions were quite promising in regards to execution time and approximation factor. Linear and exponential functions perform faster than the square root function, while the linear and square root function offer a better estimated approximation factor than the exponential one.

**Factor-3 Rule Targeting**. The Factor-3 rule will select the to be removed edge at random. We could alternatively choose an edge $e$, s.t. the removal of $e$ will allow the execution of the Vertex Domination Rule. This could potentially lower the number of Factor-3 Rule executions quite significantly.

# Applications

# Testing

Every reduction rule is tested for their correctness with unit tests. We create small graphs in these tests, that contain structures, which the rules are targeting. We then test for the elements in the partial solution and the amount of edges left in the graph.

# Branching Algorithm

## Potential Triangle Siutation

## Edge-Cover

If our problem instance is simple, we can solve it in polynomial time using an edge cover algorithm for standard graphs. We transform the hypergraph of our instance to a standard graph as follows:

-   Let $v$ be a hypergraph vertex with $\delta(v)=2$ and $e_0, e_1$ be the two edges that $v$ is incident to.
-   We then add two vertices to the auxilliary graph identified by $e_0$ and $e_1$. At last we add an edge $\{e_0, e_1\}$ identified by $v$.

## Possible Optimizations

Right now we sometimes use the Go $\texttt{apppend}$ function, to add elements to a slice. Calls to $\texttt{apppend}$ are more expensive than an ordinary assignment to a fixed size slice.

Using fixed size slices could improve performance slightly in cases where we know the maximum amount of elements we will add to a slice.

# References