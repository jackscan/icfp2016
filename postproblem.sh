#!/bin/sh

curl --compressed -L -H Expect: -H 'X-API-Key: 101-ef09387a07b469087372e29dca268d27' -F "solution_spec=@$1" -F "publish_time=$2" 'http://2016sv.icfpcontest.org/api/problem/submit'
