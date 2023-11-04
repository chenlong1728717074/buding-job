package core

import (
	"buding-job/orm/do"
	"time"
)

type Scheduler struct {
	Id              int64
	JobHandle       string
	Parser          TimeParser
	Params          string
	NextTime        time.Time
	Manager         *JobManager
	ExecuteType     int32
	RoutingPolicy   int32
	MisfireStrategy int32
	Retry           int32
}

func NewScheduler(info *do.JobInfoDo) *Scheduler {
	job := &Scheduler{
		Id: info.Id,
		//Cron:          do.Cron,
		JobHandle:       info.JobHandler,
		ExecuteType:     info.ExecuteType,
		Retry:           info.Retry,
		RoutingPolicy:   info.RoutingPolicy,
		MisfireStrategy: info.MisfireStrategy,
		Params:          info.JobParams,
	}
	//时间解析器
	job.setParser(info)
	job.NextTime = job.Next(time.Now())
	return job
}

func (scheduler *Scheduler) setParser(info *do.JobInfoDo) {
	if info.JobType == Cron {
		scheduler.Parser = NewCronTimeParser(info.Cron)
		return
	}
	scheduler.Parser = NewFixTimeParser(info.JobInterval)
}

func (scheduler *Scheduler) FlushTime() {
	scheduler.NextTime = scheduler.Parser.NextTime()
}
func (scheduler *Scheduler) Next(now time.Time) time.Time {
	return scheduler.Parser.Next(now)
}
