package service

import (
	"buding-job/common/constant"
	"buding-job/common/utils"
	"buding-job/job/grpc/to"
	"buding-job/orm"
	"buding-job/orm/do"
	"context"
	"errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

type JobService struct {
	to.UnimplementedJobServer
}

func NewJobService() *JobService {
	return &JobService{}
}

func (job *JobService) Callback(ctx context.Context, resp *to.CallbackResponse) (*emptypb.Empty, error) {
	var jobLog do.JobLogDo
	orm.DB.First(&jobLog, resp.GetId())
	if jobLog.Id == 0 {
		return nil, errors.New("entry does not exist")
	}
	job.callback(&jobLog, resp)
	return &emptypb.Empty{}, nil
}

func (job *JobService) callback(jobLog *do.JobLogDo, resp *to.CallbackResponse) {
	//解锁
	job.Unlock(jobLog.JobId)
	var jobInfo do.JobInfoDo
	orm.DB.First(&jobInfo, jobLog.JobId)
	if jobLog.ExecuteStatus != constant.Timeout {
		jobLog.ExecuteStatus = resp.Status
	}
	if resp.Status == constant.ExecutionFailed {
		jobLog.ExecuteRemark = "execute fail"
	}
	if resp.Status == constant.ExecutionSucceeded {
		jobLog.ExecuteRemark = "execute success"
	}
	startTime := resp.StartTime.AsTime()
	endTime := resp.EndTime.AsTime()
	jobLog.ExecuteStartTime = &startTime
	jobLog.ExecuteEndTime = &endTime
	jobLog.ExecuteConsumingTime = utils.ComputingTime(startTime, endTime)
	if jobLog.Retry >= jobInfo.Retry {
		jobLog.ProcessingStatus = constant.Processed
	}
	orm.DB.Updates(jobLog)
	if len(resp.Logs) > 0 {
		job.saveExecuteLog(resp.GetId(), resp.Logs)
	}
}

func (job *JobService) Unlock(id int64) {
	if id == 0 {
		return
	}
	lock := do.NewJobLock(id)
	lock.UnLock()
}

func (job *JobService) saveExecuteLog(id int64, logs []string) {
	//日志处理
	executeLogs := "--------------以下是执行日志--------------\n"
	if len(logs) != 0 {
		for _, s := range logs {
			executeLogs += s + "\n"
		}
	}
	orm.DB.Save(do.NewExecutionLog(id, executeLogs))
}
