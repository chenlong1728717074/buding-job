package core

import "sync"

type Router interface {
	GetInstance() *Instance
	AllInstance() []*Instance
	ReplaceInstance([]*Instance)
}
type CommonRouter struct {
	lock         sync.RWMutex
	instanceList []*Instance
}

func (*CommonRouter) GetInstance() *Instance {
	return nil
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

type FirstRouter struct {
	CommonRouter
}

func NewFirstRouter() *FirstRouter {
	return nil
}

func (router *FirstRouter) GetInstance() *Instance {
	if len(router.instanceList) == 0 {
		return nil
	}
	return router.instanceList[0]
}

type RandomRouter struct {
	next int32
	CommonRouter
}

func NewRandomRouter() *RandomRouter {
	return nil
}
func (router *RandomRouter) GetInstance() *Instance {
	if len(router.instanceList) == 0 {
		return nil
	}
	return nil
}

type PollingRouter struct {
	CommonRouter
}

func NewPollingRouter() *PollingRouter {
	return nil
}

func (router *PollingRouter) GetInstance() *Instance {
	if len(router.instanceList) == 0 {
		return nil
	}
	return nil
}
