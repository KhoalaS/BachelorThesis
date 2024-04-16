import pandas as pd
import argparse
from opts import rule_names

parser = argparse.ArgumentParser()
parser.add_argument("file", metavar="FILE", help="path to csv file")
parser.add_argument("out", metavar="OUT", help="path to output tex file")

args = parser.parse_args()

print("loading csv file...")
df = pd.read_csv(args.file, delimiter=";")

print("file loaded...")

opt = pd.read_csv("data_final/rome_cvd_ilphs_clp.csv", delimiter=";")
opt.drop(columns=["RatioUB", "Ratio", "HittingSet"], inplace=True)
opt.rename(columns={"Opt": "opt"}, inplace=True)
opt["File"] = opt["File"].str.replace(".txt", "")

df = pd.merge(df, opt, how="inner", on="File")

df = df.loc[df.groupby("File")["Ratio"].idxmin()]

df.drop(columns=["OVertices", "OEdges",
        "Vertices", "Edges"], inplace=True)

df["actual ratio"] = df["HittingSet"]/df["opt"]
df.drop(columns=["File"], inplace=True)

rome_stats = df.describe()
rome_stats.drop(["count", "25%", "75%"], inplace=True)
rome_stats.rename(index={"50%": "median"}, inplace=True)
rome_stats["Ratio"] = rome_stats["Ratio"].round(4)
rome_stats["HittingSet"] = rome_stats["HittingSet"].round(2)
rome_stats["Opt"] = rome_stats["Opt"].round(2)
rome_stats.rename(
    columns={"Ratio": "ratio", "HittingSet": "$|C|$", "Opt": "est. opt"}, inplace=True)

for k, v in rule_names.items():
    # rome_stats[k] = rome_stats[k].round(2)
    # rome_stats.rename(columns={k: v}, inplace=True)
    rome_stats.drop(columns=[k], inplace=True)

rome_tbl = rome_stats[["ratio", "actual ratio", "est. opt",
                      "opt", "$|C|$"]].to_latex(float_format="%.4f")
print(rome_stats)

f = open(args.out, "w+")
f.write(rome_tbl)
f.close()
