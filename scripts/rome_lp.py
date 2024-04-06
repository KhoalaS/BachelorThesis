import pandas as pd
import argparse

parser = argparse.ArgumentParser()
parser.add_argument("file", metavar="FILE", help="path to csv file")
parser.add_argument("out", metavar="OUT", help="path to output tex file")

args = parser.parse_args()

print("loading csv file...")
df = pd.read_csv(args.file, delimiter=";")

print("file loaded...")

df = df.loc[df.groupby("File")["HittingSet"].idxmin()]

rome_stats = df.describe()
rome_stats.drop(["count", "25%", "75%"], inplace=True)
rome_stats.rename(index={"50%": "median"}, inplace=True)
rome_stats["Ratio"] = rome_stats["Ratio"].round(4)
rome_stats["HittingSet"] = rome_stats["HittingSet"].round()
rome_stats["Opt"] = rome_stats["Opt"].round()
rome_stats.rename(columns={"Ratio": "ratio", "HittingSet": "$|C|$", "Opt": "$\\textnormal{Opt}^*$"}, inplace=True)

rome_tbl = rome_stats.to_latex(float_format="%.4f")

f = open(args.out, "w+")
f.write(rome_tbl)
f.close()
print(rome_stats)
