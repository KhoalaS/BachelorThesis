#!/bin/fish

set solver $argv[1]
set solverarg --$solver

if test $solver = "clp"
    set solverarg
end

set l1 ./data/rome_cvd_lphs_$solver.csv
set l2 ./data/rome_cvd_lpsc_$solver.csv

rm $l1
rm $l2

echo -e "File;RatioUB;Ratio;HittingSet;Opt" > $l1
echo -e "File;RatioUB;Ratio;HittingSet;Opt" > $l2

for file in (ls ./graphs/rome_cvd)
    for i in (seq 1 10)
        python ilp/random_hs.py --log $l1 $solverarg graphs/rome_cvd/$file > /dev/null &
        set pid1 $last_pid
        python ilp/random_sc.py --log $l2 $solverarg graphs/rome_cvd_sc/$file > /dev/null &
        set pid2 $last_pid

        wait $pid1
        wait $pid2
    end
end
