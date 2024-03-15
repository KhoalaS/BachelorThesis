from pyecharts.charts import Boxplot, Line, Page, Bar
import argparse
import os
import csv
from pyecharts import options as opts
import statistics as st
import math
from dotenv import dotenv_values
from opts import rule_names
import opts as _opts


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
parser.add_argument("-s", type=int, default=1,
                    help="If placeholder '%%' is present in dir path, replace with sequence (s..t)")
parser.add_argument("-t", type=int, default=1)
parser.add_argument("-m", type=int, action="append",
                    help="where to extract rule execution data")
parser.add_argument("-o", type=str, help="path to output file")
parser.add_argument("-st", type=int, default=1,
                    help="steps to go through from s to t")
parser.add_argument("-e",type=str, default="./scripts/.env", help="path to env file")

args = parser.parse_args()

config = dotenv_values(args.e)

if args.m == None:
    args.m = [int(args.s)]

ratios = []
marks = {}
m_empty = True

labels = []
minval = 3

for i in range(args.s, args.t+args.st, args.st):
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
                if m_empty and args.m != None:
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

seq = eval(config.get("X_AXIS_LABELS_0"))

minval = round_down(minval, 1)
avgs = []
for arr in ratios:
    avgs.append(st.mean(arr))

short_labels = [rule_names[x] for x in labels]

page = Page()

line = Line()
line.set_global_opts(toolbox_opts=_opts.img_opts())
line.add_xaxis(range(0, len(avgs))).add_yaxis("mean", avgs, yaxis_index=1)
line.set_series_opts(label_opts=opts.LabelOpts(
    is_show=False), itemstyle_opts=opts.ItemStyleOpts(color="orange"))

box = Boxplot()
box.add_xaxis(seq)
box.add_yaxis("est. appr. ratio", box.prepare_data(ratios))
box.extend_axis(yaxis=opts.AxisOpts(min_=minval))

box.set_global_opts(title_opts=opts.TitleOpts(title=config.get("TITLE_0"), subtitle=config.get("SUBTITLE_0")),
                    yaxis_opts=opts.AxisOpts(min_=minval), xaxis_opts=opts.AxisOpts(name=config.get("X_AXIS"), name_gap=30), toolbox_opts=_opts.img_opts())
box.overlap(line)
page.add(box)

subtitle_insert = eval(config.get("SUBTITLE_1_MARK"))

if args.m != None:
    i = 0
    for mark in args.m:
        r_box = Boxplot()
        r_box.set_global_opts(title_opts=opts.TitleOpts(title=config.get(
            "TITLE_1"), subtitle=config.get("SUBTITLE_1").format(subtitle_insert[i])), toolbox_opts=_opts.img_opts())
        r_box.add_xaxis(short_labels)
        rule_boxdata = r_box.prepare_data(marks[int(mark)])
        r_box.add_yaxis("rule exectutions", rule_boxdata)
        r_box.set_colors("green")
        page.add(r_box)
        i += 1

page.render(str(args.o))
