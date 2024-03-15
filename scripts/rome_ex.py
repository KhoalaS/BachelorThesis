import pandas as pd
import argparse

rule_names = {
    "kTiny": "Tiny",
    "kVertDom": "VD",
    "kEdgeDom": "ED",
    "kSmall": "Small",
    "kTri": "Tri",
    "kExtTri": "ETri",
    "kApVertDom": "AVD",
    "kApDoubleVertDom": "ADVD",
    "kSmallEdgeDegTwo": "SED2",
    "kSmallEdgeDegTwo2": "SED2[1.5]",
    "kFallback": "F3"
}

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

#df = df[df['ED'] > 150 ]
print(df.describe())

print(df.shape[0])
foi = []


for idx, row in df.iterrows():
    c_ratio = c.loc[idx]["Ratio"]
    if (row["Ratio"] - c_ratio) > 0.4:
        foi.append(idx)

print(len(foi))

for file in foi:
    print(df.loc[file])