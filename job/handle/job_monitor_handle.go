package handle

import (
	"buding-job/job/grpc/to"
	"buding-job/orm/do"
)

var JobMonitor *JobMonitorHandle

func init() {
	JobMonitor = NewJobMonitorHandle()
}

// JobMonitorHandle  这个结构体用于任务失败监听/**/
type JobMonitorHandle struct {
	failJobDone chan struct{}
	timeoutDone chan struct{}
}

func (h JobMonitorHandle) Callback(jobLog *do.JobLogDo, resp *to.CallbackResponse) {
}

func NewJobMonitorHandle() *JobMonitorHandle {
	return &JobMonitorHandle{
		failJobDone: make(chan struct{}),
		timeoutDone: make(chan struct{}),
	}
}
