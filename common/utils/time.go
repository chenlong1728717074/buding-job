package utils

import "time"

func ComputingTime(startTime time.Time, endTime time.Time) int64 {
	consumingTime := (endTime.UnixNano() / 1000000) - (startTime.UnixNano() / 1000000)
	if consumingTime == 0 {
		consumingTime = 1
	}
	return consumingTime
}
