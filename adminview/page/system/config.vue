<template>
	<div class="container">
		<div class="handle-box">
			<el-button type="primary" class="mr10" icon="el-icon-refresh" @click="handleQuery()">刷新</el-button>
			<el-button type="primary" icon="el-icon-plus" class="mr10" @click="handleAdd()">添加</el-button>
		</div>
		<div>
			<el-table :data="table_data" border class="table" max-height="700px" :cell-style="{ padding: '0' }">
				<el-table-column align="center" prop="SellerId" label="运营商id" width="100"></el-table-column>
				<el-table-column align="center" prop="SellerName" label="运营商名称" width="150"></el-table-column>
				<el-table-column align="center" label="状态" width="100">
					<template slot-scope="scope">
						<el-link :underline="false" type="primary" v-if="scope.row.State == 1">启用</el-link>
						<el-link :underline="false" type="danger" v-if="scope.row.State == 2">禁用</el-link>
					</template>
				</el-table-column>
				<el-table-column align="center" prop="Remark" label="备注" width="200"></el-table-column>
				<el-table-column align="center" prop="CreateTime" label="创建时间" width="150"></el-table-column>
				<el-table-column label="操作">
					<template slot-scope="scope">
						<el-button type="text" icon="el-icon-edit" @click="handleChangeKey(scope.$index)">更换秘钥</el-button>
						<el-button type="text" icon="el-icon-edit" @click="handleShowKey(scope.$index)">查看秘钥</el-button>
						<el-button type="text" icon="el-icon-edit" @click="handleModify(scope.$index)">修改</el-button>
						<el-button type="text" icon="el-icon-delete" class="red" @click="handleDelete(scope.$index)">删除</el-button>
					</template>
				</el-table-column>
			</el-table>
			<div class="pagination">
				<el-pagination style="margin-top: 5px" background layout="total, prev, pager, next, jumper" :hide-on-single-page="true" :total="total" @current-change="handleQuery" :page-size="pagesize"></el-pagination>
			</div>
		</div>
		<div>
			<el-dialog :title="dialog.title" :visible.sync="dialog.show" width="400px" center>
				<el-form :inline="true" label-width="100px">
					<el-form-item label="运营商名称:">
						<el-input v-model="dialog.data.SellerName"></el-input>
					</el-form-item>
					<el-form-item label="选项:">
						<el-checkbox border label="禁用" v-model="dialog.data.State" :true-label="2" :false-label="1"></el-checkbox>
					</el-form-item>
					<el-form-item label="注释:">
						<el-input type="textarea" v-model="dialog.data.Remark" :rows="4"></el-input>
					</el-form-item>
				</el-form>
				<span slot="footer" class="dialog-footer">
					<el-button type="primary" @click="handleConfirm">确 定</el-button>
				</span>
			</el-dialog>
		</div>
		<div>
			<el-dialog title="秘钥" :visible.sync="dlgkey" width="800px" center>
				<el-form :model="filters" width="800px">
					<el-form-item label="服务器业务公钥:">
						<el-input type="textarea" v-model="keydata.ApiPublicKey" style="width: 600px" :disabled="true" :rows="10"></el-input>
					</el-form-item>
					<!-- <el-form-item label="服务器风控公钥:">
						<el-input type="textarea" v-model="keydata.ApiRiskPublicKey" style="width: 600px" :disabled="true" :rows="10"></el-input>
					</el-form-item> -->
					<el-form-item label="第三方业务公钥:">
						<el-input type="textarea" v-model="keydata.ApiThirdPublicKey" style="width: 600px" :rows="10"></el-input>
					</el-form-item>
					<!-- <el-form-item label="第三方风控公钥:">
						<el-input type="textarea" v-model="keydata.ApiThirdRiskPublicKey" style="width: 600px" :rows="10"></el-input>
					</el-form-item> -->
				</el-form>
				<span slot="footer" class="dialog-footer">
					<el-button type="primary" @click="handleSetKey">确 定</el-button>
				</span>
			</el-dialog>
		</div>
	</div>
</template>
<script>
import { app } from '@/api/app.js'
import '@/assets/css/k.css'
import base from '@/api/base.js'
export default {
	extends: base,
	data() {
		return {
			dialog: {
				data: {
					SellerName: null,
					State: 1,
					Remark: null,
				},
			},
			dlgkey: false,
			textarea: '',
			keydata: {},
		}
	},
	created() {
		this.handleQuery()
	},
	methods: {
		auth2(o) {
			return app.auth2('系统管理', '运营商管理', o)
		},
		handleSetKey() {
			let data = {
				SellerId: this.table_data[this.current_row].SellerId,
				ApiThirdPublicKey: this.keydata.ApiThirdPublicKey,
				ApiThirdRiskPublicKey: this.keydata.ApiThirdRiskPublicKey,
			}
			app.post('/admin/seller/set_key', data, (result) => {
				this.dlgkey = false
			})
		},
		handleQuery(page) {
			this.page = page || 1
			if (typeof this.page == 'object') this.page = 1
			var data = {
				page: this.page,
				pagesize: this.pagesize,
			}
			app.post('/admin/seller/list', data, (result) => {
				this.table_data = result.data.data
				this.total = result.data.total
				for (var i = 0; i < this.table_data.length; i++) {
					this.table_data[i].CreateTime = this.$moment(this.table_data[i].CreateTime).format('YYYY-MM-DD hh:mm:ss')
				}
			})
		},
		handleAdd() {
			this.dialog.type = 'add'
			this.dialog.data.SellerName = null
			this.dialog.data.State = 1
			this.dialog.data.Remark = null
			this.dialog.title = '添加运营商'
			this.dialog.show = true
		},
		handleModify(index) {
			this.dialog.type = 'modify'
			this.dialog.title = '修改运营商'
			this.current_row = index
			this.dialog.data = app.clone(this.table_data[index])
			this.dialog.show = true
		},
		handleConfirm() {
			if (this.dialog.type == 'modify') {
				app.post('/admin/seller/modify', this.dialog.data, () => {
					this.dialog.show = false
					app.flushSeller()
					this.handleQuery(this.page)
				})
			}
			if (this.dialog.type == 'add') {
				app.post('/admin/seller/add', this.dialog.data, () => {
					this.dialog.show = false
					app.flushSeller()
					this.handleQuery(this.page)
				})
			}
		},
		handleDelete(index) {
			if (confirm('确定删除该运营商?')) {
				app.post('/admin/seller/delete', { SellerId: this.table_data[index].SellerId }, () => {
					app.flushSeller()
					this.handleQuery(this.page)
				})
			}
		},
		handleChangeKey(index) {
			if (confirm('确定更换秘钥?')) {
				app.post('/admin/seller/change_key', { SellerId: this.table_data[index].SellerId }, () => {
					this.$message.success('操作成功')
				})
			}
		},
		handleShowKey(index) {
			this.current_row = index
			app.post('/admin/seller/get_key', { SellerId: this.table_data[index].SellerId }, (data) => {
				this.keydata = data.data.data
				this.dlgkey = true
				console.log(this.keydata)
			})
		},
	},
}
</script>
