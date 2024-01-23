#!/bin/fish

for i in (seq 1 20)
    go run cmd/main.go -er -n=1000 -evr=$i -log=100 -d=data/eru3_"$i"_1K
end