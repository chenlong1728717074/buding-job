package bo

import "xll-job/orm/do"

type RetryJobBo struct {
	do.JobLogDo
	Enable bool
}
