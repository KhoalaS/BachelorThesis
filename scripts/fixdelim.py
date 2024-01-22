import argparse
import csv
import os

parser = argparse.ArgumentParser(prog="fixdelim")

parser.add_argument("dirname", metavar="DIR",
                    help="path to folder with csv files")
parser.add_argument("-s", default=1, type=int,
                    help="If placeholder '%%' is present in dir path, replace with sequence (s..t)")
parser.add_argument("-t", default=1, type=int)

args = parser.parse_args()


for i in range(args.s, args.t+1):
    dirs = []
    inrepl = str(args.dirname)
    inrepl = inrepl.replace("%", str(i))

    data = os.listdir(inrepl)

    for filename in data:
        fullpath = "{}/{}".format(inrepl, filename)

        file = open(fullpath, "r")
        reader = csv.reader(file, delimiter=";")
        lines = []
        for row in reader:
            line = ";".join(row[:-1])
            lines.append(line)

        file.close()
        file = open(fullpath, "w")
        for line in lines:
            file.write(line+"\n")
