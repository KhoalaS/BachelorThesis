#!/bin/bash

# dblp
python scripts/chart_single.py data/dblp -e scripts/.env.dblp
# pa
python scripts/chart.py data/pa_%_1K -s 2 -t 5 -m 2 -m 5 -e scripts/.env.pa -o out/pa_02_05.html
# er vtd
python scripts/chart.py data/er_%_1K -s 10 -t 25 -m 10 -m 20 -e scripts/.env.ervtd -o out/er_10_25_vtd.html
