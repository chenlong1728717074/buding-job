package do

import "buding-job/orm"

func init() {

}

type JobManagementDo struct {
	orm.BaseModel
	AppName string
	Name    string
}

func (JobManagementDo) TableName() string {
	return "tb_job_management"
}
