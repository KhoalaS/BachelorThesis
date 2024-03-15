import pandas as pd
import argparse
from opts import rule_names

parser = argparse.ArgumentParser()
parser.add_argument("file", metavar="FILE", help="path to csv file")

args = parser.parse_args()

df = pd.read_csv(args.file, delimiter=";")

for k, v in rule_names.items():
    df.drop(columns=[k], inplace=True)
    #df.rename(columns={k: v}, inplace=True)


dblp_stats = df.describe()
dblp_stats.drop(["count", "25%", "75%"], inplace=True)
dblp_stats.rename(index={"50%": "median"}, inplace=True)

dblp_stats.drop(columns=["Vertices", "Edges"], inplace=True)
out = open("./out/dblp_stats.md", "w+")
out.write(dblp_stats.to_latex())
