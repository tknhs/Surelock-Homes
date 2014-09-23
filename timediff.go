package main

import (
	"strconv"
	"strings"
)

func TimeDiff(ts1, ts2 string) int {
	ts1List := strings.Split(ts1, "")
	ts2List := strings.Split(ts2, "")
	ts1 = strings.Join(ts1List[:10], "")
	ts2 = strings.Join(ts2List[:10], "")
	ts1Int, _ := strconv.Atoi(ts1)
	ts2Int, _ := strconv.Atoi(ts2)
	return ts2Int - ts1Int
}
