import matplotlib.pyplot as plt
import hypernetx as hnx
import sys

file = open(sys.argv[1], "r") 

edges = {}

i = 0
for line in file:
    s = line.strip().split()
    edges.update({i : tuple(s)})
    i += 1

g = hnx.Hypergraph(edges)

plt.subplots(figsize=(5,5))
hnx.draw(g, with_edge_labels=False, with_edge_counts=False, with_node_labels=False, with_node_counts=False, node_radius=0.5)

plt.show()
