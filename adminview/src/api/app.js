export class app {}
import { Message, Loading } from 'element-ui'
import router from '../router'
import axios from 'axios'
var service = axios.create({
	baseURL: 'http://192.168.0.184:4534',
	timeout: 60000,
})
service.interceptors.request.use(
	(config) => {
		config.headers['x-token'] = sessionStorage.getItem('token')
		return config
	},
	(error) => {
		return Promise.reject(error)
	}
)

app.getInstance = (function () {
	let instance
	return function () {
		if (!instance) {
			instance = new app()
		}
		return instance
	}
})()

app.clone = function (obj) {
	return JSON.parse(JSON.stringify(obj))
}

//显示加载动画
app.showLoading = function (show) {
	if (show) {
		if (!this.loading) {
			this.loading = Loading.service({ lock: true, spinner: 'el-icon-loading', background: 'rgba(0, 0, 0, 0.7)' })
		}
	} else {
		if (this.loading) {
			this.loading.close()
			this.loading = null
		}
	}
}
//退回登录界面
app.showLoginPage = function () {
	router.push('/login')
}
//获取管理员信息
app.getInfo = function () {
	return this.info
}
//设置管理员信息
app.setInfo = function (data) {
	if (data) {
		this.info = JSON.parse(data)
	}
}
//get请求
app.get = function (url, p1, p2) {
	var data = null
	var callback = null
	if (typeof p1 == 'object') {
		data = p1
		if (typeof p2 == 'function') {
			callback = p2
		}
	} else if (typeof p1 == 'function') {
		callback = p1
	}
	if (data) {
		url += '?'
		for (var i in data) {
			url += i
			url += '='
			url += data[i]
			url += '&'
		}
	}
	if (url.charAt(url.length - 1) == '&') {
		url = url.substr(0, url.length - 1)
	}
	service({
		url: url,
		method: 'get',
	})
		.then((result) => {
			if (result.data.errmsg) {
				console.log('get:' + url + ' ' + errmsg)
			} else {
				if (callback) callback(result.data)
			}
		})
		.catch((err) => {
			console.log('get:' + url + ' ' + err)
		})
}
//post请求
app.post = function (url, data, callback, noloading) {
	noloading = false
	if (!noloading) app.showLoading(true)
	service({
		url: url,
		method: 'post',
		data,
	})
		.then((result) => {
			if (!noloading) app.showLoading(false)
			if (result.data.code != 200) {
				Message({
					message: result.data.data.errmsg,
					type: 'error',
					duration: 1000 * 3,
					showClose: true,
					center: true,
				})
			} else {
				if (callback) callback(result.data)
			}
		})
		.catch((err) => {
			if (!noloading) app.showLoading(false)
			Message({
				message: err,
				type: 'error',
				duration: 1000 * 3,
				showClose: true,
				center: true,
			})
		})
}

app.login = function (account, password, verifycode, callback) {
	app.post('/admin/user/login', { account, password, verifycode }, (result) => {
		for (let i = 0; i < result.data.MenuData.length; i++) {
			for (let j = 0; j < result.data.MenuData[i].subs.length; j++) {
				for (let k = 0; k < result.data.MenuData[i].subs[j].subs.length; k++) {
					delete result.data.MenuData[i].subs[j].subs[k].subs
				}
			}
		}
		for (let i = 0; i < result.data.MenuData.length; i++) {
			for (let j = 0; j < result.data.MenuData[i].subs.length; j++) {
				if (result.data.MenuData[i].subs[j].subs.length == 0) {
					delete result.data.MenuData[i].subs[j].subs
				}
			}
		}
		for (let i = 0; i < result.data.MenuData.length; i++) {
			if (result.data.MenuData[i].subs.length == 0) {
				delete result.data.MenuData[i].subs
			}
		}
		sessionStorage.setItem('userdata', JSON.stringify(result.data))
		sessionStorage.setItem('token', result.data.Token)
		app.setCurrentSeller(result.data.SellerId)
		if (result.data.SellerId == -1) {
			app.post('/admin/seller/name', {}, (result) => {
				sessionStorage.setItem('seller', JSON.stringify(result.data))
				this.seller = result.data
				callback()
			})
		} else {
			callback()
		}
	})
}

app.auth1 = (m, o) => {
	let info = app.getInfo()
	if (!info) return false
	let authm = info.AuthData[m]
	if (!authm) return false
	let autho = authm[o]
	if (!autho) return false
	return true
}

app.auth2 = (m, s, o) => {
	let info = app.getInfo()
	if (!info) return false
	let authm = info.AuthData[m]
	if (!authm) return false
	let auths = authm[s]
	if (!auths) return false
	let autho = auths[o]
	if (!autho) return false
	return true
}

app.auth3 = (m, s, l, o) => {
	let info = app.getInfo()
	if (!info) return false
	let authm = info.AuthData[m]
	if (!authm) return false
	let auths = authm[s]
	if (!auths) return false
	let authl = auths[l]
	if (!authl) return false
	let autho = authl[o]
	if (!autho) return false
	return true
}

app.zong = function () {
	return app.getInfo().SellerId == -1
}

app.flushSeller = function () {
	app.post('/admin/seller/name', {}, (result) => {
		sessionStorage.setItem('seller', JSON.stringify(result.data))
		this.seller = result.data
	})
}

app.getSeller = function () {
	try {
		this.seller = JSON.parse(sessionStorage.getItem('seller'))
	} catch (e) {}
	return this.seller
}
app.getSellerNoAll = function () {
	try {
		this.seller_noall = JSON.parse(sessionStorage.getItem('seller'))
		for (let i = 0; i < this.seller_noall.length; i++) {
			if (this.seller_noall[i].SellerId == 0) {
				this.seller_noall.splice(i, 1)
			}
		}
	} catch (e) {}
	return this.seller_noall
}

app.getSellerNoZong = function () {
	try {
		this.seller_nozong = JSON.parse(sessionStorage.getItem('seller'))
		for (let i = 0; i < this.seller_nozong.length; i++) {
			if (this.seller_nozong[i].SellerId == 0) {
				this.seller_nozong.splice(i, 1)
			}
		}
		for (let i = 0; i < this.seller_nozong.length; i++) {
			if (this.seller_nozong[i].SellerId == -1) {
				this.seller_nozong.splice(i, 1)
			}
		}
	} catch (e) {}
	return this.seller_nozong
}

app.currentSeller = function () {
	if (!app.zong()) {
		return app.getInfo().SellerId
	}
	if (this.current_seller == null || this.current_seller == undefined) {
		let savecurrentseller = localStorage.getItem('current_seller')
		if (savecurrentseller) {
			this.current_seller = parseInt(savecurrentseller)
		} else {
			this.current_seller = 0
		}
	}
	return this.current_seller
}

app.setCurrentSeller = function (seller) {
	if (seller == null || seller == undefined) return
	this.current_seller = seller
	localStorage.setItem('current_seller', `${seller}`)
	return this.current_seller
}
