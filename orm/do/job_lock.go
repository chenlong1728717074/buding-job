package do

import (
	"buding-job/orm"
	"time"
)

func init() {
	orm.DB.AutoMigrate(&JobLockDo{})
}

type JobLockDo struct {
	Id       int64 `gorm:"primary_key;auto_increment:false"`
	LockTime *time.Time
}

func (JobLockDo) TableName() string {
	return "tb_job_lock"
}
