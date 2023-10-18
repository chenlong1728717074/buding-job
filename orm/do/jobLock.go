package do

type JobLockDo struct {
	Id int64 `gorm:"primary_key;auto_increment:false"`
}

func (JobLockDo) TableName() string {
	return "tb_job_lock"
}
