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
				<el-table-column align="center" prop="Remark" label="备注" width="300"></el-table-column>
				<el-table-column align="center" prop="CreateTime" label="创建时间" width="150"></el-table-column>
				<el-table-column label="操作">
					<template slot-scope="scope">
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
		}
	},
	created() {
		this.handleQuery()
	},
	methods: {
		auth2(o) {
			return app.auth2('系统管理', '运营商管理', o)
		},
		handleQuery(page) {
			this.page = page || 1
			if (typeof this.page == 'object') this.page = 1
			var data = {
				page: this.page,
				pagesize: this.pagesize,
			}
			app.post('/seller/list', data, (result) => {
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
				this.dialog.data.IgnoreSeller = true
				app.post('/seller/modify', this.dialog.data, () => {
					this.dialog.show = false
					this.handleQuery(this.page)
				})
			}
			if (this.dialog.type == 'add') {
				app.post('/seller/add', this.dialog.data, () => {
					this.dialog.show = false
					this.handleQuery(this.page)
				})
			}
		},
		handleDelete(index) {
			if (confirm('确定删除该运营商?')) {
				app.post('/seller/delete', { SellerId: this.table_data[index].SellerId, IgnoreSeller: true }, () => {
					this.handleQuery(this.page)
				})
			}
		},
	},
}
</script>
