import Vue from 'vue'
import Router from 'vue-router'
import login from '@/views/login.vue'
import layout from '@/views/layout.vue'
Vue.use(Router)

const originalPush = Router.prototype.push
Router.prototype.push = function push(location) {
  return originalPush.call(this, location).catch(err => err)
}
export const constantRoutes = [
  {
    path: '/login',
    component: login,
    hidden: true
  },
  { path: '/', redirect: '/login' },
  {
    path: '/404',
    component: () => import('@/views/404'),
    hidden: true
  },
  {
    path: '/home',
    component: layout,
    hidden: false,
    redirect: '/welcome',
    children: [
      {
        path: '/welcome',
        component: () => import('@/views/welcome'),
        name: '/welcome',
      }
    ],
  }
]


const createRouter = () => new Router({
  // mode: 'history', // require service support
  scrollBehavior: () => ({ y: 0 }),
  routes: constantRoutes
})

const router = createRouter()

// Detail see: https://github.com/vuejs/vue-router/issues/1234#issuecomment-357941465
export function resetRouter() {
  const newRouter = createRouter()
  router.matcher = newRouter.matcher // reset router
}

export default router
