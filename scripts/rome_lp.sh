#!/bin/fish

rm ./data/rome_cvd_lphs.csv
rm ./data/rome_cvd_lpsc.csv

touch ./data/rome_cvd_lphs.csv
touch ./data/rome_cvd_lpsc.csv

echo -e "File;RatioUB;Ratio;HittingSet;Opt" > ./data/rome_cvd_lphs.csv
echo -e "File;RatioUB;Ratio;HittingSet;Opt" > ./data/rome_cvd_lpsc.csv

for file in (ls ./graphs/rome_cvd)
    for i in (seq 1 10)
        python ilp/random_hs.py --highs --log graphs/rome_cvd/$file > /dev/null &
        pid1=$!
        python ilp/random_sc.py --highs --log graphs/rome_cvd_sc/$file > /dev/null &
        pid2=$!

        wait $pid1
        wait $pid2
    end
end
