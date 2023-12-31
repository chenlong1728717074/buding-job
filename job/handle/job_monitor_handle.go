package handle

import (
	"buding-job/common/constant"
	"buding-job/common/log"
	"buding-job/common/utils"
	"buding-job/job/alarm"
	"buding-job/orm"
	"buding-job/orm/bo"
	"buding-job/orm/do"
	"time"
)

var JobMonitor *JobMonitorHandle

func init() {
	JobMonitor = NewJobMonitorHandle()
}

// JobMonitorHandle  这个结构体用于任务失败监听/**/
type JobMonitorHandle struct {
	expireLockDone chan struct{}
	failJobDone    chan struct{}
	timeoutDone    chan struct{}
}

func NewJobMonitorHandle() *JobMonitorHandle {
	return &JobMonitorHandle{
		expireLockDone: make(chan struct{}),
		failJobDone:    make(chan struct{}),
		timeoutDone:    make(chan struct{}),
	}
}

func (monitor *JobMonitorHandle) Start() {
	//todo 扫描过期key是为了防止系统不稳定产生bug无法及时修复,正常允许不会有这样的问题,这一步可以在配置文件上选择是否开启
	monitor.expireLockScan()
	monitor.failJobScan()
	monitor.timeoutScan()
}
func (monitor *JobMonitorHandle) Stop() {
	monitor.expireLockDone <- struct{}{}
	monitor.failJobDone <- struct{}{}
	monitor.timeoutDone <- struct{}{}
}

func (monitor *JobMonitorHandle) expireLockScan() {
	go func() {
		//sleep 30s,wait  job execute
		time.Sleep(time.Second * 30)
		for {
			select {
			case <-monitor.expireLockDone:
				return
			default:
				monitor.deleteExpireLock()
				time.Sleep(time.Second * 15)
			}
		}
	}()
}
func (monitor *JobMonitorHandle) deleteExpireLock() {
	defer utils.Recover("执行出错,原因是:")
	tx := orm.DB.Model(&do.JobLockDo{}).
		Where("expire_time < ?", time.Now()).
		Delete(&do.JobLockDo{})
	if tx.Error != nil {
		panic(tx.Error)
	}
	log.GetLog().Infoln("expire key scan complete,delete total:", tx.RowsAffected)
}

func (monitor *JobMonitorHandle) failJobScan() {
	go func() {
		//sleep 30s,wait job execute
		time.Sleep(time.Second * 30)
		for {
			select {
			case <-monitor.timeoutDone:
				return
			default:
				monitor.failJob()
				time.Sleep(time.Second * 30)
			}
		}
	}()
}

func (monitor *JobMonitorHandle) timeoutScan() {
	go func() {
		//sleep 30s,wait job execute
		time.Sleep(time.Second * 30)
		for {
			select {
			case <-monitor.timeoutDone:
				return
			default:
				monitor.timeoutJob()
				time.Sleep(time.Second * 30)
			}
		}
	}()

}
func (monitor *JobMonitorHandle) timeoutJob() {
	defer utils.Recover("执行出错,原因是:")
	var jobLogs []bo.JobTimeoutBo
	orm.DB.Raw(constant.TimeoutJob).Scan(&jobLogs)
	if len(jobLogs) == 0 {
		return
	}
	for _, jobLogBo := range jobLogs {
		//已经删除或者已经关闭的任务不需要预警也不需要重试
		if jobLogBo.DeletedAt.Valid || !jobLogBo.Enable {
			monitor.lapseTimeoutJob(&jobLogBo.JobLogDo)
			continue
		}
		monitor.effectiveTimeoutJob(&jobLogBo)
	}
}

func (monitor *JobMonitorHandle) lapseTimeoutJob(jobLog *do.JobLogDo) {
	//Tasks that do not exist do not need to be retried
	now := time.Now()
	jobLog.ExecuteStartTime = jobLog.DispatchTime
	jobLog.ExecuteEndTime = &now
	jobLog.ExecuteConsumingTime = utils.ComputingTime(*jobLog.DispatchTime, now)
	jobLog.ExecuteStatus = constant.Timeout
	jobLog.ProcessingStatus = constant.NoProcessingRequired
	orm.DB.Updates(&jobLog)
}

func (monitor *JobMonitorHandle) effectiveTimeoutJob(jobLogBo *bo.JobTimeoutBo) {
	jobLog := &jobLogBo.JobLogDo
	now := time.Now()
	jobLog.ExecuteStartTime = jobLog.DispatchTime
	jobLog.ExecuteEndTime = &now
	jobLog.ExecuteConsumingTime = utils.ComputingTime(*jobLog.DispatchTime, now)
	jobLog.ExecuteStatus = constant.Timeout
	if jobLog.Retry < jobLogBo.Retry {
		jobLog.ProcessingStatus = constant.Retry
		jobLog.Retry = jobLog.Retry + 1
	} else {
		//告警
		monitor.alarm(jobLog, jobLogBo.Author, jobLogBo.Email)
	}
	orm.DB.Updates(jobLog)
}

func (monitor *JobMonitorHandle) alarm(jobLog *do.JobLogDo, author string, email string) {
	if email == "" {
		jobLog.ProcessingStatus = constant.NoProcessingRequired
	}
	//alarm
	jobLog.ProcessingStatus = constant.WarningFailed
	status := alarm.Mail.CommonAlarm(author, email, "", "")
	if status {
		jobLog.ProcessingStatus = constant.WarnedSuccess
	}
}

func (monitor *JobMonitorHandle) failJob() {
	var logs []do.JobLogDo
	orm.DB.Model(&do.JobLogDo{}).Where("processing_status=?", 2).Find(&logs)
	for _, value := range logs {
		JobExecute.ExecuteByLog(&value)
	}
}
