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

func (instance *Instance) Equals(newInstance *Instance) bool {
	return instance.Addr == newInstance.Addr
}

func (instance *Instance) Lapse(lapseTime time.Time) bool {
	return instance.RegisterTime.After(lapseTime)
}

type ByTime []*Scheduler

func (s ByTime) Len() int      { return len(s) }
func (s ByTime) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ByTime) Less(i, j int) bool {
	// Two zero times should return false.
	// Otherwise, zero is "greater" than any other time.
	// (To sort it at the end of the list.)
	if s[i].NextTime.IsZero() {
		return false
	}
	if s[j].NextTime.IsZero() {
		return true
	}
	return s[i].NextTime.Before(s[j].NextTime)
}
