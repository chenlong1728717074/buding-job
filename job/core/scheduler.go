package core

import (
	"buding-job/common/log"
	"buding-job/orm"
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
	Timeout         int32
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
		Timeout:         info.Timeout,
	}
	//时间解析器
	job.setParser(info)
	job.NextTime = job.Next(time.Now())
	return job
}

func (scheduler *Scheduler) setParser(info *do.JobInfoDo) {
	if info.JobTimeType == Cron {
		scheduler.Parser = NewCronTimeParser(info.Cron)
		return
	}
	scheduler.Parser = NewFixTimeParser(info.JobInterval)
}

func (scheduler *Scheduler) FlushTime() {
	nextTime := scheduler.Parser.NextTime()
	scheduler.NextTime = nextTime
	orm.DB.Model(&do.JobInfoDo{}).Where("id=?", scheduler.Id).Update(
		"next_time", &nextTime)
	log.GetLog().Debugln("更新任务下次执行时间:", nextTime)
}
func (scheduler *Scheduler) Next(now time.Time) time.Time {
	return scheduler.Parser.Next(now)
}
