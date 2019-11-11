package sys

import (
	"encoding/json"
	"admin/controllers/common"
	"time"

	"admin/pkg/cache"
	"admin/pkg/convert"
	"admin/pkg/hash"
	"admin/pkg/jwt"
	"admin/pkg/util"
	"admin/models"

	linq "github.com/ahmetb/go-linq"
	"github.com/gin-gonic/gin"
)

//type User struct{}

// 用户登录
func (User) Login(c *gin.Context) {
	requestData, err := c.GetRawData()
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	var requestMap map[string]string
	err = json.Unmarshal(requestData, &requestMap)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	username := requestMap["username"]
	password := requestMap["password"]
	if username == "" || password == "" {
		common.ResFail(c, "用户名或密码不能为空")
		return
	}
	password = hash.Md5String(common.MD5_PREFIX + password)
	where := models.User{UserName: username, Password: password}
	user := models.User{}
	if username == "admin" && password == "5ae6fb77baed98bc122d9f5d304e5736" {
		user.Id = common.SUPER_ADMIN_ID

		user.Status = 1
	} else {
		notFound, err := models.First(&where, &user)
		if err != nil {
			if notFound {
				common.ResFail(c, "用户名或密码错误")
				return
			}
			common.ResErrSrv(c, err)
			return
		}
	}
	if user.Status != 1 {
		common.ResFail(c, "该用户已被禁用")
		return
	}
	// 缓存或者redis
	uuid := util.GetUUID()
	err = cache.Set([]byte(uuid), []byte(convert.ToString(user.Id)), 60*60) // 1H
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	// token jwt
	userInfo := make(map[string]string)
	userInfo["exp"] = convert.ToString(time.Now().Add(time.Hour * time.Duration(1)).Unix()) // 1H
	userInfo["iat"] = convert.ToString(time.Now().Unix())
	userInfo["uuid"] = uuid
	token := jwt.CreateToken(userInfo)
	// 发至页面
	resData := make(map[string]string)
	resData["token"] = token
	//casbin 处理
	err = common.CsbinAddRoleForUser(user.Id)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	common.ResSuccess(c, &resData)
}

// 用户登出
func (User) Logout(c *gin.Context) {
	// 删除缓存
	uuid, exists := c.Get(common.USER_UUID_Key)
	if exists {
		cache.Del([]byte(convert.ToString(uuid)))
	}
	common.ResSuccessMsg(c)
}

// 获取用户信息及可访问的权限菜单
func (User) Info2(c *gin.Context) {
	type MenuMeta struct {
		Title   string `json:"title"`
		Icon    string `json:"icon"`
		NoCache bool   `json:"noCache"`
	}
	type MenuModel struct {
		Path      string      `json:"path"`
		Component string      `json:"component"`
		Name      string      `json:"name"`
		Hidden    bool        `json:"hidden"`
		Meta      MenuMeta    `json:"meta"`
		Children  []MenuModel `json:"children"`
	}
	var menus []MenuModel
	//图标
	menu01Children01 := MenuModel{
		Path:      "/icon/index",
		Component: "icon_index", //@/views/tab/index  icon_index
		Name:      "Icons",
		Children:  []MenuModel{},
		Meta:      MenuMeta{Title: "图标管理", Icon: "icon", NoCache: true}}
	menu01Children0102 := MenuModel{
		Path:      "/icon/index2",
		Component: "icon_index", //@/views/tab/index  icon_index
		Name:      "Icons",
		Children:  []MenuModel{},
		Meta:      MenuMeta{Title: "图标管理2", Icon: "icon", NoCache: true}}
	menu01 := MenuModel{
		Path:      "/icon",
		Component: "Layout",
		Name:      "icon",
		Hidden:    false,
		Meta:      MenuMeta{Title: "图标", Icon: "icon", NoCache: true},
		Children:  []MenuModel{menu01Children01, menu01Children0102}}
	menus = append(menus, menu01)

	//文章
	menu01Children01 = MenuModel{
		Path:      "create",
		Component: "example_create",
		Name:      "CreateArticle",
		Children:  []MenuModel{},
		Meta:      MenuMeta{Title: "添加文章", Icon: "edit", NoCache: false}}
	menu01Children02 := MenuModel{
		Path:      "list",
		Component: "example_list",
		Name:      "ArticleList",
		Children:  []MenuModel{},
		Meta:      MenuMeta{Title: "文章列表", Icon: "list", NoCache: false}}
	menu01Children03 := MenuModel{
		Path:      "edit/:id",
		Component: "example_edit",
		Name:      "ArticleEdit",
		Hidden:    true,
		Children:  []MenuModel{},
		Meta:      MenuMeta{Title: "文章编辑", Icon: "edit", NoCache: false}}
	menu01 = MenuModel{
		Path:      "/example",
		Component: "Layout",
		Name:      "Article",
		Meta:      MenuMeta{Title: "文章", Icon: "example", NoCache: true},
		Children:  []MenuModel{menu01Children01, menu01Children02, menu01Children03}}
	menus = append(menus, menu01)
	type LoginModelData struct {
		Menus        []MenuModel `json:"menus"`
		Roles        []string    `json:"roles22"`
		Introduction string      `json:"introduction"`
		Avatar       string      `json:"avatar"`
		Name         string      `json:"name"`
	}
	resData := LoginModelData{Menus: menus, Roles: []string{"admin"}, Name: "Name002"}
	resData.Avatar = "https://gocn.vip/uploads/nav_menu/12.jpg"
	common.ResSuccess(c, resData)
}

