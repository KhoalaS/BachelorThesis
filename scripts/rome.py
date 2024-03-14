import pandas as pd

df = pd.read_csv("./data/rome_master.csv", delimiter=";")
fileAgg = df.agg({"File": ["mean"]})
print(fileAgg.head())