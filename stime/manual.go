package stime

import (
	"strconv"
	"time"
)

var _units = []string{"ns", "us", "ms", "s"}

func ManualDuration(du time.Duration) string {
	d := float64(du)

	i := 0
	for d > 1000 && i < len(_units)-1 {
		d = d / 1000
		i++
	}

	return strconv.FormatFloat(d, 'f', 2, 64) + _units[i]
}
