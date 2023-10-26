package do

import "buding-job/orm"

type JobInfoDo struct {
	orm.BaseModel
	ManageId      int64
	Retry         int32  `gorm:"default:0;comment:重试次数,默认0"`
	JobName       string `gorm:"type:varchar(64);comment:工作名称"`
	JobHandler    string `gorm:"type:varchar(64);comment:handler"`
	Cron          string `gorm:"type:varchar(64);comment:cron表达式"`
	Timeout       int32  `gorm:"default:600;comment:超时时间,默认十分钟"`
	Author        string `gorm:"type:varchar(64);comment:负责人"`
	Email         string `gorm:"type:varchar(64);comment:负责人邮箱"`
	RoutingPolicy int32  `gorm:"default:1;comment:路由策略 1:抛弃本次 2.并行 3.串行"`
	Enable        bool   `gorm:"column:is_enable;default:false;not null;comment:是否开启"`
}

func (JobInfoDo) TableName() string {
	return "tb_job_info"
}
