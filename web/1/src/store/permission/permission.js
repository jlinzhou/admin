import { constantRoutes } from '@/router'


// 在这里定义state 状态管理
const state = {
  routes: [],
  addRoutes: []
}

const mutations = {
  SET_ROUTES: (state, routes) => {
    state.addRoutes = routes
    state.routes = constantRoutes.concat(routes)
  }
}

const actions = {
  generateRoutes({ commit }, data) {
    return new Promise(resolve => {
      const accessedRouters = convertRouter(data)
      commit('SET_ROUTES', accessedRouters)
      resolve(accessedRouters)
    })
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}

function convertRouter(asyncRouterMap) {
  const accessedRouters = []
  if (asyncRouterMap) {
    asyncRouterMap.forEach(item => {
      var isParent = false
      if (item.children) {
        isParent = true
      }
      var parent = generateRouter(item, isParent)

      var children = []
      if (item.children) {

        item.children.forEach(child => {

          var children2 = []
          if (child.children) {
            child.children.forEach(child2 => {
              children2.push(generateRouter(child2, false))
            })
          }
          var parent2 = generateRouter(child, false)
          parent2.children = children2
          children.push(parent2)

        })
      }
      parent.children = children
      accessedRouters.push(parent)
    })
  }
  accessedRouters.push({ path: '*', redirect: '/404', hidden: true })
  return accessedRouters
}

// 对component的处理
function generateRouter(item, isParent) {
  var component = Layout // 多层嵌套时只能有一个Layout
  if (isParent !== true) {
    component = componentsMap[item.component]
  }
  var router = {
    path: item.path,
    name: item.name,
    meta: item.meta,
    hidden: item.hidden,

    component: component
  }

  return router
}

// componentsMap 需要在事先定义好
export const componentsMap = {
  example_create: () => import('@/views/example/create'), // 添加文章
  example_edit: () => import('@/views/example/edit'), // 文章编辑
  table_index: () => import('@/views/tab/index'), // 表格首页
  Icon: () => import('@/views/svg-icons/index'), // 图标管理
  Menu: () => import('@/views/app/sys/menu'), // 菜单
  Admins: () => import('@/views/app/sys/admins'), // 后台管理员
  Role: () => import('@/views/app/sys/role') // 后台角色
}
