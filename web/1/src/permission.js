import router from './router'
import store from './store'
import { Message } from 'element-ui'
import { getToken } from '@/utils/auth'
// const whiteList = ['/login', '/auth-redirect']


router.beforeEach(async (to, from, next) => {
  //const tokenStr = window.sessionStorage.getItem('token')
  const hasToken = getToken()
  if (hasToken) {
    if (to.path === '/login') {
      next()
    } else {
      const hasRoles = store.getters.roles && store.getters.roles.length > 0
      if (hasRoles) {
        next()
      } else {
        try {
          const { menus } = await store.dispatch('user/getInfo')
          const accessRoutes = await store.dispatch('permission/generateRoutes', menus)
          console.log(accessRoutes)
          router.addRoutes(accessRoutes)
          next({ ...to, replace: true })
        } catch (error) {
          await store.dispatch('user/resetToken')
          Message.error(error || 'Has Error')
          next('/login')

        }
      }
    }
  } else {
    if (to.path === '/login') {
      next()
    } {
      next('/login')
    }
  }
})
