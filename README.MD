# Building

## Requirements

-   golang version >= 1.18
-   python3 NetworkX

Clone the repo

```
git clone https://github.com/KhoalaS/BachelorThesis.git
```

In the project directory,

build an executable

```
go build cmd/main.go
```

just run the main.go program

```
go run cmd/main.go
```

# Usage

```bash
go run cmd/main.go [OPTIONS]
```

## Options

```
-c          Make charts.
-evr int    Maximum ratio |E|/|V| to compute for random graphs.
-f string   Generate a random hypergraph with fixed ratios for the edge sizes.
-fr int     Preprocess the graph with fr many Factor-3 rule executions.
-i string   Filepath to graphml input file.
-k int      The parameter k.
-m int      Number of edges if no graph file supplied. (default 20000)
-maxv int   Maximum vertices for random graphs used in charts.
-n int      Number of vertices if no graph file supplied. (default 10000)
-o string   Export the generated graph with the given string as filename. The will create a 'graphs' folder where the file is located.
-p string   Use a preconfigured chart preset. For available presets run with 'list -p'.
-pa int     Generate a random preferential attachment hypergraph
-prof       Make CPU profile
-u int      Generate a u-uniform graph.
```
