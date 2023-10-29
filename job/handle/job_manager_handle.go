package handle

import (
	"buding-job/job/core"
	"log"
	"sync"
	"time"
)

var jobManagerHandle *JobManagerHandle

// JobManagerHandle 用于服务注册以及服务管理
type JobManagerHandle struct {
	jobManagerMap map[int64]*core.JobManager
	jobList       []*core.Scheduler
	flushDone     chan struct{}
	registerDone  chan struct{}
	instanceList  []*core.Instance
	instanceLock  sync.RWMutex
	jobLock       sync.RWMutex
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
				go h.flushServer()
				time.Sleep(time.Second * 30)
			}
		}
	}()
}
func (h *JobManagerHandle) flushServer() {

}

func (h *JobManagerHandle) addInstance() {
	//添加实例
	go func() {
		log.Printf("服务注册处理器已开启")
		for {
			select {}
		}
	}()
}

// 预留设定
type instanceManager struct {
}
