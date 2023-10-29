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

func (instance *Instance) FlushRegisterTime() {
	instance.RegisterTime = time.Now()
}

func (instance *Instance) Equals(addr string) bool {
	return instance.Addr == addr
}

func (instance *Instance) Lapse(lapseTime time.Time) bool {
	return instance.RegisterTime.After(lapseTime)
}
