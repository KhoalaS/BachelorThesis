# Building

## Requirements

-   golang version 1.22

Clone the repo

```
git clone https://github.com/KhoalaS/BachelorThesis.git
```

In the project directory,

Build the binaries

```
go build cmd/main.go
go build cmd/tvd/tvd.go
```

Or use the script

```
./scripts/build.sh
```

Run the main program

```
./main
```

or use the `go run` subcommand

```
go run cmd/main.go
```

# Examples

## Main

Run with a Erdős–Rényi hypergraph with 1000 vertices and ~10000 edges.

```bash
./main -er -n=1000 -evr=10
```

Run with a preferential attachment hypergraph with 1000 vertices and a probability of 0.2 of adding a new vertex.

```bash
./main -pa=0.2 -n=1000
```

## TVD

Run with standard graph stored in `graph_file.txt` and use the frontier algorithm.

```bash
./tvd -i=graph_file.txt -fr
```

Run with standard Erdős–Rényi graph with 1000 vertices and ~20000 edges and use the frontier algorithm.

```bash
./tvd -n=1000 -er -evr=20 -fr
```

Run with standard Erdős–Rényi graph with 1000 vertices and ~20000 edges and use the frontier algorithm. Run the alorithm for 100 graphs and store the log files under `./data/er_20_1K`

```bash
./tvd -n=1000 -er -evr=20 -fr -log=100 -d=./data/er_20_1K
```

# Usage Main

```bash
./main [OPTIONS]
```

## Options

```
Usage of ./main:
  -d string
    	Output directory for log files. (default "./data")
  -er
    	Generate a andom Erdös Renyi hypergraph
  -evr int
    	Maximum ratio |E|/|V| to compute for random graphs.
  -f string
    	Generate a random hypergraph with fixed ratios for the edge sizes.
  -i string
    	Filepath to input file.
  -log int
    	Amount of logging passes.
  -m int
    	Number of edges if no graph file supplied. (default 20000)
  -n int
    	Number of vertices if no graph file supplied. (default 10000)
  -o string
    	Export the generated graph with the given string as filename. The will create a 'graphs' folder where the file is located.
  -os string
    	Export the generated graph to the given filepath.
  -pa float
    	Generate a random preferential attachment hypergraph with given float as probablity to add a new vertex.
  -pamod

  -prof
    	Make CPU profile
  -u int
    	Generate a u-uniform graph.

```

# Usage Triangle Vertex Deletion

```bash
./tvd [OPTIONS]
```

```
Usage of ./tvd:
  -d string
    	output directory (default "./data")
  -evr float
    	targetted edge/vertex ratio, takes priority over p
  -fr
    	use frontier algorithm
  -i string
    	path to input graph file
  -log int
    	log the number of rule executions, do log many runs (default 1)
  -n int
    	number of vertices (default 1000)
  -p float
    	probability of adding an edge (default 0.5)
  -prof
    	make pprof profile

```
