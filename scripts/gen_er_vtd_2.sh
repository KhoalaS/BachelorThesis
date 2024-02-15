#!/bin/bash

for i in {21..25}
do
  ./tvd -fr -n=1000 -evr=$i -log=100 -d=data/er_"$i"_1K
done
