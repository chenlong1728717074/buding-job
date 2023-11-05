package dto

type JobInfoDto struct {
	Id              int64  `json:"id" form:"id" json:"id" uri:"id" xml:"id" yaml:"id" `
	ManageId        int64  `json:"manageId" form:"manageId" json:"manageId" uri:"manageId" xml:"manageId" yaml:"manageId" binding:"required"`
	JobName         string `json:"jobName" form:"jobName" json:"jobName" uri:"jobName" xml:"jobName" yaml:"jobName" binding:"required"`
	JobDescription  string `json:"jobDescription" form:"jobDescription" json:"jobDescription" uri:"jobDescription" xml:"jobDescription" yaml:"jobDescription"`
	JobHandler      string `json:"jobHandler" form:"jobHandler" json:"jobHandler" uri:"jobHandler" xml:"jobHandler" yaml:"jobHandler" binding:"required"`
	JobParams       string `json:"jobParams" form:"jobParams" json:"jobParams" uri:"jobParams" xml:"jobParams" yaml:"jobParams"`
	JobTimeType     int32  `json:"jobTimeType" form:"jobTimeType" json:"jobTimeType" uri:"jobTimeType" xml:"jobTimeType" yaml:"jobTimeType" binding:"required"`
	JobInterval     int64  `json:"jobInterval" form:"jobInterval" json:"jobInterval" uri:"jobInterval" xml:"jobInterval" yaml:"jobInterval"`
	Cron            string `json:"core" form:"core" json:"core" uri:"core" xml:"core" yaml:"core"`
	Retry           int32  `json:"retry" form:"retry" json:"retry" uri:"retry" xml:"retry" yaml:"retry" binding:"required"`
	Timeout         int32  `json:"timeout" form:"timeout" json:"timeout" uri:"timeout" xml:"timeout" yaml:"timeout" binding:"required"`
	Author          string `json:"author" form:"author" json:"author" uri:"author" xml:"author" yaml:"author" binding:"required"`
	Email           string `json:"email" form:"email" json:"email" uri:"email" xml:"email" yaml:"email" binding:"required"`
	RoutingPolicy   int32  `json:"routingPolicy" form:"routingPolicy" json:"routingPolicy" uri:"routingPolicy" xml:"routingPolicy" yaml:"routingPolicy"`
	MisfireStrategy int32  `json:"misfireStrategy" form:"misfireStrategy" json:"misfireStrategy" uri:"misfireStrategy" xml:"misfireStrategy" yaml:"misfireStrategy"`
}
