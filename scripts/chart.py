from pyecharts.charts import Boxplot, Line
import argparse
import os
import csv
from pyecharts import options as opts
import statistics as st


def floatFormat(params):
    return str(round(params.data, 2))


parser = argparse.ArgumentParser(
    prog='chart.py',
    description='Chart est. approximation factor from csv files')


parser.add_argument("dirname", metavar="DIR",
                    help="path to folder with csv files")
parser.add_argument("-s", default=1, type=int,
                    help="If placeholder '%%' is present in dir path, replace with sequence (s..t)")
parser.add_argument("-t", default=1, type=int)

args = parser.parse_args()

ratios = []
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
                header = True
            else:
                val = float(row[0])
                if val < minval:
                    minval = val
                ratio.append(val)

        ratios.append(ratio)

seq = [args.s + x for x in range(0, len(ratios))]

minval = round(minval, 1)
avgs = []
for arr in ratios:
    avgs.append(st.mean(arr))

line = Line()
line.add_xaxis(range(0, len(avgs))).add_yaxis("mean", avgs, yaxis_index=1)
line.set_series_opts(label_opts=opts.LabelOpts(
    is_show=False), itemstyle_opts=opts.ItemStyleOpts(color="orange"))

box = Boxplot()
box.add_xaxis(seq)
box.add_yaxis("est. appr. ratio", box.prepare_data(ratios))
box.extend_axis(yaxis=opts.AxisOpts(min_=minval))

box.set_global_opts(title_opts=opts.TitleOpts(title="Triangle Vertex Deletion for ER graphs", subtitle="1000 vertices, Edge\\Vertex ratios of 10-25"),
                    yaxis_opts=opts.AxisOpts(min_=minval), xaxis_opts=opts.AxisOpts(name="EVR", name_gap=30))
box.overlap(line)
box.render()
