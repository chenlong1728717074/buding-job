package handle

import (
	"buding-job/common/constant"
	"buding-job/job/core"
	"buding-job/orm"
	"buding-job/orm/do"
	"log"
	"sync"
	"time"
)

func init() {
	JobManagerProcessor = NewJobManagerHandle()
}

var JobManagerProcessor *JobManagerHandle

// JobManagerHandle 用于服务注册以及服务管理
type JobManagerHandle struct {
	//管理器
	jobManagerMap map[int64]*core.JobManager
	//工作
	jobList []*core.Scheduler
	jobLock sync.RWMutex
	//实例
	flushDone    chan struct{}
	instanceList []*core.Instance
	registerDone chan struct{}       //暂时不用 解决高并发注册预留
	instanceChan chan *core.Instance //暂时不用 解决高并发注册预留
	instanceLock sync.RWMutex
}

func NewJobManagerHandle() *JobManagerHandle {
	return &JobManagerHandle{
		jobManagerMap: make(map[int64]*core.JobManager),
		jobList:       make([]*core.Scheduler, 0),
		jobLock:       sync.RWMutex{},
		flushDone:     make(chan struct{}),
		instanceList:  make([]*core.Instance, 0),
		registerDone:  make(chan struct{}),
		instanceChan:  make(chan *core.Instance, 100),
		instanceLock:  sync.RWMutex{},
	}
}

func (h *JobManagerHandle) Start() {
	h.init()
	h.serverInspect()
}

func (h *JobManagerHandle) init() {
	//这一步操作后期放监视器 定时检查锁
	if err := orm.DB.Exec(constant.DeleteLock).Error; err != nil {
		//log.Fatal("Failed to delete data: ", err)
	}
	//加载任务管理器
	var managers []do.JobManagementDo
	orm.DB.Model(&do.JobManagementDo{}).Find(&managers)
	//加载任务
	h.loadJob(managers)
	log.Printf("任务管理器加载成功,size=%d\n", len(h.jobManagerMap))
}

func (h *JobManagerHandle) loadJob(managers []do.JobManagementDo) {
	for _, managementDo := range managers {
		jobManager := core.NewJobManager(managementDo.Id, managementDo.Name, managementDo.AppName, core.RouterStrategy(managementDo.RoutingPolicy))
		h.jobManagerMap[managementDo.Id] = jobManager
		var jobs []do.JobInfoDo
		orm.DB.Model(&do.JobInfoDo{}).Where(&do.JobInfoDo{ManageId: jobManager.Id}).Find(&jobs)
		if len(jobs) == 0 {
			continue
		}
		for _, infoDo := range jobs {
			if infoDo.Enable {
				scheduler := core.NewScheduler(&infoDo)
				scheduler.Manager = jobManager
				h.addScheduler(scheduler)
			}
		}
	}
}

func (h *JobManagerHandle) addScheduler(scheduler *core.Scheduler) {
	h.jobLock.Lock()
	defer h.jobLock.Unlock()
	h.jobList = append(h.jobList, scheduler)
}

func (h *JobManagerHandle) serverInspect() {
	go func() {
		log.Println("服务检查处理器已开启")
		//睡十秒,等待服务注册
		time.Sleep(time.Second * 10)
		for {
			select {
			case <-h.flushDone:
				log.Println("服务检查处理器已关闭....")
				return
			default:
				go h.flush()
				time.Sleep(time.Second * 30)
			}
		}
	}()
}

func (h *JobManagerHandle) RegisterInstance(instance *core.Instance) {
	flag := true
	for index := range h.instanceList {
		if h.instanceList[index].Equals(instance) {
			h.instanceList[index].FlushRegisterTime()
			flag = false
			break
		}
	}
	if flag {
		h.addInstance(instance)
	}
	log.Printf("registration from[%s]has been refreshed\n", instance.Addr)
}

func (h *JobManagerHandle) addInstance(instance *core.Instance) {
	h.instanceLock.Lock()
	defer h.instanceLock.RUnlock()
	h.instanceList = append(h.instanceList, instance)
}

func (h *JobManagerHandle) RemoveInstance(instance *core.Instance) {
	h.instanceLock.Lock()
	defer h.instanceLock.Unlock()
	var index int
	var flag bool
	for i := range h.instanceList {
		if h.instanceList[i].Equals(instance) {
			flag = true
			index = i
			break
		}
	}
	if flag {
		h.instanceList = append(h.instanceList[:index], h.instanceList[index+1:]...)
	}
}

func (h *JobManagerHandle) flush() {
	h.instanceLock.RLock()
	defer h.instanceLock.RUnlock()
	startTime := time.Now().UnixNano() / 1000000
	log.Printf("start scrubbing service node[刷新服务]:%d\n", startTime)
	now := time.Now().Add(-time.Second * 90)
	newInstanceList := make([]*core.Instance, 0)
	//获取所有存活的服务
	for _, node := range h.instanceList {
		if now.After(node.RegisterTime) {
			node.RegisterTime = time.Now()
			newInstanceList = append(newInstanceList, node)
		}
	}
	//刷新缓存中的服务
	h.instanceList = newInstanceList
	// 分组ServiceNodeList中的节点
	temp := make(map[int64][]*core.Instance)
	for _, node := range h.instanceList {
		temp[node.JobManagerId] = append(temp[node.JobManagerId], node)
	}
	//重新分配服务
	for k := range h.jobManagerMap {
		manager := h.jobManagerMap[k]
		manager.Router.ReplaceInstance(temp[k])
	}
	endTime := time.Now().UnixNano() / 1000000
	log.Printf("service node refresh completed[刷新完成]:%d,time consuming:%d", endTime, endTime-startTime)
}

// 预留设定
type instanceManager struct {
}
