#!/bin/fish

for i in (seq 10 25)
  go run cmd/tvd/tvd.go -n=1000 -evr=$i -log=100 -d=data/er_"$i"_1K
end