package handle

import (
	"buding-job/common/log"
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
	job.start()
}
func (job *jobTriggerHandle) Stop() {
	job.JobScan <- struct{}{}
}

func (job *jobTriggerHandle) Flush() {
	job.FlushChan <- struct{}{}
}

func (job *jobTriggerHandle) start() {
	job.trigger()
}

func (job *jobTriggerHandle) trigger() {
	go func() {
		//sleep 10s,wait work register
		time.Sleep(time.Second * 10)
		for {
			var timer *time.Timer
			//flush sort
			JobManagerProcessor.flushSchedulerSort()
			if JobManagerProcessor.Permission() {
				timer = time.NewTimer(100000 * time.Hour)
			} else {
				timer = time.NewTimer(JobManagerProcessor.GetSchedulerList()[0].NextTime.Sub(time.Now()))
			}
			for {
				select {
				case <-job.JobScan:
					timer.Stop()
					return
				case <-job.FlushChan:
					timer.Stop()
					break
				case <-timer.C:
					//XXX 优化后续这一步骤可以加入时间轮而非直接执行
					start := time.Now()
					list := JobManagerProcessor.GetSchedulerList()
					for _, c := range list {
						if c.NextTime.After(start) || c.NextTime.IsZero() {
							break
						}
						go JobExecute.Execute(c, true)
						//这一部必须是同步操作
						c.FlushTime()
					}
					end := time.Now()
					consum := end.UnixMilli() - start.UnixMilli()

					log.GetLog().Debugln("本次执行耗时", consum, "ms---->", "任务大小->", len(list))
					JobManagerProcessor.flushSchedulerSort()
					break
				}
				break
			}
		}

	}()
}
