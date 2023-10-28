package handle

import (
	"buding-job/job/core"
	"sync"
	"time"
)

var jobSchedule *jobScheduleHandle

type jobScheduleHandle struct {
	lock    sync.RWMutex
	jobList []*core.Scheduler
	JobScan chan interface{}
}

func NewJobScheduleHandle() {
	jobSchedule = &jobScheduleHandle{
		lock:    sync.RWMutex{},
		JobScan: make(chan interface{}),
	}
}
func (job *jobScheduleHandle) Start() {
	//todo 获取数据
	job.jobList = make([]*core.Scheduler, 10)
	job.start()
}

func (job *jobScheduleHandle) start() {
	for {
		select {
		case <-job.JobScan:
			return
		default:
			start := time.Now()
			for _, c := range job.jobList {
				if c.NextTime.Before(start) {
					Execute(c, true)
				}
			}
			end := time.Now()
			consum := end.UnixMilli() - start.UnixMilli()

			if consum < 1000 {
				desiredSleepTime := 1000000 - consum
				time.Sleep(time.Duration(desiredSleepTime))
			}

		}
	}
}
func Execute(job *core.Scheduler, schedule bool) {
	if schedule {
		//todo 如果属于调度,那么就修改数据
	}
	go execute(job)
}
func execute(job *core.Scheduler) {

}
