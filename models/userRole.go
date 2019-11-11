package models

import (
	"github.com/jinzhu/gorm"
	"time"

)

// 用户-角色
type UserRole struct {
	Model
	UserId uint64 `DB:"column:userId;not null;" json:"userId"` // 管理员ID
	RoleId uint64 `DB:"column:roleId;not null;" json:"roleId"`   // 角色ID
}

// 表名
func (UserRole) TableName() string {
	return "UserRole"
}

// 添加前
func (m *UserRole) BeforeCreate(scope *gorm.Scope) error {
	m.CreatedTime = time.Now()
	m.UpdatedTime = time.Now()
	return nil
}

// 更新前
func (m *UserRole) BeforeUpdate(scope *gorm.Scope) error {
	m.UpdatedTime = time.Now()
	return nil
}

// 分配用户角色
func (UserRole) SetRole(userId uint64, roleids []uint64) error {
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
	if err := tx.Where(&UserRole{UserId: userId}).Delete(&UserRole{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if len(roleids) > 0 {
		for _, rid := range roleids {
			rm := new(UserRole)
			rm.RoleId = rid
			rm.UserId = userId
			if err := tx.Create(rm).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return tx.Commit().Error
}
