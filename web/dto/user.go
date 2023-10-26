package dto

type UserDto struct {
	UserName string `json:"userName" form:"userName" json:"userName" uri:"userName" xml:"userName" yaml:"userName" binding:"required" `
	Password string `json:"password" form:"password" json:"password" uri:"password" xml:"password" yaml:"password" binding:"required"`
	Salt     string `json:"salt" form:"salt" json:"salt" uri:"salt" xml:"salt" yaml:"salt" binding:"required"`
	Role     int    `json:"role" form:"role" json:"role" uri:"role" xml:"role" yaml:"role" binding:"required"`
}
