package core

import "sync"

type JobManager struct {
	Id         int64
	AppName    string
	Name       string
	ServerAddr []*ServiceNode
	JobList    map[int64]*Scheduler
	lock       sync.RWMutex
}

func NewJobManager(id int64, appName string, name string) *JobManager {
	manager := JobManager{
		Id:         id,
		AppName:    appName,
		Name:       name,
		ServerAddr: make([]*ServiceNode, 0), // 使用 make 创建新的切片
		JobList:    make(map[int64]*Scheduler),
		lock:       sync.RWMutex{},
	}
	return &manager
}

func (manager *JobManager) getJobList() {
	manager.lock.RLock()
	defer manager.lock.RUnlock()
}

func (manager *JobManager) addJob() {
	manager.lock.Lock()
	defer manager.lock.Unlock()
}

func (manager *JobManager) removeJob() {
	manager.lock.RLock()
	defer manager.lock.Unlock()
}

func (manager *JobManager) updateJob() {
	manager.lock.RLock()
	defer manager.lock.RUnlock()
}
