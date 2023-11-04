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

var JobSchedule *jobScheduleHandle

func init() {
	JobSchedule = NewJobScheduleHandle()
}

type jobScheduleHandle struct {
	lock    sync.RWMutex
	JobScan chan interface{}
}

func NewJobScheduleHandle() *jobScheduleHandle {
	return &jobScheduleHandle{
		lock:    sync.RWMutex{},
		JobScan: make(chan interface{}),
	}
}
func (job *jobScheduleHandle) Start() {
	//todo 获取数据
	job.start()
}
func (job *jobScheduleHandle) Stop() {

}

func (job *jobScheduleHandle) start() {
	go func() {
		for {
			select {
			case <-job.JobScan:
				return
			default:
				start := time.Now()
				for _, c := range JobManagerProcessor.jobList {
					if c.NextTime.Before(start) {
						Execute(c, true)
					}
				}
				end := time.Now()
				consum := end.UnixMilli() - start.UnixMilli()
				if consum < 1000 {
					desiredSleepTime := 1000 - consum
					time.Sleep(time.Duration(desiredSleepTime) * time.Millisecond)
				}
			}
		}
	}()
}
func Execute(job *core.Scheduler, schedule bool) {
	if schedule {
		//todo 判断是否自动调度,如果属于调度,那么就修改数据
		job.NextTime = job.Next(time.Now())
	}
	go execute(job)
}

// 执行逻辑
func execute(job *core.Scheduler) {
	//没有服务注册上去,不允许执行
	if !job.Manager.Permission() {
		return
	}
	// 1:单机/2:广播
	lockDo, addLog, run := permission(job)
	if !addLog {
		return
	}
	doExecute(job, lockDo, run)
}

// 执行任务
func doExecute(job *core.Scheduler, lockDo *do.JobLockDo, run bool) {
	var jobLog *do.JobLogDo
	var flag bool
	if jobLog, flag = createLog(job); !flag {
		log.Println("log add fail")
		return
	}
	if !run {
		jobLog.ProcessingStatus = constant.Serial
		orm.DB.Updates(jobLog)
	}
	initExecuteLog(jobLog)
	//调度失败就删除锁
	if !dispatch(job, jobLog) {
		orm.DB.Delete(lockDo)
	}
	//刷新下次时间
	job.FlushTime()
}

// 远程调度控制
func dispatch(job *core.Scheduler, logDo *do.JobLogDo) bool {
	instance := job.Manager.Routing(core.RouterStrategy(job.RoutingPolicy))
	flag := true
	err := doDispatch(job.JobHandle, job.Params, instance.Addr, logDo.Id)
	if err != nil && !failover(instance, job, logDo.Id) {
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
func failover(instance *core.Instance, job *core.Scheduler, logId int64) bool {
	var flag bool
	for _, newInstance := range job.Manager.Router.AllInstance() {
		if newInstance.Addr == instance.Addr {
			continue
		}
		if err := doDispatch(job.JobHandle, job.Params, newInstance.Addr, logId); err == nil {
			flag = true
			break
		}
	}
	return flag
}

// 执行调度
func doDispatch(jobHandle string, param string, addr string, logId int64) error {
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
func initExecuteLog(jobLog *do.JobLogDo) {
	now := time.Now()
	jobLog.DispatchTime = &now
	// 0失败 1成功
	jobLog.DispatchStatus = 1
	jobLog.DispatchType = 1
	jobLog.ExecuteStatus = 1
}

// 是否允许执行(管理器必须有活跃的机器&必须是允许串行)第一个bool表示是否允许添加日志,第二个bool表示是否允许执行
func permission(job *core.Scheduler) (*do.JobLockDo, bool, bool) {
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
func createLog(job *core.Scheduler) (*do.JobLogDo, bool) {
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
