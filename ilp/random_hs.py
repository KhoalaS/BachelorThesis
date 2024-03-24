from pulp import *
import argparse

parser = argparse.ArgumentParser()
parser.add_argument("input", metavar="FILE", help="path to input graph file")

args = parser.parse_args()

file = open(args.input, "r")
V_lookup = {}
V_counter = 1
E_counter = 1
V = []
E = []

inc_map = {}

for line in file:
    if V_counter % 1000 == 0:
        print(V_counter)

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

    E.append(e_tr)


print("file loaded...")

n = V_counter-1
m = len(E)

V = [x for x in range(1, n+1)]
print("graph had {} vertices and {} many edges".format(n, m))

prob = LpProblem("VC-Relax", LpMinimize)

x = LpVariable.dicts("x", range(1, n+1), 0, 1)
prob += lpSum([x[j] for j in range(1, n+1)])


for idx, e in enumerate(E):
    e_sum = lpSum([x[j] for j in e])
    prob += e_sum >= 1

print("begin solving...")
prob.solve()


print("Status:", LpStatus[prob.status])
if prob.status == LpStatusOptimal:
    print("Solution:")
    print("Sum of decision variables =", value(
        lpSum([x[j] for j in V])))

C = set()
S_0 = set()
S_1 = set()
S_gte = set()
S_l = set()

_lambda = 0.5

for j in V:
    val = value(x[j])
    if val == 0:
        S_0.add(value(x[j]))
    elif val == 1:
        S_1.add(value(x[j]))
    elif val >= 1.0/_lambda:
        S_gte.add(value(x[j]))
    else:
        S_l.add(value(x[j]))

for j in S_0:
    V.remove(j)
    for e in E:
        if j in e:
            e.remove(j)
