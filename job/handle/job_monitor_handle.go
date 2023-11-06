package handle

import (
	"buding-job/common/constant"
	"buding-job/common/utils"
	"buding-job/orm"
	"buding-job/orm/bo"
	"buding-job/orm/do"
	"log"
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
	log.Println("expire key scan complete,delete total:", tx.RowsAffected)
}

func (monitor *JobMonitorHandle) failJobScan() {

}

func (monitor *JobMonitorHandle) timeoutScan() {
	//先扫描删除的任务
	monitor.lapseTimeoutJobScan()
	//再扫描还在的任务
	monitor.effectiveTimeoutJobScan()
}

func (monitor *JobMonitorHandle) lapseTimeoutJobScan() {
	var jobLogs []*do.JobLogDo
	orm.DB.Raw(constant.LapseTimeoutJob).Scan(&jobLogs)
	if len(jobLogs) == 0 {
		return
	}
	for _, jobLog := range jobLogs {
		//Tasks that do not exist do not need to be retried
		jobLog.ExecuteStatus = constant.Timeout
		jobLog.ProcessingStatus = constant.NoProcessingRequired
		orm.DB.Updates(&jobLog)
	}
}

func (monitor *JobMonitorHandle) effectiveTimeoutJobScan() {
	var jobLogs []bo.JobTimeoutBo
	orm.DB.Raw(constant.EffectiveTimeoutJob).Scan(&jobLogs)
	if len(jobLogs) == 0 {
		return
	}
	now := time.Now()
	for _, jobLog := range jobLogs {
		dispatchTime := jobLog.DispatchTime
		orm.DB.Model(&do.JobLogDo{}).Where("id=?", jobLog.Id).Updates(map[string]interface{}{
			"execute_start_time":     dispatchTime,
			"execute_end_time":       &now,
			"execute_consuming_time": utils.ComputingTime(*dispatchTime, now),
			"execute_status":         constant.Timeout,
		})
	}
}
