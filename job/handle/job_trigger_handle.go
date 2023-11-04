package handle

import (
	"log"
	"sync"
	"time"
)

var JobSchedule *jobTriggerHandle

func init() {
	JobSchedule = NewJobTriggerHandle()
}

type jobTriggerHandle struct {
	lock      sync.RWMutex
	JobScan   chan interface{}
	FlushChan chan interface{}
}

func NewJobTriggerHandle() *jobTriggerHandle {
	return &jobTriggerHandle{
		lock:      sync.RWMutex{},
		JobScan:   make(chan interface{}),
		FlushChan: make(chan interface{}),
	}
}
func (job *jobTriggerHandle) Start() {
	//todo 获取数据
	job.start()
}
func (job *jobTriggerHandle) Stop() {
	job.JobScan <- struct{}{}
}

func (job *jobTriggerHandle) Flush() {
	job.FlushChan <- struct{}{}
}

func (job *jobTriggerHandle) start() {
	go func() {
		for {
			now := time.Now()
			var timer *time.Timer
			if JobManagerProcessor.Permission() {
				timer = time.NewTimer(100000 * time.Hour)
			} else {
				timer = time.NewTimer(JobManagerProcessor.jobList[0].NextTime.Sub(now))
			}
			for {
				select {
				case <-job.JobScan:
					timer.Stop()
					return
				case <-job.FlushChan:
					timer.Stop()
				case <-timer.C:
					start := time.Now()
					for _, c := range JobManagerProcessor.GetJobList() {
						if c.NextTime.After(start) {
							break
						}
						go JobExecute.Execute(c, true)
					}
					end := time.Now()
					consum := end.UnixMilli() - start.UnixMilli()
					log.Printf("本次执行耗时%dms", consum)
				}
				break
			}
		}

	}()
}