type MenuMeta struct {
	Title   string `json:"title"`   // 标题
	Icon    string `json:"icon"`    // 图标
	NoCache bool   `json:"noCache"` // 是不是缓存
}

type MenuModel struct {
	Path      string      `json:"path"`      // 路由
	Component string      `json:"component"` // 对应vue中的map name
	Name      string      `json:"name"`      // 菜单名称
	Hidden    bool        `json:"hidden"`    // 是否隐藏
	Meta      MenuMeta    `json:"meta"`      // 菜单信息
	Children  []MenuModel `json:"children"`  // 子级菜单
}

type UserData struct {
	Menus        []MenuModel `json:"menus"`        // 菜单
	Introduction string      `json:"introduction"` // 介绍
	Avatar       string      `json:"avatar"`       // 图标
	Name         string      `json:"name"`         // 姓名
}

// 获取用户信息及可访问的权限菜单
func (User) Info(c *gin.Context) {
	// 用户ID
	uid, isExit := c.Get(common.USER_ID_Key)
	if !isExit {
		common.ResFailCode(c, "token 无效", 50008)
		return
	}
	userID := convert.ToUint64(uid)
	// 根据用户ID获取用户权限菜单
	var menuData []models.Menu
	var err error
	//超级管理员登录的话先插入下面数据
	if userID == common.SUPER_ADMIN_ID {
		//管理员
		menuData, err = getAllMenu()
		if err != nil {
			common.ResErrSrv(c, err)
			return
		}
		if len(menuData) == 0 {
			menuModelTop := models.Menu{Status: 1, ParentId: 0, URL: "", Name: "TOP", Sequence: 1, MenuType: 1, Code: "TOP", OperateType: "none"}
			models.Create(&menuModelTop)
			menuModelSys := models.Menu{Status: 1, ParentId: menuModelTop.Id, URL: "", Name: "系统管理", Sequence: 1, MenuType: 1, Code: "Sys", Icon: "lock", OperateType: "none"}
			models.Create(&menuModelSys)
			menuModel := models.Menu{Status: 1, ParentId: menuModelSys.Id, URL: "/icon", Name: "图标管理", Sequence: 10, MenuType: 2, Code: "Icon", Icon: "icon", OperateType: "none"}
			models.Create(&menuModel)
			menuModel = models.Menu{Status: 1, ParentId: menuModelSys.Id, URL: "/menu", Name: "菜单管理", Sequence: 20, MenuType: 2, Code: "Menu", Icon: "documentation", OperateType: "none"}
			models.Create(&menuModel)
			InitMenu(menuModel)
			menuModel = models.Menu{Status: 1, ParentId: menuModelSys.Id, URL: "/role", Name: "角色管理", Sequence: 30, MenuType: 2, Code: "Role", Icon: "tree", OperateType: "none"}
			models.Create(&menuModel)
			InitMenu(menuModel)
			menuModel = models.Menu{Status: 1, ParentId: menuModel.Id, URL: "/role/setrole", Name: "分配角色菜单", Sequence: 6, MenuType: 3, Code: "RoleSetrolemenu", Icon: "", OperateType: "setrolemenu"}
			models.Create(&menuModel)
			menuModel = models.Menu{Status: 1, ParentId: menuModelSys.Id, URL: "/admins", Name: "后台用户管理", Sequence: 40, MenuType: 2, Code: "Admins", Icon: "user", OperateType: "none"}
			models.Create(&menuModel)
			InitMenu(menuModel)
			menuModel = models.Menu{Status: 1, ParentId: menuModel.Id, URL: "/admins/setrole", Name: "分配角色", Sequence: 6, MenuType: 3, Code: "AdminsSetrole", Icon: "", OperateType: "setadminrole"}
			models.Create(&menuModel)

			menuData, _ = getAllMenu()
		}
	} else {
		menuData, err = getMenusByAdminsid(userID)
		if err != nil {
			common.ResErrSrv(c, err)
			return
		}
	}
	//获取所有菜单后，对菜单列表进行排序,把子菜单依次加入
	var menus []MenuModel
	if len(menuData) > 0 {
		var topmenuid uint64 = menuData[0].ParentId
		if topmenuid == 0 {
			topmenuid = menuData[0].Id
		}
		menus = setMenu(menuData, topmenuid)
	}
	if len(menus) == 0 && userID == common.SUPER_ADMIN_ID {
		menus = getSuperAdminMenu()
	}
	resData := UserData{Menus: menus, Name: "小王"}
	resData.Avatar = "http://127.0.0.1:3003/resource/img/head_go.jpg"
	common.ResSuccess(c, &resData)
}

