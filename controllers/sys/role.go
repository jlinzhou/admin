package sys

import (
	"admin/controllers/common"
	"admin/models"

	"github.com/gin-gonic/gin"
)

type Role struct{}

// 分页数据
func (Role) List(c *gin.Context) {
	page := common.GetPageIndex(c)
	limit := common.GetPageLimit(c)
	sort := common.GetPageSort(c)
	key := common.GetPageKey(c)
	parent_id := common.GetQueryToUint64(c, "parentId")
	var whereOrder []models.PageWhereOrder
	order := "Id DESC"
	if len(sort) >= 2 {
		orderType := sort[0:1]
		order = sort[1:len(sort)]
		if orderType == "+" {
			order += " ASC"
		} else {
			order += " DESC"
		}
	}
	whereOrder = append(whereOrder, models.PageWhereOrder{Order: order})
	if key != "" {
		v := "%" + key + "%"
		var arr []interface{}
		arr = append(arr, v)
		whereOrder = append(whereOrder, models.PageWhereOrder{Where: "name like ?", Value: arr})
	}
	if parent_id > 0 {
		var arr []interface{}
		arr = append(arr, parent_id)
		whereOrder = append(whereOrder, models.PageWhereOrder{Where: "parentId = ?", Value: arr})
	}
	var total uint64
	list := []models.Role{}
	err := models.GetPage(&models.Role{}, &models.Role{}, &list, page, limit, &total, whereOrder...)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	common.ResSuccessPage(c, total, &list)
}

// 详情
func (Role) Detail(c *gin.Context) {
	id := common.GetQueryToUint64(c, "id")
	var model models.Role
	where := models.Role{}
	where.Id = id
	_, err := models.First(&where, &model)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	common.ResSuccess(c, &model)
}

// 更新
func (Role) Update(c *gin.Context) {
	model := models.Role{}
	err := c.Bind(&model)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	err = models.Save(&model)
	if err != nil {
		common.ResFail(c, "操作失败")
		return
	}
	common.ResSuccessMsg(c)
}

//新增
func (Role) Create(c *gin.Context) {
	model := models.Role{}
	err := c.Bind(&model)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	err = models.Create(&model)
	if err != nil {
		common.ResFail(c, "操作失败")
		return
	}
	common.ResSuccess(c, gin.H{"id": model.Id})
}

// 删除数据
func (Role) Delete(c *gin.Context) {
	var ids []uint64
	err := c.Bind(&ids)
	if err != nil || len(ids) == 0 {
		common.ResErrSrv(c, err)
		return
	}
	role := models.Role{}
	err = role.Delete(ids)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	go common.CsbinDeleteRole(ids)
	common.ResSuccessMsg(c)
}

// 所有角色
func (Role) AllRole(c *gin.Context) {
	var list []models.Role
	err := models.Find(&models.Role{}, &list, "parentId asc", "sequence asc")
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	common.ResSuccess(c, &list)
}

// 获取角色下的菜单ID列表
func (Role) RoleMenuIDList(c *gin.Context) {
	roleid := common.GetQueryToUint64(c, "roleId")
	menuIDList := []uint64{}
	where := models.RoleMenu{RoleId: roleid}
	err := models.PluckList(&models.RoleMenu{}, &where, &menuIDList, "menuId")
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	common.ResSuccess(c, &menuIDList)
}

// 设置角色菜单权限
func (Role) SetRole(c *gin.Context) {
	roleid := common.GetQueryToUint64(c, "roleId")
	var menuids []uint64
	err := c.Bind(&menuids)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	rm := models.RoleMenu{}
	err = rm.SetRole(roleid, menuids)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	go common.CsbinSetRolePermission(roleid)
	common.ResSuccessMsg(c)
}
