package core

import "time"

type Instance struct {
	Addr           string
	JobManagerId   int64
	JobManagerName string
	RegisterTime   time.Time
}

func NewInstance(addr string, jobManagerId int64, jobManagerName string) *Instance {
	return &Instance{
		Addr:           addr,
		JobManagerId:   jobManagerId,
		JobManagerName: jobManagerName,
		RegisterTime:   time.Now(),
	}
}
