---
    title: Implementation and evaluation of a self-monitoring approximation algorithm for 3-Hitting-Set
    link-citations: true
    author: Khoa Le
    abstract: |
        In this thesis we implemented an approximation algorithm for the \textsc{$3$-Hitting Set} problem. The algorithm is based on various reduction rules and the approximation ratio is estimated with the local ratio technique. We show test results on real-world networks that can be interpreted as derivative problems of \textsc{Hitting Set}. At last, we compare these results with two linear programming based algorithms.
...

\RestyleAlgo{ruled}
\DontPrintSemicolon
\SetAlgoVlined
\LinesNumbered
\SetKwFunction{FRecurs}{FnRecursive}
\SetKwFunction{ExpandFrontier}{ExpandFrontier}
\SetKwProg{Fn}{func}{\string:}{}

# Introduction

The \textsc{$d$-Hitting Set} problem is basically the \textsc{Vertex Cover} problem on hypergraphs, which are a generalization of standard graphs by allowing an edge to have an arbitrary amount of endpoints. The problem can be formulated as follows.

$$
\begin{center}
\fbox{
    \begin{minipage}[h]{0.75\textwidth}
        \textbf{Problem Name:} \textsc{$d$-Hitting Set}

        \textbf{Given:} A hypergraph $G=(V,E)$ with edges of size at most $d$ and a non-negative integer $k$.

        \textbf{Output:} Is there a \textit{hitting set} $C\subseteq V$ of size at most $k$, such that for all $e \in E$ it holds that $e\cap C \neq \emptyset$.
    \end{minipage}
}
\end{center}
$$

The problem is $\NP$-hard, meaning that there exist no reasonably fast algorithms to solve the problem, unless $\P = \NP$. One way to tackle these kinds of problems is the use of approximation algorithms. Such algorithms will compute a solution to a (minimization) problem, that is at most $r$ times the size of an optimal solution, for some $r \geq 1$. 

We are going to implement a 3-\textsc{Hitting Set} 3-approximation algorithm, which is based on the work of Brankovic and Fernau in [@BraFer]. Note that this is a preliminary version of the paper, the latest version contains full proofs and is still unpublished. The algorithm is based on exact and approximative reduction rules and the analysis is made via techniques reminiscent of the _local ratio theorem_ by Bar-Yehuda and Even. The idea behind this theorem is to view an approximate solution to a problem as a decomposition of its cost and lower bound. We direct to [@DBLP:journals/csur/Bar-YehudaBFR04] for a more in-depth presentation. An algorithm using local ratio techniques is practically building this composition step by step, which is essentially what the reduction rules are doing. Recording the number of executions for each rule during the run of the algorithm, allows us to make an estimation of the actual approximation ratio. This method of analyzing the approximation ratio has been described as _a posteriori analysis_ or _self-monitoring_ by Fernau in [@WS:fernau/fptapprox].

We are going to use hypergraphs obtained from real world networks and randomly generated hypergraphs for our testing. The primary metrics we are focusing on are execution time, hitting set size and estimated ratio. We will then compare our algorithm to two linear programming based algorithms.

## Preliminaries

A hypergraph $G=(V,E)$ consists of a set of vertices $V$ and a set of edges $E$. Each element of $E$ is a subset of $V$. The degree of a vertex $v$ is the number of edges that are incident to $v$. We also denote the degree of a vertex $v$ by $deg(v)$. The maximum degree of a hypergraph is denoted by $\Delta$. An edge $e$ has size $l$, if it contains $l$ elements. A hypergraph is said to be $d$-_uniform_ if all edges of the hypergraph are of size $d$.

We also need to define what a minimization problem and an $\alpha$-preserving reduction is. We are going to follow the definitions in [@WS:fernau/fptapprox]. 

__Definition.__ A minimization problem $\mathcal{P}$ can be specified by a triple $(I_\mathcal{P} , \textnormal{SOL}_\mathcal{P} , m_\mathcal{P} )$, where 

1. $I_\mathcal{P}$ is the set of input instances of $\mathcal{P}$;
2. $\textnormal{SOL}_\mathcal{P}(x)$ is a set of feasible solutions of $x$, where $x \in I_\mathcal{P}$;
3. $m_\mathcal{P}$ is a function that maps $(x, y)$ to a non-negative integer $m_\mathcal{P} (x, y)$, which denotes the size of the solution $y$ to instance $x$. 

An optimum solution $y^*$ to $x$ satisfies: (i) $y^* \in \text{SOL}_\mathcal{P} (x)$, and (ii) $m_\mathcal{P} (y^*) = \min\{m_\mathcal{P} (y) \mid y \in \text{SOL}_\mathcal{P} (x)\}$. The value $m_\mathcal{P} (x, y^*)$ is also referred to as $m_\mathcal{P}^*(x)$ for brevity.

Given a minimization problem $\mathcal{P}$, a factor-$\alpha$ approximation, $\alpha \geq 1$, associates to each $x \in I_\mathcal{P}$ some $y \in \textnormal{SOL}_\mathcal{P} (x)$ such that $m_\mathcal{P} (x, y) \leq \alpha \cdot m_\mathcal{P}^* (x)$.
A solution $y \in \textnormal{SOL}_\mathcal{P} (x)$ satisfying $m_\mathcal{P} (x, y) \leq \alpha \cdot m_\mathcal{P}^* (x)$ is also called an $\alpha$-approximate solution for $x$.

__Definition.__ Let $\mathcal{P} = (I_\mathcal{P} , \textnormal{SOL}_\mathcal{P} , m_\mathcal{P} )$ be a minimization problem. An $\alpha$-preserving reduction, with $\alpha \geq 1$, is a pair of mappings $\texttt{inst}_\mathcal{P} : I_\mathcal{P} \rightarrow I_\mathcal{P}$ and $\texttt{sol}_\mathcal{P} : \textnormal{SOL}_\mathcal{P} (\texttt{inst}_\mathcal{P} (x)) \rightarrow \textnormal{SOL}_\mathcal{P} (x)$ satisfying the following conditions. Function $\texttt{inst}_\mathcal{P}$ maps $x$ into $\texttt{inst}_\mathcal{P} (x)$, where $x, \texttt{inst}_\mathcal{P} (x) \in I_\mathcal{P}$. Function $\texttt{sol}_\mathcal{P}$ maps $y' \in \textnormal{SOL}_\mathcal{P} (\texttt{inst}_\mathcal{P} (x))$ into some $y \in \textnormal{SOL}_\mathcal{P} (x)$ such that there are constants $a, b \geq 0$ satisfying the following inequalities:

