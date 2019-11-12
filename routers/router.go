package routers

import (
	"admin/controllers/sys"
	"github.com/gin-gonic/gin"
	"net/http"
	"admin/middleware"
)
func RegisterRouterSys(app *gin.RouterGroup) {
	menu := sys.Menu{}
	app.GET("/menu/list", menu.List)
	app.GET("/menu/detail", menu.Detail)
	app.GET("/menu/allmenu", menu.AllMenu)
	app.GET("/menu/menubuttonlist", menu.MenuButtonList)
	app.POST("/menu/delete", menu.Delete)
	app.POST("/menu/update", menu.Update)
	app.POST("/menu/create", menu.Create)
	user := sys.User{}
	app.GET("/user/info", user.Info)
	app.POST("/user/login", user.Login)
	app.POST("/user/logout", user.Logout)
	app.POST("/user/editpwd", user.EditPwd)
	userSet := sys.User{}
	app.GET("/userSet/list", userSet.List)
	app.GET("/userSet/detail", userSet.Detail)
	app.GET("/userSet/adminsroleidlist", userSet.AdminsRoleIDList)
	app.POST("/userSet/delete", userSet.Delete)
	app.POST("/userSet/update", userSet.Update)
	app.POST("/userSet/create", userSet.Create)
	app.POST("/userSet/setrole", userSet.SetRole)
	role := sys.Role{}
	app.GET("/role/list", role.List)
	app.GET("/role/detail", role.Detail)
	app.GET("/role/rolemenuidlist", role.RoleMenuIDList)
	app.GET("/role/allrole", role.AllRole)
	app.POST("/role/delete", role.Delete)
	app.POST("/role/update", role.Update)
	app.POST("/role/create", role.Create)
	app.POST("/role/setrole", role.SetRole)
}

func RegisterRouter(app *gin.Engine) {
	//首页
	app.GET("/", func(c *gin.Context) { c.HTML(http.StatusOK, "index.html", nil) })
	apiPrefix:="/api"
	g := app.Group(apiPrefix)
	// 登录验证 jwt token 验证 及信息提取
	var notCheckLoginUrlArr []string
	notCheckLoginUrlArr = append(notCheckLoginUrlArr, apiPrefix+"/user/login")
	notCheckLoginUrlArr = append(notCheckLoginUrlArr, apiPrefix+"/user/logout")
	g.Use(middleware.UserAuthMiddleware(
		middleware.AllowPathPrefixSkipper(notCheckLoginUrlArr...),
	))
	// 权限验证
	var notCheckPermissionUrlArr []string
	notCheckPermissionUrlArr = append(notCheckPermissionUrlArr, notCheckLoginUrlArr...)
	notCheckPermissionUrlArr = append(notCheckPermissionUrlArr, apiPrefix+"/menu/menubuttonlist")
	notCheckPermissionUrlArr = append(notCheckPermissionUrlArr, apiPrefix+"/menu/allmenu")
	notCheckPermissionUrlArr = append(notCheckPermissionUrlArr, apiPrefix+"/admins/adminsroleidlist")
	notCheckPermissionUrlArr = append(notCheckPermissionUrlArr, apiPrefix+"/user/info")
	notCheckPermissionUrlArr = append(notCheckPermissionUrlArr, apiPrefix+"/user/editpwd")
	notCheckPermissionUrlArr = append(notCheckPermissionUrlArr, apiPrefix+"/role/rolemenuidlist")
	notCheckPermissionUrlArr = append(notCheckPermissionUrlArr, apiPrefix+"/role/allrole")
	g.Use(middleware.CasbinMiddleware(
		middleware.AllowPathPrefixSkipper(notCheckPermissionUrlArr...),
	))
	//sys
	RegisterRouterSys(g)

}
