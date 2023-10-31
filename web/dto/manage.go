package dto

type JobManagementDto struct {
	AppName       string `json:"appName" form:"appName" json:"appName" uri:"appName" xml:"appName" yaml:"appName" binding:"required"`
	Name          string `json:"name" form:"name" json:"name" uri:"name" xml:"name" yaml:"name" binding:"required"`
	Description   string `json:"description" form:"description" json:"description" uri:"description" xml:"description" yaml:"description"`
	RoutingPolicy int32  `json:"routingPolicy" form:"routingPolicy" json:"routingPolicy" uri:"routingPolicy" xml:"routingPolicy" yaml:"routingPolicy" binding:"required"`
}
