package core

type JobManager struct {
	Id         int64
	AppName    string
	Name       string
	ServerAddr []*ServiceNode
	JobList    map[int64]*Job
}

func NewJobManager(id int64, appName string, name string) *JobManager {
	manager := JobManager{
		Id:         id,
		AppName:    appName,
		Name:       name,
		ServerAddr: make([]*ServiceNode, 0), // 使用 make 创建新的切片
		JobList:    make(map[int64]*Job),
	}
	return &manager
}
