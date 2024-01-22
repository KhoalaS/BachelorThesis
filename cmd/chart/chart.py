from pyecharts.charts import Boxplot
import argparse
import os
import csv
from pyecharts import options as opts
import math

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

box = Boxplot()
box.add_xaxis(seq)
box.add_yaxis("Ratios", box.prepare_data(ratios))

minval = round(minval, 1)
box.set_global_opts(yaxis_opts=opts.AxisOpts(min_= minval))

box.render()
