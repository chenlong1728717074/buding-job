package handle

import (
	"buding-job/common/constant"
	"buding-job/common/log"
	"buding-job/job/core"
	"encoding/json"
	transport "github.com/Jille/raft-grpc-transport"
	"github.com/hashicorp/raft"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"sync"
)

var JobFsm *JobFsmHandle

func init() {
	JobFsm = NewJobFsmHandle()
}

type JobFsmHandle struct {
	lock       sync.RWMutex
	kvStore    map[string]string
	raftNode   *raft.Raft
	status     int
	statusLock sync.RWMutex
}

func NewJobFsmHandle() *JobFsmHandle {
	return &JobFsmHandle{
		kvStore: make(map[string]string),
		lock:    sync.RWMutex{},
	}
}

func NodeRaftNode(config *raft.Config, fsm *JobFsmHandle, server *grpc.Server) (*raft.Raft, error) {
	// 使用内存存储
	logStore := raft.NewInmemStore()
	stableStore := raft.NewInmemStore()
	snapshotStore := raft.NewInmemSnapshotStore()
	// 配置传输
	manager := transport.New("localhost:8082", []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	manager.Register(server)
	// 创建 Raft 节点
	raftNode, err := raft.NewRaft(config, fsm, logStore, stableStore, snapshotStore, manager.Transport())
	if err != nil {
		return nil, err
	}
	return raftNode, nil
}

func (fsm *JobFsmHandle) GetKv() map[string]string {
	fsm.lock.RLock()
	fsm.lock.Unlock()
	return fsm.kvStore
}

func (fsm *JobFsmHandle) InitRaftNode(server *grpc.Server) {
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID("localhost:8082") // 设置本地节点的唯一标识
	config.Logger = core.NewRaftLog()
	node, err := NodeRaftNode(config, fsm, server)
	if err != nil {
		panic(err)
	}
	fsm.raftNode = node
}

func (fsm *JobFsmHandle) GetStatus() int {
	fsm.statusLock.RLock()
	fsm.statusLock.Unlock()
	return fsm.status
}
func (fsm *JobFsmHandle) SetStatus(status int) {
	fsm.statusLock.Lock()
	fsm.statusLock.Unlock()
	fsm.status = status
}

func (fsm *JobFsmHandle) Start() {
	go func() {
		if err := fsm.raftNode.BootstrapCluster(raft.Configuration{Servers: []raft.Server{{ID: "localhost:8082", Address: "localhost:8082"}}}).Error(); err != nil {
			log.GetLog().Fatalf("Error bootstrapping Raft cluster: %v\n", err)
		} else {
			fsm.status = 1
		}
	}()
}

// Apply 接受到同步命令
func (fsm *JobFsmHandle) Apply(logEntry *raft.Log) interface{} {
	fsm.lock.Lock()
	fsm.lock.Unlock()
	cmd := &core.Command{}
	if err := json.Unmarshal(logEntry.Data, cmd); err != nil {
		log.GetLog().Errorln("Failed to unmarshal command: %v\n", err)
	}
	if cmd.Cmd == constant.Put {
		fsm.kvStore[cmd.Key] = cmd.Value
		return nil
	}
	if cmd.Cmd == constant.Delete {
		fsm.kvStore[cmd.Key] = cmd.Value
		return nil
	}
	return nil
}

// Snapshot 目前项目没有快照需求
func (fsm *JobFsmHandle) Snapshot() (raft.FSMSnapshot, error) {
	// 首先对当前状态进行深拷贝以避免在序列化过程中的并发修改问题
	//fsmStateCopy := make(map[string]string)
	//for k, v := range fsm.kvStore {
	//	fsmStateCopy[k] = v
	//}
	return nil, nil
}

func (fsm *JobFsmHandle) Restore(serializedSnapshot io.ReadCloser) error {
	// 在这里恢复快照
	return nil
}
