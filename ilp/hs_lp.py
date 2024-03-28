from pulp import *
import argparse
from random import random

parser = argparse.ArgumentParser()
parser.add_argument("input", metavar="FILE", help="path to input graph file")
parser.add_argument("--highs", action='store_true',
                    help="use the HiGHS solver")
parser.add_argument("--glpk", action='store_true',
                    help="use the GLPK solver")
parser.add_argument("-l", action='store_true', help="keep log files")

args = parser.parse_args()

file = open(args.input, "r")
V_lookup = {}
V_counter = 1
E_counter = 1
V = []
E = {}

for line in file:
    line_str = line.strip()
    e = [int(ep) for ep in line_str.split()]
    e_tr = []

    for v in e:
        if v not in V_lookup:
            V_lookup.update({v: V_counter})
            e_tr.append(V_counter)
            V_counter += 1
        else:
            e_tr.append(V_lookup.get(v))

    E.update({E_counter: e_tr})
    E_counter += 1

print("file loaded...")

n = V_counter-1
m = len(E)

V = [x for x in range(1, n+1)]
print("graph had {} vertices and {} many edges".format(n, m))

prob = LpProblem("VC-Relax", LpMinimize)

x = LpVariable.dicts("x", range(1, n+1), 0, 1)
prob += lpSum([x[j] for j in range(1, n+1)])


for idx, e in E.items():
    e_sum = lpSum([x[j] for j in e])
    prob += e_sum >= 1

print("begin solving...")

keep_logs = args.l

if args.highs:
    prob.solve(HiGHS_CMD(mip=False, msg="using HiGHS", keepFiles=keep_logs,
                         path="/usr/local/bin/highs", threads=os.cpu_count()))
elif args.glpk:
    prob.solve(GLPK(msg="using GLPK solver", keepFiles=keep_logs))
else:
    prob.solve(PULP_CBC_CMD(keepFiles=keep_logs))

opt = 0.0
print("Status:", LpStatus[prob.status])
if prob.status == LpStatusOptimal:
    print("Solution:")
    opt = value(lpSum([x[j] for j in V]))
    print("Sum of decision variables =", opt)


S_1 = []
counter = 0
out = open("/tmp/s1.sol", "ab+")

for key, val  in V_lookup.items():
    if value(x[val]) == 1:
        print(key, val)
        S_1.append(key)
        counter += 1
        if counter == 8192:
            out.write(b''.join([j.to_bytes(4, byteorder='big', signed=False) for j in S_1]))  # Assuming 32-bit integers
            counter = 0
            S_1.clear()

if counter > 0:
    out.write(b''.join([j.to_bytes(4, byteorder='big', signed=False) for j in S_1]))

out.close()
