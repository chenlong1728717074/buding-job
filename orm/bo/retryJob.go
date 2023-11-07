package bo

import "buding-job/orm/do"

type RetryJobBo struct {
	do.JobLogDo
	Enable bool
}
