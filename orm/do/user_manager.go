package do

import "buding-job/orm"

func init() {
	orm.DB.AutoMigrate(&UserManagerDo{})
}

type UserManagerDo struct {
	orm.BaseModel
	UserId    int64
	ManagerId int64
}

func (*UserManagerDo) TableName() string {
	return "tb_user_manager"
}
