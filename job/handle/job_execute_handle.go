package handle

import (
	"buding-job/common/constant"
	"buding-job/job/core"
	"buding-job/job/grpc/to"
	"buding-job/orm"
	"buding-job/orm/do"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync"
	"time"
)

var JobExecute *jobExecuteHandle

func init() {
	JobExecute = NewJobExecuteHandle()
}

type jobExecuteHandle struct {
	lock sync.RWMutex
}

func NewJobExecuteHandle() *jobExecuteHandle {
	return &jobExecuteHandle{
		lock: sync.RWMutex{},
	}
}

// 执行逻辑
func (jobExecute *jobExecuteHandle) Execute(job *core.Scheduler, schedule bool) {
	//没有服务注册上去,不允许执行
	if job.Manager.Permission() {
		// 1:单机/2:广播
		lockDo, addLog, run := jobExecute.permission(job)
		if addLog {
			jobExecute.doExecute(job, lockDo, run)
		}
	}
	//刷新下次时间
	if schedule {
		job.FlushTime()
	}
}

// 执行任务
func (jobExecute *jobExecuteHandle) doExecute(job *core.Scheduler, lockDo *do.JobLockDo, run bool) {
	var jobLog *do.JobLogDo
	var flag bool
	if jobLog, flag = jobExecute.createLog(job); !flag {
		log.Println("log add fail")
		return
	}
	if !run {
		jobLog.ProcessingStatus = constant.Serial
		orm.DB.Updates(jobLog)
	}
	jobExecute.initExecuteLog(jobLog)
	//调度失败就删除锁
	if !jobExecute.dispatch(job, jobLog) {
		orm.DB.Delete(lockDo)
	}
}

// 远程调度控制
func (jobExecute *jobExecuteHandle) dispatch(job *core.Scheduler, logDo *do.JobLogDo) bool {
	instance := job.Manager.Routing(core.RouterStrategy(job.RoutingPolicy))
	flag := true
	err := jobExecute.doDispatch(job.JobHandle, job.Params, instance.Addr, logDo.Id)
	if err != nil && !jobExecute.failover(instance, job, logDo.Id) {
		//调度失败
		flag = false
		logDo.DispatchStatus = 2
		logDo.ExecuteStatus = -1
		logDo.Remark = fmt.Sprintf("任务调度失败%s", err.Error())
	}
	orm.DB.Updates(logDo)
	return flag
}

// 故障转移
func (jobExecute *jobExecuteHandle) failover(instance *core.Instance, job *core.Scheduler, logId int64) bool {
	var flag bool
	for _, newInstance := range job.Manager.Router.AllInstance() {
		if newInstance.Addr == instance.Addr {
			continue
		}
		if err := jobExecute.doDispatch(job.JobHandle, job.Params, newInstance.Addr, logId); err == nil {
			flag = true
			break
		}
	}
	return flag
}

// 执行调度
func (jobExecute *jobExecuteHandle) doDispatch(jobHandle string, param string, addr string, logId int64) error {
	dial, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer dial.Close()
	con := to.NewInstanceClient(dial)
	//dispatch
	_, callErr := con.Dispatch(context.Background(), &to.Request{
		JobHandler: jobHandle,
		CallbackId: logId,
		JobParams:  param,
	})
	return callErr
}

// 调度前初始化日志
func (jobExecute *jobExecuteHandle) initExecuteLog(jobLog *do.JobLogDo) {
	now := time.Now()
	jobLog.DispatchTime = &now
	// 0失败 1成功
	jobLog.DispatchStatus = 1
	jobLog.DispatchType = 1
	jobLog.ExecuteStatus = 1
}

// 是否允许执行(管理器必须有活跃的机器&必须是允许串行)第一个bool表示是否允许添加日志,第二个bool表示是否允许执行
func (jobExecute *jobExecuteHandle) permission(job *core.Scheduler) (*do.JobLockDo, bool, bool) {
	if !job.Manager.Permission() {
		return nil, false, false
	}
	//2串行
	if job.MisfireStrategy == 2 {
		return nil, true, true
	}
	now := time.Now()
	lock := &do.JobLockDo{
		Id:       job.Id,
		LockTime: &now,
	}
	tx := orm.DB.Create(lock)
	//加锁失败&&丢弃
	if tx.RowsAffected == 0 && job.MisfireStrategy == 1 {
		return nil, false, false
	}
	//加锁失败&&串行
	if job.MisfireStrategy == 3 && tx.RowsAffected == 0 {
		return nil, true, false
	}
	return lock, true, true
}

// 创建日志
func (jobExecute *jobExecuteHandle) createLog(job *core.Scheduler) (*do.JobLogDo, bool) {
	logDo := &do.JobLogDo{
		JobId:                job.Id,
		ManageId:             job.Manager.Id,
		DispatchHandler:      job.JobHandle,
		Retry:                0,
		ProcessingStatus:     constant.NoProcessingRequired,
		ExecuteConsumingTime: -1,
	}
	tx := orm.DB.Create(logDo)
	if tx.RowsAffected == 0 {
		return nil, false
	}
	return logDo, true
}
