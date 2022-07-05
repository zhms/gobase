import Vue from 'vue'
import App from './App.vue'
import router from './router'
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
import './assets/css/icon.css'
import './components/directives'
import 'babel-polyfill'
import { app } from './api/app.js'
import axios from 'axios'
import moment from 'moment'
import md5 from 'js-md5'
Vue.prototype.$moment = moment
Vue.prototype.$md5 = md5
Vue.config.productionTip = false
Vue.use(ElementUI, {
	size: 'small',
})

//使用钩子函数对路由进行权限跳转
router.beforeEach((to, from, next) => {
	document.title = `后台演示系统`
	app.setInfo(sessionStorage.getItem('userdata'))
	const token = sessionStorage.getItem('token')
	if (!token && to.path !== '/login') {
		next('/login')
	} else {
		axios.defaults.headers.common['x-token'] = token
		next()
	}
})

new Vue({
	router,
	render: (h) => h(App),
}).$mount('#app')
