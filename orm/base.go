package orm

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	Id        int64          `gorm:"primary_key;auto_increment:false"`
	CreatedAt *time.Time     `gorm:"column:add_time"`
	UpdatedAt *time.Time     `gorm:"column:update_time"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_time"`
}

func (base *BaseModel) GetId() int64 {
	return base.Id
}
func (base *BaseModel) SetId(id int64) {
	base.Id = id
}

func (base *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	// 在创建操作之前，设置 CreatedAt 和 UpdatedAt 字段
	now := time.Now()
	base.CreatedAt = &now
	base.UpdatedAt = &now
	base.Id = SnowflakeGenerator.Generate()
	return
}

func (base *BaseModel) BeforeUpdate(tx *gorm.DB) (err error) {
	// 在更新操作之前，更新 UpdatedAt 字段
	now := time.Now()
	base.UpdatedAt = &now
	return
}

// 自定义逻辑删除逻辑,目前用原生DeletedAt,框架限制用gorm自带的更好
/*var ErrCancelled = errors.New("Deletion cancelled")
func (b *BaseModel) BeforeDelete(tx *gorm.DB) (err error) {
	if b.IsDeleted {
		return gorm.ErrRecordNotFound
	}
	now := time.Now()
	b.DeletedAt = &now
	b.IsDeleted = true
	// skip actual delete
	tx.Session(&gorm.Session{SkipHooks: true}).Model(b).Updates(b)
	return ErrCancelled
}
func WithDeletedFalse() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("is_deleted = ?", false)
	}
}*/
