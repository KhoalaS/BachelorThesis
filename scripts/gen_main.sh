#!/bin/bash

for i in {1..15}
do
    ./main -er -n=1000 -evr=$i -log=100 -d=data/eru3_"$i"_1K
done