import networkx as nx
import sys

G = nx.Graph()

gFile = open(sys.argv[1])

for line in gFile:
    v = line.strip().split(",")
    G.add_edge(int(v[1]), int(v[2]), id=int(v[0]))

cover = nx.min_edge_cover(G)
ids = []

for e in cover:
    ids.append(G.edges[e]['id'])
    
print(ids)

sys.exit(0)