//查询所有菜单
func getAllMenu() (menus []models.Menu, err error) {
	models.Find(&models.Menu{}, &menus, "parentId asc", "sequence asc")
	return
}

//获取超级管理员初使菜单
func getSuperAdminMenu() (out []MenuModel) {
	menuTop := MenuModel{
		Path:      "/sys",
		Component: "Sys",
		Name:      "Sys",
		Meta:      MenuMeta{Title: "系统管理", NoCache: false},
		Children:  []MenuModel{}}
	menuModel := MenuModel{
		Path:      "/icon",
		Component: "Icon",
		Name:      "Icon",
		Meta:      MenuMeta{Title: "图标管理", NoCache: false},
		Children:  []MenuModel{}}
	menuTop.Children = append(menuTop.Children, menuModel)
	menuModel = MenuModel{
		Path:      "/menu",
		Component: "Menu",
		Name:      "Menu",
		Meta:      MenuMeta{Title: "菜单管理", NoCache: false},
		Children:  []MenuModel{}}
	menuTop.Children = append(menuTop.Children, menuModel)
	menuModel = MenuModel{
		Path:      "/role",
		Component: "Role",
		Name:      "Role",
		Meta:      MenuMeta{Title: "角色管理", NoCache: false},
		Children:  []MenuModel{}}
	menuTop.Children = append(menuTop.Children, menuModel)
	menuModel = MenuModel{
		Path:      "/admins",
		Component: "Admins",
		Name:      "Admins",
		Meta:      MenuMeta{Title: "用户管理", NoCache: false},
		Children:  []MenuModel{}}
	menuTop.Children = append(menuTop.Children, menuModel)
	out = append(out, menuTop)
	return
}

