from pulp import *
import argparse
from math import exp, log
from random import random
import os

parser = argparse.ArgumentParser()
parser.add_argument("input", metavar="FILE", help="path to input file")
parser.add_argument("--highs", action='store_true',
                    help="use the HiGHS solver")
parser.add_argument("--glpk", action='store_true',
                    help="use the GLPK solver")
parser.add_argument("--cplex", action='store_true',
                    help="use the CPLEX solver")
parser.add_argument("-l", action='store_true', help="keep log files")
parser.add_argument("--log")
parser.add_argument("--ipm", action='store_true')

args = parser.parse_args()

file = open(args.input, "r")

U_lookup = {}
S_lookup = {}

U_counter = 1
S_counter = 1

U = []
S = {}

inc_map = {}
delta = 0

for line in file:
    line_spl = line.split(":")
    id = int(line_spl[0].strip())

    S_lookup.update({S_counter: id})

    line_str = line_spl[1].strip()
    e = [int(ep) for ep in line_str.split()]

    if len(e) > delta:
        delta = len(e)

    e_tr = []

    for v in e:
        if v not in U_lookup:
            U_lookup.update({v: U_counter})
            e_tr.append(U_counter)
            U_counter += 1
        else:
            e_tr.append(U_lookup[v])

        if U_lookup[v] not in inc_map:
            inc_map.update({U_lookup[v]: [S_counter]})
        else:
            inc_map[U_lookup[v]].append(S_counter)

    S.update({S_counter: e_tr})
    S_counter += 1

print("file loaded...")

n = U_counter-1
m = len(S)

U = [x for x in range(1, n+1)]

print("instance had {} elements and {} many sets".format(n, m))

prob = LpProblem("SC-Relax", LpMinimize)

x = LpVariable.dicts("x", range(1, m+1), 0, 1)
prob += lpSum([x[j] for j in range(1, m+1)])

for i in range(1, n+1):
    # at least one set containing i needs to be in the cover
    e_sum = lpSum([x[j] for j in inc_map[i]])
    prob += e_sum >= 1

print("begin solving...")

keep_logs = args.l
opts = []
if args.highs and args.ipm:
    opts = ["--solver", "ipm"]

if args.highs:
    prob.solve(HiGHS_CMD(mip=False, msg="using HiGHS", keepFiles=keep_logs, options=opts,
                         path="/usr/local/bin/highs", threads=os.cpu_count()))
elif args.glpk:
    prob.solve(GLPK(msg="using GLPK solver", keepFiles=keep_logs))
elif args.cplex:
    prob.solve(CPLEX_CMD(keepFiles=keep_logs))
else:
    prob.solve(PULP_CBC_CMD(keepFiles=keep_logs))

opt = 0.0
print("Status:", LpStatus[prob.status])
if prob.status == LpStatusOptimal:
    # print("Solution:")
    # for j in range(1, m+1):
    #    print(f"{S_lookup[j]} =", value(x[j]))
    opt = value(lpSum([x[j] for j in range(1, m+1)]))
    print("Sum of decision variables =", opt)

R_1 = []
R_2 = []
C = set()

k = max([len(inc) for _, inc in inc_map.items()])
alpha = 1 - exp(-1.0 * (log(delta)/float(k-1)))

print("k =", k)
print("delta =", delta)
print("alpha =", alpha)

for j in range(1, m+1):
    p_j = min([1.0, alpha*k*value(x[j])])
    r = random()
    if r <= p_j:
        covered = True
        for i in S[j]:
            if i not in C:
                covered = False
            C.add(i)
        if not covered:
            R_1.append(j)
        

# get elements not covered by sets
I_r = set()
for i in range(1, n+1):
    if i not in C:
        I_r.add(i)

for i in I_r:
    if i in C:
        continue
    for j in inc_map[i]:
        if j not in R_2:
            R_2.append(j)
            for v in S[j]:
                C.add(v)
            break

sc = [S_lookup[j] for j in R_1 + R_2]

ratio_ub = (k-1)*(1-exp(-1*(log(delta)/(k-1)))) + 1
ratio = len(sc)/opt

print("ratio upper bound:", ratio_ub)
print("actual ratio:", ratio)
print("found set-cover of size", len(sc))

if args.log != None:
    logfile = open(args.log, "a+")
    filename = str(args.input).split("/")[-1]
    logfile.write("{};{};{};{};{}\n".format(filename, ratio_ub, ratio, len(sc), opt))