1. $b \leq \alpha \cdot a$,
2. $m_\mathcal{P}^* (\texttt{inst}_\mathcal{P} (x)) + a \leq m_\mathcal{P}^* (x)$, and
3. $\forall y' \in \textnormal{SOL}_\mathcal{P} (\texttt{inst}_\mathcal{P} (x))( m_\mathcal{P} (\texttt{inst}_\mathcal{P} (x), y') + b \geq m_\mathcal{P} (x, \texttt{sol}_\mathcal{P} (y')))$.

We will now illustrate how one can formulate $\alpha$-preserving reduction rules for the 3-Hitting Set problem, for $\alpha = 3$. Let $G=(V,E)$ be a hypergraph and $a,b \geq 1$ with $b \leq \alpha \cdot a$. Let $R$ be a reduction rule which will put $b$ many vertices into the hitting set, by removing these vertices from $V$ and removing all edges incident to them from $E$. We now argue that $R$ is $\alpha$-preserving, iff at least $a$ many of these vertices are always present in any minimum hitting set of $G$. 

_Proof._ Let $Y$ be the vertices that $R$ puts into the hitting set and let $G'$ be the hypergraph obtained by applying $R$. For every hitting set $C'$ of $G'$, by adding the vertices in $Y$ to $C'$, we obtain a hitting set $C$ for $G$ with $|C'| + |Y| \geq |C|$. 

Let $C^*$ be a minimum hitting set of $G$. Let $G''$ be the hypergraph obtained by removing the vertices in $Z = \{ v \mid v \in Y,\; v \in C^*\}$ and all edges incident to $Z$ from $G$. Since $C^*\setminus Z$ is a hitting set for $G''$, it holds that $|C^*| - |Z| \geq |C''^*|$, with $C''^*$ being a minimal hitting set of $G''$. Since $Z$ is a subset of $Y$, $G'$ has to be a subgraph of $G''$. The size of $C'^*$ is then at most the size of $C''^*$. It then holds that,

$$
\begin{align*}
\: & |C'^*| \leq |C''^*| \quad \wedge \quad  |C^*| - |Z| \geq |C''^*|\\ \Leftrightarrow \: & |C'^*| \leq |C''^*| \quad \wedge \quad |C''^*| + |Z| \leq |C^*|\\
\Leftrightarrow  \: & |C'^*| + |Z| \leq |C''^*| + |Z| \leq |C^*|\\
\end{align*}
\qed
$$

The composition of two $\alpha$-preserving reduction rules is also an $\alpha$-preserving reduction, which is necessary for the a posteriori analysis. For a more comprehensive overview, we direct to [@WS:fernau/fptapprox] section 2.2 and 2.3.

## Programming Language

We chose the [_Go_](https://go.dev/) programming language for our implementation. It is a statically typed, compiled language, with a C-like syntax and built-in concurrency primitives. It uses a garbage collector to handle memory management, which is off-putting for an application like this. Since Go's GC is very efficient, we are not worried about that fact. Go also provides great developer tooling without the need for additional third party tools. The built-in profiler _pprof_ can visualize captured performance profiles as flamegraphs or standard graphs, which makes it easy to spot performance bottlenecks. Another strength of Go is its simplicity, which makes writing idiomatic Go code very easy. 

# Data Structures

## Vertex

```go
type Vertex struct {
	Id int32
	Data int
}
```

The $\texttt{Vertex}$ datatype has two fields. The field $\texttt{Id}$ is an arbitrary identifier and $\texttt{Data}$ serves as a placeholder for actual data associated with the vertex.

## Edge

```go
type Edge struct {
	V map[int32]bool
}
```

The $\texttt{Edge}$ datatype has one field. The field $\texttt{V}$ is a map with keys of type $\texttt{int32}$ and values of type $\texttt{bool}$. This map will store the ids of vertices contained in the edge. Using \texttt{bool}s as values simulates a _Set_ datatype, allowing faster access times than simple arrays and slices.

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

The $\texttt{HyperGraph}$ datatype has five fields. Both fields $\texttt{Vertices}$ and $\texttt{Edges}$ are maps with keys of type $\texttt{int32}$ and values of type $\texttt{Vertex}$ and $\texttt{Edge}$ respectively. We chose this set-like data structure over lists again because of faster access times, but also operations that remove edges/vertices are built-in to the map type. The field $\texttt{edgeCounter}$ is an internal counter used to assign ids to edges. The field $\texttt{IncMap}$ is a map of maps, essentially storing the hypergraph as a sparse incidence matrix. We will also derive vertex degrees from this map. And at last the $\texttt{AdjCount}$ map, which will associate every vertex $v \in V$ with all vertices adjacent to $v$. Additionally, this map will also store the number of times such a vertex is adjacent to $v$. We also provide various struct methods to do basic operations on the hypergraph. All of these methods handle mutations to the \texttt{IncMap} and \texttt{AdjCount} fields of the receiving hypergraph struct. For example, calling \texttt{RemoveEdge(0)} on a hypergraph struct \texttt{g} will remove the edge with id $0$ from \texttt{g.Edges} and will remove \texttt{0} from all entries \texttt{g.IncMap[v]}, \texttt{v} being an endpoint of the removed edge. 

# Misc. Algorithms/Utilities

## Edge Hashing

Given a set $S$ of size $n$, compute a unique hash of $S$.

Time Complexity: $O(n + n\cdot \log(n))$, assuming an average of $O(n \cdot \log(n))$ comparisons for $\textsc{Quick-Sort}$

We start by sorting $S$ with a $\textsc{Quick-Sort}$-algorithm. We then join the elements of $S$ with the delimiter $"|"$, returning a string of the form $"|i_0|i_1|\dots|i_n|"$. Whenever we refer to \textit{the hash of an edge} we refer to the output of this function, using the endpoints of the edge as the input set.

## Compute Subsets of Size $k$

Given an array $A$ and an integer $k$, compute all subsets of $A$ of size $k$.

Time Complexity: $O(n^k)$, where $n$ denotes the size of $A$. 

$$
\begin{algorithm}[H]
\label{subset}

\KwIn{An array $A$, the subset size $k$ and a list $subsets$.}
\KwOut{None}
\BlankLine

$data\gets [\;]$\;
$n\gets |A|$\;
$last\gets k-1$\;
$\texttt{FnRecursive}(0,0)$\;

\Fn{\FRecurs{i, next}}{
    \For{$j\gets next$ \KwTo $n$}{
        $data[i] \gets A[j]$\;
        \eIf{$i=last$}{
            $subsets.push(data)$\;
        }{
            $\texttt{FnRecursive}(i+1, j+1)$\;
        }
    }
}
\caption{An algorithm to compute all subsets of size $k$}
\end{algorithm}
$$

Lists in Go are not very memory efficient, but since we exclusively call this function with $A$ representing the vertices in an edge, the value is usually fixed at 3. The raised memory problems occur at values of $n>10000$, justifying the continued usage of lists. For the case where we have to compute a lot of subsets, we provide a slightly different version of this function. Instead of passing in the $subsets$ list, we pass in a callback function that is called whenever we find a subset, using the found subset as an argument.

## Two-Sum

Given an array $A$ of integers and an integer target $t$, return indices of the two numbers such that they add up to $t$.

Time Complexity: $O(n)$, where $n$ denotes the size of $A$.

$$
\begin{algorithm}
\label{twosum}

\KwIn{An array of integers $A$, a target value $t$}
\KwOut{Two indices $a,b$, such that $A[a] + A[b] = t$, a boolean indicating if a solution was found}
\BlankLine

$lookup \gets \texttt{map}[\mathbb{N}]\mathbb{N}$\;

\For{$i\gets 0$ \KwTo $len(A)$}{
    \lIf{$lookup[t-A[i]]$ exists}{
        \KwRet $(i, lookup[t-A[i]]), true$\;
    }
    \Else{
        $lookup[A[i]] \gets i$\;
    }
}
\KwRet $nil, false$\;

\caption{An algorithm for the \textssc{Two-Sum} problem}
\end{algorithm}
$$

We start by creating a map called $lookup$. Iterating other the elements of $A$, we check if the entry $lookup[t-A[i]]$ exists. If the entry exists, we return a pair $(i, lookup[t-A[i]])$ and the boolean value $true$ since we found a solution. If the entry does not exist, we add a new entry to the $lookup$ map using $A[i]$ as key and $i$ as value. If no solution was found, we return $nil$ and the boolean value $false$.

This algorithm is an ingredient for the implementation of one of the reduction rules, specifically the approximative vertex domination rule. The actual implementation accepts a map instead of an array as its first parameter. We also implemented a version that finds all solutions, which accepts an additional callback function as a parameter. Instead of returning the solution, we call the callback function with the solution as the first argument. This version is an ingredient for the implementation of the approximative double vertex domination rule. We refer to this version as $\texttt{Two-SumAll}$ in the algorithm listings.

# Hypergraph Models

## First Testing Model

```go
func GenerateTestGraph(n int32, m int32, tinyEdges bool) *HyperGraph
```

-   $\texttt{n}$ are the number of vertices the hypergraph will have
-   $\texttt{m}$ is the number of edges the hypergraph will at most have
-   $\texttt{tinyEdges}$ when $\texttt{false}$ indicates that we do not want to generate edges of size 1.

We first create an empty Hypergraph struct and add $\texttt{n}$ many vertices to that hypergraph. We then compute a random $\texttt{float32}$ value $\texttt{r}$ in the half-open interval $[ 0.0, 1.0 )$. This value will be used to determine the size of an edge $e$. The edges are distributed based on their size as follows.

$$
size(\texttt{r})=
\begin{cases}
    1 & \quad \texttt{r} < 0.01\\
    2 & \quad 0.01 \leq \texttt{r} < 0.60\\
    3 & \quad \text{else}
\end{cases}
$$

The result of $size(\texttt{r})$ is stored in a variable $\texttt{d}$. We then randomly pick vertices in the half-open interval $[ 0, \texttt{n} )$, until we have picked $\texttt{d}$ many distinct vertices. If an edge with these endpoints does not exist, then we add it to our graph. The hypergraph has at most $\texttt{m}$ edges and not exactly $\texttt{m}$, since we remove duplicate edges. The advantage of this model over a model like the Erdős–Rényi model is, that we can compute hypergraphs with large vertex count in time $O(m)$. The model was primarily used during the implementation phase of the algorithm to gather profiling data.

## Preferential Attachment Hypergraph Model

In the Preferential Attachment Model, one will add edges to an existing hypergraph, with a probability proportional to the degree of the endpoints of that edge. This edge will either contain a newly added vertex or will be comprised of vertices already part of the hypergraph. We will use an implementation by Antelmi et al. as reference [@SimpleHypergraphs.js], which is part of their work on _SimpleHypergraphs.jl_ [@DBLP:journals/im/SpagnuoloCSPSKA20], a hypergraph software library written in the Julia language. The implementation is based on a preferential attachment model proposed by Avin et al. in [@DBLP:conf/asunam/AvinLNP19].

```go
func GeneratePrefAttachmentGraph(n int, p float64, maxEdgesize int32) *HyperGraph
```

-   $\texttt{n}$ is the number of vertices of the hypergraph
-   $\texttt{p}$ is the probability of adding a new vertex to the hypergraph
-   $\texttt{maxEdgesize}$ is the maximum size of a generated edge

We start by generating a random hypergraph with five vertices and five edges, using the first model for the generation. With probability $\texttt{p}$ add a new vertex $v$ to the hypergraph and select at most $\texttt{maxEdgesize}-1$ other vertices that will form a new edge together with $v$. Else, with probability $1-\texttt{p}$ select at most $\texttt{maxEdgesize}$ many vertices that will form a new edge. Whenever we select a vertex $u$ to be included in an edge, we do this with probability

$$
P[\text{u is chosen}] = \frac{deg(u)}{\sum_{w\in V}deg(w)}.
$$

## Erdős–Rényi Model

The Erdős–Rényi model was published by Paul Erdős and Alfréd Rényi in [@erdren/random]. The model requires two inputs, the amount of vertices $n$ and a probability $p$. Let $E'$ be the set of all possible edges on $n$ vertices. To construct a random graph, for all edges $e \in E'$ include $e$ in the graph with probability $p$, independently for all $e$. We can generalize this model for $d$-uniform hypergraphs quite easily by defining $E'$ as the set of all size $d$ subsets of $\{ 1,\dots,n \}$. This immediately results in a lower bound of $\binom{n}{d}$ for any generation algorithm. and the reason why it is not feasible to compute hypergraphs with large $n$ with this model. 

```go
func UniformERGraph(n int, p float64, evr float64, d int) *HyperGraph
```

- \texttt{n} is number of vertices of the hypergraph
- \texttt{p} is the probability of adding an edge
- \texttt{evr} is the desired edge to vertex ratio
- \texttt{d} is the size of an edge

Note that our implementation takes an additional parameter called $\texttt{evr}$. If $\texttt{p}=0$, then we use the value of $\texttt{evr}$ to compute $\texttt{p}$ as follows,

$$
\texttt{p} = \texttt{n} \cdot \frac{\texttt{evr}}{\binom{\texttt{n}}{\texttt{d}}}
$$

The generated hypergraph will have $\texttt{evr}$ times as many edges as vertices. The subset generation is done with our subset algorithm, using a callback instead of a list.  
 
# Reduction Rules

The usual signature of an implemented reduction rule looks as follows.

```go
func NameRule(g *HyperGraph, c map[int32]bool) int32
```

We take both a pointer to a \texttt{HyperGraph} struct \texttt{g} and a map \texttt{c} as arguments and mutate them. We then return the number of rule executions. We prioritize time complexity over memory complexity when implementing rules, which does not equate to ignoring memory complexity completely.

## Executions

The proposed rules are meant to be applied exhaustively. A one-to-one implementation of a rule will only find one of the structures the rule is targeting. Calling such a rule implementation exhaustively will take polynomial time but will be very inefficient regarding memory writes and execution time. It is therefore advantageous to design the algorithms, with the aspect of exhaustive application in mind.

The general outline of an algorithm will look as follows.

1. (If needed) Construct auxiliary data structures that are used to find parts of the hypergraph, for which the rule can be applied.
2. If we can apply rules do:

    Iterate over the auxiliary data structure, or the vertices/edges themselves.

    - Identify targets of the rule.
    - Mutate the hypergraph according to the rule.
    - Mutate the auxiliary data structure according to the rule.

This way we minimize memory writes by reusing existing data structures we built in step 1, but also save on execution time, since we mutate and iterate the auxiliary data structure at the same time. We can also alter these algorithms with minimal changes, to only execute a rule once. 

## Algorithms

The algorithm descriptions occasionally omit implementation details. We do this to try and keep these descriptions as concise as possible. The real implementations correspond to the pseudocode listings. We do not impose the use of our proposed hypergraph model in those listings. Instead we use general language, for example "$e$ incident to $v$" instead of $e \in G.IncMap[v]$. Proofs are provided in the case, where the algorithm does not simply follow the description of a rule.

We further introduce a $map$ data structure in our pseudocode with syntax $\texttt{map}[A]B$, which describes a mapping from $A$ into $B$. We use square brackets to indicate access and mutation of the mapping, e.g. $\gamma[0] \gets 1$ maps the value $1$ to the key $0$ in map $\gamma$. The map type also exposes a primitive function with the signature $delete(\gamma, x)$, which simply means that we want to delete the entry with key $x$ from our map. When iterating over a map in a \texttt{for}-loop we destructure the entries into a $(key, value)$-pair, e.g. $\texttt{for}\; (\_,v) \in map\; \texttt{do}$. Unused values are omitted with the underscore symbol.

### Tiny/Small Edge Rule

-   tiny edges: Delete all hyperedges of size one and place the corresponding vertices into the hitting set.
-   small edges: If $e$ is a hyperedge of size two, i.e. $e = \{ x,y \}$, then put both $x$
    and $y$ into the hitting set.

$O(|E|\cdot\Delta)$ **Algorithm**. We iterate over all edges of the hypergraph. If the current edge $e$ is of size $t$, put $e$ into the partial hitting set and remove $e$ and all edges incident to vertices in $e$ from the hypergraph.

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

This rule can be modified to work for hypergraphs with an edge size $d>3$. In that case, the time complexity of this algorithm could become exponential in the input size, if $d \geq \min \{|V|, |E|\}$.

\begin{lemma}
Algorithm \ref{edom} finds all edges of G that are dominated, iff G has no size one edges.
\end{lemma}

\begin{proof} 
Let $e$ be a dominated edge. Since there are no size one edges, edges with size two cannot be dominated. Thus $e$ has to be of size three. Then simply removing $e$ will not create or eliminate an edge domination situation. It is therefore sufficient to only check size three edges for the domination condition. 
\end{proof}

This also allows us to parallelize the main part of the algorithm, where we check each edge in our $dom$ set. We can achieve a speedup of approximately $2$ on a six-core CPU and a pseudo-random hypergraph with one million vertices and two million edges.

### Vertex Domination Rule

-   A vertex $x$ is dominated by a vertex $y$ if, whenever $x$
    belongs to some hyperedge $e$, then $y$ also belongs to $e$. Then, we can simply
    delete $x$ from the vertex set and from all edges it belongs to.

$O(|V| \cdot \Delta)$ **Algorithm**
Iterate over all vertices $v\in V$. If one of the entries in $AdjCount[v]$ is equal to $deg(v)$, then $v$ is dominated. In that case we remove $v$ from all edges and our vertex set.

\clearpage

$$
\begin{algorithm}[h]

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

$O(|V| \cdot \Delta)$ **Algorithm**. We start by iterating other the $AdjCount$ map of the hypergraph, referring to the current value in the iteration as $AdjCount[x]$. We then use the $\textsc{Two-Sum}$ algorithm to compute and return the first pair $(y,z)$ in $AdjCount[x]$, s.t. for $(y,z)$ holds,

$$
AdjCount[x][y] + AdjCount[x][z] = deg(x)+1
$$

If such a pair exists, then we conclude that for every edge $f$ incident to $x$, it holds that either $y\in f$ or $z\in f$.

\begin{lemma} Algorithm \ref{apvdom} is correct, under the assumption that the hypergraph does not contain any duplicate edges.
\end{lemma}

\begin{proof}
Let $G$ be a hypergraph without duplicate edges. Let $x$ be an entry in $AdjCount$ and $sol=(y,z)$ the result of calling our $\textsc{Two-Sum}$ algorithm on $AdjCount[x]$ with a target sum of $n=deg(x)+1$.

\textit{Proposition.} If $sol$ is non-empty, then the edge $\{ x,y,z \}$ exists.

Let $sol=(y,z)$ be the solution obtained by calling our $\textsc{Two-Sum}$ algorithm on $AdjCount[x]$ with a target sum of $n=deg(x)+1$. For the sake of contradiction let us assume that the edge $\{ x,y,z \}$ does not exist. Since our hypergraph does not contain duplicate edges and not the edge $\{ x,y,z \}$, there exist $deg(x)+1$ many edges that contain either $\{y,x\}$ or $\{z,x\}$. This however contradicts that there only exist $deg(x)$ many edges containing $x$. Therefore the assumption that $\{ x,y,z \}$ does not exist, must be false.

Since $\{ x,y,z \}$ exists, $y$ and $z$ can only occur $n-2 = deg(x) - 1$ times in other edges containing $x$. Since duplicate edges of $\{ x,y,z \}$ cannot exist, we know that every other edge containing $x$ also contains $y$ or $z$, but not both simultaneously.
\end{proof}

We then add the two vertices in $sol$ to our partial solution $C$.

The proof did not specifically mention why we have to assume that duplicate edges do not exist in the hypergraph. Let us construct a simple example to illustrate the importance of this assumption. Let $G$ be a hypergraph without duplicate edges. Let $e= \{ x,y,z \}$ be an edge $G$ for which the approximative vertex domination rule applies, i.e. all edges incident to $x$ contain either $y$ or $z$. Since there are no duplicate edges in $G$, it follows that $AdjCount[x][y] + AdjCount[x][z] = deg(x)+1$. We now add the edge $e_d = \{ x,y,z \}$ to our hypergraph. This will increase the degree of $x$ by 1, but will also increase both $AdjCount[x][y]$ and $AdjCount[x][z]$ by 1 each. Our algorithm would now fail to identify the approximative vertex domination situation at $e$.

$$
\begin{algorithm}[t]

\KwIn{A hypergraph $G=(V,E)$, a set $C$}
\KwOut{An integer denoting the number of rule applications.}
\BlankLine

$exec \gets 0$\;

$outer \gets true$\;
\While{$outer$}{
    $outer \gets false$\;
    \For{$(x, count) \in AdjCount$}{
        $sol, ex \gets \texttt{TwoSum}(count, deg(x)+1)$\;
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

**Idea**: The initial idea for this algorithm involved the usage of a complete incidence matrix, where edges are identified by the rows and the vertices are identified by the columns. To check the domination condition for a vertex $v$, the algorithm would select all edges/columns that contain $v$ and then add up the columns. Now let $n$ be the number of edges containing $v$. If there exist two entries in the resulting column that have a combined value of $n+1$, then the rule applies for $v$ under the assumption that there are no duplicate edges.

### Approximative Double Vertex Domination Rule

-   approximative double vertex domination: Assume there is a hyperedge $e =
\{x, y, a\}$ and another vertex $b$ such that, whenever $x$ or $y$ belong to some
    hyperedge $h$, then $a$ or $b$ also belong to $h$. Then, we put $a$ and $b$ together
    into the hitting set that we produce.

$O(|V|^2 \cdot \Delta^2)$ **Algorithm**. We start by creating a map called $tsHashes$, which maps subsets of $E$ to subsets of vertices. We then iterate over all vertices in $V$. For the current vertex $x$, compute all \textsc{Two-Sum} solutions with input array/map $AdjCount[x]$ and target $deg(x)$. If there exists a solution $sol = \{ z_0, z_1 \}$, such that $tsHashes[sol]\neq \emptyset$, then for all $y \in tsHashes[sol], y\neq x$ construct two edges $\{ x,y,z_0 \}$ and $\{ x,y,z_1 \}$. If one of these two edges exists in $E$, then we found a approximative double vertex domination situation. If $tsHashes[sol] = \emptyset$, map $tsHashes[sol]$ to $tsHashes[sol]\cup \{ x \}$.

\begin{lemma}
Algorithm \ref{apdvdom} is correct, under the assumption that the hypergraph does not contain any size one edges.
\end{lemma}

\begin{proof}
Let $G=(V,E)$ be a hypergraph. Let $x$ be the current vertex in the algorithm's iteration over $V$. If $T = \textsc{Two-SumAll}(AdjCount[x], deg(x))$ is empty, then $x$ will not be able to trigger an approximative double vertex domination situation. If $T$ is not empty, then we know that there exist sets $\{ a,b \}$, such that all edges incident to $x$ either contain $a$ or $b$. If $tsHashes[\{ a,b \}]$ is empty, then there are currently no other vertices for which $\{ a,b \}$ is a \textsc{Two-Sum} solution. In that case we add $x$ to $tsHashes[\{ a,b \}]$. Otherwise there exist vertices $y$, such that $\{ a,b \}$ is a \textsc{Two-Sum} solution for $y$. All edges incident to $x$ and $y$ either contain $a$ or $b$. If for one of these vertices $y$ there exists a size three edge $\{ x,y,a \}$, or without loss of generality $\{ x,y,b \}$, then we found a approximative double vertex domination situation at $\{ x,y,a \}$ or $\{ x,y,b \}$ respectively. If none of these edges exist, then $x$ could still trigger another approximative double vertex situation with another vertex. Thus we need to add $x$ to $tsHashes[\{ a,b \}]$.
\end{proof}

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

$O(|E| + |V| \cdot \Delta^4)$ **Algorithm** We start by constructing an adjacency list $adjList$ for all edges of size two. We then iterate over the entries of the list. For the current entry $adjList[v]$ we compute all subsets of size two of the entry. If there exists a subset $s$ such that $s \in E$, then we found a small triangle situation. If we find a triangle situation, we put the corresponding vertices in our partial solution and alter the adjacency list to reflect these changes. We do this by iterating over all vertices that are adjacent to the triangle. For every vertex $w$ of these vertices we delete all vertices of the triangle from the entry $adjList[w]$.

$$
\begin{algorithm}[h]
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

\clearpage

### Extended Triangle Rule

-   Assume that the hypergraph contains a small edge $e = \{y, z\}$. Moreover,
    there are hyperedges $f, g$ such that $e \cap f = \{ y \}$, $e \cap g = \{z\}$, $f \cup g=\{v, x, y, z\}$ and $|f| = 3$. Then, put all of $f \cup g$ into the hitting set.

$O(|E| \cdot (\Delta^2 + \Delta) )$ **Algorithm** We start by iterating over $E$ until we find a size two edge $e$. We then iterate over the vertices of $e$. Let $a$ be the vertex in the current iteration. Assign $a$ to the variable $y$ and $e\setminus \{ y \}$ to $\{ z \}$. Then we iterate over the edges $f$ that are incident to $y$. If $f$ is not of size three or $z\in f$, continue with the iteration. Else, iterate over the edges $g$ that are incident to $z$. If $g\setminus \{ z \} \subset f$, then we found an extended triangle situation. In that case, save $f$ in a variable $f_0$ and break out of the loop that iterates over the edges incident to $y$. If $f_0 \neq nil$, put $f_0 \cup \{ z \}$ into the partial hitting set $C$ and break out of the loop iterating over the endpoints of $e$.

$$
\begin{algorithm}

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
                put  $f_0 \cup \{ z \}$ into the partial hitting set $C$\;
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

$O(|V|\cdot \Delta)$ **Algorithm** We start by iterating other all vertices $v\in V$. If $v$ is of degree two, check if there is a size two edge $e$ and a size three edge $f$ incident to $v$. If these two edges exist, save $e\setminus \{ v \}$ in a variable $x$. Then iterate over the vertices $w$ in $f\setminus{v}$ and check if there exists an edge $h$ incident to $w$, such that $h\neq f$ and $x\notin h$. If such an edge exists, then we found a small edge degree 2 situation. We then put both $x$ and all vertices of $h$ into the partial hitting set.

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
        $h \gets nil$\;

        \For{$w\in s3\setminus \{ v \} $}{
            \For{$f$ \textnormal{incident to} $w$}{
                \eIf{$x \in f \textnormal{\textbf{ or }} s3 = f$}{
                    \textbf{continue}
                }{
                    $h\gets f$\;
                    \textbf{break}\;
                }
            }
            \If{$h \neq nil$}{
                \textbf{break}\;
            }
        }
        \If{$h \neq nil$}{
            $outer \gets true$\;
            $exec \gets exec +1$\;
            \For{$a \in \{ x \}\cup h$}{
                $C\gets C\cup \{ a \}$\;
                \For{$f$ \textnormal{incident to} $a$}{
                    $E\gets E \setminus \{ f \}$\;
                }
            }
        }
    }
}

\KwRet{$exec$}\;

\caption{Algorithm for exhaustive application of Small Edge Degree 2 Rule\label{sed2}}
\end{algorithm}
$$

### F3 Low Degree Rule

-   F3 low degree: Let $v$ be a vertex with degree 2 and let $e$ be an edge that is incident to $v$. If there exists an edge $f$ that is incident to a vertex in $e$ and does not contain $v$, put $f$ into the hitting set.

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

\clearpage

## Self-Monitoring

Each of the reduction rule functions returns a $\texttt{int32}$ value, which indicates the number of rule executions. We store the $a$ and $b$ values for each rule in a map called $\texttt{Ratios}$. The vertex domination and edge domination rule do not have entries in this map, since they do not put any vertex into the partial hitting set.

```go
var Ratios = map[string]pkg.IntTuple{
	"kTiny":             {A: 1, B: 1},
	"kSmall":            {A: 1, B: 2},
	"kTri":              {A: 2, B: 3},
	"kExtTri":           {A: 2, B: 4},
	"kApVertDom":        {A: 1, B: 2},
	"kApDoubleVertDom":  {A: 1, B: 2},
	"kSmallEdgeDegTwo":  {A: 2, B: 4},
	"kSmallEdgeDegTwo2": {A: 2, B: 3},
	"kFallback":         {A: 1, B: 3},
}
```

The number of vertices put into the partial hitting set by a single rule execution is denoted by $\texttt{B}$. And $\texttt{A}$ denotes the number of these vertices that are present in any optimal solution. We can then use these values to calculate the estimated approximation factor as follows.

```go
g := NewHypergraph()
c := make(map[int32]bool)

execs := ThreeHS_F3ApprPoly(g, c)

var num float64 = 0
var denom float64 = 0

for key, val := range execs {
	num += float64(Ratios[key].B * val)
	denom += float64(Ratios[key].A * val)
}

ratio := num / denom
```

# Algorithms

## Main Algorithm

It is crucial to apply the rules in a given order. Some rules expect that other rules cannot be applied. As for the exact rules, one can save execution time, if a specific order is enforced when executing them. We use the order proposed in [@BraFer], also referred to as _precedence of rule executions_. We will just use the shorter term _precedence_ going forward.

1. Exact Rules: vertex domination $\rightarrow$ tiny edge $\rightarrow$ edge domination
2. approximate (double) vertex domination rules
3. small edge degree 2 rule
4. small triangle rule
5. extended triangle rule
6. small edge rule

The main algorithm will apply all rules exhaustively according to the precedence, possibly mutating the input hypergraph. If the hypergraph has no more edges, then we are done. If not, apply the fallback rule by putting a size three edge into the partial hitting set and start over with the rule application. Note that the step of "applying all rules exhaustively" leaves room for interpretation and that the produced hitting sets are highly dependent on the actions taken during this step. We will refer to the procedure in this step as _rule strategy_ or just _strategy_ if the context allows it. We now present three strategies to apply the rules in the precedence.

The base strategy is quite simple. Execute a rule exhaustively and then do the same for the next rule in the precedence. Do this until no more rules can be applied. This is obviously not very optimal, since we do not check whether a "better" rule can be applied again, before possibly executing a worse rule. We can however add some heuristics that preserve the fast execution time and improve on effectiveness. Instead of executing the vertex domination $\rightarrow$ tiny edge $\rightarrow$ edge domination cascade only once, we execute it three times consecutively. By just doing this, we can recoup a lot of cascades we might have missed with the base strategy. Iteration values higher than three did not yield better results. We should also execute this cascade after rules that are applied the most. Since these rules are the ones most likely to make the exact rules applicable. That would be the approximative (double) vertex domination rules.
The second strategy is also derived from the same base strategy. The only modification is that we start over with the vertex domination rule, whenever one of the rules was applied at least once. This should in theory lead to more exact rules being executed, while sacrificing a bit of execution time.
The third strategy tries to maximize the exact rule executions. We start by applying the exact rules exhaustively. We then execute the other rules at most once, according to the precedence. If a rule has been applied, we execute the exact rules exhaustively and start over with the first non-exact rule in the precedence.

## Incremental Frontier Algorithm

There are some hypergraphs that do not work well with the prior algorithm. These instances do not admit any rule executions after applying the fallback rule. Running the algorithm on the whole hypergraph again, knowing that only small parts of the hypergraph have changed is not very time efficient. The algorithm should only look at the part of the hypergraph where the fallback rule, or in fact any rule, was applied. We therefore propose a new approach, which involves the usage of a _vertex frontier_. Such frontiers are common in search algorithms such as BFS, where each vertex in the frontier has the same distance to the root vertex. But let us explain the general structure of the algorithm first.

We first apply all reduction rules exhaustively for the entire hypergraph $G=(V,E)$. If $G$ has no more edges, we are done. Else we apply the fallback rule, but only select a size three edge $e$ that will be put into the partial hitting set. We then build up our initial frontier. We put all vertices in $\{ f \mid f\cap e \neq \emptyset,\; f\in E\}$ into the frontier, excluding vertices in $e$ itself. We then put $e$ into the partial hitting set. Next the frontier will be expanded by adding all edges incident to the frontier to a new hypergraph $H$. All vertices in $H$ that were not part of the frontier will form the next frontier. The amount of expansion steps can be set by the user. The fields \texttt{AdjCount} and \texttt{IncMap} of $G$ will be reused by $H$. Thus, changes in $H$ will be reflected in $G$ and vice versa. The main loop will be explained next.

We apply the rules as usual, but only on $H$. The rules are modified, such that they also return the vertices of edges that were edge-adjacent to a modified/removed edge. If there are any, then $H$ will be expanded at these vertices and the loop will be continued. If not, apply the fallback rule, considering the edges in the original hypergraph $G$ and expand accordingly. This will be repeated until $G$ has no more edges.

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
    $H.IncMap \gets G.IncMap$\;
    $H.AdjCount \gets G.AdjCount$\;
    \KwRet $H$\;
}

\caption{Incremental 3-approximation algorithm for 3-HS\label{3hs_frontier}}
\end{algorithm}
$$

Algorithm \ref{3hs_frontier} was designed around both rule strategy 1 and 2. It is not compatible with strategy 3, since the auxiliary hypergraph could at some point be empty and thus $|exp| = 0$, with the original hypergraph still admitting rule applications. In such a case, we apply the first applicable rule in the precedence and use the neighborhood of the removed vertices as an "entry point" to rebuild the auxiliary hypergraph.

\clearpage

# Applications and Results {#appl}

All tests were run on a machine with an AMD Ryzen 5 3600 3.6 GHz six-core, 12-thread CPU and 32 gigabytes of 3200 MT/s DDR4 RAM. The machine was running Ubuntu Server 22.04.04 LTS and the binaries were compiled with Go version 1.22.1. Measured execution times do not include the time it takes to load in a graph from disk or the time needed to transform the graph to a new problem instance. We provide results for $\textsc{Triangle Vertex Deletion}$ and $\textsc{Cluster Vertex Deletion}$ problems using real world networks, as well as randomly generated hypergraphs. The code is available on a public GitHub repository[^1].

[^1]: [https://github.com/KhoalaS/BachelorThesis](https://github.com/KhoalaS/BachelorThesis)

## F3 Low Degree Rule/Fallback

The goal of the F3 low degree rule is to reduce the degree of a vertex, which preferably has a degree of two. Applying the rule would lower the degree of the vertex to one, making the vertex domination rule applicable. We tested the effectiveness of the F3 low degree rule on random ER hypergraphs. First, we generated 100 random ER hypergraphs with 1000 vertices and approximately 3000 edges. We ran our algorithm on every hypergraph 10 times each to account for run-to-run variance. We collected data for two datasets. For the first dataset we selected a random size three edge when applying the fallback rule. For the second dataset we used the F3 low degree rule instead. The results can be seen in table \ref{er3_evr3_stats}.

We do in fact see a slight increase of vertex domination rule executions, with the mean number of executions increasing from about 345 to 349. More interesting observations can be made about the fallback rule and extended triangle rule. The mean number of fallback rule executions decreased from about 44 to 25 when not selecting random edges. The number of extended triangle rule executions increased from about 47 to 63. The mean values for the other rules stay mostly the same. 

To investigate these changing values, we reran the algorithm with both fallback variants on some of the hypergraphs. This time keeping a log of all the rules that are applied, see \ref{A:f3_logs} for the first 30 lines of such a log. Every time the F3 low degree rule is executed, the $F3\rightarrow VDom\rightarrow ETri$ cascade is most likely triggered. For the random variant, one can see batches of fallback rule executions and the occasional $F3\rightarrow VDom\rightarrow ETri$ cascade in the logs. 

So what happens is, that by using the F3 low degree rule, we substitute multiple fallback rule executions with a single $F3\rightarrow VDom\rightarrow ETri$ cascade. The desired effects on the hitting set size are underwhelming, with the mean hitting set size decreasing from about 592 to 590. The ratio did improve from 1.92 to 1.87, caused by a higher estimated optimum. Since we did not regress and even improved the ratio, we chose to keep using this variant of the rule instead of the random variant.

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

\textsc{Triangle Vertex Deletion} is a problem that can be reduced to the \textsc{3-Hitting Set}. Triangles in the input graph will be represented as size 3 edges in a hypergraph, containing the vertices that make up the triangle. A hitting set of size at most $k$ on this hypergraph can be used to transform the original standard graph into a triangle-free graph, removing at most $k$ vertices.

### DBLP Coauthor Graph

$$
\begin{figure}[t]
    \centering
    \includegraphics[width=\textwidth]{./img/dblp_tvd_rules_strat.png}
    \caption{Mean number of rule executions for \textsc{Triangle Vertex Deletion} on DBLP coauthor graph per rule strategy, $n=100$}
    \label{dblp_rules_strat}
\end{figure}
$$

It was considered to use the complete DBLP coauthor graph, which was quite large. We ultimately decided to use a dataset of the largest connected component. The network is made available by the [Stanford Network Analysis Project](https://snap.stanford.edu/data/com-DBLP.html) and part of a paper by J. Yang and J. Leskovec [@DBLP:journals/kais/YangL15] about community detection in real-world networks. The network contains 2224385 triangles, resulting in a \textsc{3-Hitting Set} instance of the same size. We now present the collected data during a run of the algorithm on this hypergraph, which can be found in table \ref{stats_dblp}.

Rule strategy 1 has the fastest execution time with a mean of 29 seconds, but also the worst ratio with a mean of $1.5708$. Rule strategy 3 has the best ratio with a mean of $1.4195$, but also the worst execution time of about 8 minutes. We used the frontier technique in all the tests, since we could observe a speedup of about 90, just using the base rule strategy. This speedup can be attributed to the fact that the original coauthor graph exerts a community structure. These communities are preserved when converting it to a $\textsc{Triangle Vertex Deletion}$ and thus $\textsc{3-Hitting Set}$ instance. The expansion steps will only expand the auxiliary hypergraph $H$ inside these communities. This drastically reduces the number of edges the algorithm has to process per iteration. One can see the opposite happen in dense hypergraphs. Due to the high density, even a small expansion two neighborhoods deep will add all edges of $G$ to $H$, eliminating possible gains. In fact, constant rebuilding of the auxiliary hypergraph $H$ will even worsen the execution time. This can be observed when using any rule strategy besides strategy 3 on a dense 3-uniform ER hypergraph. On inspection of the generated CPU profile, around 30% of the time is spent computing the auxiliary hypergraph $H$. See figure \ref{A:flame_base} for a flamegraph obtained from one of those profiles.

Figure \ref{dblp_rules_strat} displays the number of rule executions per strategy. Note that the there is a rule called SED2* in this chart. Recall that the small edge degree 2 rule can either put four or three vertices into the hitting set. Both cases were counted separately, where SED2* refers to the version that puts three vertices into the hitting set. There are several interesting observations to be made, for example the decrease of edge domination rule executions for strategy 2 compared to strategy 1. Quite frankly we do not know why this happens. We can also see that strategy 3 does indeed increase the amount of exact rule executions. Another observation can be made about the number of small edge rule executions. Strategy 2 does not apply the small edge rule at all, which is surprising since it is conceptually similar to strategy 3. This is probably a result of only applying rules one at a time using strategy 3. The exact rules that are executed after each rule application will possibly destroy situations that are favorable for other rules. The approximative vertex domination rule, which is executed a lot more with strategy 2 compared to 3, might have removed all the size two edges as a byproduct. The number of fallback rule executions is also interesting. All strategies are executing this rule around 3000 times. This leads us to the hypothesis that as long as the order in the precedence is obeyed, a "worse" strategy will not result in significantly more fallback rule executions.

$$
\begin{table}[t]
    \begin{subtable}[b]{0.45\textwidth}
        \centering
        \begin{tabular}{lrrrr}
            \toprule
            & est. ratio & $|C|$ & est. opt & time \\
            \midrule
            mean & 1.7594 & 115672 & 65745 & 17 sec \\
            std & 0.0005 & 79 & 41 & 1 sec\\
            min & 1.7577 & 115469 & 65635 & 16 sec\\
            median & 1.7594 & 115677 & 65741 & 17 sec\\
            max & 1.7610 & 115998 & 65904 & 19 sec\\
            \bottomrule
        \end{tabular}
        \caption{base rule strategy\label{dblp_str_base}}
    \end{subtable}
    \hfill
    \begin{subtable}[b]{0.45\textwidth}
        \centering
        \begin{tabular}{lrrrr}
            \toprule
            & est. ratio & $|C|$ & est. opt & time\\
            \midrule
            mean & 1.5708 & 106339 & 67697 & 29 sec\\
            std & 0.0006 & 65 & 40 & 1 sec\\
            min & 1.5691 & 106183 & 67597 & 26 sec\\
            median & 1.5708 & 106342 & 67699 & 29 sec\\
            max & 1.5727 & 106500 & 67797 & 31 sec\\
            \bottomrule
        \end{tabular}
        \caption{rule strategy 1\label{dblp_str_1}}
    \end{subtable}
    \newline
    \vspace{4mm}
    \newline
     \begin{subtable}[b]{0.45\textwidth}
        \centering
        \begin{tabular}{lrrrr}
            \toprule
            & est. ratio & $|C|$ & est. opt & time \\
            \midrule
            mean & 1.5437 & 105084 & 68071 & 104 sec \\
            std & 0.0006 & 63 & 35 & 2 sec\\
            min & 1.5422 & 104922 & 67956 & 99 sec\\
            median & 1.5438 & 105084 & 68072 & 104 sec\\
            max & 1.5453 & 105261 & 68146 & 110 sec\\
            \bottomrule
        \end{tabular}
        \caption{rule strategy 2\label{dblp_str_2}}
    \end{subtable}
    \hfill
    \begin{subtable}[b]{0.45\textwidth}
        \centering
        \begin{tabular}{lrrrr}
            \toprule
            & est. ratio & $|C|$ & est. opt & time \\
            \midrule
            mean & 1.4195 & 99456 & 70066 & 474 sec\\
            std & 0.0009 & 52 & 31 & 7 sec\\
            min & 1.4171 & 99323 & 69987 & 457 sec\\
            median & 1.4196 & 99450 & 70069 & 474 sec\\
            max & 1.4215 & 99575 & 70133 & 493 sec\\
            \bottomrule
        \end{tabular}
        \caption{rule strategy 3\label{dblp_str_3}}
    \end{subtable}
    \caption{Results for \textsc{Triangle Vertex Deletion} on DBLP coauthor graph; $n=100$}
    \label{stats_dblp}
\end{table}
$$

### Amazon Product Co-Purchasing Graph

To validate that our previous findings were not just due to the structure of the DBLP coauthor graph, we conduct the same tests with an unrelated network. We chose an Amazon product co-purchasing graph, which had a similar number of vertices and edges compared to the DBLP coauthor graph. The resulting \textsc{Triangle-Vertex-Deletion} instance only contains 667129 edges and thus should be structurally distinct, compared to the instance obtained from the DBLP coauthor graph. As seen in table \ref{stats_amzn} and figure \ref{amzn_rules_strat}, the algorithm achieves a mean estimated ratio of $1.3136$ and behaves the same as for the DBLP coauthor graph.

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

We tested our algorithm on the _Rome Graphs_, a collection of 11534 undirected graphs used and made available by Di Battista et al. in [@DiBattista1997b]. We ran the algorithm ten times for each derived problem instance and picked the run with the lowest ratio to include during calculation of the statistics. The results can be seen in table \ref{cvd_rome_stats}. The algorithm achieved a mean ratio of $1.1082$ with a standard deviation of $0.0851$. We also recorded a ratio of at most $1.75$ for the selected problem instances. It can be observed and inferred from the low ratio, that most of the applied rules were exact rules.

$$
\begin{table}[h]
\makebox[\textwidth][c]{
    \input{"out/rome_cvd_minratio"}
}
    \caption{Results for \textsc{Cluster Vertex Deletion} on graphs from the Rome Graphs collection. Hitting sets were computed ten times for each graph and the run with the lowest ratio was chosen for final calculation.}
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

$$
\begin{table}[t]
    \input{out/amazon_cvd_lp}
\end{table}
$$

We used PuLP [@GH:pulp] to model the linear programs in Python and implemented the algorithms in Python as well. Both algorithms were tested with three open source solvers GLPK [@GH:glpk], CLP [@GH:clp] and HiGHS [@GH:highs] a comparatively newer solver that uses a parallel dual simplex method developed by Q. Huangfu and J. A. J. Hall [@DBLP:journals/mpc/HuangfuH18]. Table \ref{lp_stats} shows the results for both algorithms on the \textsc{Triangle Vertex Deletion} instance of the Amazon product co-purchasing graph, which is a 3-uniform hypergraph with maximum vertex degree of 551. Algorithm (1) promises a hitting set size of $|C| \leq d(1-\frac{d-1}{8\Delta})\cdot \textnormal{Opt}^*$ with probability $\frac{3}{4}$, under the assumption that the hypergraph is $d$-uniform and $3 \leq d \leq \frac{16}{3}\Delta$. Our hypergraph fulfills both of these requirements, which should result in a ratio of at most $3(1-\frac{3-1}{8 \cdot 551})\cdot \textnormal{Opt}^* = 2.9986 \cdot \textnormal{Opt}^*$. This upper bound is only valid for $\lambda = l(1-\epsilon)$ where $\epsilon = \frac{l\textnormal{Opt}^*-|S_1|}{2m}$. Algorithm (2) promises a hitting set size of at most $|C| \leq ((\Delta -1)(1-e^{-\frac{ln d}{\Delta -1}})+1) \cdot \textnormal{Opt}^*$. Applied to our set cover instance, this yields an upper bound for the ratio of $((3-1)(1-e^{-\frac{ln 551}{3-1}})+1)\cdot \textnormal{Opt}^* = 2.9148\cdot \textnormal{Opt}^*$.

Algorithm (1) slightly outperforms algorithm (2) in terms of hitting set size. The LP solution obtained with the HiGHS solver using the interior point method yields the smallest hitting set for both algorithms, with an execution time of about 1 minute. The hitting set obtained with the GLPK solver comes close but at a way higher execution time of about 1 hour. On inspection of the logs for each solver, both the GLPK and HiGHS interior point method solutions contain significantly more variables with a value of 1 or 0 than the over solvers. The CLP solver is not very well suited for the algorithms on this specific problem instance. Execution time as well as the resulting hitting set size are worse compared to HiGHS and GLPK. Compared to our algorithm, only algorithm (1) using the HiGHS solver with the interior point method yielded a hitting set with comparable size.

We already mentioned how different solutions of the LPs can influence the output of the algorithms. We have seen this for solutions obtained with GLPK and HiGHS (ipm), which assigned exact integer values to more variables compared to CLP and HiGHS (simplex). Such solutions with more variables close to or equal 1 seem to be "more optimal" than a solution with less of such variables. This intuition is reasonable but not true in general. For our application, these solutions do lead to better results.

Also note that our implementations of the algorithms are not strict one-to-one implementations. Both algorithms do not check if a vertex $v$ is isolated, before adding it to the hitting set. It would not make sense to add $v$ to the hitting set, since $v$ does not cover any edges. We assume that this check is not mentioned in the algorithm descriptions to make the analysis easier, since omitting this check does not influence the upper bound of the approximation ratio. We did run the benchmark on naive implementations of the algorithms, the results can be seen in table \ref{amzn_lp_stats_naive}. Not adding this simple check resulted in hitting sets with about 20 thousand more vertices.

$$
\begin{table}[t]
    \centering
    \input{out/rome_cvd_opt_multi}
    \caption{Results for \textsc{Cluster Vertex Deletion} with LP rounding based algorithms on graphs from the Rome Graphs collection. Same method as previous CVD benchmark.}
    \label{cvd_rome_multi}
\end{table}
$$

We also conducted the same benchmark from the Cluster Vertex Deletion section with both algorithms using the same solvers and methods. Both algorithms perform very similarly, with algorithm (1) slightly outperforming algorithm (2) regarding mean hitting set size. This time the GLPK solver was best suited for the problem, the results for algorithm (1) can be seen in table \ref{cvd_rome_multi}. For results from the other solvers see table \ref{A:cvd_rome_all} in the appendix. We also added two additional columns "opt" and "actual ratio" to table \ref{cvd_rome_multi}, where "opt" denotes the size of an optimal solution and "actual ratio" denotes $\frac{|C|}{\text{opt}}$. The values for "opt" were obtained by solving the \textsc{Hitting Set} ILP for each \textsc{Cluster Vertex Deletion} instance in the dataset. We can see that both the Hitting Set LP algorithm and our algorithm overestimate the ratio. Our algorithm achieves a mean actual ratio of 1.04 while the mean estimated ratio is about 1.10. We can also see that our algorithm never exceeded an actual ratio of 1.5.

### Random Preferential Attachment Hypergraphs

The last benchmark dataset consists of 100 random preferential attachment hypergraphs. We used a value of $p=0.2$ for the probability of adding a new vertex to the hypergraph. Data was collected ten times each for every hypergraph in the dataset. The solver with the best ratio was HiGHS using the interior point method. The results can be seen in table \ref{pa02_stats}. Similar to the Rome Graphs benchmark, our algorithm yields smaller mean hitting sets and has a smaller mean ratio.

$$
\begin{table}[h]
    \centering
    \begin{subtable}[t]{0.5\textwidth}
        \input{out/pa02/hs_pa02_highs_ipm}
        \caption{Hitting Set LP (1), HiGHS (ipm)}
    \end{subtable}
    \hfill
    \begin{subtable}[t]{0.4\textwidth}
        \input{out/pa02/master_std}
        \caption{3-HS 3-approx, strategy 3}
    \end{subtable}
    \caption{Results for HiGHS (ipm) and our algorithm on 100 random preferential attachment hypergraphs, $p=0.2$}
    \label{pa02_stats}
\end{table}
$$

# Conclusion

The implemented \textsc{Hitting Set} algorithm can compute hitting sets for large inputs in a reasonable amount of time. In order to achieve this, it is important that every rule used in the algorithm is reasonably fast. A "bad" implementation of a single rule can slow down the entire algorithm. We experienced this for both the approximative (double) vertex domination rules. A naive implementation is fast for smaller hypergraphs, but performance will degrade significantly for large hypergraphs. We switched to an incremental algorithm design, to cope with long running times on large hypergraphs. Our algorithm can also compete with LP randomized rounding based algorithms, and even manages to outperform them on select hypergraphs. We could also observe that the algorithm never exceeded an estimated ratio of 2 on our benchmarks. 

## Further Questions {-}

The only part of the algorithm that is parallelized is the edge domination rule. The incremental variant of the main algorithm could also be parallelized. The hypergraph could at some point during the run of the algorithm consist of multiple disjoint connected components. These components could then be processed in parallel. Every component would yield its own number of rule executions and hitting sets, which are aggregated at the end. The question is, if the overhead of finding these components is greater than the possible speed up.

There already exist parallel algorithms for the vertex and edge domination rule, which run on a GPU as shown by Bevern et al. in [@DBLP:journals/jcss/BevernKSST24]. It would be interesting to know if such algorithms also exist for the other domination type reduction rules, namely the approximate (double) vertex domination rule. Our implementation of these rules rely on arguments about vertex degrees and adjacency between vertices, which is similar to the standard vertex domination implementation. 

Could our hypergraph data structure be further improved? Right now we only use "primitive" data structures like maps to model our hypergraph. Improvements could be made for both the incidence and adjacency map struct fields. When accessing elements of a map with the \texttt{range} operator, Go does not guarantee that the order of elements is always the same. Using a real sparse matrix could improve performance and make the algorithm more deterministic.

We could technically use optimal solutions of the hitting set LP to assess the quality of our estimated optimum. The value of an optimal solution of the ILP is at least the value of an optimal solution of the LP relaxation. If our estimated optimum is lower than an optimal solution of the LP, then we know that we are underestimating the actual optimum by at least $\text{Opt}^* - \text{est. opt}$.

Do rule strategies exist that are better than strategy 3? One might have noticed a small flaw in rule strategy 3. Let us highlight this flaw with a small example. Let $G$ be a hypergraph, where none of the exact rules are applicable. According to rule strategy 3 we now apply the next approximative rule $R_0$ in the precedence. We then construct our auxiliary graph $H$. Let $R_1$ be the only applicable rule in $H$. The algorithm would apply $R_1$, even if there exist other approximative rules lower in the precedence that are applicable in $G$.
