package handle

import (
	"buding-job/job/core"
	"log"
	"sync"
	"time"
)

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

func NewJobManagerHandle() {

}

func (h *JobManagerHandle) Start() {
	h.serverInspect()
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
		if h.instanceList[index].Equals(instance.Addr) {
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

func (h *JobManagerHandle) RemoveInstance() {
	h.instanceLock.Lock()
	defer h.instanceLock.RUnlock()

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
