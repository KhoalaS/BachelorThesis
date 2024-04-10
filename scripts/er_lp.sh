#!/bin/fish

set solver $argv[1]
set graphdir $argv[2]
set solverarg --$solver

if test $solver = "clp"
    set solverarg
end

mkdir ./data/$graphdir

set l1 ./data/$graphdir/hs_$solver.csv
set l2 ./data/$graphdir/sc_$solver.csv

rm $l1
rm $l2

echo -e "File;RatioUB;Ratio;HittingSet;Opt" > $l1
echo -e "File;RatioUB;Ratio;HittingSet;Opt" > $l2

for file in (ls ./graphs/$graphdir)
    for i in (seq 1 10)
        python ilp/random_hs.py --log $l1 $solverarg graphs/$graphdir/$file > /dev/null
        #set p1 $last_pid
        python ilp/random_sc.py --log $l2 $solverarg graphs/$graphdir\_sc/$file > /dev/null
        #set p2 $last_pid

        #wait $p1
        #wait $p2
    end
end
