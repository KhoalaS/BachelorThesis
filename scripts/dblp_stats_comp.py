import os
import pandas as pd
import argparse
from dotenv import dotenv_values
from pyecharts.charts import Bar, Page
import opts as _opts
from opts import rule_names
import pyecharts.options as opts

parser = argparse.ArgumentParser()
parser.add_argument("file", metavar="FILE", help="path to first csv file")
parser.add_argument("comp", metavar="FILE", help="path to second csv file")
parser.add_argument(
    "-e", type=str, default="./scripts/dblp.env", help="path to env file")

args = parser.parse_args()

stock = pd.read_csv(args.file, delimiter=";")
stock_stats = stock.describe()
stock_stats.drop(["count", "25%", "75%"], inplace=True)
stock_stats.rename(index={"50%": "median"}, inplace=True)
stock_stats.drop(columns=["Vertices", "Edges"], inplace=True)


noedom = pd.read_csv(args.comp, delimiter=";")
noedom_stats = noedom.describe()
noedom_stats.drop(["count", "25%", "75%"], inplace=True)
noedom_stats.rename(index={"50%": "median"}, inplace=True)
noedom_stats.drop(columns=["Vertices", "Edges"], inplace=True)

config = dotenv_values(args.e)

short_labels = []
vals_stock = []
vals_noedom = []
vals_diff = []


for rule, short in rule_names.items():
    short_labels.append(short)
    v_0 = stock_stats.loc["mean", rule]
    v_1 = noedom_stats.loc["mean", rule]
    vals_stock.append(v_0)
    vals_noedom.append(v_1)
    d = ((v_1 - v_0)/v_0)*100
    print(d, rule)
    vals_diff.append(d)

diff_df = pd.DataFrame([vals_diff], columns=short_labels)
#diff_table = diff_df.to_latex()
#diff_out = open("./out/diff_table.md", "w+")
#diff_out.write(diff_table)

page = Page()
bar = Bar().set_global_opts(title_opts=opts.TitleOpts(
    title=config.get("TITLE_1")), toolbox_opts=_opts.img_opts())
bar.add_xaxis(short_labels)
bar.add_yaxis("stock", vals_stock, color="#00ff99",
              label_opts=opts.LabelOpts(is_show=False))
bar.add_yaxis("no edge domination", vals_noedom,
              color="#0099ff", label_opts=opts.LabelOpts(is_show=False))


diff = Bar().set_global_opts(yaxis_opts=opts.AxisOpts(axislabel_opts=opts.LabelOpts()), title_opts=opts.TitleOpts(
    title="#Rule Executions, Percent Change of Mean"), toolbox_opts=_opts.img_opts())
diff.add_xaxis(short_labels)
diff.add_yaxis("percent change", vals_diff, color="#00ff99",
               label_opts=opts.LabelOpts(is_show=False))

page.add(bar)
page.add(diff)
page.render("./out/test.html")
os.system("sed -i 's/https:\/\/assets.pyecharts.org\/assets\/v5\/echarts.min.js/.\/echarts.min.js/g' ./out/test.html")
