package handle

import (
	"buding-job/job/core"
	"sync"
	"time"
)

var JobSchedule *jobScheduleHandle

func init() {
	JobSchedule = NewJobScheduleHandle()
}

type jobScheduleHandle struct {
	lock    sync.RWMutex
	jobList []*core.Scheduler

	JobScan chan interface{}
}

func NewJobScheduleHandle() *jobScheduleHandle {
	return &jobScheduleHandle{
		lock:    sync.RWMutex{},
		JobScan: make(chan interface{}),
	}
}
func (job *jobScheduleHandle) Start() {
	//todo 获取数据
	job.start()
}
func (job *jobScheduleHandle) Stop() {

}

func (job *jobScheduleHandle) start() {
	go func() {
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
					desiredSleepTime := 1000 - consum
					time.Sleep(time.Duration(desiredSleepTime) * time.Millisecond)
				}
			}
		}
	}()
}
func Execute(job *core.Scheduler, schedule bool) {
	if schedule {
		//todo 判断是否自动调度,如果属于调度,那么就修改数据
		job.NextTime = job.Next(time.Now())
	}
	go execute(job)
}

// todo 没有设定故障转移逻辑
func execute(job *core.Scheduler) {
	//没有服务注册上去,不允许执行
	if !job.Manager.Permission() {
		return
	}
	// 1:单机/2:广播
	//routing := job.Manager.Routing(core.RouterStrategy(job.RoutingPolicy))

}
