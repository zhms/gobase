import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

export default new Router({
	routes: [
		{ path: '/', redirect: '/home' },
		{ path: '/login', component: () => import('../components/login.vue'), meta: { title: '登录' } },
		{ path: '*', redirect: '/404' },
		{
			path: '/',
			component: () => import('../components/home.vue'),
			meta: { title: '自述文件' },
			children: [
				//首页
				{ path: '/home', component: () => import('../../page/home.vue'), meta: { title: '系统首页' } },
				//玩家管理
				{ path: '/user_list', component: () => import('../../page/user/user.vue'), meta: { title: '账号管理' } },
				//游戏管理
				{ path: '/game_list', component: () => import('../../page/game/game.vue'), meta: { title: '游戏列表' } },
				//系统管理
				{ path: '/system_seller', component: () => import('../../page/system/seller.vue'), meta: { title: '运营商管理' } },
				{ path: '/system_login_log', component: () => import('../../page/system/loginlog.vue'), meta: { title: '登录日志' } },
				{ path: '/system_account', component: () => import('../../page/system/account.vue'), meta: { title: '账号管理' } },
				{ path: '/system_role', component: () => import('../../page/system/roles.vue'), meta: { title: '角色管理' } },
				{ path: '/system_log', component: () => import('../../page/system/log.vue'), meta: { title: '操作日志' } },
			],
		},
	],
})
