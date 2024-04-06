#!/bin/fish

set solver $argv[1]

set l1 ./data/rome_cvd_lphs_$solver.csv
set l2 ./data/rome_cvd_lpsc_$solver.csv
rm $l1
rm $l2
touch $l1
touch $l2

echo -e "File;RatioUB;Ratio;HittingSet;Opt" > ./data/rome_cvd_lphs.csv
echo -e "File;RatioUB;Ratio;HittingSet;Opt" > ./data/rome_cvd_lpsc.csv

for file in (ls ./graphs/rome_cvd)
    for i in (seq 1 10)
        python ilp/random_hs.py --log $l1 --$solver graphs/rome_cvd/$file > /dev/null &
        set pid1 $last_pid
        python ilp/random_sc.py --log $l2 --$solver graphs/rome_cvd_sc/$file > /dev/null &
        set pid2 $last_pid

        wait $pid1
        wait $pid2
    end
end
