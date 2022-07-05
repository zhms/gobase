<template>
	<div class="container">
		<!-- 筛选 -->
		<div>
			<el-form :inline="true" :model="filters">
				<el-form-item label="账号:">
					<el-input v-model="filters.Account" style="width: 150px" :clearable="true"></el-input>
				</el-form-item>
				<el-form-item label="运营商:" v-show="zong">
					<el-select v-model="filters.SellerId" placeholder="运营商" style="width: 130px" @change="handleSelectSeller">
						<el-option v-for="item in seller" :key="item.SellerId" :label="item.SellerName" :value="item.SellerId"> </el-option>
					</el-select>
				</el-form-item>
				<el-form-item>
					<el-button type="primary" icon="el-icon-refresh" class="mr10" @click="handleQuery">查询</el-button>
					<el-button type="primary" icon="el-icon-plus" class="mr10" v-show="auth2('增')" @click="handleAdd">添加</el-button>
				</el-form-item>
			</el-form>
		</div>
		<!--表-->
		<div>
			<el-table :data="table_data" border max-height="620px" class="table" :cell-style="{ padding: '0' }">
				<el-table-column align="center" prop="Id" label="id" width="80"></el-table-column>
				<el-table-column align="center" prop="Account" label="账号" width="100"></el-table-column>
				<el-table-column align="center" prop="SellerName" label="运营商" width="130"></el-table-column>
				<el-table-column align="center" prop="RoleSellerName" label="角色运营商" width="130"></el-table-column>
				<el-table-column align="center" prop="RoleName" label="角色" width="150"></el-table-column>
				<el-table-column align="center" label="状态" width="100">
					<template slot-scope="scope">
						<el-link :underline="false" type="primary" v-if="scope.row.State == 1">启用</el-link>
						<el-link :underline="false" type="danger" v-if="scope.row.State == 2">禁用</el-link>
					</template>
				</el-table-column>
				<el-table-column align="center" prop="LoginTime" label="登录时间" width="160"></el-table-column>
				<el-table-column align="center" prop="LoginIp" label="ip" width="120"></el-table-column>
				<el-table-column align="center" prop="LoginCount" label="登录次数" width="100"></el-table-column>
				<el-table-column align="center" prop="Remark" label="备注" width="200"></el-table-column>
				<el-table-column label="操作">
					<template slot-scope="scope">
						<el-button type="text" icon="el-icon-eleme" class="primary" v-show="auth2('改')" @click="handleChangeGoogle(scope.$index)">谷歌</el-button>
						<el-button type="text" icon="el-icon-edit" v-show="handleShowEdit(scope.$index)" @click="handleModify(scope.$index)">编辑</el-button>
						<el-button type="text" icon="el-icon-delete" class="red" v-show="auth2('删')" @click="handleDel(scope.$index)">删除</el-button>
					</template>
				</el-table-column>
			</el-table>
			<div class="pagination">
				<el-pagination style="margin-top: 5px" background layout="total, prev, pager, next, jumper" :hide-on-single-page="true" :total="total" @current-change="handleQuery" :page-size="pagesize"></el-pagination>
			</div>
		</div>
		<!--对话框-->
		<div>
			<el-dialog :title="dialog.title" :visible.sync="dialog.show" width="415px" center>
				<el-form :inline="true" label-width="100px">
					<el-form-item label="账号:">
						<el-input v-model="dialog.data.Account" :disabled="dialog.type == 'modify'"></el-input>
					</el-form-item>
					<el-form-item label="运营商:" v-show="zong">
						<el-select v-model="dialog.data.SellerId" placeholder="运营商" style="width: 150px" :disabled="dialog.type == 'modify'">
							<el-option v-for="item in seller_noall" :key="item.SellerId" :label="item.SellerName" :value="item.SellerId"> </el-option>
						</el-select>
					</el-form-item>
					<el-form-item label="密码:">
						<el-input v-model="dialog.data.Password" show-password style="width: 200px"></el-input>
					</el-form-item>
					<el-form-item label="角色运营商:" v-show="zong">
						<el-select v-model="dialog.data.RoleSellerId" placeholder="运营商" style="width: 150px" @change="handleSelectRoleSeller">
							<el-option v-for="item in seller_noall" :key="item.SellerId" :label="item.SellerName" :value="item.SellerId"> </el-option>
						</el-select>
					</el-form-item>
					<el-form-item label="角色:">
						<el-select v-model="dialog.data.RoleName" placeholder="请选择" style="width: 200px">
							<el-option v-for="item in dialog.options.role" :key="item.RoleName" :label="item.RoleName" :value="item.RoleName"></el-option>
						</el-select>
					</el-form-item>
				</el-form>
				<el-form :inline="true" label-width="100px">
					<el-form-item label="选项:">
						<el-checkbox border label="禁用" v-model="dialog.data.State" :true-label="2" :false-label="1"></el-checkbox>
					</el-form-item>
				</el-form>
				<el-form :inline="true" label-width="100px">
					<el-form-item label="备注:">
						<el-input type="textarea" v-model="dialog.data.Remark" :rows="4"></el-input>
					</el-form-item>
				</el-form>
				<span slot="footer" class="dialog-footer">
					<el-button type="primary" @click="handleConfirm">确 定</el-button>
				</span>
			</el-dialog>
		</div>
		<div>
			<el-dialog :title="dialog_google.title" :visible.sync="dialog_google.show" width="350px" center>
				<vueqr :text="dialog_google.url" :size="300"> </vueqr>
			</el-dialog>
		</div>
	</div>
