package constant

const (
	NotStarted         = -1
	InProgress         = 1
	ExecutionSucceeded = 2
	ExecutionFailed    = 3
	Timeout            = 4
)
const (
	DispatchSuccess = 1
	DispatchFailed  = 2
)

const (
	ManualTriggering    = 1
	AutomaticTriggering = 2
)

// 0 无需处理 1 已处理 2告警 3告警失败 4串行
const (
	NoProcessingRequired = 0
	Processed            = 1
	Warned               = 2
	WarningFailed        = 3
	Serial               = 4
)