// 递归菜单
func setMenu(menus []models.Menu, parentID uint64) (out []MenuModel) {
	var menuArr []models.Menu
	linq.From(menus).Where(func(c interface{}) bool {
		return c.(models.Menu).ParentId == parentID
	}).OrderBy(func(c interface{}) interface{} {
		return c.(models.Menu).Sequence
	}).ToSlice(&menuArr)
	if len(menuArr) == 0 {
		return
	}
	noCache := false
	for _, item := range menuArr {
		menu := MenuModel{
			Path:      item.URL,
			Component: item.Code,
			Name:      item.Code,
			Meta:      MenuMeta{Title: item.Name, Icon: item.Icon, NoCache: noCache},
			Children:  []MenuModel{}}
		if item.MenuType == 3 {
			menu.Hidden = true
		}
		//查询是否有子级
		menuChildren := setMenu(menus, item.Id)
		if len(menuChildren) > 0 {
			menu.Children = menuChildren
		}
		if item.MenuType == 2 {
			// 添加子级首页，有这一级NoCache才有效
			menuIndex := MenuModel{
				Path:      "index",
				Component: item.Code,
				Name:      item.Code,
				Meta:      MenuMeta{Title: item.Name, Icon: item.Icon, NoCache: noCache},
				Children:  []MenuModel{}}
			menu.Children = append(menu.Children, menuIndex)
			menu.Name = menu.Name + "index"
			menu.Meta = MenuMeta{}
		}
		out = append(out, menu)
	}
	return
}

//查询登录用户权限菜单
func getMenusByAdminsid(adminsid uint64) (ret []models.Menu, err error) {
	menu := models.Menu{}
	var menus []models.Menu
	err = menu.GetMenuByAdminsid(adminsid, &menus)
	if err != nil || len(menus) == 0 {
		return
	}
	allmenu, err := getAllMenu()
	if err != nil || len(allmenu) == 0 {
		return
	}
	menuMapAll := make(map[uint64]models.Menu)
	for _, item := range allmenu {
		menuMapAll[item.Id] = item
	}
	menuMap := make(map[uint64]models.Menu)
	for _, item := range menus {
		menuMap[item.Id] = item
	}
	for _, item := range menus {
		_, exists := menuMap[item.ParentId]
		if exists {
			continue
		}
		setMenuUp(menuMapAll, item.ParentId, menuMap)
	}
	for _, m := range menuMap {
		ret = append(ret, m)
	}
	linq.From(ret).OrderBy(func(c interface{}) interface{} {
		return c.(models.Menu).ParentId
	}).ToSlice(&ret)
	return
}

// 向上查找父级菜单
func setMenuUp(menuMapAll map[uint64]models.Menu, menuid uint64, menuMap map[uint64]models.Menu) {
	menuModel, exists := menuMapAll[menuid]
	if exists {
		mid := menuModel.Id
		_, exists = menuMap[mid]
		if !exists {
			menuMap[mid] = menuModel
			setMenuUp(menuMapAll, menuModel.ParentId, menuMap)
		}
	}
}

// 用户修改密码
func (User) EditPwd(c *gin.Context) {
	// 用户ID
	uid, isExit := c.Get(common.USER_ID_Key)
	if !isExit {
		common.ResFailCode(c, "token 无效", 50008)
		return
	}
	userID := convert.ToUint64(uid)
	reqData := make(map[string]string)
	err := c.Bind(&reqData)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	old_password := reqData["old_password"]
	old_password = hash.Md5String(common.MD5_PREFIX + old_password)
	new_password := reqData["new_password"]
	if len(new_password) < 6 || len(new_password) > 20 {
		common.ResFail(c, "密码长度在 6 到 20 个字符")
		return
	}
	new_password = hash.Md5String(common.MD5_PREFIX + new_password)
	where := models.User{}
	where.Id = userID
	modelOld := models.User{}
	_, err = models.First(&where, &modelOld)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	if old_password != modelOld.Password {
		common.ResFail(c, "原密码输入不正确")
		return
	}
	modelNew := models.User{Password: new_password}
	err = models.Updates(&modelOld, &modelNew)
	if err != nil {
		common.ResFail(c, "操作失败")
		return
	}
	common.ResSuccessMsg(c)
}
