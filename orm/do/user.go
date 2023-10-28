package do

import (
	"buding-job/orm"
)

func init() {

}

type UserDo struct {
	orm.BaseModel
	UserName string `gorm:"type:varchar(100);"`
	Password string `gorm:"type:varchar(100);not null"`
	Salt     string `gorm:"column:salt;;type:varchar(512)"`
	Role     int    `gorm:"column:role;comment:权限 1管理员 2:普通;not null"`
}

func (UserDo) TableName() string {
	return "tb_user"
}