</template>
<script>
import { app } from '@/api/app.js'
import base from '@/api/base.js'
import '@/assets/css/k.css'
import vueqr from 'vue-qr'
export default {
	extends: base,
	data() {
		return {
			filters: {
				Account: null,
			},
			dialog_google: {
				show: false,
				title: '谷歌验证码',
				url: '',
			},

			dialog: {
				options: {
					role: [],
				},
				data: {
					Account: null,
					SellerId: null,
					Password: null,
					RoleSellerId: null,
					RoleName: null,
					State: 1,
					Remark: null,
				},
			},
		}
	},
	components: {
		vueqr,
	},
	created() {
		this.handleQuery(1)
	},
	methods: {
		auth2(o) {
			return app.auth2('系统管理', '账号管理', o)
		},
		handleSelectSeller() {
			app.setCurrentSeller(this.filters.SellerId)
		},
		handleShowEdit(index) {
			return this.auth2('改') && this.table_data[index].RoleName != '运营商超级管理员'
		},
		handleSelectRoleSeller() {
			this.dialog.data.RoleName = null
			app.post('/admin/role/listall', { SellerId: this.dialog.data.RoleSellerId }, (result) => {
				this.dialog.options.role = []
				for (let i = 0; i < result.data.length; i++) {
					this.dialog.options.role.push({ RoleName: result.data[i] })
				}
			})
		},
		handleChangeGoogle(index) {
			if (confirm('确定更换谷歌验证码?')) {
				let data = {
					Account: this.table_data[index].Account,
					SellerId: this.table_data[index].SellerId,
				}
				app.post('/admin/user/google', data, (result) => {
					this.dialog_google.url = result.data
					this.dialog_google.show = true
				})
			}
		},
		handleQuery(page) {
			this.page = page || 1
			if (typeof this.page == 'object') this.page = 1
			var data = {
				page: this.page,
				pagesize: this.pagesize,
				Account: this.filters.Account || '',
				SellerId: parseInt(this.filters.SellerId || 0),
			}
			app.post('/admin/user/list', data, (result) => {
				this.total = result.data.total
				this.table_data = result.data.data
				for (var i = 0; i < this.table_data.length; i++) {
					this.table_data[i].CreateTime = this.$moment(this.table_data[i].CreateTime).format('YYYY-MM-DD hh:mm:ss')
					this.table_data[i].LoginTime = this.$moment(this.table_data[i].LoginTime).format('YYYY-MM-DD hh:mm:ss')
					for (let j = 0; j < this.seller.length; j++) {
						if (this.seller[j].SellerId == this.table_data[i].SellerId) this.table_data[i].SellerName = this.seller[j].SellerName
						if (this.seller[j].SellerId == this.table_data[i].RoleSellerId) this.table_data[i].RoleSellerName = this.seller[j].SellerName
					}
				}
			})
		},
		handleAdd() {
			this.dialog.data.Account = null
			this.dialog.data.SellerId = null
			this.dialog.data.Password = null
			this.dialog.data.RoleSellerId = null
			this.dialog.data.RoleName = null
			this.dialog.data.Remark = null
			this.dialog.data.State = 1
			this.dialog.title = '添加账号'
			this.dialog.type = 'add'
			this.dialog.show = true
		},
		handleModify(index) {
			this.current_row = index
			this.dialog.title = '修改账号'
			this.dialog.data = app.clone(this.table_data[index])
			this.dialog.type = 'modify'
			let v = this.dialog.data.RoleName
			this.handleSelectRoleSeller()
			this.dialog.data.RoleName = v
			this.dialog.show = true
		},
		handleDel(index) {
			if (confirm('确定删除该账号?')) {
				let data = {
					Id: this.table_data[index].Id,
					Account: this.table_data[index].Account,
					SellerId: this.table_data[index].SellerId,
				}
				app.post('/admin/user/delete', data, () => {
					this.table_data.splice(index, 1)
					this.table_data = app.clone(this.table_data)
					this.$message.success('操作成功')
				})
			}
		},
		handleConfirm() {
			if (this.dialog.type == 'modify') {
				let data = {
					Account: this.dialog.data.Account,
					SellerId: this.dialog.data.SellerId,
					Remark: this.dialog.data.Remark,
					RoleSellerId: this.dialog.data.RoleSellerId,
					State: this.dialog.data.State,
					RoleName: this.dialog.data.RoleName,
				}
				if (this.dialog.data.Password && this.dialog.data.Password.length > 0) {
					data.Password = this.$md5(this.dialog.data.Password)
				}
				app.post('/admin/user/modify', data, () => {
					this.dialog.show = false
					this.$message.success('操作成功')
					this.handleQuery(this.page)
				})
			}
			if (this.dialog.type == 'add') {
				let data = {
					Account: this.dialog.data.Account,
					SellerId: this.dialog.data.SellerId,
					Remark: this.dialog.data.Remark,
					RoleSellerId: this.dialog.data.RoleSellerId,
					State: this.dialog.data.State,
					RoleName: this.dialog.data.RoleName,
				}
				if (this.dialog.data.Password && this.dialog.data.Password.length > 0) {
					data.Password = this.$md5(this.dialog.data.Password)
				}
				app.post('/admin/user/add', data, () => {
					this.dialog.show = false
					this.$message.success('操作成功')
					this.handleQuery(this.page)
				})
			}
		},
	},
}
</script>
