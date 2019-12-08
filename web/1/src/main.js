import Vue from 'vue'
import App from './App.vue'
import router from './router'
// import './plugins/element.js'
import store from './store'
import './permission'
//导入全局样式表
import '@/assets/css/global.css'
import '@/assets/fonts/iconfont.css'
// import axios from 'axios'
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'

Vue.use(ElementUI);
//配置请求的genglu

Vue.config.productionTip = false


new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
