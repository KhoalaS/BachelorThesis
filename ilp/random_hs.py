from pulp import *
import argparse
from random import random

parser = argparse.ArgumentParser()
parser.add_argument("input", metavar="FILE", help="path to input graph file")
parser.add_argument("--highs", action='store_true',
                    help="use the HiGHS solver")
parser.add_argument("--glpk", action='store_true',
                    help="use the GLPK solver")
parser.add_argument("--cplex", action='store_true',
                    help="use the CPLEX solver")
parser.add_argument("-l", action='store_true', help="keep log files")
parser.add_argument("--log")
parser.add_argument("--ipm", action='store_true')
parser.add_argument("--sol", action='store_true')


args = parser.parse_args()

file = open(args.input, "r")
V_lookup = {}
V_lookup_inv = {}
V_counter = 1
E_counter = 1
V = []
E = {}

inc_map = {}

for line in file:
    line_str = line.strip()
    e = [int(ep) for ep in line_str.split()]
    e_tr = []

    for v in e:
        if v not in V_lookup:
            V_lookup.update({v: V_counter})
            V_lookup_inv.update({V_counter: v})
            e_tr.append(V_counter)
            V_counter += 1
        else:
            e_tr.append(V_lookup.get(v))

        if V_lookup[v] not in inc_map:
            inc_map.update({V_lookup[v]: [E_counter]})
        else:
            inc_map[V_lookup[v]].append(E_counter)

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
opts = []
if args.highs and args.ipm:
    opts = ["--solver", "ipm"]

if args.highs:
    prob.solve(HiGHS_CMD(msg="using HiGHS", keepFiles=keep_logs, options=opts,
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
    print("Solution:")
    opt = value(lpSum([x[j] for j in V]))
    print("Sum of decision variables =", opt)

if args.sol:
    buffer = []
    counter = 0
    out = open("VC-Relax.sol", "w+")

    for v in V:
        out.write("{} {}\n".format(V_lookup_inv[v], value(x[v])))

    out.close()
    sys.exit(0)

C = set()
S_0 = set()
S_1 = set()
S_gte = set()
S_l = set()

for j in V:
    val = value(x[j])
    if val == 0:
        S_0.add(j)
    elif val == 1:
        S_1.add(j)

l = max([len(e) for _, e in E.items()])
e = (l * (opt-len(S_1)))/(2.0 * m)
delta = max([len(inc) for _, inc in inc_map.items()])
_lambda = l*(1.0-e)

print("l =", l)
print("e =", e)
print("delta =", delta)
print("lambda =", _lambda)

for j in V:
    val = value(x[j])
    if val != 1 and val >= 1.0/_lambda:
        S_gte.add(j)
    elif val < 1.0/_lambda and val != 0:
        S_l.add(j)

print("|S_0| =", len(S_0))
print("|S_1| =", len(S_1))
print("|S_≥| =", len(S_gte))
print("|S_<| =", len(S_l))

# print("skip step 4 of algorithm, not removing vertices in S_0")
for j in S_0:
    V.remove(j)
    for e in inc_map[j]:
        E[e].remove(j)
    inc_map[j] = []

for j in S_1:
    if len(inc_map[j]) == 0:
        continue
    C.add(j)
    V.remove(j)
    rem_e = []
    for e in inc_map[j]:
        rem_e.append(e)

    for e in rem_e:
        for v in E[e]:
            inc_map[v].remove(e)
        del E[e]

for j in S_gte:
    if len(inc_map[j]) == 0:
        continue
    C.add(j)
    V.remove(j)
    rem_e = []
    for e in inc_map[j]:
        rem_e.append(e)

    for e in rem_e:
        for v in E[e]:
            inc_map[v].remove(e)
        del E[e]

for j in S_l:
    if len(inc_map[j]) == 0:
        continue
    p = _lambda*value(x[j])
    r = random()
    if r <= p:
        C.add(j)
        V.remove(j)
        rem_e = []
        for e in inc_map[j]:
            rem_e.append(e)

        for e in rem_e:
            for v in E[e]:
                inc_map[v].remove(e)
            del E[e]

if len(E) == 0:
    pass
else:
    while len(E) > 0:
        for _, e in E.items():
            rem = e[0]

            C.add(rem)
            V.remove(rem)

            rem_e = []
            for h in inc_map[rem]:
                rem_e.append(h)

            for h in rem_e:
                for v in E[h]:
                    inc_map[v].remove(h)
                del E[h]

            break

ratio_ub = l*(1-((l-1)/(8*delta)))
ratio = len(C)/opt

print("ratio upper bound:", ratio_ub)
print("actual ratio:", ratio)
print("found hitting-set of size", len(C))

if args.log != None:
    logfile = open(args.log, "a+")
    filename = str(args.input).split("/")[-1]
    logfile.write("{};{};{};{};{}\n".format(
        filename, ratio_ub, ratio, len(C), opt))
