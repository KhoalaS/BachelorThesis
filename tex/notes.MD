---
    title: Implementation and evaluation of a self-monitoring approximation algorithm for 3-Hitting-Set
    link-citations: true
    author: Khoa Le
    abstract: |
        In this thesis we implemented an approximation algorithm for the \textsc{Hitting Set} problem. The algorithm is based on various reduction rules and the approximation ratio is estimated with the local ratio method. We show test results on real-world networks that can be interpreted as derivative problems of \textsc{Hitting Set}. At last, we compare these results with two linear programming based algorithms.
...

\RestyleAlgo{ruled}
\DontPrintSemicolon
\SetAlgoVlined
\LinesNumbered
\listoftodos
\SetKwFunction{FRecurs}{FnRecursive}
\SetKwFunction{ExpandFrontier}{ExpandFrontier}
\SetKwProg{Fn}{func}{\string:}{}

# Introduction


## Preliminaries

A hypergraph $G=(V,E)$ consists of a set of vertices $V$ and a set of edges $E$. Each element of $E$ is a subset of $V$. The degree of a vertex $v$ is the number of edges that are incident to $v$. We also denote the degree of a vertex $v$ with $deg(v)$. An edge $e$ has size $m$, if it contains $m$ elements.

# Programming Language

We chose the [_Go_](https://go.dev/) programming language. It is a statically typed, compiled language, with a C-like syntax. It is using a garbage collector to handle memory management, which at first seems off-putting for an application like this. Since Go's GC is very efficient, we are not worried about that fact. Go also provides great developer tooling without the need for additional third party tools. The built-in profiler _pprof_ can visualize captured performance profiles as flame graphs or standard graphs, which makes it easy to spot performance bottlenecks.

# Data Structures

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

The $\texttt{Edge}$ datatype has one field. The field $\texttt{v}$ is a map with keys of type $\texttt{int32}$ and values of type $\texttt{bool}$. When working with the endpoints of an edge, we are usually not interested in the associated values, since we never mutate the edges. This simulates a _Set_ datatype while allowing faster access times than simple arrays/slices.

## Hypergraph

```go
type HyperGraph struct {
	Vertices       map[int32]Vertex
	Edges          map[int32]Edge
	edgeCounter    int32
	IncMap         map[int32]map[int32]bool
	AdjCount       map[int32]map[int32]int32
}

func (g *HyperGraph) AddVertex(id int32, data int)
func (g *HyperGraph) RemoveVertex(id int32) bool
func (g *HyperGraph) RemoveElem(elem int32) bool
func (g *HyperGraph) AddEdge(eps ...int32)
func (g *HyperGraph) RemoveEdge(e int32) bool
func (g *HyperGraph) Deg(v int32) int
func (g *HyperGraph) RemoveDuplicate()
```

The $\texttt{HyperGraph}$ datatype has five fields. Both fields $\texttt{Vertices}$ and $\texttt{Edges}$ are maps with keys of type $\texttt{int32}$ and values of type $\texttt{Vertex}$ and $\texttt{Edge}$ respectively. We chose this set-like data structure over lists again because of faster access times, but also operations that remove edges/vertices are built-in to the map type. The field $\texttt{edgeCounter}$ is an internal counter used to assign ids to added edges. The field $\texttt{IncMap}$ is a map of maps, essentially storing the hypergraph as a sparse incidence matrix. We will also derive vertex degrees from this map. And at last the $\texttt{AdjCount}$ map, which will associate every vertex $v \in V$ to all vertices adjacent to $v$. Additionally, this map will also store the number of times such a vertex is adjacent to $v$. We also provide various struct methods to do basic operations on the hypergraph. All of these methods handle mutations to the \texttt{IncMap} and \texttt{AdjCount} fields of the receiving hypergraph struct. For example, calling \texttt{RemoveEdge(0)} on a hypergraph struct \texttt{g} will remove the edge with id 0 from \texttt{g.E} and will remove \texttt{0} from all entries \texttt{g.IncMap[v]}, where \texttt{v} is an endpoint of the removed edge. 

# Misc. Algorithms/Utilities

## Edge Hashing

Given a set $S$ of size $n$, compute a unique hash of $S$.

Time Complexity: $n + n\cdot \log(n)$

We start by sorting $S$ with a $\textsc{Quick-Sort}$-Algorithm. We then join the elements of $S$ with the delimiter $"|"$, returning a string of the form $"|i_0|i_1|\dots|i_n|"$. Whenever we refer to \textit{the hash of an edge} we refer to the output of this function, using the endpoints of the edge as the input set. This hash value is useful when we want to check the existence of specific edges.

## Compute Subsets of Size $s$

Given an array $arr$ and an integer $s$, compute all subsets of $arr$ of size $s$.

$$
\begin{algorithm}[H]
\label{subset}

\KwIn{An array $arr$, the subset size $s$ and a list $subsets$.}
\KwOut{}
\BlankLine

$data\gets [\;]$\;
$n\gets |arr|$\;
$last\gets s-1$\;
$\texttt{FnRecursive}(0,0)$\;

\Fn{\FRecurs{i, next}}{
    \For{$j\gets next$ \KwTo $n$}{
        $data[i] \gets arr[j]$\;
        \eIf{$i=last$}{
            $subsets.push(data)$\;
        }{
            $\texttt{FnRecursive}(i+1, j+1)$\;
        }
    }
}
\caption{An algorithm to compute all subsets of size $s$}
\end{algorithm}
$$

Lists in Go are not very memory efficient, but since we exclusively call this function with $arr$ representing the vertices in an edge, the value is usually fixed at 3. The raised memory problems occur at values of $n>10000$, justifying the continued usage of lists. For the case where we have to compute a lot of subsets, we provide a slightly different version of this function. Instead of passing in the $subsets$ list, we pass in a callback function that is called whenever we find a subset, using the found subset as an argument.

## Two-Sum

Given an array $\texttt{items}$ of integers and an integer target $t$, return indices of the two numbers such that they add up to $t$.

Time Complexity: $n$, where $n$ denotes the size of $\texttt{items}$.

$$
\begin{algorithm}
\label{twosum}

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

We start by creating a map called $lookup$. Iterating other the elements of $arr$, we check if the entry $lookup[t-arr[i]]$ exists. If the entry exists, we return a pair $(i, lookup[t-arr[i]])$ and the boolean value $true$ since we found a solution. If the entry does not exist, we add a new entry to the $lookup$ map using $arr[i]$ as key and $i$ as value. If no solution was found, we return $nil$ and the boolean value $false$.

This algorithm is an ingredient for the implementation of one of the reduction rules, specifically the approximative vertex domination rule. The actual implementation accepts a map instead of an array as its first parameter. We also implemented a version that finds all solutions, which accepts an additional callback function as a parameter. Instead of returning the solution, we call the callback function with the solution as the first argument. This version is an ingredient for the implementation of the approximative double vertex domination. We refer to this version as $\textssc{Two-SumAll}$ inside the algorithm listings.

# Hypergraph Models

## First Testing Model

```go
func GenerateTestGraph(n int32, m int32, tinyEdges bool) *HyperGraph
```

Let us explain the arguments first:

-   $\texttt{n}$ are the number of vertices the graph will have
-   $\texttt{m}$ is the number of edges the graph will at most have
-   $\texttt{tinyEdges}$ when $\texttt{false}$ indicates that we do not want to generate edges of size 1.

We use a very naive approach for generating (pseudo-)random graphs. We first create an empty Hypergraph struct and add $\texttt{n}$ many vertices to that graph. We then compute a random $\texttt{float32}$ value $\texttt{r}$ in the half-open interval $[ 0.0, 1.0 )$. This value will be used to determine the size of an edge $e$. The edges are distributed based on their size as follows.

$$
size(\texttt{r})=
\begin{cases}
    1 & \quad \texttt{r} < 0.01\\
    2 & \quad 0.01 \leq \texttt{r} < 0.60\\
    3 & \quad \text{else}
\end{cases}
$$

The result of $size(\texttt{r})$ is stored in a variable $\texttt{d}$. We then randomly pick vertices in the half-open interval $[ 0, \texttt{n} )$, until we have picked $\texttt{d}$ many distinct vertices. If an edge with these endpoints does not exist, we add it to our graph.

That results in the graph having at most $\texttt{m}$ edges and not exactly $\texttt{m}$, since we did not want to artificially saturate the graph with edges. One could also look into generating random bipartite graphs that translate back to a hypergraph with the desired vertex and edge numbers. The advantage of this model over a model like the Erdős–Rényi model is that we can compute hypergraphs with large vertex count fast. The model was primarily used during the implementation phase of the algorithm to gather profiling data.

## Preferential Attachment Hypergraph Model

In the Preferential Attachment Model, one will add edges to an existing graph, with a probability proportional to the degree of the endpoints of that edge. This edge will either contain a newly added vertex or will be comprised of vertices already part of the graph. We will use an implementation by Antelmi et al. as reference [@SimpleHypergraphs.js], which is part of their work on _SimpleHypergraphs.jl_ [@DBLP:journals/im/SpagnuoloCSPSKA20], a hypergraph software library written in the Julia language. The implementation is based on a preferential attachment model proposed by Avin et al. in [@DBLP:conf/asunam/AvinLNP19].

```go
func GeneratePrefAttachmentGraph(n int, p float64, maxEdgesize int32)
```

-   $\texttt{n}$ is the amount of vertices the graph will have
-   $\texttt{p}$ is the probability of adding a new vertex to the graph
-   $\texttt{maxEdgesize}$ is the maximum size of a generated edge

# Reduction Rules

The ususal signature of an implemented reduction rule looks as follows.

```go
func NameRule(g *HyperGraph, c map[int32]bool) int32
```

We take both a pointer to a \texttt{HyperGraph} struct \texttt{g} and a map \texttt{c} as arguments and mutate them. We then return the number of rule executions. We prioritize time complexity over memory complexity when implementing rules, which does not equate to ignoring memory complexity completely.

## Executions

The proposed rules are meant to be applied exhaustively. A one-to-one implementation of a rule will only find one of the structures the rule is targeting. Calling such a rule implementation exhaustively will take polynomial time but will be very inefficient regarding memory writes and execution time. It is therefore advantageous to design the algorithms for the rules, with the aspect of exhaustive application in mind.

The general outline of an algorithm will look as follows.

1. (If needed) Construct auxiliary data structures that are used to find parts of the graph, for which the rule can be applied.
2. If we can apply rules do:

    Iterate over the auxiliary data structure, or the vertices/edges themselves.

    - Identify targets of the rule.
    - Mutate the graph according to the rule.
    - Mutate the auxiliary data structure according to the rule.

This way we minimize memory writes by reusing existing data structures we built in step 1, but also save on execution time, since we mutate and iterate the auxiliary data structure at the same time. We can also alter these algorithms with minimal changes, to only execute a rule once. 

## Algorithms

The algorithm descriptions occasionally omit implementation details. We do this to try and keep these descriptions as concise as possible. The real implementations correspond to the pseudocode listings. We do not impose the use of our proposed hypergraph model in those listings. Instead we use general language, for example $e \textnormal{incident to} v$ instead of $e \in H.IncMap[v]$. Proofs are provided in the case, where the algorithm does not simply follow the definition of a rule.

We further introduce a $map$ data structure in our pseudocode with syntax $\texttt{map}[A]B$, which describes a mapping from $A$ into $B$. We use square brackets to indicate access and mutation of the mapping, e.g. $\gamma[0] \gets 1$ maps the value $1$ to the key $0$ in map $\gamma$. The map type also exposes a primitive function with the signature $delete(\gamma, x)$, which simply means that we want to delete the entry with key $x$ from our map. When iterating over a map in a \texttt{for}-loop we destructure the entries into a $(key, value)$-pair, e.g. $\texttt{for}\; (\_,v) \in map\; \texttt{do}$. Unused values are omitted with the underscore symbol.

### Tiny/Small Edge Rule

-   tiny edges: Delete all hyperedges of size one and place the corresponding vertices into the hitting set.
-   small edges: If $e$ is a hyperedge of size two, i.e. $e = \{ x,y \}$, then put both $x$
    and $y$ into the hitting set.

$O(|E|^2)$ **Algortihm**. We iterate over all edges of the graph. If the current edge $e$ is of size $t$, put $e$ into the partial hitting set and remove $e$ and all edges adjacent to vertices in $e$ from the graph.

$$
\begin{algorithm}

\KwIn{A hypergraph $G=(V,E)$, a set $C$, an integer $t$ denoting the size of the edges to be removed}
\KwOut{An integer denoting the number of rule applications.}

\BlankLine
$rem \gets \emptyset$\;
$exec \gets 0$\;

\For{$e \in E$}{
    \If{$|e| = t$}{
        $exec \gets exec+1$\;
        \For{$v \in e$}{
            $C \gets C \cup \{v\}$\;
            $V \gets V\setminus \{v\}$\;
            \For{$f$ \textnormal{incident to} $v$}{
                $E \gets E\setminus \{f\}$\;
            }
        }
    }
}

\KwRet $exec$\;

\caption{Algorithm for exhaustive application of Tiny/Small Edge Rule\label{TinySmall}}
\end{algorithm}
$$

### Edge Domination Rule

-   (hyper)edge domination: A hyperedge $e$ is _dominated_ by another hyperedge $f$ if $f\subset e$. In that case, delete $e$.

$O(|E|)$ **Algorithm**. We partition our set of edges into two disjoint sets $sub$ and $dom$. The set $dom$ will contain edges that could be dominated. The set $sub$ will contain hashes of edges $e$ that could dominate another edge. We then iterate over the set $dom$ and compute every strict subset of the current edge $f$. For each of these subsets, we test if the hash of the subset is present in our set $sub$. If it is then $f$ is dominated by another edge.

$$
\begin{algorithm}
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
\caption{Algorithm for exhaustive application of Edge Domination Rule\label{edom}}
\end{algorithm}
$$

The exact time complexity is as follows.

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
Algorithm \ref{edom} finds all edges of G that are dominated, iff G has no size one edges.
\end{lemma}

_Proof._ Let $e$ be a dominated edge. Since there are no size one edges, edges with size two cannot be dominated. Thus $e$ has to be of size three. Then simply removing $e$ will not create or eliminate an edge domination situation. It is therefore sufficient to only check size three edges for the domination condition. $\square$

This also allows us to parallelize the main part of the algorithm, where we check each edge in our $dom$ set. We can achieve a speedup of approximately $2$ on a six-core CPU and a pseudo-random graph with one million vertices and two million edges.

### Vertex Domination Rule

-   A vertex $x$ is dominated by a vertex $y$ if, whenever $x$
    belongs to some hyperedge $e$, then $y$ also belongs to $e$. Then, we can simply
    delete $x$ from the vertex set and from all edges it belongs to.

$O(|V|^2 \cdot |E|)$ **Algorithm**
A vertex $v$ is dominated, if one of the entries in $AdjCount[v]$ is equal to $deg(v)$. In that case we remove $v$ from all edges and our vertex set.

$$
\begin{algorithm}

\KwIn{A hypergraph $G=(V,E)$ }
\KwOut{An integer denoting the number of rule applications.}
\BlankLine

$exec \gets 0$\;

$outer \gets true$\;
\While{$outer$}{
    $outer \gets false$\;
    \For{$v \in V$}{
        $dom \gets false$\;
        \For{$(\_, val) \in AdjCount[v]$}{
            \If{$val = deg(v)$}{
                $dom = true$\;
                $\textbf{break}$\;
            }
        }
        \If{$dom$}{
            $outer = true$\;
            \For{$e$ \textnormal{incident to} $v$}{
                $e \gets e\setminus \{ v \}$;
            }
            $V \gets V\setminus \{ v \}$\;
            $exec \gets exec+1$\;
        }
    }
}
\KwRet $exec$\;

\caption{Algorithm for exhaustive application of Vertex Domination Rule\label{vdom}}
\end{algorithm}
$$

### Approximative Vertex Domination Rule

-   approximative vertex domination: Assume there is a hyperedge $e = \{ x,y,z \}$ such that, whenever $x$ belongs to some hyperedge $h$, then $y$ or $z$ also belong to $h$. Then, we put $y$ and $z$ together into the hitting set that we produce.

$O(|V|^2 \cdot |E|)$ **Algorithm**. The additional factor $|E|$ looks scary at first, but will only occur in the worst case, if there exists a vertex $v$ that is incident to all edges in $G$. We start by iterating other the $AdjCount$ map of the graph, referring to the current value in the iteration as $AdjCount[v]$. We then use the $\textsc{Two-Sum}$ algorithm to compute and return the first pair $(a,b)$ in $AdjCount[v]$, s.t. for $(a,b)$ holds,

$$
AdjCount[v][a] + AdjCount[v][b] = deg(v)+1
$$

If such a pair exists, then we conclude that for every edge $f$ incident to $v$, it holds that either $a\in f$ or $b\in f$.

\begin{lemma} The outlined procedure above is correct, under the assumption that the underlying graph does not contain any duplicate edges or size one edges.
\end{lemma}

_Proof._ Let $G$ be a hypergraph without size one edges. Let $v$ be an entry in $AdjCount$ and $sol=(a,b)$ the result of calling our $\textsc{Two-Sum}$ algorithm on $AdjCount[v]$ with a target sum of $n=deg(v)+1$.

_Proposition._ If $sol$ is non-empty, then the edge $\{ v,a,b \}$ exists.

Let $sol=(a,b)$ be the solution obtained by calling our $\textsc{Two-Sum}$ algorithm on $AdjCount[v]$ with a target sum of $n=deg(v)+1$. For the sake of contradiction let us assume that the edge $\{ v,a,b \}$ does not exist. Since our graph does not contain duplicate edges, size one edges and not the edge $\{ v,a,b \}$, there exist $deg(v)+1$ many edges that contain either $\{a,v\}$ or $\{b,v\}$. This however contradicts that there only exist $deg(v)$ many edges containing $v$. Therefore the assumption that $\{ v,a,b \}$ does not exist, must be false.

Since $\{ v,a,b \}$ exists, $a$ and $b$ can only occur $n-2 = deg(v) - 1$ times in other edges containing $v$. Since duplicate edges of $\{ v,a,b \}$ cannot exist, we know that every other edge containing $v$ also contains $a$ or $b$, but not both simultaneously. $\square$

We then add the two vertices in the solution $sol$ to our partial solution $C$.

$$
\begin{algorithm}

\KwIn{A hypergraph $G=(V,E)$, a set $C$}
\KwOut{An integer denoting the number of rule applications.}
\BlankLine

$exec \gets 0$\;

$outer \gets true$\;
\While{$outer$}{
    $outer \gets false$\;
    \For{$(v, count) \in AdjCount$}{
        $sol, ex \gets \texttt{TwoSum}(count, deg(v)+1)$\;
        \If{\textbf{not} $ex$}{
            \textbf{continue}\;
        }
        $outer \gets true$\;
        $exec \gets exec+1$\;

        \For{$w \in sol$}{
            $C\gets C\cup \{ w \}$\;
            $V \gets V\setminus \{ w \}$\;
            \For{$e$ \textnormal{incident to} $w$}{
                $E \gets E\setminus \{ e \}$\;
            }
        }
    }
}

\KwRet $exec$\;

\caption{Algorithm for exhaustive application of Approximative Vertex Domination Rule\label{apvdom}}
\end{algorithm}
$$

**Idea**: The initial idea for this algorithm involved the usage of a complete incidence matrix, where edges are identified by the rows and the vertices are identified by the columns. To check the _Domination Condition_ for a vertex $v$, the algorithm would select all edges/columns that contain $v$ and then add up the columns. Now let $n$ be the number of edges containing $v$. If there exist two entries in the resulting column that have a combined value of $n+1$, then the rule applies for $v$ under the assumption that there are no duplicate edges. This would result in an algorithm with a worse time complexity of $|V|+|V|^2\cdot|E|$.

### Approximative Double Vertex Domination Rule

-   approximative double vertex domination: Assume there is a hyperedge $e =
\{x, y, a\}$ and another vertex $b$ such that, whenever $x$ or $y$ belong to some
    hyperedge $h$, then $a$ or $b$ also belong to $h$. Then, we put $a$ and $b$ together
    into the hitting set that we produce.

$\mathcal{O}(|V|^2 + |E|)$ **Algorithm**. We start by creating a map called $tsHashes$, which maps subsets of $E$ to subsets of vertices. We then iterate over all vertices in $V$. For the current vertex $x$, compute all \textsc{Two-Sum} solutions with input array/map $AdjCount[x]$ and target $deg(x)$. If there exists a solution $sol = \{ z_0, z_1 \}$, such that $tsHashes[sol]\neq \emptyset$, then for all $y \in tsHashes[sol], y\neq x$ construct two edges $\{ x,y,z_0 \}$ and $\{ x,y,z_1 \}$. If one of these two edges exists in $E$, then we found a approximative double vertex domination situation. If $tsHashes[sol] = \emptyset$, map $tsHashes[sol]$ to $tsHashes[sol]\cup \{ x \}$.

\begin{lemma}
Algorithm \ref{apdvdom} is correct.
\end{lemma}

_Proof_.
Let $G=(V,E)$ be a hypergraph. Let $x$ be the current vertex in the algorithm's iteration over $V$. If $T = \textsc{Two-SumAll}(AdjCount[x], deg(x))$ is empty, then $x$ will not be able to trigger a approximative double vertex domination situation. If $T$ is not empty, then we know that there exist sets $\{ a,b \}$, such that all edges incident to $x$ either contain $a$ or $b$. If $tsHashes[\{ a,b \}]$ is empty, then there are currently no other vertices for which $\{ a,b \}$ is a \textsc{Two-Sum} solution. In that case we add $x$ to $tsHashes[\{ a,b \}]$. Otherwise there exist vertices $y$, such that $\{ a,b \}$ is a \textsc{Two-Sum} solution for $y$. All edges incident to $x$ and $y$ either contain $a$ or $b$. If for one of these vertices $y$ there exists a size three edge $\{ x,y,a \}$, or without loss of generality $\{ x,y,b \}$, then we found a approximative double vertex domination situation at $\{ x,y,a \}$ or $\{ x,y,b \}$ respectively. If none of these edges exist, then $x$ could still trigger another approximative double vertex situation with another vertex. Thus we need to add $x$ to $tsHashes[\{ a,b \}]$.
$\square$

$$
\begin{algorithm}

\KwIn{A hypergraph $G=(V,E)$, a set $C$}
\KwOut{An integer denoting the number of rule applications.}
\BlankLine

$exec \gets 0$\;

$outer \gets true$\;
\For{$outer$}{
    $outer \gets false$\;
    $tsHashes \gets \texttt{map}[2^V]2^V$\;

    \For{$x \in V$}{
        \For{$sol \in \textnormal{\texttt{TwoSumAll}}(AdjCount[x], deg(x))$}{
            $\{ z_0, z_1 \} \gets sol$\;
            \eIf{$tsHashes[sol] \neq \emptyset$}{
                \For{$y \in tsHashes[sol]$}{
                    \If{$y = x$}{
                        \textbf{continue}\;
                    }
                    $f_0 \gets \{ x,y,z_0 \}$\;
                    $f_1 \gets \{ x,y,z_1 \}$\;

                    $found \gets false$\;
                    \For{$e$ \textnormal{incident to} $y$}{
                        \If{$e=f_0 \textbf{ or } e=f_1$}{
                            $found \gets true$\;
                            \textbf{break}\;
                        }
                    }

                    \If{found}{
                        $exec\gets exec + 1$\;
                        $outer\gets true$\;
                        $C\gets C \cup sol$\;
                        \For{$a \in sol$}{
                            \For{$e$ \textnormal{incident to} $a$}{
                                $E\gets E\setminus \{ e \}$\;
                            }
                        }
                        \textbf{break}\;
                    }
                    $tsHashes[sol] \gets tsHashes[sol] \cup \{ x \}$\;
                }
            }{
                $tsHashes[sol] \gets \{ x \}$\;
            }
        }
    }
}

\KwRet $exec$\;

\caption{Algorithm for exhaustive application of Approximative Double Vertex Domination Rule\label{apdvdom}}
\end{algorithm}
$$

### Small Triangle Rule

-   small triangle situation: Assume there are three small hyperedges $e = \{y, z\}$, $f = \{x, y\}$, $g = \{x, z\}$. This describes a triangle situation $(e, f, g)$. Then, we put $\{x, y, z\}$ together into the hitting set, and we can even choose another hyperedge of size three to worsen the ratio.

$O(|E|+|V|^2)$ **Algorithm** We start by constructing an adjacency list $adjList$ for all edges of size two. We then iterate over the entries of the list. For the current entry $adjList[v]$ we compute all subsets of size two of the entry. If there exists a subset $s$ such that $s \in E$, then we found a small triangle situation. If we find a triangle situation, we put the corresponding vertices in our partial solution and alter the adjacency list to reflect these changes. We do this by iterating over all vertices that are adjacent to the triangle. For every vertex $w$ of these vertices we delete all vertices of the triangle from the entry $adjList[w]$.

This last step will introduce the quadratic complexity, since in the worst case, for a vertex $v$ in a triangle, there could exist $|V|$ many size two edges that contain $v$. This worst case occurs very rarely, which justifies using this quadratic algorithm. We could alternatively move the last step of the algorithm outside of the loop and wrap both procedures with an outer loop which breaks if we do not find any more triangles. This simulates calling the rule exhaustively, while achieving a linear time complexity.

$$
\begin{algorithm}
\SetKwFunction{subsets}{getSubsetsRec}
\KwIn{A hypergraph $G = (V, E)$, a set $C$}
\KwOut{An integer denoting the number of rule applications.}
\BlankLine

$adjList \gets \texttt{map}[V]2^V$\;
$rem \gets \emptyset$\;
$exec \gets 0$\;

\For{$e\in E$}{
    \If{$|e| \neq 2$}{
        \textbf{continue}\;
    }
    $\{ x,y \} \gets e$\;
    $adjList[x] \gets adjList[x] \cup \{ y \}$\;
    $adjList[y] \gets adjList[y] \cup \{ x \}$\;
}

\For{$(z, val) \in adjList$}{
    \If{$|val| < 2$}{
        \textbf{continue}\;
    }
    $subsets \gets \subsets{val, 2}$\;
    \For{$s\in subsets$}{
        $\{x,y\} \gets s$\;
        \If{$y \in adjList[x]$ \textnormal{\textbf{or}} $x\in adjList[y]$}{
            $exec \gets exec +1$\;
            $C\gets C \cup \{ x,y,z \}$\;
            $rem \gets rem \cup \{ x,y,z \}$\;
            \For{$u \in \{ x,y,z \}$}{
                \For{$v\in adjList[u]$}{
                    $adjList[v] \gets adjList[v]\setminus \{ u \}$\;
                }
                $delete(adjList, u)$
            }
            \textbf{break};
        }
    }
}

\For{$v \in rem$}{
    $V\gets V\setminus \{ v \}$\;
    \For{$e$ \textnormal{incident to} $v$}{
        $E\gets E\setminus \{ e \}$\;
    }
}

\KwRet $exec$\;

\caption{: Algorithm for exhaustive application of Small Triangle Rule\label{tri}}
\end{algorithm}
$$

### Extended Triangle Rule

-   Assume that the hypergraph contains a small edge $e = \{y, z\}$. Moreover,
    there are hyperedges $f, g$ such that $e \cap f = \{ y \}$, $e \cap g = \{z\}$, $f \cup g=\{v, x, y, z\}$ and $|f| = 3$. Then, put all of $f \cup g$ into the hitting set.

$O(|E|^3)$ **Algorithm** We start by iterating over $E$ until we find a size two edge $e$. We then iterate over the vertices of $e$. Let $a$ be the vertex in the current iteration. Assign $a$ to the variable $y$ and $e\setminus \{ y \}$ to $\{ z \}$. Then we iterate over the edges $f$ that are incident to $y$. If $f$ is not of size three or $z\in f$, continue with the iteration. Else, iterate over the edges $g$ that are incident to $z$. If $g\setminus \{ z \} \subset f$, then we found a extended triangle situation. In that case, save $f$ in a variable $f_0$ and break out of the loop that iterates over the edges incident to $y$. If $f_0 \neq nil$ then put $f_0 \cup \{ z \}$ into the partial hitting set $C$ and break out of the loop iterating over the endpoints of $e$.

$$
\begin{algorithm}[H]

\KwIn{A hypergraph $G = (V, E)$, a set $C$}
\KwOut{An integer denoting the number of rule applications.}
\BlankLine

$exec \gets 0$\;
$outer \gets true$\;
\For{$outer$}{
    $outer \gets false$\;

    \For{$e \in E$}{
        \If{$|e| \neq 2$}{
            \textbf{continue}\;
        }

        \For{$a\in e$}{
            $y\gets a$\;
            $\{ z \} \gets e\setminus \{ y \}$\;
            $f_0 \gets nil$\;

            \textbf{incy:}\;
            \For{$f$ \textnormal{incident to} y}{
                \If{$|f|\neq 3 \textnormal{\textbf{ or }} z\in f$}{
                    \textbf{continue}\;
                }

                \For{$g$ \textnormal{incident to} $z$}{
                    $cond \gets true$\;
                    \For{$b \in g$}{
                        \If{$b = z$}{
                            \textbf{continue}\;
                        }
                        \If{$b\notin f$}{
                            $cond \gets false$\;
                            \textbf{break}\;
                        }
                    }
                    \If{cond}{
                        $f_0 \gets f$\;
                        \textbf{break incy}\;
                    }
                }
            }
            \If{$f_0 \neq nil$}{
                $outer \gets true$\;
                $exec \gets exec +1$\;
                $C\gets C\cup f \cup \{ z \}$\;
                $V\gets V\setminus f\cup \{ z \}$\;
                \For{$h$ \textnormal{incident to} $f\cup \{ z \}$}{
                    $E\gets E\setminus \{ h \}$\;
                }
                \textbf{break}\;
            }
        }
    }
}


\KwRet exec\;

\caption{Algorithm for exhaustive application of Extended Triangle Rule\label{etri}}
\end{algorithm}
$$

### Small Edge Degree 2 Rule

-   small edge degree 2: Let $v$ be a vertex of degree 2, and let the two hyperedges
    containing $v$ be $e = \{x, v\}$ and $f = \{v, y, z\}$. Then we can select a hyperedge $g$
    that contains one of the neighbors of $v$ in $f$ but not $x$, for example $g = \{u, w, z\}$
    (when $y = w$ is possible as a special case) or $g = \{u, z\}$. We put $x, u$ and $z$
    and $w$ (when existing) into the hitting set.

$O(|V|\cdot |E|)$ **Algorithm** We start by iterating other all vertices in $V$. If the vertex $v$ in the current iteration is of degree two, check if there is a size two edge $e$ and a size three edge $f$ incident to $v$. If these two edges exist, save $e\setminus \{ v \}$ in a variable $x$. Then iterate over the vertices $w$ in $f\setminus{v}$ and check if there exists an edge $h$ incident to $w$, such that $h\neq f$ and $x\notin h$. If such an edge exists, then we found a small edge degree 2 situation. We then put both $x$ and all vertices of $h$ into the partial hitting set.

$$
\begin{algorithm}[H]

\KwIn{A hypergraph $G=(V,E)$, a set $C$}
\KwOut{An integer denoting the number of rule applications.}
\BlankLine

$exec \gets 0$\;
$outer \gets true$\;

\For{outer}{
    \For{$v \in V$}{
        \If{$deg(v) \neq 2$}{
            \textbf{continue}\;
        }

        $s2,s3 \gets nil$\;
        \For{$e$ \textnormal{incident to} $v$}{
            \uIf{$|e| = 3$}{
                $s3 \gets e$\;
            }
            \uElseIf{$|e| = 2$}{
                $s2 \gets e$\;
            }
        }

        \If{$s2 = nil \textnormal{\textbf{ or }} s3 =nil$}{
            \textbf{continue}
        }

        $\{ x,\_ \} \gets s2$\;
        $found \gets false$\;
        $rem \gets nil$\;

        \For{$w\in s3\setminus \{ v \} $}{
            \For{$f$ \textnormal{incident to} $w$}{
                \eIf{$x \in f \textnormal{\textbf{ or }} s3 = f$}{
                    \textbf{continue}
                }{
                    $found \gets true$\;
                    $rem\gets f$\;
                    \textbf{break}\;
                }
            }
            \If{$found$}{
                \textbf{break}\;
            }
        }
        \If{$found$}{
            $outer \gets true$\;
            $exec \gets exec +1$\;
            \For{$a \in \{ x \}\cup rem$}{
                $C\gets C\cup \{ a \}$\;
                \For{$h$ \textnormal{incident to} $a$}{
                    $E\gets E \setminus \{ h \}$\;
                }
            }
        }
    }
}

\caption{Algorithm for exhaustive application of Small Edge Degree 2 Rule\label{sed2}}
\end{algorithm}
$$

### F3 Low Degree Rule

-   F3 low degree: Let $v$ be a vertex with degree 2 and let $e$ be an edge that is incident to $x$. If there exists an edge $f$ that is incident to a vertex in $e$ and does not contain $v$, put $f$ into the hitting set.

As seen in algorithm \ref{f3target} we extended this rule in the final implementation. If there is no vertex for which the rule applies, then we check the condition of the rule once again with the vertex $v_{min}$. If this attempt does not work, then we put a random size three edge into the hitting set. The effectiveness of this rule will be discussed in section \ref{appl}.

$$
\begin{algorithm}[H]

\KwIn{A hypergraph $G=(V,E)$, a set $C$}
\KwOut{An integer denoting the number of rule applications.}
\BlankLine

$min \gets \infty$\;
$v_{min} \gets nil$\;
$f_0 \gets nil$;

\textbf{check:}\;
\For{$x \in V$}{
    \uIf{$deg(x) = 2$}{
        $found \gets false$\;
        \textbf{incx:}\;
        \For{$e$ \textnormal{incident to} $x$}{
            \For{$v \in e\setminus \{ x \}$}{
                \For{$f\neq e$ \textnormal{incident to} $v$}{
                    \If{$x\notin f$ \textnormal{and} $|f| = 3$}{
                        $found\gets true$\;
                        $f_0 \gets f$\;
                        \textbf{break incx}\;
                    }
                }
            }
        }
        \If{found}{
            $C\gets C\cup f_0$\;
            \For{$v \in f_0$}{
                \For{$e$ \textnormal{incident to} $e$}{
                    $E \gets E\setminus \{ e \}$\;
                }
            }
            \KwRet 1\;
        }
    }
    \uElseIf{$deg(x) < min$ \textbf{and} $deg(x) > 1$}{
        $min \gets deg(x)$\;
        $v_{min} \gets x$
    }
    \BlankLine
    jump to label \textbf{incx} with $x \gets v_{min}$\;

    \If{$f_0 = nil$}{
        $f_0 \gets$ random size three edge\;
    }

    \KwRet $f_0 = nil$\;
}

\caption{Algorithm that selects a single size three edge to put into the hitting set\label{f3target}}
\end{algorithm}
$$

## Self-Monitoring

Each of the reduction rule functions returns a $\texttt{int32}$ value, which indicates the number of rule executions. We store the ratios for each rule in a map of the form:

```go
var Ratios = map[string]pkg.IntTuple{
	"kRulename": {A:1, B:1},
}
```

Where $\texttt{A}$ denotes the number of vertices put into the partial solution by a single rule execution. And $\texttt{B}$ denotes the number of vertices present in an optimal solution. We can then use these values to calculate the estimated approximation factor as follows.

```go
g := NewHypergraph()
c := make(map[int32]bool)
execs := make(map[string]int)

ApplyRules(g, c, execs)

var num float64 = 0
var denom float64 = 0

for key, val := range execs {
	num += float64(Ratios[key].A * val)
	denom += float64(Ratios[key].B * val)
}

ratio := num / denom
```

# Algorithms

## Main Algorithm

It is crucial to apply the rules in a given order. Some rules expect that other rules cannot be applied. As for the exact rules, one can save execution time, if a specific order is enforced when executing them. We use the order proposed in the original paper [@10.1007/978-3-642-29116-6_6], also referred to as precedence of rule executions. We will just use the shorter term _precedence_ going forward.

1. Exact Rules: vertex domination $\rightarrow$ tiny edge $\rightarrow$ edge domination
2. approximate (double) vertex domination rules
3. small edge degree 2 rule
4. small triangle rule
5. extended triangle rule
6. small edge rule

The main algorithm will apply all rules exhaustively according to the precedence, possibly mutating the input graph. If the graph has no more edges, then we are done. If not, put a size three edge into the partial hitting set and start over with the rule application. Note that the step of "applying all rules exhaustively" leaves room for interpretation and that the produced hitting sets are highly dependant on the actions taken during this step. We will refer to the procedure in this step as _rule strategy_ or just _strategy_ if the context allows it. We now present three strategies to apply the rules in the precedence.

The base strategy is quite simple. Execute a rule exhaustively and then do the same for the next rule in the precedence. Do this until no more rules can be applied. This is obviously not very optimal, since we do not check, wether a "better" rule can be applied again, before possibly executing a worse rule. We can however add some heuristics that preserve the fast execution time and improve on effectiveness. Instead of executing the vertex domination $\rightarrow$ tiny edge $\rightarrow$ edge domination cascade only once, we execute it three times. By just doing this, we can recoup a lot of cascades we might have missed with the base strategy. Iteration values higher than three did not yield better results. We should also execute this cascade after rules that are applied the most. Since these rules are the ones most likely to make the extact rules applicable. That would be the approximative (double) vertex domination rules.
The second strategy is also derived from the same base strategy. The only modification is that we start over with the vertex domination rule, whenever one of the rules was applied at least once. This should in theory lead to more exact rules being executed, while sacrificing a bit of execution time.
The third strategy tries to maximize the exact rule executions. We start by applying the exact rules exhaustively. We then execute the other rules at most once, according to the precedence. If a rule has been applied, we execute the exact rules exhaustively and start over with the first non-exact rule in the precedence.

## Incremental Frontier Algorithm

As experienced with the DBLP coauthor graph, there are some graph instances that do not work well with the prior algorithm. These instances do not admit any rule executions after applying the fallback rule. Running the algorithm on the whole graph again, knowing that only small parts of the graph have changed is not very time efficient. The algorithm should only look at the part of the graph where the fallback rule, or in fact any rule, was applied. We therefore propose a new approach, which involves the usage of a _vertex frontier_. Such frontiers are common in search algorithms such as BFS, where each vertex in the frontier has the same distance to the root vertex. But let us explain the general structure of the algorithm first.

We first apply all reduction rules exhaustively for the entire graph $G$. If the graph has no more edges, we are done. Else we try to find a size three edge $e$, which will trigger a vertex domination situation if removed. We choose a random size three edge if there is none. We then build up our initial frontier. We put all vertices into the frontier that are adjacent to a vertex in $e$, excluding vertices in $e$ itself. Then remove all edges that are incident to vertices in $e$. Next the frontier will be expanded, by adding all edges incident to the frontier to a new graph $H$. All vertices in $H$ that were not part of the frontier will form the next frontier. The amount of expansion steps can be set by the user\todo{proof that expansion step of 1 is sufficient ?}. The fields \texttt{AdjCount} and \texttt{IncMap} of $G$ will be reused by $H$. Thus, changes in $H$ will be reflected in $G$ and vice versa. The main loop will be explained next.

We apply the rules as usual, but only on $H$. The rules are modified, such that they also return the vertices adjacent to vertices in a modified/removed edge. If there are any, then $H$ will be expanded at these vertices and the loop will be continued. If not, apply the targeted fallback rule, considering the edges in the original graph $G$ and expand accordingly. This will be repeated until $G$ has no more edges.

$$
\begin{algorithm}[H]

\KwIn{Hypergraph $G = (V, E)$, and a partial hitting set $C$}
\KwOut{A hitting set $C$}
\BlankLine

Apply all reduction rules exhaustively mutating $G$ and $C$\;
$l\gets 2$\;

\If{$|E| = 0$}{
    \KwRet $C$\;
}

$H \gets$ empty hypergraph\;
\tcc{$H$ will act as a mask over $G$, the rules will only be applied to vertices/edges in $H$}

\For{$|E| > 0$}{
    Apply all reduction rules exhaustively on $H$ mutating $H$,$G$ and $C$\;
    $exp\gets \textnormal{ vertices of edges adjacent to removed or modified edge }$\;
    \If{$|exp| > 0$}{
        $H\gets \texttt{ExpandFrontier}(G, l, exp)$\;
        \textbf{continue}\;
    }

    $e\gets \texttt{F3LowDegree}(G)$\;
    \If{$e = \emptyset$}{
        \textbf{continue}\;
    }
    $H\gets \texttt{ExpandFrontier}(G, l, e)$\;
}

\KwRet $C$\;

\Fn{\ExpandFrontier{G, l, exp}}{
    $H\gets$ empty hypergraph\;

    \For{$i\gets 0$ \KwTo $l$}{
        $next \gets \emptyset$\;
        \For{$v\in exp$}{
            \For{$e \in G.E$ \textnormal{incident to} $v$}{
                \If{$e \notin H.E$}{
                    $H.E\gets H.E \cup \{ e \}$\;
                    \For{$w\in e$}{
                        \If{$w \notin H.V$}{
                            $H.V \gets H.V \cup \{ w \}$\;
                            $next \gets next \cup \{ w \}$\;
                        }
                    }
                }
            }
        }
        \If{$|next | = 0$}{
            \textbf{break}\;
        }
        $exp\gets next$\;
    }
    $H.IncMap = G.IncMap$\;
    $H.AdjCount = G.AdjCount$\;
    \KwRet $H$\;
}

\caption{Incremental 3-approximation algorithm for 3-HS\label{3hs_frontier}}
\end{algorithm}
$$

Algorithm \ref{3hs_frontier} was designed around both rule strategy 1 and 2. It is not compatible with strategy 3, since the auxiliary graph could at some point be empty and thus $|exp| = 0$, with the original graph still admitting rule applications. In such a case, we apply the first applicable rule in the precedence and use the neighborhood of the removed vertices to rebuild the auxiliary graph.

# Applications and Results {#appl}

All tests were run on a machine with a AMD Ryzen 5 3600 3.6 GHz six-core, 12-thread CPU and 32 gigabytes of 3200 MT/s DDR4 RAM. The machine was running Ubuntu Server 22.04.04 LTS and the binaries were compiled with Go version 1.22.1. Measured execution times do not include the time it takes to load in a graph from disk or the time needed to transform the graph to a new problem instance. We provide results for $\textsc{Triangle Vertex Deletion}$ and $\textsc{Cluster Vertex Deletion}$ problems using real world networks, as well as randomly generated hypergraphs.

## F3 Low Degree Rule

The goal of the F3 low degree rule is to reduce the degree of a vertex, which preferably has a degree of two. Applying the rule would lower the degree of the vertex to one, making the vertex domination rule applicable. We tested the effectiveness of the F3 low degree rule on random ER hypergraphs. First, we generated 100 random ER hypergraphs with 1000 vertices and approximately 3000 edges. We ran our algorithm on every graph 10 times each to account for run-to-run variance. We collected data for two datasets. For the first dataset we selected a random size three edge when applying the fallback rule. For the second dataset we used the F3 low target degree rule instead. The results can be seen in table \ref{er3_evr3_stats}.

We do in fact see a slight increase of vertex domination rule executions, with the mean number of executions increasing from about 345 to 349. More interesting observations can be made about the fallback rule and extended triangle rule. The mean number of fallback rule executions decreased from about 44 to 25 when not selecting random edges. The number of extended triangle rule executions increased from about 47 to 63. The mean values for the over rules stay mostly the same. 

To investigate these changing values, we reran the algorithm with both fallback variants on some of the graphs. This time keeping a log of all the rules that are applied. Every time the F3 low degree rule is executed, the $F3\rightarrow VDom\rightarrow ETri$ cascade is most likely triggered. For the random variant, one can see batches of fallback rule executions and the occasional $F3\rightarrow VDom\rightarrow ETri$ cascade in the logs. 

So what happens is, that the second variant substitutes multiple fallback rule executions with a single $F3\rightarrow VDom\rightarrow ETri$ cascade. The desired effects are underwhelming, with the mean hitting set size decreasing from about 592 to 590. Since we did not regress, we chose to keep using this variant of the rule over the random variant.

## Triangle Vertex Deletion

$$
\begin{center}
\fbox{
    \begin{minipage}[h]{0.75\textwidth}
        \textbf{Problem Name:} \textsc{Triangle Vertex Deletion}

        \textbf{Given:} A graph $G=(V,E)$ and a non-negative integer $k$.

        \textbf{Output:} Is there a set $C\subseteq V$ of size at most $k$, such that $G\setminus C$ is a triangle-free graph.
    \end{minipage}
}
\end{center}
$$

\textsc{Triangle Vertex Deletion} can be reduced to the \textsc{3-Hitting Set} problem. Triangles in the input graph will be represented as size 3 edges in a hypergraph, containing the vertices that make up the triangle. A hitting set of size at most $k$ on this hypergraph can be used to transform the original standard graph into a triangle-free graph, removing at most $k$ vertices.

### DBLP Coauthor Graph

$$
\begin{figure}[t]
    \centering
    \includegraphics[width=\textwidth]{./img/dblp_tvd_rules_strat.png}
    \caption{Mean number of rule executions for \textsc{Triangle Vertex Deletion} on DBLP coauthor graph per rule strategy, $n=100$}
    \label{dblp_rules_strat}
\end{figure}
$$

We evaluated the algorithms on a DBLP coauthor graph and ER graphs of varying densities. It was considered to use the complete DBLP coauthor graph, which was quite large. We ultimately decided to use a dataset of the largest connected component. The network is made available by the [Stanford Network Analysis Project](https://snap.stanford.edu/data/com-DBLP.html) and part of a paper by J. Yang and J. Leskovec [@DBLP:journals/kais/YangL15] about community detection in real-world networks. The network contains 2224385 triangles, resulting in a \textsc{3-Hitting Set} instance of the same size. We now present the collected data during a run of the algorithm on this hypergraph, which can be found in table \ref{stats_dblp}.

Rule strategy 1 has the fastest execution time with a mean of 29 seconds, but also the worst ratio with a mean of $1.5708$. Rule strategy 3 has the best ratio with a mean of $1.4195$, but also the worst execution time of about 8 minutes. We used the frontier technique in all the tests, since we could observe a speedup of about 90, just using the base rule strategy. This speedup can be attributed to the fact that the original coauthor graph exerts a community structure. These communities are preserved when converting it to a $\textsc{Triangle Vertex Deletion}$ and thus $\textsc{3-Hitting Set}$ instance. The expansion steps will only expand the auxiliary graph $H$ inside these communities. This drastically reduces the number of edges the algorithm has to process per iteration. One can see the opposite happen in dense hypergraphs. Due to the high density, even a small expansion two neighborhoods deep will add all edges of $G$ to $H$, eliminating possible gains. In fact, constant rebuilding of the auxiliary graph $H$ will even worsen the execution time. This can be observed when using any rule strategy besides strategy 3 on a dense 3-uniform ER hypergraph. On inspection of the generated CPU profile, around 30% of the time is spent computing the auxiliary graph $H$. See figure \ref{flamegraph} for a flamegraph obtained from one of those profiles.

Figure \ref{dblp_rules_strat} displays the number of rule executions per strategy. Note that the there is a rule called SED2* in this chart. Recall that the small edge degree 2 rule can either put four or three vertices into the hitting set. Both cases were counted seperately, where SED2* refers to the version that puts three vertices into the hitting set. There are several interesting observations to be made, for example the decrease of edge domination rule executions for strategy 2 compared to strategy 1.\todo{explanation ?} Quite frankly we do not know why this happens. We can also see that strategy 3 does indeed increase the amount of exact rule executions. Another observation can be made about the number of small edge rule executions. Strategy 2 does not apply the small edge rule at all, which is surprising since it is conceptually similar to strategy 3. This is probably a result of only applying rules one at a time using strategy 3. The exact rules that are executed after each rule application will possibly destroy situations that are favorable for other rules. The approximative vertex domination rule, which is executed a lot more with strategy 2 compared to 3, might have removed all the size two edges as a byproduct. The number of fallback rule executions is also interesting. All strategies are executing this rule around 3000 times. This leads us to the hypothesis that as long as the order in the precedence is obeyed, a "worse" strategy will not result in significantly more fallback rule executions. It could also be true that this number stays the same for any precedence. \todo{get data for random rule application to compare}

$$
\begin{table}[t]
    \begin{subtable}[b]{0.45\textwidth}
        \centering
        \caption{base rule strategy\label{dblp_str_base}}
        \begin{tabular}{lrrrr}
            \toprule
            & ratio & $|C|$ & est. opt & time \\
            \midrule
            mean & 1.7594 & 115672 & 65745 & 17 sec \\
            std & 0.0005 & 79 & 41 & 1 sec\\
            min & 1.7577 & 115469 & 65635 & 16 sec\\
            median & 1.7594 & 115677 & 65741 & 17 sec\\
            max & 1.7610 & 115998 & 65904 & 19 sec\\
            \bottomrule
        \end{tabular}
    \end{subtable}
    \hfill
    \begin{subtable}[b]{0.45\textwidth}
        \centering
        \caption{rule strategy 1\label{dblp_str_1}}
        \begin{tabular}{lrrrr}
            \toprule
            & ratio & $|C|$ & est. opt & time\\
            \midrule
            mean & 1.5708 & 106339 & 67697 & 29 sec\\
            std & 0.0006 & 65 & 40 & 1 sec\\
            min & 1.5691 & 106183 & 67597 & 26 sec\\
            median & 1.5708 & 106342 & 67699 & 29 sec\\
            max & 1.5727 & 106500 & 67797 & 31 sec\\
            \bottomrule
        \end{tabular}
    \end{subtable}
    \newline
    \vspace{4mm}
    \newline
     \begin{subtable}[b]{0.45\textwidth}
        \centering
        \caption{rule strategy 2\label{dblp_str_2}}
        \begin{tabular}{lrrrr}
            \toprule
            & ratio & $|C|$ & est. opt & time \\
            \midrule
            mean & 1.5437 & 105084 & 68071 & 104 sec \\
            std & 0.0006 & 63 & 35 & 2 sec\\
            min & 1.5422 & 104922 & 67956 & 99 sec\\
            median & 1.5438 & 105084 & 68072 & 104 sec\\
            max & 1.5453 & 105261 & 68146 & 110 sec\\
            \bottomrule
        \end{tabular}
    \end{subtable}
    \hfill
    \begin{subtable}[b]{0.45\textwidth}
        \centering
        \caption{rule strategy 3\label{dblp_str_3}}
        \begin{tabular}{lrrrr}
            \toprule
            & ratio & $|C|$ & est. opt & time \\
            \midrule
            mean & 1.4195 & 99456 & 70066 & 474 sec\\
            std & 0.0009 & 52 & 31 & 7 sec\\
            min & 1.4171 & 99323 & 69987 & 457 sec\\
            median & 1.4196 & 99450 & 70069 & 474 sec\\
            max & 1.4215 & 99575 & 70133 & 493 sec\\
            \bottomrule
        \end{tabular}
    \end{subtable}
    \caption{Results for \textsc{Triangle Vertex Deletion} on DBLP coauthor graph; $n=100$}
    \label{stats_dblp}
\end{table}
$$

### Amazon Product Co-Purchasing Graph

To validate that our previous findings were not just due to the structure of the DBLP coauthor graph, we conduct the same tests with an unrelated network. We chose an Amazon product co-purchasing graph, which had a similar number of vertices and edges compared to the DBLP coauthor graph. The resulting \textsc{Triangle-Vertex-Deletion} instance only contains 667129 edges and thus should be structurally distinct, compared to the instance obtained from the DBLP coauthor graph. As seen in table \ref{stats_amzn} and figure \ref{amzn_rules_strat}, the algorithm achieves a mean estimated ratio of $1.3136$ and behaves the same as for the DBLP coauthor graph.

$$
\begin{figure}[t]
    \centering
    \includegraphics[width=\textwidth]{./img/amazon_tvd_rules_strat.png}
    \caption{Mean number of rule executions for \textsc{Triangle Vertex Deletion} on Amazon product co-purchasing graph per rule strategy, $n=100$}
    \label{amzn_rules_strat}
\end{figure}
$$

$$
\begin{table}[t]
    \begin{subtable}[b]{0.45\textwidth}
        \centering
        \caption{base rule strategy\label{amzn_str_base}}
        \begin{tabular}{lrrrr}
            \toprule
            & ratio & $|C|$ & est. opt & time \\
            \midrule
            mean & 1.7466 & 95256 & 54539 & 3 sec\\
            std & 0.0003 & 82 & 41 & 0 sec\\
            min & 1.7458 & 95032 & 54433 & 3 sec\\
            median & 1.7466 & 95269 & 54543 & 3 sec\\
            max & 1.7474 & 95441 & 54632 & 4 sec\\
            \bottomrule
        \end{tabular}
    \end{subtable}
    \hfill
    \begin{subtable}[b]{0.45\textwidth}
        \centering
        \caption{rule strategy 1\label{amzn_str_1}}
        \begin{tabular}{lrrrr}
            \toprule
            & ratio & $|C|$ & est. opt & time \\
            \midrule
            mean & 1.4731 & 86359 & 58624 & 6 sec\\
            std & 0.0005 & 61 & 36 & 0 sec\\
            min & 1.4719 & 86213 & 58534 & 6 sec\\
            median & 1.4731 & 86359 & 58627 & 6 sec\\
            max & 1.4741 & 86496 & 58723 & 6 sec\\
            \bottomrule
        \end{tabular}
    \end{subtable}
    \newline
    \vspace{4mm}
    \newline
    \begin{subtable}[b]{0.45\textwidth}
        \centering
        \caption{rule strategy 2\label{amzn_str_2}}
        \begin{tabular}{lrrrr}
            \toprule
            & ratio & $|C|$ & est. opt & time \\
            \midrule
            mean & 1.4215 & 84813 & 59666 & 38 sec\\
            std & 0.0005 & 58 & 32.3300 & 0 sec\\
            min & 1.4206 & 84670 & 59586 & 37 sec\\
            median & 1.4215 & 84814 & 59664 & 38 sec\\
            max & 1.4227 & 84927 & 59740 & 39 sec\\
            \bottomrule
        \end{tabular}
    \end{subtable}
    \hfill
    \begin{subtable}[b]{0.45\textwidth}
        \centering
        \caption{rule strategy 3\label{amzn_str_3}}
        \begin{tabular}{lrrrr}
            \toprule
            & ratio & $|C|$ & est. opt & time \\
            \midrule
            mean & 1.3136 & 80829 & 61532 & 166 sec\\
            std & 0.0006 & 59 & 34.6300 & 2 sec\\
            min & 1.3121 & 80684 & 61444 & 160 sec\\
            median & 1.3136 & 80836 & 61535 & 166 sec\\
            max & 1.3149 & 80943 & 61616 & 171 sec\\
            \bottomrule
        \end{tabular}
    \end{subtable}
    \caption{Results for \textsc{Triangle Vertex Deletion} on Amazon product co-purchasing graph; $n=100$}
    \label{stats_amzn}
\end{table}
$$

## Cluster Vertex Deletion

Another derivative problem of \textsc{Hitting Set} is the \textsc{Cluster Vertex Deletion} problem. The problem can be formulated as follows,

$$
\begin{center}
\fbox{
    \begin{minipage}[h]{0.75\textwidth}
        \textbf{Problem Name:} \textsc{Cluster Vertex Deletion}

        \textbf{Given:} A graph $G=(V,E)$ and a non-negative integer $k$.

        \textbf{Output:} Is there a set $C\subseteq V$ of size at most $k$, such that $G\setminus C$ is a cluster graph.
    \end{minipage}
}
\end{center}
$$

A cluster graph is a union of disjoint complete graphs which is equivalent to a graph that does not contain an induced path on three vertices, also called a $P_3$. These $P_3$ can be interpreted as edges in a hypergraph. A hitting set of size at most $k$ on this hypergraph can be used to transform the original standard graph into a cluster graph, removing at most $k$ vertices.

We tested our algorithm on the _Rome Graphs_, a collection of 11534 undirected graphs used and made available by Di Battista et al. in [@DiBattista1997b]. We ran the algorithm ten times for each derived problem instance and picked the run with the lowest hitting set size to include during calculation of the statistics. The results can be seen in table \ref{cvd_rome_stats}. The algorithm achieved a mean ratio of $1.1188$ with a standard deviation of about $0.1$. We also recorded a ratio of at most 2 for the selected problem instances. It can be observed and inferred from the low ratio that the most applied rules are exact rules.

$$
\begin{table}[h]
\makebox[\textwidth][c]{
    \begin{tabular}{lrrrrrrrrrrrrrr}
    \toprule
    & ratio & Tiny & VD & ED & Small & Tri & ETri & AVD & ADVD & SED2 & SED2* & F3 & $|C|$ & est. opt \\
    \midrule
    mean & 1.1188 & 13.77 & 31.58 & 84.01 & 0.21 & 0.05 & 0.01 & 1.38 & 0.07 & 0 & 0.08 & 0 & 18 & 16 \\
    std & 0.0962 & 7.20 & 15.75 & 56.79 & 0.45 & 0.23 & 0.09 & 1.39 & 0.26 & 0.07 & 0.27 & 0.01 & 9 & 8 \\
    min & 1 & 0 & 0 & 0 & 0 & 0 & 0 & 0 & 0 & 0 & 0 & 0 & 2 & 2 \\
    median & 1.1111 & 12 & 28 & 75 & 0 & 0 & 0 & 1 & 0 & 0 & 0 & 0 & 16 & 14 \\
    max & 2 & 32 & 69 & 310 & 5 & 2 & 1 & 9 & 2 & 1 & 2 & 1 & 40 & 33 \\
    \bottomrule
    \end{tabular}
}
    \caption{Results for \textsc{Cluster Vertex Deletion} on graphs from the Rome Graphs collection. Hitting sets were computed ten times for each graph and the run with the lowest hitting set size was chosen for final calculation.}
    \label{cvd_rome_stats}
\end{table}
$$

## Comparison with LP rounding based Algorithms

There exist randomized approximation algorithms for $d-\textsc{Hitting Set}$ that rely on LP rounding, which promise quite competitive approximation ratios. They are also fast since most of the work happens during the solving of the LP. We will begin by explaining the ILP/LP formulation of $\textsc{Hitting Set}$ and its dual problem $\textsc{Set Cover}$, since one of the algorithms is actually a $\textsc{Set Cover}$ algorithm. After that we will briefly examine the two proposed algorithms and end with a comparison to our algorithm. Note that we will sometimes use slightly different notation and variable names than the original papers to avoid confusion (for example $\Delta$ is defined in both, in one as maximum vertex degree and the other as maximum edge size).

Let $G=(V,E)$ be a hypergraph with $n:=|V|$ and $m:=|E|$. Let $V= \left\{ v_1, v_2, \dots \right\}$ and $E = \left\{ e_1, e_2, \dots \right\}$ be the set of vertices and edges respectively. Abusing notation, we also allow referring to a vertex or edge using its index, i.e. vertex $1$ and vertex $v_1$ shall mean the same. The $\textsc{Hitting Set}$ problem can now be formulated as an integer linear program as follows,

$$
\begin{alignat*}{2}
    & \textnormal{minimize}\quad   && \sum_{i=1}^{n} x_i \\
    & \textnormal{subject to}\quad && \sum_{i \in e_j} x_i \geq 1 \quad \forall e_j \in E \\
    &                              && x_i \in \left\{0,1\right\} \quad \forall i \in [n]\\
\end{alignat*}
$$

The LP relaxation of this ILP would allow for $x_i$ to be any value in the interval $[0,1]$ for all $i \in [n]$. The dual to the problem, the $\textsc{Set Cover}$ problem, can be formulated similarly for the unweighted case,

$$
\begin{alignat*}{2}
    & \textnormal{minimize}\quad   && \sum_{j=1}^{m} x_j \\
    & \textnormal{subject to}\quad && \sum_{j: i\in e_j} x_j \geq 1 \quad \forall i \in [n] \\
    &                              && x_j \in \left\{0,1\right\} \quad \forall j \in [m]\\
\end{alignat*}
$$

And again, the LP relaxation of this ILP would allow for $x_j$ to be any value in the interval $[0,1]$ for all $j \in [m]$.

It might not obvious at first that \textsc{Set Cover} is the dual of \textsc{Hitting Set}, but formulating both problems informally gives it away in an understandable way. One problem tries to cover sets with elements, while the over aims at covering elements with sets. For a formal reduction from \textsc{Hitting Set} to \textsc{Set Cover} let $G=(V,E)$ be a hypergraph and $C$ be a subset of $V$. Let $S_v = \left\{ e \in E \mid v \in e \right\}$ for $v\in V$, i.e. the set of edges that are incident to a vertex $v$. Now $C$ is a hitting set if and only if $S = \left\{ e \in S_v \mid v \in C\right\}$ is a set cover. Let us assume that $C$ is a hitting set for $G$, then it holds that,

$$
\forall e \in E \mid e \cap C \neq\emptyset \quad\Leftrightarrow\quad \forall e \in E\; \exists v \in C \;\; \textnormal{s.t.} \;\; e \in S_v \quad \Leftrightarrow \quad \bigcup_{v\in C} S_v = E
$$

The first algorithm we look at was proposed by Ouali et. al. in [@DBLP:journals/tcs/OualiFS14]. Let $G=(V,E)$ be a hypergraph. The algorithm starts by solving the LP relaxation of the hitting set ILP for $G$. Let $\textnormal{Opt}^* = \sum_{i=1}^{n} x_i^*$ be an optimal solution to the LP relaxation. With that solution construct the following four sets:

$$
\begin{align*}
     &S_0 = \left\{ i \in [n] \mid x_i^* = 0 \right\}
     &S_\geq = \left\{ i \in [n] \mid 1\neq x_i^* \geq \tfrac{1}{\lambda} \right\}\\                                                                           \\
     &S_1 = \left\{ i \in [n] \mid x_i^* = 1 \right\}
     &S_{<} = \left\{ i \in [n] \mid 0\neq x_i^* < \tfrac{1}{\lambda} \right\}\\
\end{align*}
$$

Remove all vertices in $S_0$ from $V$ and remove the vertices in $S_0$ from all edges. Then put all vertices in $S_1$ and $S_\geq$ into the hitting set $C$, by removing all vertices in $S_1 \cup S_\geq$ from $V$ and removing all edges incident to $S_1 \cup S_\geq$ from $E$. Next is the randomized rounding step. For all vertices $i \in S_<$ include $i$ in the hitting set with probability $\lambda \cdot x_i^*$, independently for all $i$. If $|E|=0$ return the hitting set $C$. Else, as long as $|E|>0$, select a vertex from an uncovered edge and put it into the hitting set $C$. We will further elaborate on the definition of $\lambda$ when discussing the test results of the algorithms.

Next, we look at the algorithm proposed by Saket and Sviridenko in [@DBLP:conf/approx/SaketS12]. Let $V$ be a ground set and $E \subseteq 2^V$ a set system. We also want to define two variables, $\Delta$ the maximum number of sets any element in $V$ is contained in and $d$ the maximum size of any set in $E$. The algorithm starts by solving the LP relaxation to the $\textsc{Set Cover}$ ILP for $(V,E)$. Let $\textnormal{Opt}^* = \sum_{j=1}^{m} x_j^*$ be an optimal solution to the LP relaxation. For all $j \in [m]$, choose to include set $e_j$ with probability $p_j=\min\left\{1, \alpha\Delta \cdot x_j^*\right\}$, where $\alpha = 1-e^{-\frac{\ln d}{\Delta-1}}$. Let $I^r$ be the set of elements of $V$ that are uncovered. For all $v \in I^r$, choose a set with the lowest weight in $E$ that contains $v$ and include it in the cover.

We used PuLP[@GH:pulp] to model the linear programs in Python and implemented the algorithms in Python as well. Both algorithms were tested with three open source solvers GLPK[@GH:glpk], CLP[@GH:clp] and HiGHS[@GH:highs] a comparatively newer solver that uses a parallel dual simplex method developed by Q. Huangfu and J. A. J. Hall [@DBLP:journals/mpc/HuangfuH18]. Table \ref{lp_stats} shows the results for both algorithms on the \textsc{Triangle Vertex Deletion} instance of the Amazon product co-purchasing graph, which is a 3-uniform hypergraph with maximum vertex degree of 551. Algorithm (1) promises a hitting set size of $|C| \leq d(1-\frac{d-1}{8\Delta})\cdot \textnormal{Opt}^*$ with probability $\frac{3}{4}$, under the assumption that the hypergraph is $d$-uniform and $3 \leq d \leq \frac{16}{3}\Delta$. Our hypergraph fulfills both of these requirements, which should result in a ratio of at most $3(1-\frac{3-1}{8 \cdot 551})\cdot \textnormal{Opt}^* = 2.9986 \cdot \textnormal{Opt}^*$. This upper bound is only valid for $\lambda = l(1-\epsilon)$ where $\epsilon = \frac{l\textnormal{Opt}^*-|S_1|}{2m}$. Algorithm (2) promises a hitting set size of $|C| \leq ((\Delta -1)(1-e^{-\frac{ln d}{\Delta -1}})+1) \cdot \textnormal{Opt}^*$. Applied to our set cover instance, this yields a ratio of $((3-1)(1-e^{-\frac{ln 551}{3-1}})+1)\cdot \textnormal{Opt}^* = 2.9148\cdot \textnormal{Opt}^*$.

$$
\begin{table}[t]
     \centering
     \begin{tabular}{c|rrrrrrr}
          \toprule
          algorithm                      & solver & method & ratio UB  & actual ratio &
          $\textnormal{Opt}^*$           & $|C|$  & solver time                                               \\
          \midrule
          \multirow{3}{4em}{Hitting Set LP (1)} & GLPK & simplex  & \multirow{3}{3em}{2.9986}         & 1.5814       & 64452.6941 &
          101928                         & 1 hour                                                             \\
                                         & CLP    &  simplex    &      & 1.7838       & 64452.6992 & 114972 & 50 sec \\
                                         & HiGHS & simplex &            & 1.6709       & 64452.6993 & 107697 & 17 sec \\
                                         & HiGHS & ipm &           & 1.5713       & 64452.6993 & 101276 & 76 sec \\
          \midrule
          \multirow{3}{4em}{Set Cover LP (2)}   & GLPK & simplex  & \multirow{3}{3em}{2.9148} & 1.5944       & 64452.6940 &
          102762                         & 1 hour                                                             \\
                                         & CLP    & simplex &           & 1.8068       & 64452.6992 & 116455 & 47 sec \\
                                         & HiGHS & simplex &            & 1.6826       & 64452.6993 & 108447 & 19 sec \\
                                         & HiGHS & ipm &            & 1.5797       & 64452.6993 & 101820 & 68 sec \\
     \end{tabular}
     \caption{Results for LP based rounding algorithms for \textsc{Hitting Set} and \textsc{Set Cover} on \textsc{Triangle Vertex Deletion} instance of Amazon product co-purchasing graph}
     \label{lp_stats}
\end{table}
$$

Algorithm (1) slightly outperforms algorithm (2) in terms of hitting set size. The LP solution obtained with the HiGHS solver using the interior point method yields the smallest hitting set for both algorithms, with a execution time of about 1 minute. The hitting set obtained with the GLPK solver comes close but at a way higher execution time of about 1 hour. The CLP solver is not very well suited for the algorithms on this specific problem instance. Execution time as well as the resulting hitting set size are worse compared to HiGHS and GLPK. Compared to our algorithm, both algorithm (1) and (2) yield hitting sets with bigger size.

$$
\begin{table}[h]
    \begin{subtable}[b]{0.45\textwidth}
        \begin{tabular}{lrrrr}
            \toprule
            & ratio UB & ratio & $|C|$ & $\textnormal{Opt}^*$ \\
            \midrule
            mean & 2.9697 & 1.4246 & 25 & 17 \\
            std & 0.0197 & 0.4268 & 16 & 8 \\
            min & 2.7500 & 1 & 2 & 2 \\
            median & 2.9758 & 1.3565 & 21 & 15 \\
            max & 2.9925 & 2.5714 & 80 & 35 \\
            \bottomrule
        \end{tabular}
        \caption{Hitting Set LP (1)}
    \end{subtable}
    \hfill
    \begin{subtable}[b]{0.45\textwidth}
        \begin{tabular}{lrrrr}
            \toprule
            & ratio UB & ratio & $|C|$ & $\textnormal{Opt}^*$ \\
            \midrule
            mean & 2.6132 & 1.4050 & 24 & 17 \\
            std & 0.1089 & 0.4101 & 16 & 8 \\
            min & 1.8453 & 1 & 2 & 2 \\
            median & 2.6408 & 1.3220 & 21 & 15 \\
            max & 2.8000 & 2.5385 & 76 & 35 \\
            \bottomrule
        \end{tabular}
        \caption{Set Cover LP (2)}
    \end{subtable}
    \caption{Results for \textsc{Cluster Vertex Deletion} using LP rounding based algorithms on graphs from the Rome Graphs collection. Same method as previous CVD benchmark.}
    \label{cvd_rome_stats_lp}
\end{table}
$$

\todo{update with new test results using other solvers}
We also conducted the same test from the Cluster Vertex Deletion section with both algorithms using the the same solvers and methods. This time the GLPK solver was best suited for the problem. The results can be seen in table \ref{cvd_rome_stats_lp}. Both algorithms perform very similarly, with almost identical values for $\textnormal{Opt}^*$ and the hitting set sizes. Both have a high maximum hitting set size, which is almost double the size compared to our algorithm. The $\textnormal{Opt}^*$ values for both are very similar to the estimated optimum values that our algorithm computes, but the difference between optimum and actual hitting set size is smaller with our algorithm. This is can be observed for every statistic we calculated. The most extreme case is the standard deviation for hitting set size and $\textnormal{Opt}^*$. We can see that the standard deviation for the hitting set size with a value of 16, is twice as big as the standard deviation for $\textnormal{Opt}^*$, for both algorithms.

## Preferential Attachment Hypergraphs

## 3-Uniform ER Hypergraphs

# Testing

Every reduction rule is tested for their correctness with unit tests. We create small graphs in these tests, which contain structures, which the rules are targeting. We then test for the elements in the partial solution, number of edges/vertices left in the graph and degree of the vertices left.

# Conclusion

The implmented \textsc{Hitting Set} algorithm can compute hitting sets for large inputs in a reasonable amount of time. The algorithm also yields smaller hitting sets compared to LP randomized rounding based algorithms. We also indirectly highlighted the weaknesses of these LP based algorithms. In order to guarantee the algorithm's approximation ratio, the probabilities during the randomized rounding have to be dependant on parameters of the hypergraph. Using constants for these probabilities will most likely lead to smaller hitting sets. 

## Further Questions

The only part of the algorithm that is parallelized is the edge domination rule. The incremental variant of the main algorithm could also be parallelized. The graph could at some point during the run of the algorithm consist of multiple connected components. These components could then be processed in parallel. Every component would yield its own number of rule executions and hitting sets, which are aggregated at the end. The question is, if the overhead of finding these components is greater than the possible speed up.

There already exist parallel algorithms for the vertex and edge domination rule, which run on a GPU as shown by Bevern et al. in [@DBLP:journals/jcss/BevernKSST24]. It would be interesting to know if such algorithms also exist for the other domination type reduction rules, namely the approximate (double) vertex domination rule. Our implementation of these rules rely on arguments about vertex degrees and adjacency between vertices, which is similar to the standard vertex domination implementation. 

Could our hypergraph data structure be further improved? Right now we only use Go primitives to model our hypergraph. Improvements could be made for both the incidence and adjacency map struct fields. When accessing elements of a map with the \texttt{range} operator, Go does not guarantee that the order of elements is always the same. Using a real sparse matrix could improve performance and make the algorithm more deterministic.
