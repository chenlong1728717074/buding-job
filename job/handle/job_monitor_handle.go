package handle

import (
	"buding-job/common/constant"
	"buding-job/common/utils"
	"buding-job/job/grpc/to"
	"buding-job/orm"
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

func NewJobMonitorHandle() *JobMonitorHandle {
	return &JobMonitorHandle{
		failJobDone: make(chan struct{}),
		timeoutDone: make(chan struct{}),
	}
}

func (h *JobMonitorHandle) Callback(jobLog *do.JobLogDo, resp *to.CallbackResponse) {
	var job do.JobInfoDo
	orm.DB.First(&job, jobLog.JobId)
	//已经删除的任务不需要进行更新
	if job.Id == 0 {
		return
	}
	h.Unlock(job.Id)
	if jobLog.ExecuteStatus != constant.Timeout {
		jobLog.ExecuteStatus = resp.Status
	}
	startTime := resp.StartTime.AsTime()
	endTime := resp.EndTime.AsTime()
	jobLog.ExecuteStartTime = &startTime
	jobLog.ExecuteEndTime = &endTime
	jobLog.ExecuteConsumingTime = utils.ComputingTime(startTime, endTime)
}

func (h *JobMonitorHandle) Unlock(id int64) {
	if id == 0 {
		return
	}
	lock := do.NewJobLock(id)
	lock.UnLock()
}
