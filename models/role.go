package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// 角色
type Role struct {
	Model
	Memo     string `gorm:"column:memo;size:64;" json:"memo" `                // 备注
	Name     string `gorm:"column:name;size:32;not null;" json:"name"`       // 名称
	Sequence int    `gorm:"column:sequence;not null;" json:"sequence" `   // 排序值
	ParentId uint64 `gorm:"column:parentId;not null;" json:"parentId"` // 父级ID
}

// 表名
func (Role) TableName() string {
	return "Role"
}

// 添加前
func (m *Role) BeforeCreate(scope *gorm.Scope) error {
	m.CreatedTime = time.Now()
	m.UpdatedTime = time.Now()
	return nil
}

// 更新前
func (m *Role) BeforeUpdate(scope *gorm.Scope) error {
	m.UpdatedTime = time.Now()
	return nil
}

// 删除角色及关联数据
func (Role) Delete(roleids []uint64) error {
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("id in (?)", roleids).Delete(&Role{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("roleId in (?)", roleids).Delete(&RoleMenu{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
