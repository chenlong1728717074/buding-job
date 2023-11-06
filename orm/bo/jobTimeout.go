package bo

import "buding-job/orm/do"

type JobTimeoutBo struct {
	do.JobLogDo
	Timeout int32
}
