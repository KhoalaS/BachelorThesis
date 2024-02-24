import pandas as pd
import argparse

parser = argparse.ArgumentParser()
parser.add_argument("file", metavar="FILE", help="path to csv file")

args = parser.parse_args()

df = pd.read_csv(args.file)

print(df)