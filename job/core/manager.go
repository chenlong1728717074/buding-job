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

func NewJobManager(id int64, appName string, name string) *JobManager {
	var manager = JobManager{
		Id:      id,
		AppName: appName,
		Name:    name,
		Router:  NewFirstRouter(),
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
func (manager *JobManager) RouterInstance() {
	manager.Router.GetInstance()
}
