import pandas as pd
print("loading csv file...")
df = pd.read_csv("./data/rome_master.csv", delimiter=";")

print("file loaded...")
df.drop(columns=["File", "OVertices", "OEdges",
        "Vertices", "Edges"], inplace=True)
rome_stats = df.describe()
rome_stats.drop(["count"], inplace=True)
rome_tbl = rome_stats.to_latex()

f = open("./out/rome_stats.tex", "w+")
f.write(rome_tbl)
f.close()
print(rome_stats)