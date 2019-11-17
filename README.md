# 基于go和vue的后台管理系统





## 功能介绍

```go
1：用户管理 
    管理可以使用此系统的用户
    功能：查询、新增、修改（可以直接重置密码）删除、分配用户角色
2：角色管理
    主要是用于用户分组和权限分组
    功能：查询、新增、修改、删除、关联权限（设置这个组拥有哪些权限）
3：菜单管理
    用于管理导航菜单中的菜单目录以及定义所有系统中api的名称，便于分配权限。
    功能：查询、新增、修改、删除
```



## 技术选型

1. web:gin
2. orm:gorm
3. database: mysql
4. 权限管理: casbin 
5. 前后端分离
6. 前端: vue-element-admin 



## 项目结构

```go
-admin
	|-cmd 程序执行入口
    |-confs 配置文件目录
    |-controllers 控制器目录
    |-middleware 中间件目录
    |-models 数据库访问目录
    |-pkg 公用程序包目录
    |-routers 路由目录
    |-test 测试目录
	|-web vue前端目录
		|-admin
			|-node_modules项目开发依赖模块
			|-public
			|-src
				|-api 所有后端路由接口
				|-assets 静态文件
				|-components 公用的调用组件
				|-router 路由配置
				|-store 全局变量数据存储
				|-utils 公用的js函数
				|-views 所有的前端页面
				|-App.vue 入口组件
				|-main.js 入口js文件
				|-permission.js 路由导航权限控制拦截
			|-.eslintrc.js
			|-package.json
```



## 已实现功能

- [x] 后端用户、角色、菜单三个模块的curd开发
- [x] 后端 JWT Authorization token认证 
- [x] 后端基于 Casbin 的 RBAC 访问控制模型 
- [x] 前端动态路由生成
- [x] 前端所有按钮显示与否依据后端传来的权限数据来控制