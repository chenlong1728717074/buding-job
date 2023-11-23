package service

import (
	"buding-job/common/constant"
	"buding-job/common/log"
	"buding-job/common/utils"
	"buding-job/job/alarm"
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
	defer func() {
		if err := recover(); err != nil {
			log.GetLog().Errorln("操作出现异常:", err)
		}
	}()
	//解锁
	job.Unlock(jobLog.JobId)
	var jobInfo do.JobInfoDo
	orm.DB.First(&jobInfo, jobLog.JobId)
	//设置无需处理,即使时任务超时,只要进行了回调,那就算有结果,所以系统判断的超时不需要重试
	startTime := resp.StartTime.AsTime()
	endTime := resp.EndTime.AsTime()
	jobLog.ExecuteStartTime = &startTime
	jobLog.ExecuteEndTime = &endTime
	executeTime := utils.ComputingTime(startTime, endTime)
	jobLog.ExecuteConsumingTime = executeTime
	//处理状态
	jobLog.ProcessingStatus = constant.NoProcessingRequired
	jobLog.ExecuteRemark = "execute success"
	//超时
	if jobLog.ExecuteStatus != constant.Timeout {
		jobLog.ExecuteStatus = resp.Status
	}
	//执行失败
	if resp.Status == constant.ExecutionFailed {
		job.failJob(jobLog, jobInfo)
	}
	//执行成功,但是超时
	if resp.Status == constant.ExecutionSucceeded && executeTime > int64(jobInfo.Timeout*1000) {
		jobLog.ProcessingStatus = constant.Timeout
		jobLog.ExecuteRemark = "timeout,execute success"
	}

	orm.DB.Updates(jobLog)
	if len(resp.Logs) > 0 {
		job.saveExecuteLog(resp.GetId(), resp.Logs)
	}
}

func (job *JobService) failJob(jobLog *do.JobLogDo, jobInfo do.JobInfoDo) {
	jobLog.ExecuteRemark = "execute fail"
	jobLog.ProcessingStatus = constant.NoProcessingRequired
	//删除了就不需要处理了
	if jobInfo.Id == 0 {
		return
	}
	//开启就需要重试,关闭了就用上面的无需处理
	if jobInfo.Enable {
		jobLog.ProcessingStatus = constant.Retry
	}
	//只要任务存在不管是开启状态还是关闭都要告警
	if jobLog.Retry >= jobInfo.Retry && jobLog.ProcessingStatus != constant.WarnedSuccess &&
		jobLog.ProcessingStatus != constant.WarningFailed {
		jobLog.ProcessingStatus = constant.WarningFailed
		status := alarm.Mail.CommonAlarm(jobInfo.Author, jobInfo.Email, "", jobInfo.JobName)
		if status {
			jobLog.ProcessingStatus = constant.WarnedSuccess
		}
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
