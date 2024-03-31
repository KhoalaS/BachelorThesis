import pandas as pd
from pyecharts.charts import Bar, Grid
import pyecharts.options as opts
import opts as _opts
from opts import rule_names
import os

df_base = pd.read_csv(
    "./data/dblp_base/master_CUSTOM_1711049002.csv", delimiter=";")
df_str1 = pd.read_csv(
    "./data/dblp_str1/master_ER_1710866151.csv", delimiter=";")
df_str2 = pd.read_csv(
    "./data/dblp_str2/master_ER_1710773299.csv", delimiter=";")
df_str3 = pd.read_csv(
    "./data/dblp_str3/master_CUSTOM_1710947310.csv", delimiter=";")

frames = [df_base, df_str1, df_str2, df_str3]
strat = ["base", "str1", "str2", "str3"]
colors = ["#3398DB",
          "#FFC107",
          "#FF5722",
          "#9C27B0"]

exclude = ["kTiny",
           "kVertDom",
           "kEdgeDom",
           "kApVertDom",
           "kApDoubleVertDom",
           "kFallback"]

out = open("out/dblp_rules.md", "w+")

i = 0
for df in frames:
    df = df.describe()
    df.drop(["count", "25%", "75%"], inplace=True)
    df.rename(index={"50%": "median"}, inplace=True)
    df.drop(columns=["Vertices", "Edges"], inplace=True)
    df["Ratio"] = df["Ratio"].round(4)
    df["HittingSet"] = df["HittingSet"].round(2)
    df["Opt"] = df["Opt"].round(2)
    df["Time"] = df["Time"].round()

    frames[i] = df
    i += 1
    out.write(df.to_markdown())
    out.write("\n\n")

bar = Bar().set_global_opts(toolbox_opts=_opts.img_opts())

short_labels = [label for key, label in rule_names.items() if key in exclude]
bar.add_xaxis(short_labels)

i = 0
for df in frames:
    arr = []
    for col in df.columns:
        if col in exclude:
            arr.append(df.loc["mean", col])

    bar.add_yaxis(strat[i], arr, color=colors[i])
    print(arr)
    i += 1

bar.set_global_opts(yaxis_opts=opts.AxisOpts(name="rule executions"))
bar.set_series_opts(label_opts=opts.LabelOpts(is_show=False))

for df in frames:
    df.drop(columns=exclude, inplace=True)

bar_s = Bar().set_global_opts(toolbox_opts=_opts.img_opts())

short_labels = [label for key, label in rule_names.items()
                if key not in exclude]
bar_s.add_xaxis(short_labels)

i = 0
for df in frames:
    arr = [x for x in df.loc["mean"].values[1:-3]]
    bar_s.add_yaxis(strat[i], arr, color=colors[i])
    i += 1

bar_s.set_series_opts(label_opts=opts.LabelOpts(is_show=False))

grid = Grid()

grid.add(bar_s, grid_opts=opts.GridOpts(pos_left="60%"))
grid.add(bar, grid_opts=opts.GridOpts(pos_right="50%"))

grid.render("out/dblp_rules.html")
os.system("sed -i 's/https:\/\/assets.pyecharts.org\/assets\/v5\/echarts.min.js/.\/echarts.min.js/g' out/dblp_rules.html")
