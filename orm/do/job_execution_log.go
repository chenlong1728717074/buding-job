package do

import "buding-job/orm"

func init() {
	orm.DB.AutoMigrate(&ExecutionLog{})
}

type ExecutionLog struct {
	Id          int64  `gorm:"primary_key;auto_increment:false"`
	ExecuteLogs string `gorm:"comment:执行日志"`
}

func NewExecutionLog(id int64, log string) *ExecutionLog {
	return &ExecutionLog{
		id, log,
	}
}

func (*ExecutionLog) TableName() string {
	return "tb_job_execution_log"
}
