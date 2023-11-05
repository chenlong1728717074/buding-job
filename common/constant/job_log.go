package constant

// -1:没有执行/1:进行中/2:执行成功/3:执行失败/4:超时/5:串行
const (
	NotStarted         = -1
	InProgress         = 1
	ExecutionSucceeded = 2
	ExecutionFailed    = 3
	Timeout            = 4
	Serial             = 5
)

// 1调度成功/2调度失败
const (
	DispatchSuccess = 1
	DispatchFailed  = 2
)

// 触发类型 1:手动/2:自动
const (
	ManualTriggering    = 1
	AutomaticTriggering = 2
)

// 0:无需处理/1:已处理/2:告警/3:告警失败
const (
	NoProcessingRequired = 0
	Processed            = 1
	Warned               = 2
	WarningFailed        = 3
)
