package utils

import "time"

func ComputingTime(startTime time.Time, endTime time.Time) int64 {
	consumingTime := endTime.UnixMilli() - startTime.UnixMilli()
	if consumingTime == 0 {
		consumingTime = 1
	}
	return consumingTime
}
