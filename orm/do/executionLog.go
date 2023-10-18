package do

type ExecutionLog struct {
	Id          int64  `gorm:"primary_key;auto_increment:false"`
	ExecuteLogs string `gorm:"comment:执行日志"`
}

func (ExecutionLog) TableName() string {
	return "tb_job_execution_log"
}
