package main

import (
	"strconv"
	"strings"
)

type courseTimes struct {
	Start int
	End   int
	Days  []string
}

func parseTimes(s string) ([]courseTimes, error) {
	content := strings.Split(s, ";")
	times := make([]courseTimes, 0, 4)
	for i := 0; i < len(content)/3; i++ {
		start, end, days := content[3*i], content[3*i+1], content[3*i+2]
		stime, err := strconv.Atoi(start)
		if err != nil {
			return nil, err
		}
		etime, err := strconv.Atoi(end)
		if err != nil {
			return nil, err
		}
		times = append(times, courseTimes{stime, etime, strings.Split(days, "")})
	}
	return times, nil
}
