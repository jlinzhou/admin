package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

// 角色-菜单
type RoleMenu struct {
	Model
	RoleId uint64 `DB:"column:roleId;unique_index:uk_role_menu_role_id;not null;" json:"roleId"` // 角色ID
	MenuId uint64 `DB:"column:menuId;unique_index:uk_role_menu_role_id;not null;" json:"menuId"` // 菜单ID
}

// 表名
func (RoleMenu) TableName() string {
	return "RoleMenu"
}

// 添加前
func (m *RoleMenu) BeforeCreate(scope *gorm.Scope) error {
	m.CreatedTime = time.Now()
	m.UpdatedTime = time.Now()
	return nil
}

// 更新前
func (m *RoleMenu) BeforeUpdate(scope *gorm.Scope) error {
	m.UpdatedTime = time.Now()
	return nil
}

// 设置角色菜单权限
func (RoleMenu) SetRole(roleid uint64, menuids []uint64) error {
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
	if err := tx.Where(&RoleMenu{RoleId: roleid}).Delete(&RoleMenu{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if len(menuids) > 0 {
		for _, mid := range menuids {
			rm := new(RoleMenu)
			rm.RoleId = roleid
			rm.MenuId = mid
			if err := tx.Create(rm).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return tx.Commit().Error
}
