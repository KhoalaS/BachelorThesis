import pandas as pd
import argparse

parser = argparse.ArgumentParser()
parser.add_argument("file", metavar="FILE", help="path to csv file")
parser.add_argument("out", metavar="OUT", help="path to output tex file")

args = parser.parse_args()

print("loading csv file...")
df = pd.read_csv(args.file, delimiter=";")

print("file loaded...")
df.drop(columns=["File", "OVertices", "OEdges",
        "Vertices", "Edges"], inplace=True)
rome_stats = df.describe()
rome_stats.drop(["count"], inplace=True)
rome_tbl = rome_stats.to_latex()

f = open(args.out, "w+")
f.write(rome_tbl)
f.close()
print(rome_stats)