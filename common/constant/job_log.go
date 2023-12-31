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

// 1:无需处理/2:需要重试/3:告警成功/4:告警失败
const (
	NoProcessingRequired = 1
	Retry                = 2
	WarnedSuccess        = 3
	WarningFailed        = 4
)
