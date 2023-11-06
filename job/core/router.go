package core

import (
	"buding-job/common/utils"
	"sync"
)

type RouterStrategy int32

const (
	First   RouterStrategy = 1
	Random  RouterStrategy = 2
	Polling RouterStrategy = 3
	Common  RouterStrategy = 4
)

type Router interface {
	GetInstance() *Instance
	AllInstance() []*Instance
	ReplaceInstance([]*Instance)
}

// abstractRouter 通用路由器 用于实现Router接口
type abstractRouter struct {
	lock         sync.RWMutex
	instanceList []*Instance
}

func (*abstractRouter) GetInstance() *Instance {
	return nil
}
func (router *abstractRouter) AllInstance() []*Instance {
	router.lock.RLock()
	defer router.lock.RUnlock()
	return router.instanceList
}
func (router *abstractRouter) ReplaceInstance(newList []*Instance) {
	router.lock.Lock()
	defer router.lock.Unlock()
	router.instanceList = newList
}

// CommonRouter 通用的
type CommonRouter struct {
	PollingRouter
}

func NewCommonRouter() *CommonRouter {
	router := &CommonRouter{}
	router.lock = sync.RWMutex{}
	router.instanceList = make([]*Instance, 0)
	router.next = 0
	return router
}

func (router *CommonRouter) GetInstance() *Instance {
	return router.PollingRouter.GetInstance()
}
func (router *CommonRouter) GetStrategyInstance(strategy RouterStrategy) *Instance {
	if len(router.instanceList) == 0 {
		return nil
	}
	var result *Instance
	switch strategy {
	case First:
		result = router.instanceList[0]
		break
	case Random:
		result = router.instanceList[1]
		break
	default:
		result = router.PollingRouter.GetInstance()
		break
	}
	return result
}
func (router *CommonRouter) AllInstance() []*Instance {
	router.lock.RLock()
	defer router.lock.RUnlock()
	return router.instanceList
}
func (router *CommonRouter) ReplaceInstance(newList []*Instance) {
	router.lock.Lock()
	defer router.lock.Unlock()
	router.instanceList = newList
}

func NewStrategyRouter(strategy RouterStrategy) Router {
	var router Router
	switch strategy {
	case First:
		router = NewFirstRouter()
		break
	case Random:
		router = NewRandomRouter()
		break
	case Polling:
		router = NewPollingRouter()
		break
	default:
		router = NewCommonRouter()
	}
	return router
}

// FirstRouter 第一个实例路由器
type FirstRouter struct {
	abstractRouter
}

func NewFirstRouter() *FirstRouter {
	result := &FirstRouter{}
	result.lock = sync.RWMutex{}
	result.instanceList = make([]*Instance, 0)
	return result
}

func (router *FirstRouter) GetInstance() *Instance {
	if len(router.instanceList) == 0 {
		return nil
	}
	return router.instanceList[0]
}

// RandomRouter 随机路由器
type RandomRouter struct {
	abstractRouter
}

func NewRandomRouter() *RandomRouter {
	result := &RandomRouter{}
	result.lock = sync.RWMutex{}
	result.instanceList = make([]*Instance, 0)
	return result
}
func (router *RandomRouter) GetInstance() *Instance {
	if len(router.instanceList) == 0 {
		return nil
	}
	index := utils.RandI64(len(router.instanceList))
	return router.instanceList[index]
}

// PollingRouter 轮询路由器
type PollingRouter struct {
	next int
	abstractRouter
}

func NewPollingRouter() *PollingRouter {
	router := &PollingRouter{}
	router.lock = sync.RWMutex{}
	router.instanceList = make([]*Instance, 0)
	router.next = 0
	return router
}

func (router *PollingRouter) GetInstance() *Instance {
	if len(router.instanceList) == 0 {
		return nil
	}
	router.lock.Lock()
	defer router.lock.Unlock()
	next := router.next
	router.next = 0
	if !((next + 1) >= len(router.instanceList)) {
		router.next = next + 1
	}
	return router.instanceList[next]
}
