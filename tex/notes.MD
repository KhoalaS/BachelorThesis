---
    title: Implementation and evaluation of a self-monitoring approximation algorithm for 3-Hitting-Set
    link-citations: true
    author: Khoa Le
...

\RestyleAlgo{ruled}
\DontPrintSemicolon
\SetAlgoVlined
\LinesNumbered
\listoftodos
\SetKwFunction{FRecurs}{FnRecursive}
\SetKwFunction{ExpandFrontier}{ExpandFrontier}
\SetKwProg{Fn}{func}{\string:}{}

# Programming Language

We chose the [_Go_](https://go.dev/) programming language. It is a statically typed, compiled language, with a C-like syntax. It is using a garbage collector to handle memory management, which at first seems off-putting for an application like this. Since Go's GC is very efficient, we are not worried about that fact. If the situation arises in which we do need to handle memory manually, we can utilize Go's $\texttt{unsafe}$ package in conjunction with C-interop.

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
```

The $\texttt{HyperGraph}$ datatype has five fields. Both fields $\texttt{Vertices}$ and $\texttt{Edges}$ are maps with keys of type $\texttt{int32}$ and values of type $\texttt{Vertex}$ and $\texttt{Edge}$ respectively. We chose this $\texttt{Set}$-like datastructure over lists again because of faster access times, but also operations that remove edges/vertices are built-in to the map type. The field $\texttt{edgeCounter}$ is an internal counter used to assign ids to added edges. The field $\texttt{IncMap}$ is a map of maps, essentially storing the hypergraph as a sparse incidence matrix. We will also derive vertex degrees from this map. And at last the $\texttt{AdjCount}$ map, which will associate every vertex $v \in V$ all vertices adjacent to $v$. Additionally this map will also store the amount of times such a vertex is adjacent to $v$.

# Misc. Algorithms/Utilities

## Edge Hashing

Time Complexity: $n + n\cdot \log(n)$, where $n$ denotes the size of $\texttt{arr}$.

We start by sorting $\texttt{arr}$ with a $\textsc{Quick-Sort}$-Algorithm. We then join the elements of $\texttt{arr}$ with the delimiter $"|"$, returning a string of the form $"|id_0|id_1|\dots|id_n|"$. Whenever we refer to $\textit{the hash of an edge}$ we refer to the output of this function, using the endpoints of the edge as the $\texttt{arr}$ argument.

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

Lists in Go are not very memory efficient, but since we exclusively call this function with $arr$ representing the vertices in an edge, the value is usually fixed at 3. The raised memory problems occur at values of $n>10000$, thus justifying the continued usage of lists. For the case where we have to compute a lot of subsets, we provide a slightly different version of this function. Instead of passing in the $subsets$ list, we pass in a callback function that shall be called whenever we find a subset, using the found subset as an argument.

## Two-Sum

Given an array of integers and an integer target $t$, return indices of the two numbers such that they add up to $t$.

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

We start by creating a map called $lookup$. Iterating other the elements of $arr$, we check if the entry $lookup[t-arr[i]]$ exists.

-   If the entry exists, we return a pair $(i, lookup[t-arr[i]])$ and the boolean value $true$ since we found a solution.
-   If the entry does not exist, we add a new entry to the $lookup$ map using $arr[i]$ as key and $i$ as value.

If no solution was found, we return $nil$ and the boolean value $false$.

This algorithm is an ingredient for the implementation of one of the reduction rules, specifically the approximative vertex domination rule. The actual implementation accepts a map instead of an array as its first parameter. We also implemented a version that finds all solutions, which accepts an additional callback function as a parameter. Instead of returning the solution, we call the callback function with the solution as the first argument. This version is an ingredient for the implementation of the approximative double vertex domination. We refer to this version as $\textssc{Two-SumAll}$ inside the algorithm listings.

# Hypergraph Models

## First Testing Model

```go
func GenerateTestGraph(n int32, m int32, tinyEdges bool) *HyperGraph
```

Let us explain the arguments first:

-   $\texttt{n}$ are the number of vertices the graph will have
-   $\texttt{m}$ is the amount of edges the graph will at most have
-   $\texttt{tinyEdges}$ when $\texttt{false}$ indicates that we do not want to generate edges of size 1.

We use a very naive approach for generating (pseudo-)random graphs. We first create an empty Hypergraph struct and add $\texttt{n}$ many vertices to that graph. We then compute a random $\texttt{float32}$ value $\texttt{r}$ in the half-open interval $[ 0.0, 1.0 )$. This value will be used to determine the size of an edge $e$. The edges are distributed based on their size as follows:

$$
size(\texttt{r})=
\begin{cases}
    1 & \quad \texttt{r} < 0.01\\
    2 & \quad 0.01 \leq \texttt{r} < 0.60\\
    3 & \quad \text{else}
\end{cases}
$$

The result of $size(\texttt{r})$ is stored in a variable $\texttt{d}$. We then randomly pick vertices in the half-open interval $[ 0, \texttt{n} )$, until we have picked $\texttt{d}$ many distinct vertices. If an edge with these endpoints does not exist, we add it to our graph.

That results in the graph having at most $\texttt{m}$ edges and not exactly $\texttt{m}$, since we did not want to artificially saturate the graph with edges. One could also look into generating random bipartite graphs that translate back to a hypergraph with the desired vertex and edge numbers.

## Preferential Attachment Hypergraph Model

In the Preferential Attachment Model, one will add edges to an existing graph, with a probability proportional to the degree of the endpoints of that edge. This edge will either contain a newly added vertex, or will be comprised of vertices already part of the graph.

We will use an implementation by Antelmi et al. as reference [@SimpleHypergraphs.js], which is part of their work on _SimpleHypergraphs.jl_ [@DBLP:journals/im/SpagnuoloCSPSKA20], a hypergraph software library written in the Julia language. The implementation is based on a preferential attachment model proposed by Avin et al. in [@DBLP:journals/corr/AvinLP15].

```go
func GeneratePrefAttachmentGraph(n int, p float64, maxEdgesize int32)
```

-   $\texttt{n}$ is the amount of vertices the graph will have
-   $\texttt{p}$ is the probability of adding a new vertex to the graph
-   $\texttt{maxEdgesize}$ is the maximum size of a generated edge

# Reduction Rules

The ususal signature of a reduction rule looks as follows:

```go
func NameRule(g HyperGraph, c map[int32]bool) int32
```

We take both a \texttt{HyperGraph} struct \texttt{g} and a _Set_ \texttt{c} as arguments and mutate them. We then return the number of rule executions. We prioritize time complexity over memory complexity when implementing rules, which does not equate to ignoring memory complexity completely.

## Executions

The proposed rules are meant to be applied exhaustively. A one to one implementation of a rule will only find one of the structures the rule is targeting. Calling such a rule implementation exhaustively will take polynomial time, but will be very inefficient in regards to memory writes and execution time. It is therefor necessary to design the algorithms for the rules with the aspect of exhaustive application in mind.

The general outline of an algorithm will look as follows,

1. Construct auxiliary data structures that are used to find parts of the graph, for which the rule can be applied.
2. As long as we can apply rules do:

    Iterate over the auxiliary data structure, or the vertices/edges themselves.

    - Identify targets of the rule.
    - Mutate the graph according to the rule.
    - Mutate the auxiliary data structure according to the rule.

This way we minimize memory writes by reusing existing data structures we built in step 1, but also save on execution time, since we mutate and iterate the auxiliary data structure at the same time.

## Algorithms

The algorithm descriptions occasionally omit implementation details. We do this to try and keep these descriptions as concise as possible. The real implementations correspond to the pseudocode listings. Proofs are provided in the case, where the algorithm does not simply follow the definition of a rule.

We further introduce a $map$ data structure in our pseudocode with syntax $\texttt{map}[A]B$, which describes a mapping from $A$ into $B$. We use square brackets to indicate access and mutation of the mapping, e.g. $\gamma[0] \gets 1$. The map type also exposes a primitive function with the signature $delete(\gamma, x)$, which simply means that we want to delete the entry with key $x$ from our map. When iterating over a map in a \texttt{for}-loop we destructure the entries into a $(key, value)$-pair, e.g. $\texttt{for}\; (\_,v) \in map\; \texttt{do}$. Note that unused values are omitted with the underscore symbol.

### Tiny/Small Edge Rule

-   tiny edges: Delete all hyperedges of size one and place the corresponding vertices into the hitting set.
-   small edges: If $e$ is a hyperedge of size two, i.e., $e = \{ x,y \}$, then put both $x$
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

$O(|E|)$ **Algorithm**. We partition our set of edges into two disjoint sets $sub$ and $dom$. The set $dom$ will contain edges that could possibly be dominated. The set $sub$ will contain hashes of edges $e$ that could dominate another edge. We then iterate over the set $dom$ and compute every strict subset of the current edge $f$. For each of these subsets, we test if the hash of the subset is present in our set $sub$. If it is then $f$ is dominated by another edge.

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
Algorithm \ref{edom} finds all edges of G that are dominated, iff G has no size one edges.
\end{lemma}

_Proof._ Let $e$ be a dominated edge. Since there are no size one edges, edges with size two cannot be dominated. Thus $e$ has to be of size three. Then simply removing $e$ will not create or eliminate an edge domination situation. It is therefore sufficient to only check size three edges for the domination condition. $\square$

This also allows us to parallelize the main part of the algorithm, where we check each edge in our $dom$ set. We can achieve a speedup of approximatly $2$ on a six-core CPU an a pseudo-random graph with one million vertices and two million edges.

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

$O(|V|^2 \cdot |E|)$ **Algorithm**. The additional factor $|E|$ looks scary at first, but will only occur in the worst case, if there exists a vertex $v$, that is incident to all edges in $G$.

We start by iterating other the $AdjCount$ map of the graph, refering to the current value in the iteration as $AdjCount[v]$. We then use the $\textsc{Two-Sum}$-Algorithm to compute and return the first pair $(a,b)$ in $AdjCount[v]$, s.t. for $(a,b)$ holds,

$$
AdjCount[v][a] + AdjCount[v][b] = deg(v)+1
$$

If such a pair exists, then we conclude that for every edge $f$ such that $v \in f$, it holds that, either $a\in f$ or $b\in f$.

\begin{lemma} The outlined procedure above is correct, under the assumption that the underlying graph does not contain any duplicate edges.
\end{lemma}

_Proof._ Let $G$ be a hypergraph. We first remove all edges of size one with the _Tiny Edge Rule_. We then construct our map $AdjCount$. Let $v$ be an entry in $AdjCount$ and $sol=(a,b)$ the result of calling our $\textsc{Two-Sum}$ implementation on $AdjCount[v]$ with a target sum of $n=deg(v)+1$.

_Proposition._ If $sol$ is non-empty, then the edge $\{ v,a,b \}$ exists.

Let $sol=(a,b)$ be the solution obtained by calling our $\textsc{Two-Sum}$ algorithm on $AdjCount[v]$ with a target sum of $n=deg(v)+1$. For the sake of contradiction let us assume that the edge $\{ v,a,b \}$ does not exist. Since our graph does not contain duplicate edges and does not contain $\{ v,a,b \}$, there exist $deg(v)+1$ many edges that contain either $\{a,v\}$ or $\{b,v\}$. This however contradicts that there only exist $deg(v)$ many edges containing $v$. Therefore it must be, that the assumption that $\{ v,a,b \}$ does not exist, is false.

Since $\{ v,a,b \}$ exists, $a$ and $b$ can only occur $n-2 = deg(v) - 1$ times in other edges containing $v$. Since duplicate edges of $\{ v,a,b \}$ can not exist, we know that every other edge containing $v$ also contains $a$ or $b$, but not both simultaniously. $\square$

We then add the two vertices in the solution $sol$ to our partial solution $c$.

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

**Idea**: The initial idea for this algorithm involved the usage of a complete incidence matrix, where edges are identified by the rows and the vertices are identified by the columns. To check the _Domination Condition_ for a vertex $v$, the algortihm would select all edges/columns that contain $v$ and then add up the columns. Now let $n$ be the amount of edges containing $v$. If there exist two entries in the resulting column that have a combined value of $n+1$, then the rule applies for $v$ under the assumption that there are no duplicate edges. This would result in an algorithm with a worse time complexity of $|V|+|V|^2\cdot|E|$.

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

$O(|E|+|V|^2)$ **Algorithm** We start by constructing an adjacency list $\texttt{adjList}$ for all edges of size two. We then iterate over the entries of the list. For the current entry $adjList[v]$ we compute all subsets of size two of the entry. If there exists a subset $s$ such that $s \in E$, then we found a small triangle situation. If we find a triangle situation we put the corresponding vertices in our partial solution and alter the adjacency list to reflect these changes. We do this by iterating over all vertices that are adjacent to the triangle. For every vertex $w$ of these vertices we delete all vertices of the triangle from the entry $adjList[w]$.

This last step will introduce the quadratic complexity, since in the worst case, for a vertex $v$ in a triangle, there could exist $|V|$ many size two edges that contain $v$. This worst case occurs very rarely, which justifies using this quadratic algorithm. We could alternatively move the last step of the algorithm outside of the loop, and wrap both procedures with an outer loop which breaks if we dont find any more triangles. This simulates calling the rule exhaustively, while achieving a linear time complexity.

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
            }
        }
    }
}


\KwRet exec\;

\caption{Algorithm for exhaustive application of Extended Triangle Rule\label{etri}}
\end{algorithm}
$$

### Small Edge Degree 2 Rule

-   small edge degree 2: Let v be a vertex of degree 2, and let the two hyperedges
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

## Self-Monitoring

Each of the reduction rule functions returns a $\texttt{int32}$ value, which indicates the number of rule executions.

We store the ratios for each rule in a map of the form:

```go
var ratios = map[string]pkg.IntTuple{
	"kRulename": {A:1, B:1},
}
```

Where $\texttt{A}$ denotes the amount of vertices put into the partial solution by a single rule execution. And $\texttt{B}$ denotes the amount of vertices present in an optimal solution. We can then use these values to calculate the estimated approximation factor as follows:

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

![Estimated Approximation Factor for Graphs with 10, 100, 1K, 10K Vertices for $r\in[1,20]$](img/line_fac.png)

We also looked at the number of rule executions for $r=1$ and $r=10$.

![Number of Rule Executions, $r=1$](img/bar_ex_1.png)

![Number of Rule Executions, $r=10$](img/bar_ex_10.png)

Note that the reduction rules alone are sufficient to compute a Hitting-Set for our graphs. The est. ratio is quite low, since we do not need to put whole size 3 edges into our hitting set. This is due to the way our random-graph model works. This will not work for all graph classes, namely 3-uniform graphs.

We can observe a large performance hit, once our graph is near 3-uniform. We present some ideas that could potentially speedup the execution for these graphs.

\todo[backgroundcolor=red]{remove Factor-3 Rule function paragraph}
**Factor-3 Rule Targeting**. The Factor-3 rule will select the to be removed edge at random. We could alternatively choose an edge $e$, s.t. the removal of $e$ will allow the execution of the Vertex Domination Rule. This could potentially lower the number of Factor-3 Rule executions quite significantly.

# Algorithms

## Main Algorithm

It is crucial to apply the rules in a given order. Some rules expect that other rules can not be applied. As for the exact rules, one can save execution time, if a specific order is enforced when executing them. We use the order proposed in the original paper [@10.1007/978-3-642-29116-6_6], also refered to as precedence of rule executions. We will just use the shorter term _precedence_ going forward.

1. Exact Rules: vertex domination $\rightarrow$ tiny edge $\rightarrow$ edge domination
2. approximate (double) vertex domination rules
3. small edge degree 2 rule
4. small triangle rule
5. extended triangle rule
6. small edge rule

The main algorithm will apply all rules exhaustively according to the precedence, possibly mutating the input graph. If the graph has no more edges, then we are done. If not, put a size three edge into the partial hitting set and start over with the rule application. Note that the step of "applying all rules exhaustively" leaves room for interpretation and that the produced hitting sets are highly dependant on the actions taken during this step. We will refer to the procedure in this step as _rule strategy_ or just _strategy_ if the context allows it. We now present three strategies to apply the rules in the precedence.

The base strategy is quite simple. Execute a rule exhaustively and then do the same for the next rule in the precedence. Do this until no more rules can be applied. This is obviously not very optimal, since we do not check, wether a "better" rule can be applied again, before possibly executing a worse rule. We can however add some heuristics, that preserve the fast execution time and improve on effectiveness. Instead of executing the vertex domination $\rightarrow$ tiny edge $\rightarrow$ edge domination cascade only once, we execute it three times. By just doing this, we can recoup a lot of cascades we might have missed with the base strategy. Iteration values higher than three did not yield better results. We should also execute this cascade after rules, that are applied the most. Since these rules are the ones most likely to make the extact rules applicable. That would be the approximative (double) vertex domination rules.
The second strategy is also derived from the same base strategy. The only modification is that, we start over with the vertex domination rule, whenever one of the rules was applied at least once. This should in theory lead to more exact rules being executed, while sacrificing a bit of execution time.
The third strategy tries to maximize the exact rule executions. We start by applying the exact rules exhaustively. We then execute the other rules at most once, according to the precedence. If a rule has been applied, we execute the exact rules exhaustively and start over with the first non-exact rule in the precedence.

## Incremental Frontier Algorithm

As experienced with the DBLP coauthor graph, there are some graph instances that do not work well with the prior algorithm. These instances do not admit a lot of rule executions after applying the fallback rule. Running the algorithm on the whole graph again, knowing that only small parts of the graph have changed is not very time efficient. The algorithm should only look at the part of the graph where the fallback rule, or in fact any rule, was applied. We therefore propose a new approach, which involves the usage of a _vertex frontier_. Such frontiers are common in search algorithms such as BFS, where each vertex in the frontier has the same distance to the root vertex. But let us explain the general structure of the algorithm first.

We first apply all reduction rules exhaustively for the entire graph $G$. If the graph has no more edges we are done. Else we try to find a size three edge $e$, which will trigger a vertex domination situation if removed. We choose a random size three edge if there is none. We then build up our initial frontier. We put all vertices into the frontier that are adjacent to a vertex in $e$, excluding vertices in $e$ itself. Then remove all edges that are incident to vertices in $e$. Next the frontier will be expanded, by adding all edges incident to the frontier to a new graph $H$. All vertices in $H$ that were not part of the frontier will form the next frontier. The amount of expansion steps can be set by the user\todo{proof that expansion step of 1 is sufficient ?}. The fields \texttt{AdjCount} and \texttt{IncMap} of $G$ will be reused by $H$. Thus changes in $H$ will be reflected in $G$ and vice versa. The main loop will be explained next.

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

    $e\gets \texttt{F3TargetLowDegree}(G)$\;
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

Algorithm \ref{3hs_frontier} was designed around both rule strategy 1 and 2. It is not compatible with strategy 3, since the auxiliary graph could at some point be empty and thus $|exp| = 0$, with the original graph still admitting rule applications. In such a case, we have to apply the first applicable rule in the precedence and use the neighborhood of the removed vertices to rebuild the auxiliary graph.

# Applications and Results {#appl}

All tests were run on a machine with a AMD Ryzen 5 3600 3.6 GHz six-core, 12-thread CPU and 32 gigabytes of 3200 MT/s DDR4 RAM. The machine was running Ubuntu Server 22.04.04 LTS and the binaries were compiled with Go version 1.22. Measured execution times do not include the time it takes to load in a graph from disk or the time needed to transform the graph to a new problem instance.

## Triangle Vertex Deletion

[PROBLEM DEFINITION]

\textsc{Triangle Vertex Deletion} can be reduced to the \textsc{3-Hitting Set} problem. Triangles in the input graph will be represented as size 3 edges in a hypergraph, containing the vertices that make up the triangle.

### DBLP Coauthor Graph

We tested the algorithms on a DBLP coauthor graph and ER graphs of varying densities. It was considered to use the complete DBLP coauthor graph, which was quite large. We ultimately decided to use a dataset of the largest connected component. The network is made available by the [Stanford Network Analysis Project](https://snap.stanford.edu/data/com-DBLP.html) and part of a paper by J. Yang and J. Leskovec [@DBLP:journals/corr/abs-1205-6233] about community detection in real-world networks. The network contains 2224385 triangles, resulting in a \textsc{3-Hitting Set} instance of the same size. We now present the collected data during a run of the algorithm on this hypergraph.

![Triangle Vertex Deletion on a DBLP coauthor graph, amount of rule executions](img/dblp_tvd_rules.png)

Rule strategy 1 has the fastest execution time with a mean of 29 seconds, but also the worst ratio with a mean of $1.5708$. Rule strategy 3 has the best ratio with a mean of $1.4199$\todo{replace with real value}, but also the worst execution time of about 8 minutes. We used the frontier technique in all of the tests, since we could observe a speedup of about 90, just using the base rule strategy. This speedup can be attributed to the fact, that the original coauthor graph exerts a community structure. These communities are preserved when converting it to a $\textsc{Triangle Vertex Deletion}$ and thus $\textsc{3-Hitting-Set}$ instance. The expansion steps will only expand the auxiliary graph $H$ inside these communities. This drastically reduces the amount of edges the algorithm has to process per iteration. One can see the opposite happen in really dense hypergraphs. Due to the high density, even a small expansion two neighborhoods deep will add all edges of $G$ to $H$, eliminating possible gains. In fact, constant rebuilding of the auxiliary graph $H$ will even worsen the execution time.\todo{needs proof?}

$$
\begin{table}[t]
    \label{stats_dblp}
    \begin{subtable}[b]{0.45\textwidth}
        \centering
        \caption{base rule strategy\label{str_base}}
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
        \caption{rule strategy 1\label{str_1}}
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
        \caption{rule strategy 2\label{str_2}}
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
        \caption{rule strategy 3\label{str_3}}
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
\end{table}
$$

### Amazon Product Co-Purchasing Graph

\todo[backgroundcolor=green]{add amazon graph section}

To validate that our previous findings were not just due to the structure of the DBLP coauthor graph, we conduct the same tests with an unrelated network. We chose an Amazon product co-purchasing graph, which had a similar amount of vertices and edges compared to the DBLP coauthor graph. The resulting \textsc{Triangle-Vertex-Deletion} instance only contains 667129 edges and thus should be structurally distinct, compared to the instance obtained from the DBLP coauthor graph. As seen in table \ref{stats_amzn} \todo{interpretation of results}

\begin{table}[t]
\begin{subtable}[b]{0.3\textwidth}
\centering
\caption{unmodified}
\begin{tabular}{lrr}
\toprule
& Ratio & $|C|$ \\
\midrule
mean & 1.7469 & 95252.10 \\
std & 0.0003 & 74.38 \\
min & 1.7460 & 95081.00 \\
median & 1.7469 & 95249.50 \\
max & 1.7477 & 95435.00 \\
\bottomrule
\end{tabular}
\label{amzn_stock}
\end{subtable}
\hfill
\begin{subtable}[b]{0.3\textwidth}
\caption{no edge domination rule}
\begin{tabular}{lrr}
\toprule
& Ratio & $|C|$ \\
\midrule
mean & 1.7404 & 92614.51 \\
std & 0.0003 & 80.23 \\
min & 1.7397 & 92382.00 \\
median & 1.7404 & 92621.00 \\
max & 1.7413 & 92802.00 \\
\bottomrule
\end{tabular}
\label{amzn_noedom}
\end{subtable}
\hfill
\begin{subtable}[b]{0.3\textwidth}
\caption{no edge domination rule and small edge rule$(*)$}
\begin{tabular}{lrr}
\toprule
& Ratio & $|C|$ \\
\midrule
mean & 1.6894 & 89951.35 \\
std & 0.0005 & 83.60 \\
min & 1.6876 & 89742.00 \\
median & 1.6894 & 89953.00 \\
max & 1.6905 & 90153.00 \\
\bottomrule
\end{tabular}
\label{amzn_noedomsmall}
\end{subtable}
\newline
\vspace{4mm}
\newline
\begin{subtable}[b]{\linewidth}
\centering
\caption{Number of rule executions, percent change of mean after not using some rules}
\begin{tabular}{l|rrrrrrrrrr}
\toprule
rules disabled & Tiny & VD & ED & Small & Tri & ETri & AVD & ADVD & SED2 & F3 \\
\midrule
EDom & 0.99 & 0.42 & -100 & -52.48 & -99.67 & 19.02 & 5.49 & -1.84 & -75.23 & 2.68 \\
Edom, Small & 20.20 & 4.86 & -100 & -99.89 & -99.58 & 18.93 & 5.52 & -2.38 & -75.56 & 2.85 \\
\bottomrule
\end{tabular}
\label{rules_diff_amazon}
\end{subtable}
\caption{Results for \textsc{Triangle Vertex Deletion} on Amazon co-purchasing graph; $n=100$}
\label{stats_amzn}
\end{table}

Next we looked at random ER graphs. We conducted the tests on graphs with $1000$ vertices and an edge to vertex ratio of 10 to 25. Data was collected for $100$ graphs per EVR.

![Triangle Vertex Deletion on ER graphs, est. appr. ratios](img/er_tvd_ratio.png)

![Triangle Vertex Deletion on ER graphs, amount of rule executions for EVR=10](img/er_tvd_rules_1.png)

![Triangle Vertex Deletion on ER graphs, amount of rule executions for EVR=25](img/er_tvd_rules_2.png)

## Cluster Vertex Deletion

\begin{table}[h]
\centering
\begin{tabular}{llrrrrrrrrrrrr}
\toprule
& & Ratio & Tiny & VD & ED & Small & Tri & ETri & AVD & ADVD & SED2 & F3 & $|C|$ \\
\midrule
\multirow{5}{4em}{stock} & mean & 1.7529 & 3.07 & 11.69 & 35.54 & 3.10 & 0.03 & 0 & 7.77 & 0.03 & 0.05 & 0 & 25.13 \\
& std & 0.1703 & 2 & 6.11 & 29.80 & 2.35 & 0.18 & 0.07 & 4.59 & 0.18 & 0.23 & 0.01 & 13.24 \\
& min & 1 & 0 & 0 & 0 & 0 & 0 & 0 & 0 & 0 & 0 & 0 & 2 \\
& median & 1.7778 & 3 & 11 & 29 & 3 & 0 & 0 & 7 & 0 & 0 & 0 & 23 \\
& max & 2.6667 & 13 & 37 & 198 & 16 & 2 & 4 & 25 & 3 & 4 & 2 & 62 \\
\hline
\multirow{5}{4em}{no EDom and small} & mean & 1.7109 & 3.47 & 12.10 & 0 & 0 & 0 & 0 & 9.75 & 0 & 0 & 0 & 22.98 \\
& std & 0.1807 & 2.18 & 6.32 & 0 & 0 & 0.06 & 0 & 5.59 & 0.02 & 0 & 0.01 & 12.03 \\
& min & 1 & 0 & 0 & 0 & 0 & 0 & 0 & 0 & 0 & 0 & 0 & 2 \\
& median & 1.7391 & 3 & 11 & 0 & 0 & 0 & 0 & 9 & 0 & 0 & 0 & 21 \\
& max & 2.50 & 15 & 42 & 0 & 0 & 2 & 1 & 27 & 2 & 1 & 2 & 56 \\
\bottomrule
\end{tabular}
\caption{Results for $\textsc{Cluster Vertex Deletion}$ on the Rome-Graph dataset, $n=100$}
\label{cvd_rome_stats}
\end{table}

\todo{interpretation of results}

## Preferential Attachment Hypergraphs

![](img/pa_ratio.png)
![](img/pa_rules_02.png)
![](img/pa_rules_05.png)

## 3-Uniform ER Hypergraphs

![](img/eru3_ratio.png)
![](img/eru3_rules_1.png)
![](img/eru3_rules_2.png)

# Testing

Every reduction rule is tested for their correctness with unit tests. We create small graphs in these tests, that contain structures, which the rules are targeting. We then test for the elements in the partial solution, amount of edges/vertices left in the graph and degree of the vertices left.

# References
