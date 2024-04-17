import pandas as pd
import argparse
from opts import rule_names

parser = argparse.ArgumentParser()
parser.add_argument("file", metavar="FILE", help="path to csv file")
parser.add_argument("out", metavar="OUT", help="path to output tex file")
parser.add_argument("--opt", action="store_true")
parser.add_argument("--short", action="store_true")
parser.add_argument("--deriv", action="store_true")



args = parser.parse_args()

print("loading csv file...")
df = pd.read_csv(args.file, delimiter=";")
filename = str(args.file).split("/")[-1]

print("file loaded...")

if args.opt:
    opt = pd.read_csv("data_final/rome_cvd_ilphs_clp.csv", delimiter=";")
    opt.drop(columns=["RatioUB", "Ratio", "HittingSet"], inplace=True)
    opt.rename(columns={"Opt": "opt"}, inplace=True)
    opt["File"] = opt["File"].str.replace(".txt", "")

    df = pd.merge(df, opt, how="inner", on="File")

df = df.loc[df.groupby("Filename")["Ratio"].idxmin()]

df.drop(columns=["Vertices", "Edges"], inplace=True)

if args.deriv:
    df.drop(columns=["OVertices", "OEdges"], inplace=True)

if args.opt:
    df["actual ratio"] = df["HittingSet"]/df["opt"]

rome_stats = df.describe()
rome_stats.drop(["count", "25%", "75%"], inplace=True)
rome_stats.rename(index={"50%": "median"}, inplace=True)
rome_stats["Ratio"] = rome_stats["Ratio"].round(4)
rome_stats["HittingSet"] = rome_stats["HittingSet"].round(2)
rome_stats["Opt"] = rome_stats["Opt"].round(2)
rome_stats.rename(
    columns={"Ratio": "ratio", "HittingSet": "$|C|$", "Opt": "est. opt"}, inplace=True)

for k, v in rule_names.items():
    if args.short:
        rome_stats.drop(columns=[k], inplace=True)
    else:
        rome_stats[k] = rome_stats[k].round(2)
        rome_stats.rename(columns={k: v}, inplace=True)

if args.opt:
    rome_tbl = rome_stats[["ratio", "actual ratio", "est. opt",
                      "opt", "$|C|$"]].to_latex(float_format="%.4f")
else:
    rome_tbl = rome_stats.to_latex(float_format="%.4f")

print(rome_stats)

f = open("{}/{}".format(args.out, filename.replace(".csv", ".tex")), "w+")
f.write(rome_tbl)
f.close()
