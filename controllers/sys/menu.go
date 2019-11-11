package sys

import (
	"admin/controllers/common"
	"admin/pkg/convert"
	"admin/models"

	"github.com/gin-gonic/gin"
)

type Menu struct{}

// 分页数据
func (Menu) List(c *gin.Context) {
	page := common.GetPageIndex(c)
	limit := common.GetPageLimit(c)
	sort := common.GetPageSort(c)
	key := common.GetPageKey(c)
	menuType := common.GetQueryToUint(c, "type")
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
		arr = append(arr, v)
		whereOrder = append(whereOrder, models.PageWhereOrder{Where: "name like ? or code like ?", Value: arr})
	}
	if menuType > 0 {
		var arr []interface{}
		arr = append(arr, menuType)
		whereOrder = append(whereOrder, models.PageWhereOrder{Where: "menuType = ?", Value: arr})
	}
	if parent_id > 0 {
		var arr []interface{}
		arr = append(arr, parent_id)
		whereOrder = append(whereOrder, models.PageWhereOrder{Where: "parentId = ?", Value: arr})
	}
	var total uint64
	list := []models.Menu{}
	err := models.GetPage(&models.Menu{}, &models.Menu{}, &list, page, limit, &total, whereOrder...)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	common.ResSuccessPage(c, total, &list)
}

// 详情
func (Menu) Detail(c *gin.Context) {
	id := common.GetQueryToUint64(c, "id")
	var menu models.Menu
	where := models.Menu{}
	where.Id = id
	_, err := models.First(&where, &menu)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	common.ResSuccess(c, &menu)
}

// 更新
func (Menu) Update(c *gin.Context) {
	//logs.Info(c.Request.Method, c.Request.PostForm)
	//
	//logs.Info(binding.Form)
	//for _, value := range c.Request.PostForm {
	//	logs.Info(value)
	//}
	//logs.Info(c.Request.Body)
	model := models.Menu{}
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
func (Menu) Create(c *gin.Context) {
	menu := models.Menu{}
	err := c.Bind(&menu)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	err = models.Create(&menu)
	if err != nil {
		common.ResFail(c, "操作失败")
		return
	}
	// 新增菜单后自动添加菜单下的常规操作
	go InitMenu(menu)
	common.ResSuccess(c, gin.H{"id": menu.Id})
}

// 删除数据
func (Menu) Delete(c *gin.Context) {
	var ids []uint64
	err := c.Bind(&ids)
	if err != nil || len(ids) == 0 {
		common.ResErrSrv(c, err)
		return
	}
	menu := models.Menu{}
	err = menu.Delete(ids)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	common.ResSuccessMsg(c)
}

// 所有菜单
func (Menu) AllMenu(c *gin.Context) {
	var menus []models.Menu
	err := models.Find(&models.Menu{}, &menus, "parentId asc", "sequence asc")
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	common.ResSuccess(c, &menus)
}

// 新增菜单后自动添加菜单下的常规操作
func InitMenu(model models.Menu) {
	if model.MenuType != 2 {
		return
	}
	add := models.Menu{Status: 1, ParentId: model.Id, URL: model.URL + "/create", Name: "新增", Sequence: 1, MenuType: 3, Code: model.Code + "Add", OperateType: "add"}
	models.Create(&add)
	del := models.Menu{Status: 1, ParentId: model.Id, URL: model.URL + "/delete", Name: "删除", Sequence: 2, MenuType: 3, Code: model.Code + "Del", OperateType: "del"}
	models.Create(&del)
	view := models.Menu{Status: 1, ParentId: model.Id, URL: model.URL + "/detail", Name: "查看", Sequence: 3, MenuType: 3, Code: model.Code + "View", OperateType: "view"}
	models.Create(&view)
	update := models.Menu{Status: 1, ParentId: model.Id, URL: model.URL + "/update", Name: "编辑", Sequence: 4, MenuType: 3, Code: model.Code + "Update", OperateType: "update"}
	models.Create(&update)
	list := models.Menu{Status: 1, ParentId: model.Id, URL: model.URL + "/list", Name: "分页api", Sequence: 5, MenuType: 3, Code: model.Code + "List", OperateType: "list"}
	models.Create(&list)
}

// 获取菜单有权限的操作列表
func (Menu) MenuButtonList(c *gin.Context) {
	// 用户ID
	uid, isExit := c.Get(common.USER_ID_Key)
	if !isExit {
		common.ResFailCode(c, "token 无效", 50008)
		return
	}
	userID := convert.ToUint64(uid)
	menuCode := common.GetQueryToStr(c, "menucode")
	if userID == 0 || menuCode == "" {
		common.ResFail(c, "err")
		return
	}
	btnList := []string{}
	if userID == common.SUPER_ADMIN_ID {
		//管理员
		btnList = append(btnList, "add")
		btnList = append(btnList, "del")
		btnList = append(btnList, "view")
		btnList = append(btnList, "update")
		btnList = append(btnList, "setrolemenu")
		btnList = append(btnList, "setadminrole")
	} else {
		menu := models.Menu{}
		err := menu.GetMenuButton(userID, menuCode, &btnList)
		if err != nil {
			common.ResErrSrv(c, err)
			return
		}
	}
	common.ResSuccess(c, &btnList)
}
