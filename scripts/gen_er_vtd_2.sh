#!/bin/bash

for i in {21..25}
do
  go run cmd/tvd/tvd.go -n=1000 -evr=$i -log=100 -d=data/er_"$i"_1K
done
