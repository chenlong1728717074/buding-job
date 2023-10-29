package do

import (
	"buding-job/orm"
	"time"
)

func init() {
	orm.DB.AutoMigrate(&JobInfoDo{})
}

type JobInfoDo struct {
	orm.BaseModel
	ManageId        int64
	JobName         string     `gorm:"type:varchar(64);comment:工作名称"`
	JobDescription  string     `gorm:"type:varchar(256);comment:描述"`
	JobHandler      string     `gorm:"type:varchar(64);comment:handler"`
	JobParams       string     `gorm:"type:varchar(256);comment:参数"`
	JobType         int32      `gorm:"default:1;comment:工作类型 1:cron/2:固定时间"`
	Cron            string     `gorm:"type:varchar(64);comment:cron表达式"`
	JobInterval     int32      `gorm:"default:0;comment:间隔时间 单位:s"`
	NextTime        *time.Time `gorm:"comment:下次执行时间"`
	Timeout         int32      `gorm:"default:600;comment:超时时间,默认十分钟"`
	ExecuteType     int32      `gorm:"default:1;comment:调度类型 1:单机/2:广播"`
	RoutingPolicy   int32      `gorm:"default:1;comment:路由策略 1:第一个/2:随机/3:轮询"`
	MisfireStrategy int32      `gorm:"default:1;comment:过期策略 1:抛弃本次结果(调度成功后不记入调度结果)/2.并行/3.串行"`
	Retry           int32      `gorm:"default:0;comment:重试次数,默认0"`
	Enable          bool       `gorm:"column:is_enable;default:false;not null;comment:是否开启"`
	Author          string     `gorm:"type:varchar(64);comment:负责人"`
	Email           string     `gorm:"type:varchar(64);comment:负责人邮箱"`
}

func (JobInfoDo) TableName() string {
	return "tb_job_info"
}
