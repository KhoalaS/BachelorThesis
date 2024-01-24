#!/bin/bash

# dblp
python scripts/chart_single.py data/dblp -e scripts/.env.dblp
# pa
python scripts/chart.py data/pa_%_1K -s 2 -t 5 -m 2 -m 5 -e scripts/pa.env -o out/pa_02_05.html
# er vtd
python scripts/chart.py data/er_%_1K -s 10 -t 25 -m 10 -m 20 -e scripts/ervtd.env -o out/er_10_25_vtd.html
# er vtd
python scripts/chart.py data/eru3_%_1K -s 1 -t 20 -m 10 -m 20 -e scripts/eru3.env -o out/eru3_1_20.html
