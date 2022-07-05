<template>
	<div class="login-wrap">
		<div class="ms-login">
			<div class="ms-title">演示后台系统</div>
			<el-form :model="form_data" :rules="rules" ref="login" label-width="0px" class="ms-content">
				<el-form-item prop="account">
					<el-input v-model="form_data.account" placeholder="账号">
						<el-button slot="prepend" icon="el-icon-lx-people"></el-button>
					</el-input>
				</el-form-item>
				<el-form-item prop="password">
					<el-input type="password" placeholder="密码" v-model="form_data.password" @keyup.enter.native="dologin()">
						<el-button slot="prepend" icon="el-icon-lx-lock"></el-button>
					</el-input>
				</el-form-item>
				<el-form-item prop="verifycode">
					<el-input v-model="form_data.verifycode" placeholder="谷歌验证码">
						<el-button slot="prepend" icon="el-icon-suitcase"></el-button>
					</el-input>
				</el-form-item>
				<div class="login-btn">
					<el-button type="primary" @click="dologin()">登录</el-button>
				</div>
			</el-form>
		</div>
	</div>
</template>

<script>
import { app } from '@/api/app.js'
export default {
	data: function () {
		return {
			form_data: {
				account: 'admin',
				password: 'admin',
				verifycode: '1',
			},
			rules: {
				account: [{ required: true, message: '请输入用户名' }],
				password: [{ required: true, message: '请输入密码' }],
				verifycode: [{ required: true, message: '请输入谷歌验证码' }],
			},
		}
	},
	methods: {
		dologin() {
			this.$refs.login.validate((valid, object) => {
				if (valid) {
					app.login(this.form_data.account, this.$md5(this.form_data.password), this.form_data.verifycode, () => {
						this.$router.push('/home')
						this.$message.success('登录成功')
					})
				}
			})
		},
	},
}
</script>

<style scoped>
.login-wrap {
	position: relative;
	width: 100%;
	height: 100%;
	background-image: url(../assets/img/login-bg.jpg);
	background-size: 100%;
}
.ms-title {
	width: 100%;
	line-height: 50px;
	text-align: center;
	font-size: 20px;
	color: #fff;
	border-bottom: 1px solid #ddd;
}
.ms-login {
	position: absolute;
	left: 50%;
	top: 50%;
	width: 350px;
	margin: -190px 0 0 -175px;
	border-radius: 5px;
	background: rgba(255, 255, 255, 0.3);
	overflow: hidden;
}
.ms-content {
	padding: 30px 30px;
}
.login-btn {
	text-align: center;
}
.login-btn button {
	width: 100%;
	height: 36px;
	margin-bottom: 10px;
}
.login-tips {
	font-size: 12px;
	line-height: 30px;
	color: #fff;
}
</style>
