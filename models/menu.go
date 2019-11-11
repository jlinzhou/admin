package models

import (

	"time"
	"admin/pkg/logs"
	"github.com/jinzhu/gorm"
)

type Model struct {
	Id        		uint64    `gorm:"column:id;primary_key;auto_increment;" json:"id"`                  // 主键
	CreatedTime 	time.Time `gorm:"column:createdTime;type:datetime;not null;" json:"createdTime" ` // 创建时间
	UpdatedTime 	time.Time `gorm:"column:updatedTime;type:datetime;not null;" json:"updatedTime" ` // 更新时间
	CreatedUserId 	uint64    `gorm:"column:createdUserId;default:0;not null;" json:"createdUserId"`     // 创建人
	UpdatedUserId 	uint64    `gorm:"column:updatedUserId;default:0;not null;" json:"updatedUserId" `     // 更新人
}

// 菜单
type Menu struct {
	Model
	Status      uint8  `gorm:"column:status;type:tinyint(1);not null;" json:"status" form:"status"`        // 状态(1:启用 2:不启用)
	Memo        string `gorm:"column:memo;size:64;" json:"memo" form:"memo"`                               // 备注
	ParentId    uint64 `gorm:"column:parentId;not null;" json:"parentId" form:"parentId"`                  // 父级ID
	URL         string `gorm:"column:url;size:72;" json:"url" form:"url"`                                  // 菜单URL
	Name        string `gorm:"column:name;size:32;not null;" json:"name" form:"name"`                      // 菜单名称
	Sequence    int    `gorm:"column:sequence;not null;" json:"sequence" form:"sequence"`                  // 排序值
	MenuType    uint8  `gorm:"column:menuType;type:tinyint(1);not null;" json:"menuType" form:"menuType"`  // 菜单类型 1模块2菜单3操作
	Code        string `gorm:"column:code;size:32;not null;" json:"code" form:"code"`                      // 菜单代码
	Icon        string `gorm:"column:icon;size:32;" json:"icon" form:"icon"`                               // icon
	OperateType string `gorm:"column:operateType;size:32;not null;" json:"operateType" form:"operateType"` // 操作类型 none/add/del/view/update
}

// 表名
func (Menu) TableName() string {
	return "Menu"
}

// 添加前
func (m *Menu) BeforeCreate(scope *gorm.Scope) error {
	m.CreatedTime = time.Now()
	m.UpdatedTime = time.Now()
	return nil
}

// 更新前
func (m *Menu) BeforeUpdate(scope *gorm.Scope) error {
	m.UpdatedTime = time.Now()
	return nil
}

// 获取菜单有权限的操作列表
func (Menu) GetMenuButton(userId uint64, menuCode string, btns *[]string) (err error) {
	sql := `select operateType from Menu
	      where id in (
					select menuId from RoleMenu where 
					menuId in (select id from Menu where parentId in (select id from Menu where code=?))
					and roleId in (select roleId from UserRole where userId=?)
				)`
	err = DB.Raw(sql, menuCode, userId).Pluck("operateType", btns).Error
	logs.Info(btns)
	return
}

// 获取管理员权限下所有菜单
func (Menu) GetMenuByAdminsid(userId uint64, menus *[]Menu) (err error) {
	sql := `select * from Menu
	      where id in (
					select menuId from RoleMenu where 
				  roleId in (select roleId from UserRole where userId=?)
				)`
	err = DB.Raw(sql, userId).Find(menus).Error
	return
}

// 删除菜单及关联数据
func (Menu) Delete(menuIds []uint64) error {
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
	for _, menuid := range menuIds {
		if err := deleteMenuRecurve(tx, menuid); err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Where("menuId in (?)", menuIds).Delete(&RoleMenu{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("id in (?)", menuIds).Delete(&Menu{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func deleteMenuRecurve(db *gorm.DB, parentID uint64) error {
	where := &Menu{}
	where.ParentId = parentID
	var menus []Menu
	dbslect := db.Where(&where)
	if err := dbslect.Find(&menus).Error; err != nil {
		return err
	}
	for _, menu := range menus {
		if err := db.Where("menuId = ?", menu.Id).Delete(&RoleMenu{}).Error; err != nil {
			return err
		}
		if err := deleteMenuRecurve(db, menu.Id); err != nil {
			return err
		}
	}
	if err := dbslect.Delete(&Menu{}).Error; err != nil {
		return err
	}
	return nil
}
