package sys

import (
	"admin/controllers/common"
	"admin/models"
	"admin/pkg/hash"

	"github.com/gin-gonic/gin"
)

type User struct{}

// 分页数据?page=1&limit=20&sort=""&key=""&status=0
//sort= + or -  
func (User) List(c *gin.Context) {
	//默认是1
	page := common.GetPageIndex(c)
	//默认是20
	limit := common.GetPageLimit(c)
	//默认是"",sort= +** or  -**
	sort := common.GetPageSort(c)
	//默认是""
	key := common.GetPageKey(c)
	//默认是0
	status := common.GetQueryToUint(c, "status")
	var whereOrder []models.PageWhereOrder
	order := "id DESC"
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
	//key用于模糊查询用户名
	if key != "" {
		v := "%" + key + "%"
		var arr []interface{}
		arr = append(arr, v)
		arr = append(arr, v)
		whereOrder = append(whereOrder, models.PageWhereOrder{Where: "username like ? or realname like ?", Value: arr})
	}
	//查找状态
	if status > 0 {
		var arr []interface{}
		arr = append(arr, status)
		whereOrder = append(whereOrder, models.PageWhereOrder{Where: "status = ?", Value: arr})
	}
	var total uint64
	list := []models.User{}
	err := models.GetPage(&models.User{}, &models.User{}, &list, page, limit, &total, whereOrder...)
	if err != nil {
		//返回到前端的数据格式是{code: ,message:}
		common.ResErrSrv(c, err)
		return
	}
	//返回到前端的数据格式是{code: ,message:,data:}
	common.ResSuccessPage(c, total, &list)
}

// 详情
func (User) Detail(c *gin.Context) {
	id := common.GetQueryToUint64(c, "id")
	var model models.User
	where := models.User{}
	where.Id = id
	_, err := models.First(&where, &model)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}

	common.ResSuccess(c, &model)
}

// 更新,更改redis
func (User) Update(c *gin.Context) {
	model := models.User{}
	err := c.Bind(&model)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	where := models.User{}
	where.Id = model.Id
	modelOld := models.User{}
	_, err = models.First(&where, &modelOld)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	model.UserName = modelOld.UserName
	model.Password = modelOld.Password
	err = models.Save(&model)
	if err != nil {
		common.ResFail(c, "操作失败")
		return
	}
	common.ResSuccessMsg(c)
}

//新增 更改redis
func (User) Create(c *gin.Context) {
	model := models.User{}
	err := c.Bind(&model)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	model.Password = hash.Md5String(common.MD5_PREFIX + model.Password)
	err = models.Create(&model)
	if err != nil {
		common.ResFail(c, err.Error())
		return
	}
	common.ResSuccess(c, gin.H{"id": model.Id})
}

// 删除数据,userId是一个数组
func (User) Delete(c *gin.Context) {
	var ids []uint64
	err := c.Bind(&ids)
	if err != nil || len(ids) == 0 {
		common.ResErrSrv(c, err)
		return
	}
	user:=models.User{}
	err = user.Delete(ids)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	common.ResSuccessMsg(c)
}

// 获取用户下的角色ID列表
func (User) AdminsRoleIDList(c *gin.Context) {
	adminsid := common.GetQueryToUint64(c, "adminsId")
	roleList := []uint64{}
	where := models.UserRole{UserId: adminsid}
	err := models.PluckList(&models.UserRole{}, &where, &roleList, "roleId")
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	common.ResSuccess(c, &roleList)
}

// 分配用户角色权限
func (User) SetRole(c *gin.Context) {
	adminsid := common.GetQueryToUint64(c, "adminsId")
	var roleids []uint64
	err := c.Bind(&roleids)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	ar := models.UserRole{}
	err = ar.SetRole(adminsid, roleids)
	if err != nil {
		common.ResErrSrv(c, err)
		return
	}
	//另外开一个线程让Casbin处理权限
	go common.CsbinAddRoleForUser(adminsid)
	common.ResSuccessMsg(c)
}
