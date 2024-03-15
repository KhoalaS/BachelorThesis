import pandas as pd
import argparse
from opts import rule_names


parser = argparse.ArgumentParser()
parser.add_argument("file", metavar="FILE", help="path to csv file")
parser.add_argument("comp", metavar="COMP", help="file to compare to")

args = parser.parse_args()

print("loading csv file...")
df = pd.read_csv(args.file, delimiter=";")
c = pd.read_csv(args.comp, delimiter=";")
df.set_index("File", inplace=True)
c.set_index("File", inplace=True)

print("file loaded...")
#df.drop(columns=["OVertices", "OEdges",
#        "Vertices", "Edges"], inplace=True)

for k, v in rule_names.items():
    df.rename(columns={k: v}, inplace=True)
    c.rename(columns={k: v}, inplace=True)

df = df.groupby('File').mean()
c = c.groupby('File').mean()


print(df.shape[0])
foi = []

for idx, row in df.iterrows():
    c_ratio = c.loc[idx]["Ratio"]
    if (row["Ratio"] - c_ratio) > 0.4:
        foi.append(idx)

print(len(foi))

for file in foi:
    pass
    print(c.loc[file])
    print(df.loc[file])