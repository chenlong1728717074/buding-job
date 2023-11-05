package do

import (
	"buding-job/orm"
	"time"
)

func init() {
	orm.DB.AutoMigrate(&JobLockDo{})
}

type JobLockDo struct {
	Id         int64 `gorm:"primary_key;auto_increment:false"`
	ExpireTime time.Time
}

func NewJobLock(id int64) *JobLockDo {
	return &JobLockDo{
		Id: id,
	}
}

func (lock *JobLockDo) Lock(expire int32) bool {
	lock.ExpireTime = time.Now().Add(time.Second * time.Duration(expire))
	if tx := orm.DB.Create(lock); tx.RowsAffected == 0 || tx.Error != nil {
		lock.UnLock()
		return false
	}
	return true
}

func (lock *JobLockDo) UnLock() {
	orm.DB.Delete(lock)
}

func (*JobLockDo) TableName() string {
	return "tb_job_lock"
}
