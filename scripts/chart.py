from pyecharts.charts import Boxplot, Line, Page, Bar
import argparse
import os
import csv
from pyecharts import options as opts
import statistics as st
import math

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

def round_down(x: float, d: int) -> float:
    mul = math.pow(10, d)
    y = x*mul
    z = math.floor(y)
    return float(z)/mul

parser = argparse.ArgumentParser(
    prog='chart.py',
    description='Chart est. approximation factor from csv files')


parser.add_argument("dirname", metavar="DIR",
                    help="path to folder with csv files")
parser.add_argument("-s", default=1, type=int,
                    help="If placeholder '%%' is present in dir path, replace with sequence (s..t)")
parser.add_argument("-t", default=1, type=int)
parser.add_argument("-m", type=int, action="append")

args = parser.parse_args()

ratios = []
marks = {}
m_empty = True


labels = []
minval = 3

for i in range(args.s, args.t+1):
    dirs = []
    inrepl = str(args.dirname)
    inrepl = inrepl.replace("%", str(i))

    data = os.listdir(inrepl)

    for filename in data:
        if filename.find("master") == -1:
            continue

        fullpath = "{}/{}".format(inrepl, filename)
        file = open(fullpath, "r")
        reader = csv.reader(file, delimiter=';')

        header = False
        ratio = []

        for row in reader:
            if not header:
                labels = row[1:]
                if m_empty:
                    for k in args.m:
                        marks.update({int(k): [[] for _ in labels]})
                    m_empty = False
                header = True
            else:
                val = float(row[0])
                if val < minval:
                    minval = val
                ratio.append(val)
                if i in marks.keys():
                    for j in range(1, len(row)):
                        marks[i][j-1].append(int(row[j]))

        ratios.append(ratio)

seq = [args.s + x for x in range(0, len(ratios))]

minval = round_down(minval, 1)
avgs = []
for arr in ratios:
    avgs.append(st.mean(arr))

short_labels = [rule_names[x] for x in labels]

page = Page()

line = Line()
line.add_xaxis(range(0, len(avgs))).add_yaxis("mean", avgs, yaxis_index=1)
line.set_series_opts(label_opts=opts.LabelOpts(
    is_show=False), itemstyle_opts=opts.ItemStyleOpts(color="orange"))

box = Boxplot()
box.add_xaxis(seq)
box.add_yaxis("est. appr. ratio", box.prepare_data(ratios))
box.extend_axis(yaxis=opts.AxisOpts(min_=minval))

box.set_global_opts(title_opts=opts.TitleOpts(title="Triangle Vertex Deletion for ER graphs", subtitle="1000 vertices, Edge\\Vertex ratios of 10-25, 100 graphs per EVR"),
                    yaxis_opts=opts.AxisOpts(min_=minval), xaxis_opts=opts.AxisOpts(name="EVR", name_gap=30))
box.overlap(line)
page.add(box)

for mark in args.m:
    r_box = Boxplot()
    r_box.set_global_opts(title_opts=opts.TitleOpts(title="#Rule Executions", subtitle="EVR={} for input ER graph".format(mark)))
    r_box.add_xaxis(short_labels)
    rule_boxdata = r_box.prepare_data(marks[int(mark)])
    r_box.add_yaxis("rule exectutions", rule_boxdata)
    r_box.set_colors("green")
    page.add(r_box)

page.render("./out/er_10_25_vtd.html")

