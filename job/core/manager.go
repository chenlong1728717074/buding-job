package core

import "sync"

type JobManager struct {
	Id      int64
	AppName string
	Name    string
	Router  Router
	JobList map[int64]*Scheduler
	lock    sync.RWMutex
}

func NewJobManager(id int64, appName string, name string, strategy RouterStrategy) *JobManager {
	var manager = JobManager{
		Id:      id,
		AppName: appName,
		Name:    name,
		Router:  NewStrategyRouter(strategy),
		JobList: make(map[int64]*Scheduler),
		lock:    sync.RWMutex{},
	}
	return &manager
}

func (manager *JobManager) GetJobList() {
	manager.lock.RLock()
	defer manager.lock.RUnlock()
}

func (manager *JobManager) AddJob() {
	manager.lock.Lock()
	defer manager.lock.Unlock()
}

func (manager *JobManager) RemoveJob() {
	manager.lock.RLock()
	defer manager.lock.Unlock()
}

func (manager *JobManager) UpdateJob() {
	manager.lock.RLock()
	defer manager.lock.RUnlock()
}

func (manager *JobManager) Routing(strategy RouterStrategy) *Instance {
	commonRouter, ok := manager.Router.(*CommonRouter)
	if ok {
		return commonRouter.GetStrategyInstance(strategy)
	} else {
		return manager.Router.GetInstance()
	}
}

func (manager *JobManager) Permission() bool {
	return len(manager.Router.AllInstance()) > 0
}
