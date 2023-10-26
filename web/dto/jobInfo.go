package dto

type JobInfoDto struct {
	Id            int64  `json:"id" form:"id" json:"id" uri:"id" xml:"id" yaml:"id" `
	ManageId      int64  `json:"manageId" form:"manageId" json:"manageId" uri:"manageId" xml:"manageId" yaml:"manageId" binding:"required"`
	JobName       string `json:"jobName" form:"jobName" json:"jobName" uri:"jobName" xml:"jobName" yaml:"jobName" binding:"required"`
	JobHandler    string `json:"jobHandler" form:"jobHandler" json:"jobHandler" uri:"jobHandler" xml:"jobHandler" yaml:"jobHandler" binding:"required"`
	Cron          string `json:"core" form:"core" json:"core" uri:"core" xml:"core" yaml:"core" binding:"required"`
	Retry         int32  `json:"retry" form:"retry" json:"retry" uri:"retry" xml:"retry" yaml:"retry" binding:"required"`
	Timeout       int32  `json:"timeout" form:"timeout" json:"timeout" uri:"timeout" xml:"timeout" yaml:"timeout" binding:"required"`
	Author        string `json:"author" form:"author" json:"author" uri:"author" xml:"author" yaml:"author" binding:"required"`
	Email         string `json:"email" form:"email" json:"email" uri:"email" xml:"email" yaml:"email" binding:"required"`
	RoutingPolicy int32  `json:"routingPolicy" form:"routingPolicy" json:"routingPolicy" uri:"routingPolicy" xml:"routingPolicy" yaml:"routingPolicy" binding:"required"`
}
