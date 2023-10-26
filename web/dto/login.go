package dto

type LoginDto struct {
	UserName string `json:"userName" form:"userName" json:"userName" uri:"userName" xml:"userName" yaml:"userName" binding:"required" `
	Password string `json:"password" form:"password" json:"password" uri:"password" xml:"password" yaml:"password" binding:"required"`
}
