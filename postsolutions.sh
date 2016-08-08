#!/bin/sh

#curl --compressed -L -H Expect: -H 'X-API-Key: 101-ef09387a07b469087372e29dca268d27' -F 'problem_id=1' -F 'solution_spec=@work/solution.txt' 'http://2016sv.icfpcontest.org/api/solution/submit'
for f in solution-*.txt; do
    i=${f%.txt}
    i=${i#solution-}
    if [ ! -e "pushed-$i" ]; then
        echo "posting $i"
        sleep 1
        if curl --compressed -L -H Expect: -H 'X-API-Key: 101-ef09387a07b469087372e29dca268d27' -F "problem_id=$i" -F "solution_spec=@$f" 'http://2016sv.icfpcontest.org/api/solution/submit' ; then
            touch "pushed-$i"
            echo ""
        else
            echo "failed to push $i"
        fi
    fi
done
