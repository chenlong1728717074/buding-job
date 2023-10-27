package core

import "time"

type ServiceNode struct {
	Addr           string
	JobManagerId   int64
	JobManagerName string
	RegisterTime   time.Time
}

func NewServiceNode(addr string, jobManagerId int64, jobManagerName string) *ServiceNode {
	return &ServiceNode{
		Addr:           addr,
		JobManagerId:   jobManagerId,
		JobManagerName: jobManagerName,
		RegisterTime:   time.Now(),
	}
}
