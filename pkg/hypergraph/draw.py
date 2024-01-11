import matplotlib.pyplot as plt
import hypernetx as hnx
import sys

edges = {}
data = {}

file = open(sys.argv[1], "r")

datafile = None
if len(sys.argv) > 2:
    datafile = open(sys.argv[2], "r")
    for line in datafile:
        s = line.strip().split()
        data.update({s[0]: {"com": s[1]}})
    


i = 0
for line in file:
    s = line.strip().split()
    edges.update({i: tuple(s)})
    i += 1

g = hnx.Hypergraph(edges, node_properties=data)

colors_map = {
    "0": "red",
    "1": "green",
    "2": "blue",
    "3": "yellow",
    "4": "pink",
    "5": "orange",
    "6": "purple",
    "7": "brown",
    "8": "gray",
    "9": "olive",
    "10":"cyan",
}

plt.subplots(figsize=(5, 5))

if datafile == None:
    hnx.draw(g, with_edge_labels=False,
             with_edge_counts=False,
             with_node_labels=False,
             with_node_counts=False,
             node_radius=0.5)
else:
    colors = [
        colors_map[g.get_properties(id=v, level=1, prop_name="com")]
        for v in list(g.nodes)
    ]
    hnx.draw(g, with_edge_labels=False,
             with_edge_counts=False,
             with_node_labels=False,
             with_node_counts=False,
             node_radius=0.5,
             nodes_kwargs={
                 'facecolors': colors
             })

plt.show()
