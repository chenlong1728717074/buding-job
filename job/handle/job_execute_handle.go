package handle

import (
	"buding-job/common/constant"
	"buding-job/common/log"
	"buding-job/job/core"
	"buding-job/job/grpc/to"
	"buding-job/orm"
	"buding-job/orm/do"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

// Execute 执行逻辑
func (jobExecute *jobExecuteHandle) Execute(job *core.Scheduler, triggerType bool) {
	defer func() {
		if err := recover(); err != nil {
			log.GetLog().Errorln("执行出错,原因是:", err)
		}
	}()
	//没有服务注册上去,不允许执行
	if !job.Manager.Permission() {
		return
	}
	if job.ExecuteType == 1 {
		jobExecute.dispatchBroadcast(job, triggerType)
		return
	}
	jobExecute.dispatchClustering(job, triggerType)
}

// Execute 执行逻辑
func (jobExecute *jobExecuteHandle) ExecuteByLog(jobLog *do.JobLogDo) {
	defer func() {
		if err := recover(); err != nil {
			log.GetLog().Errorln("重试任务执行出错,原因是:", err)
		}
	}()
	var jobInfo do.JobInfoDo
	orm.DB.First(&jobInfo, jobLog.JobId)
	if jobInfo.Id == 0 || !jobInfo.Enable {
		jobLog.DispatchRemake = "任务关闭/删除,无需重试"
		jobLog.ProcessingStatus = constant.NoProcessingRequired
		orm.DB.Updates(jobLog)
		return
	}
	scheduler := JobManagerProcessor.GetScheduler(jobLog.JobId)
	if scheduler == nil || scheduler.Id == 0 {
		log.GetLog().Infoln("task not loaded[this is a serious error]")
		jobLog.DispatchRemake = "任务没有被正确加载,请联系开发者或者重启服务"
		jobLog.ProcessingStatus = constant.NoProcessingRequired
		orm.DB.Updates(jobLog)
		return
	}
	//没有服务注册上去,不允许执行
	if !scheduler.Manager.Permission() {
		jobLog.DispatchRemake = "该任务所在的任务管理器没有注册,无法进行重试"
		jobLog.Retry = jobLog.Retry + 1
		jobLog.ProcessingStatus = constant.Retry
		orm.DB.Updates(jobLog)
		return
	}
	//以上校验都通过,即可以进行重试
	jobLog.Retry = jobLog.Retry + 1
	if scheduler.ExecuteType == 1 {
		//lock, allow := jobExecute.permission(scheduler)
		//jobExecute.dispatchBroadcast(scheduler, false)
		return
	}
	jobExecute.dispatchClustering(scheduler, false)
}

// 集群模式下路由调用
func (jobExecute *jobExecuteHandle) dispatchBroadcast(job *core.Scheduler, triggerType bool) {
	lock, allow := jobExecute.permission(job)
	//加锁失败并且抛弃
	if !allow && job.MisfireStrategy == 1 {
		return
	}
	jobLog := jobExecute.createLog(job, triggerType)
	//串行的话需要添加日志进行重试
	if !allow && job.MisfireStrategy == 3 {
		jobLog.ProcessingStatus = constant.Serial
		orm.DB.Updates(jobLog)
		return
	}
	//路由
	var flag bool
	instance := job.Manager.Routing(core.RouterStrategy(job.RoutingPolicy))
	err := jobExecute.dispatch(job.JobHandle, job.Params, instance.Addr, jobLog.Id)
	if !(err == nil || jobExecute.failover(instance, job, jobLog)) {
		//调度失败
		flag = true
		jobLog.DispatchStatus = 2
		jobLog.ExecuteStatus = -1
		jobLog.ExecuteRemark = fmt.Sprintf("任务调度失败%s", err.Error())
	}
	orm.DB.Updates(jobLog)
	//调度失败并且不为并行
	if !flag && job.MisfireStrategy != 2 {
		lock.UnLock()
	}
}

func (jobExecute *jobExecuteHandle) dispatchClustering(job *core.Scheduler, triggerType bool) {
	instance := job.Manager.Router.AllInstance()
	for _, nowInstance := range instance {
		fmt.Println(nowInstance)
	}
}

// 故障转移
func (jobExecute *jobExecuteHandle) failover(instance *core.Instance, job *core.Scheduler, jobLog *do.JobLogDo) bool {
	var flag bool
	for _, nowInstance := range job.Manager.Router.AllInstance() {
		if nowInstance.Addr == instance.Addr {
			continue
		}
		if err := jobExecute.dispatch(job.JobHandle, job.Params, nowInstance.Addr, jobLog.GetId()); err == nil {
			flag = true
			jobLog.DispatchAddress = nowInstance.Addr
			break
		}
	}
	return flag
}

// 执行调度
func (jobExecute *jobExecuteHandle) dispatch(jobHandle string, param string, addr string, logId int64) error {
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

// 是否允许执行(管理器必须有活跃的机器&必须是允许串行)第一个bool表示是否允许添加日志,第二个bool表示是否允许执行
func (jobExecute *jobExecuteHandle) permission(job *core.Scheduler) (*do.JobLockDo, bool) {
	//2串行
	if job.MisfireStrategy == 2 {
		return nil, true
	}
	lock := do.NewJobLock(job.Id)
	//加锁失败&&丢弃
	isLock := lock.Lock(job.Timeout)
	return lock, isLock
}

// 创建日志
func (jobExecute *jobExecuteHandle) createLog(job *core.Scheduler, schedule bool) *do.JobLogDo {
	now := time.Now()
	logDo := &do.JobLogDo{
		JobId:                job.Id,
		ManageId:             job.Manager.Id,
		DispatchHandler:      job.JobHandle,
		DispatchTime:         &now,
		DispatchStatus:       1,
		DispatchType:         1,
		ExecuteStatus:        1,
		Retry:                0,
		ProcessingStatus:     constant.NoProcessingRequired,
		ExecuteConsumingTime: -1,
	}
	logDo.DispatchType = constant.ManualTriggering
	if schedule {
		logDo.DispatchType = constant.AutomaticTriggering
	}
	tx := orm.DB.Create(logDo)
	if tx.RowsAffected == 0 || tx.Error != nil {
		log.GetLog().Infoln("log insert fail")
		panic(tx.Error)
	}
	return logDo
}
