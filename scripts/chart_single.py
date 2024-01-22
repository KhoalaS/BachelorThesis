from pyecharts.charts import Page, Line, Bar
import pyecharts.options as opts 
import csv
import argparse
import os
import math

def round_down(x: float, d: int) -> float:
    mul = math.pow(10, d)
    y = x*mul
    z = math.floor(y)
    return float(z)/mul

rule_names = {
    "kTiny": "Tiny",
    "kVertDom": "VDom",
    "kEdgeDom": "EDom",
    "kSmall": "Small",
    "kTri": "Tri",
    "kExtTri": "ETri",
    "kApVertDom": "AVDom",
    "kApDoubleVertDom": "ADVDom",
    "kSmallEdgeDegTwo": "SETwo",
    "kFallback": "F3"
}

parser = argparse.ArgumentParser(prog="chart_single")
parser.add_argument("folder")
args = parser.parse_args()

ratios = []
min_ratio = 3
rules = []
labels = []

for filename in os.listdir(str(args.folder)):
    fullpath = "{}/{}".format(args.folder, filename)
    file = open(fullpath, "r")
    reader = csv.reader(file, delimiter=";")

    header = True
    is_master_file = filename.find("master") != -1

    for row in reader:
        if header:
            header = False
            labels = row[1:]
            continue

        if is_master_file:
            rules = row[1:]
        else:
            ratio = float(row[0])
            if ratio < min_ratio:
                min_ratio = ratio
            ratios.append(ratio)

short_labels = [rule_names[x] for x in labels]

page = Page()
line = Line().add_xaxis(range(0, len(ratios)))
line.add_yaxis("est. appr. ratio", ratios, symbol="none", color="orange")
line.set_global_opts(title_opts=opts.TitleOpts(title="Est. appr. ratio", subtitle="The ratio was calculated after every F3 rule execution"), yaxis_opts=opts.AxisOpts(min_=round_down(min_ratio, 1)))

bar = Bar().set_global_opts(title_opts=opts.TitleOpts(title="#Rule Executions"))
bar.add_xaxis(short_labels)
bar.add_yaxis("rule exectutions", rules)

page.add(line)
page.add(bar)
page.render("./out/dblp_vtd.html")