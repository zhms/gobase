package server

var MenuDataStr string = `[
	{
		"icon": "el-icon-lx-home",
		"index": "home",
		"title": "系统首页"
	},
	{
		"icon": "el-icon-user-solid",
		"index": "0",
		"title": "玩家管理",
		"subs": [
			{
				"title": "账号管理",
				"index": "user_list"
			}
		]
	},
	{
		"icon": "el-icon-setting",
		"index": "2",
		"title": "系统管理",
		"subs":
		[
			{
				"index": "system_seller",
				"title": "运营商管理"
			},
			{
				"index": "system_account",
				"title": "账号管理"
			},
			{
				"index": "system_role",
				"title": "角色管理"
			},
			{
				"index": "system_login_log",
				"title": "登录日志"
			},
			{
				"index": "system_log",
				"title": "操作日志"
			}
		]
	}
]`

var AuthDataStr = `{
	"系统首页": { "查" : 1},
	"玩家管理": {
		"账号管理": { "查": 1,"增": 1,"删": 1,"改": 1}
	},
	"系统管理": {
		"运营商管理": { "查": 1,"增": 1,"删": 1,"改": 1},
		"账号管理": { "查": 1,"增": 1,"删": 1,"改": 1},
		"角色管理": { "查": 1,"增": 1,"删": 1,"改": 1},
		"登录日志": { "查": 1},
		"操作日志": { "查": 1}
	}
}`

var SellerAuthDataStr = `{
	"系统首页": { "查" : 1 },
	"玩家管理": {
		"账号管理": { "查": 1,"增": 1,"删": 1,"改": 1}
	},
	"系统管理": {
		"运营商管理": { "查": 1,"增": 1,"删": 1,"改": 1},
		"账号管理": { "查": 1,"增": 1,"删": 1,"改": 1},
		"角色管理": { "查": 1,"增": 1,"删": 1,"改": 1},
		"登录日志": { "查": 1},
		"操作日志": { "查": 1}
	}
}`
