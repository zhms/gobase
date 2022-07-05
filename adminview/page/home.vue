<template>
	<div>
		<el-row :gutter="20">
			<el-col :span="8">
				<el-card shadow="hover" class="mgb20" style="height: 330px; width: 400px">
					<div class="user-info">
						<img src="../src/assets/img/img.jpg" class="user-avator" alt />
						<div class="user-info-cont">
							<div class="user-info-name">{{ info.Account }}</div>
							<div>{{ info.role }}</div>
						</div>
					</div>
					<div class="user-info-list"><span style="margin-left: 59px">账号：</span> {{ info.Account }}</div>
					<div class="user-info-list"><span style="margin-left: 29px">登录次数：</span> {{ info.LoginCount }}</div>
					<div class="user-info-list"><span style="margin-left: 43px">登录IP：</span> {{ info.Ip }}</div>
					<div class="user-info-list"><span>上次登录时间：</span> {{ info.LoginTime }}</div>
					<div class="user-info-list"><span style="margin-left: 29px">系统版本：</span> {{ info.Version }}</div>
				</el-card>
			</el-col>
		</el-row>
	</div>
</template>

<script>
import { app } from '@/api/app.js'
export default {
	data() {
		return {
			info: app.getInfo(),
		}
	},
	components: {},
	computed: {},
	created() {},
	methods: {},
}
</script>

<style scoped>
.el-row {
	margin-bottom: 20px;
}

.grid-content {
	display: flex;
	align-items: center;
	height: 100px;
}

.grid-cont-right {
	flex: 1;
	text-align: center;
	font-size: 14px;
	color: #999;
}

.grid-num {
	font-size: 30px;
	font-weight: bold;
}

.grid-con-icon {
	font-size: 50px;
	width: 100px;
	height: 100px;
	text-align: center;
	line-height: 100px;
	color: #fff;
}

.grid-con-1 .grid-con-icon {
	background: rgb(45, 140, 240);
}

.grid-con-1 .grid-num {
	color: rgb(45, 140, 240);
}

.grid-con-2 .grid-con-icon {
	background: rgb(100, 213, 114);
}

.grid-con-2 .grid-num {
	color: rgb(45, 140, 240);
}

.grid-con-3 .grid-con-icon {
	background: rgb(242, 94, 67);
}

.grid-con-3 .grid-num {
	color: rgb(242, 94, 67);
}

.user-info {
	display: flex;
	align-items: center;
	padding-bottom: 20px;
	border-bottom: 2px solid #ccc;
	margin-bottom: 20px;
}

.user-avator {
	width: 120px;
	height: 120px;
	border-radius: 50%;
}

.user-info-cont {
	padding-left: 50px;
	flex: 1;
	font-size: 14px;
	color: #999;
}

.user-info-cont div:first-child {
	font-size: 30px;
	color: #222;
}

.user-info-list {
	font-size: 14px;
	color: #999;
	line-height: 25px;
}

.user-info-list span {
	margin-left: 1px;
}

.mgb20 {
	margin-bottom: 20px;
}

.todo-item {
	font-size: 14px;
}

.todo-item-del {
	text-decoration: line-through;
	color: #999;
}

.schart {
	width: 100%;
	height: 300px;
}
</style>
