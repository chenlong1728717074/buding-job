package do

import "buding-job/orm"

func init() {
	orm.DB.AutoMigrate(&JobManagementDo{})
}

type JobManagementDo struct {
	orm.BaseModel
	AppName         string `gorm:"type:varchar(128);comment:app名称"`
	Name            string `gorm:"type:varchar(128);comment:执行器名称"`
	Description     string `gorm:"type:varchar(256);comment:描述"`
	RoutingPolicy   int32  `gorm:"default:1;comment:路由策略 1:指定第一个/2:指定随机/3:指定轮询/4:跟随任务本身"`
	BelongingServer string `gorm:"type:varchar(32);comment:归属服务,属于预留字段，为后期集群做准备"`
}

func (*JobManagementDo) TableName() string {
	return "tb_job_management"
}
