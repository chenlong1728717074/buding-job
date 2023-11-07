package bo

import (
	"buding-job/orm/do"
	"gorm.io/gorm"
)

type JobTimeoutBo struct {
	do.JobLogDo
	Enable    bool           `gorm:"column:is_enable;"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_time"`
	Retry     int32
	Author    string
	Email     string
}
