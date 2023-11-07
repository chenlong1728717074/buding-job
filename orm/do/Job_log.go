package do

import (
	"buding-job/orm"
	"time"
)

func init() {
	orm.DB.AutoMigrate(&JobLogDo{})
}

type JobLogDo struct {
	orm.BaseModel
	ManageId             int64      `gorm:"column:manage_id;comment:管理器id"`
	JobId                int64      `gorm:"column:job_id;comment:任务id"`
	DispatchHandler      string     `gorm:"type:varchar(64);comment:调度handler"`
	DispatchTime         *time.Time `gorm:"column:dispatch_time;comment:调度时间"`
	DispatchAddress      string     `gorm:"column:dispatch_address;comment:调度地址;type:varchar(64)"`
	DispatchStatus       int64      `gorm:"comment:调度状态 1:调度成功 2:调度失败"`
	DispatchType         int64      `gorm:"comment:调度类型 1:自动 2:手动"`
	DispatchRemake       string     `gorm:"type:varchar(256);comment:调度说明"`
	ExecuteStartTime     *time.Time `gorm:"comment:执行开始时间"`
	ExecuteEndTime       *time.Time `gorm:"comment:执行结束时间"`
	ExecuteConsumingTime int64      `gorm:"comment:执行耗时"`
	ExecuteParams        string     `gorm:"type:varchar(256);comment:参数"`
	ExecuteStatus        int32      `gorm:"comment:执行状态 -1:未开始 0:串行 1:进行中 2:执行成功 3:执行出现异常 4:执行超时"`
	ExecuteRemark        string     `gorm:"type:varchar(128);"`
	Retry                int32      `gorm:"default:0;comment:重试次数"`
	ProcessingStatus     int32      `gorm:"default:0;comment:处理状态  1:无需处理/2:需要重试/3:告警成功/4:告警失败"`
}

func (*JobLogDo) TableName() string {
	return "tb_job_log"
}
