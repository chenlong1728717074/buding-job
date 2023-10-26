package core

import (
	"buding-job/orm/do"
	"github.com/gorhill/cronexpr"
	"time"
)

type Job struct {
	Id            int64
	Cron          string
	JobHandle     string
	Expression    *cronexpr.Expression
	NextTime      time.Time
	Manager       string
	RoutingPolicy int32
	Retry         int32
}

func NewJob(do do.JobInfoDo) *Job {
	job := &Job{
		Id:            do.Id,
		Cron:          do.Cron,
		JobHandle:     do.JobHandler,
		Retry:         do.Retry,
		RoutingPolicy: do.RoutingPolicy,
	}
	parse, _ := cronexpr.Parse(do.Cron)
	job.Expression = parse
	job.NextTime = parse.Next(time.Now())
	//job.manager
	return job
}
