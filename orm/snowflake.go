package orm

import (
	"math/rand"
	"sync"
	"time"
)

var SnowflakeGenerator *Snowflake

func init() {
	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)
	uid := generator.Intn(256)
	SnowflakeGenerator = NewSnowflake(uid)
}

// Snowflake 结构体用于生成唯一ID
type Snowflake struct {
	mu       sync.Mutex
	epoch    int64 // 起始时间戳，单位毫秒
	workerID int   // 当前worker ID
	lastTime int64 // 上一次生成ID的时间戳
	sequence int   // 当前序列号
}

const (
	workerBits      = 8                         // 扩展worker ID的位数，占用8位
	sequenceBits    = 4                         // 剩余位数为序列号，占用4位
	workerMax       = -1 ^ (-1 << workerBits)   // worker ID最大值，255
	sequenceMax     = -1 ^ (-1 << sequenceBits) // 序列号最大值，15
	timeShiftBits   = workerBits + sequenceBits // 时间戳左移位数，12
	workerShiftBits = sequenceBits              // worker ID左移位数，4
)

// NewSnowflake 创建一个 Snowflake 实例
// NewSnowflake 返回一个新的雪花算法对象
func NewSnowflake(workerID int) *Snowflake {
	return &Snowflake{
		epoch:    1630550400000, // 自定义起始时间戳，这里设置为2021-09-02 00:00:00的毫秒数
		workerID: workerID,
		lastTime: 0,
		sequence: 0,
	}
}

// NextID 生成下一个唯一ID
// Generate 生成一个新的ID
func (sf *Snowflake) Generate() int64 {
	sf.mu.Lock()
	defer sf.mu.Unlock()

	currTime := time.Now().UnixNano() / 1e6 // 当前时间戳，单位毫秒

	if currTime < sf.lastTime {
		// 如果当前时间小于上一次生成ID的时间戳，则可能是时钟回拨，暂停生成ID，直到时间追上
		time.Sleep(time.Millisecond * time.Duration(sf.lastTime-currTime))
		currTime = time.Now().UnixNano() / 1e6 // 重新获取当前时间戳
	}

	if currTime == sf.lastTime {
		sf.sequence = (sf.sequence + 1) & sequenceMax
		if sf.sequence == 0 {
			// 如果当前序列号达到最大值，等待下一毫秒
			currTime = sf.waitNextMilli()
		}
	} else {
		sf.sequence = 0
	}

	sf.lastTime = currTime
	id := ((currTime - sf.epoch) << timeShiftBits) | (int64(sf.workerID) << workerShiftBits) | int64(sf.sequence)
	return id
}

func (sf *Snowflake) waitNextMilli() int64 {
	currTime := time.Now().UnixNano() / 1e6
	for currTime <= sf.lastTime {
		currTime = time.Now().UnixNano() / 1e6
	}
	return currTime
}
