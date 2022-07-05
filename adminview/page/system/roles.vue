<template>
	<div class="container">
		<!-- 筛选 -->
		<div>
			<el-form :inline="true" :model="filters">
				<el-form-item label="运营商:" v-show="zong">
					<el-select v-model="filters.SellerId" placeholder="运营商" style="width: 130px" @change="handleSelectSeller">
						<el-option v-for="item in seller" :key="item.SellerId" :label="item.SellerName" :value="item.SellerId"> </el-option>
					</el-select>
				</el-form-item>
				<el-form-item>
					<el-button type="primary" icon="el-icon-plus" class="mr10" v-show="auth2('增')" @click="handleAdd">添加</el-button>
					<el-button type="primary" icon="el-icon-refresh" v-on:click="handleQuery">查询</el-button>
				</el-form-item>
			</el-form>
		</div>
		<!--表-->
		<div>
			<el-table :data="table_data" border max-height="620px" class="table" :cell-style="{ padding: '0' }">
				<el-table-column align="center" prop="SellerName" label="运营商" width="200"></el-table-column>
				<el-table-column align="center" prop="RoleName" label="角色名" width="200"></el-table-column>
				<el-table-column align="center" prop="ParentSellerName" label="上级角色运营商" width="200"></el-table-column>
				<el-table-column align="center" prop="Parent" label="上级角色" width="200"></el-table-column>
				<el-table-column label="操作">
					<template slot-scope="scope">
						<el-button type="text" icon="el-icon-edit" v-show="auth2('改')" @click="handleModify(scope.$index)">编辑</el-button>
						<el-button type="text" icon="el-icon-delete" class="red" v-show="auth2('删')" @click="handleDel(scope.$index)">删除</el-button>
					</template>
				</el-table-column>
			</el-table>
		</div>
		<!-- 对话框 -->
		<div>
			<el-dialog :title="dialog.title" :visible.sync="dialog.show" width="400px" center>
				<div>
					<el-form :inline="true" :model="filters">
						<el-form-item label="角色名:" style="margin-left: 15px">
							<el-input v-model="dialog.data.RoleName" :disabled="dialog.type == 'modify'"></el-input>
						</el-form-item>
						<el-form-item label="运营商:" v-show="zong" style="margin-left: 15px">
							<el-select v-model="dialog.data.SellerId" placeholder="请选择" style="" :disabled="dialog.type == 'modify'" @change="handleDialogSelectSellerId">
								<el-option v-for="item in seller_noall" :key="item.SellerId" :label="item.SellerName" :value="item.SellerId"> </el-option>
							</el-select>
						</el-form-item>
						<el-form-item label="上级角色:" v-show="zong">
							<el-select v-model="dialog.data.Parent" placeholder="请选择" style="" :disabled="dialog.type == 'modify'" @change="handleDialogSelectRole">
								<el-option v-for="item in dialog.options.Parents" :key="item.RoleName" :label="item.RoleName" :value="item.RoleName"> </el-option>
							</el-select>
						</el-form-item>
					</el-form>
				</div>
				<el-tree :default-checked-keys="dialog.selected" node-key="path" ref="authtree" :props="dialog.tree_props" show-checkbox v-show="dialog.show_tree"> </el-tree>
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
				selected: [],
				show_tree: false,
				tree_props: {
					label: 'name',
					children: 'children',
				},
				data: {
					RoleName: null,
					Parent: null,
					SellerId: null,
				},
				options: {
					Parents: [],
				},
			},
		}
	},
	components: {},
	computed: {},
	created() {
		this.handleQuery()
	},
	methods: {
		auth2(o) {
			return app.auth2('系统管理', '角色管理', o)
		},
		handleSelectSeller() {
			app.setCurrentSeller(this.filters.SellerId)
		},
		handleDialogSelectSellerId() {
			this.dialog.show_tree = false
			this.dialog.data.Parent = null
			app.post('/admin/role/listall', { SellerId: this.dialog.data.SellerId }, (result) => {
				let parents = []
				for (let i = 0; i < result.data.length; i++) {
					parents.push({ RoleName: result.data[i] })
				}
				this.dialog.options.Parents = parents
			})
		},
		handleDialogSelectRole() {
			app.post('/admin/role/roledata', { SellerId: this.dialog.data.SellerId, RoleName: this.dialog.data.Parent }, (result) => {
				this.dialog.parentroledata = JSON.parse(result.data.RoleData)
				this.dialog.superroledata = JSON.parse(result.data.SuperRoleData)
				this.dialog.show_tree = true
				this.dialog.roledata = {}
				let treedata = this.getTreeData()
				this.$refs.authtree.root.setData(treedata.menu)
				this.dialog.selected = treedata.selected
			})
		},
		handleQuery(page) {
			this.page = page || 1
			if (typeof this.page == 'object') this.page = 1
			var data = {
				SellerId: parseInt(this.filters.SellerId || 0),
				page: this.page,
				pagesize: this.pagesize,
			}
			app.post('/admin/role/list', data, (result) => {
				this.table_data = result.data.data
				this.total = result.data.total
				for (var i = 0; i < this.table_data.length; i++) {
					for (let j = 0; j < this.seller.length; j++) {
						if (this.seller[j].SellerId == this.table_data[i].SellerId) this.table_data[i].SellerName = this.seller[j].SellerName
						if (this.seller[j].SellerId == this.table_data[i].ParentSellerId) this.table_data[i].ParentSellerName = this.seller[j].SellerName
					}
				}
			})
		},
		handleAdd() {
			this.dialog.title = `添加角色`
			this.dialog.type = 'add'
			this.dialog.data.SellerId = null
			this.dialog.data.ParentSellerId = null
			this.dialog.data.Parent = null
			this.dialog.data.RoleName = null
			this.dialog.show_tree = false
			this.dialog.show = true
		},
		handleModify(index) {
			this.current_row = index
			if (this.table_data[this.current_row].Parent == 'god') {
				this.$message.error('该角色不可修改')
				return
			}
			this.dialog.data = app.clone(this.table_data[this.current_row])
			this.dialog.title = `修改角色`
			this.dialog.type = 'modify'
			this.dialog.show = true
			setTimeout(() => {
				this.$refs.authtree.root.setData([])
			}, 10)
			app.post('/admin/role/roledata', { SellerId: this.dialog.data.SellerId, RoleName: this.dialog.data.Parent }, (resulta) => {
				this.dialog.parentroledata = JSON.parse(resulta.data.RoleData)
				this.dialog.superroledata = JSON.parse(resulta.data.SuperRoleData)
				app.post('/admin/role/roledata', { SellerId: this.dialog.data.SellerId, RoleName: this.dialog.data.RoleName }, (resultb) => {
					this.dialog.roledata = JSON.parse(resultb.data.RoleData)
					this.dialog.show_tree = true
					let treedata = this.getTreeData()
					this.$refs.authtree.root.setData(treedata.menu)
					this.dialog.selected = treedata.selected
				})
			})
		},
		handleDel(index) {
			this.current_row = index
			this.dialog.data = this.table_data[this.current_row]
			if (this.dialog.data.Parent == 'god') {
				this.$message.error('该角色不可删除')
				return
			}
			if (confirm('确定删除该角色?')) {
				let data = {
					SellerId: this.dialog.data.SellerId,
					RoleName: this.dialog.data.RoleName,
				}
				app.post('/admin/role/delete', data, () => {
					this.dialog.show = false
					this.$message.success('操作成功')
					this.handleQuery(this.page)
				})
			}
		},
		handleConfirm() {
			if (!this.dialog.data.RoleName) return this.$message.error('请填写角色名')
			if (!this.dialog.data.Parent) return this.$message.error('请选择上级角色')
			if (!this.dialog.data.SellerId) return this.$message.error('请选择运营商')
			let setdisable = (node) => {
				for (let n in node) {
					if (typeof node[n] == 'object') {
						setdisable(node[n])
					} else {
						node[n] = 0
					}
				}
			}
			let newroledata = app.clone(this.dialog.superroledata)
			setdisable(newroledata)
			let selected = this.$refs.authtree.getCheckedNodes()
			for (let i = 0; i < selected.length; i++) {
				if (!selected[i].leaf) continue
				let path = selected[i].path.split('.')
				let pn = newroledata
				for (let i = 0; i < path.length - 1; i++) {
					pn = pn[path[i]]
				}
				pn[path[path.length - 1]] = 1
			}
			if (this.dialog.type == 'modify') {
				let data = {
					SellerId: this.dialog.data.SellerId,
					RoleName: this.dialog.data.RoleName,
					RoleData: JSON.stringify(newroledata),
				}
				app.post('/admin/role/modify', data, () => {
					this.dialog.show = false
					this.$message.success('操作成功')
					this.handleQuery(this.page)
				})
			}
			if (this.dialog.type == 'add') {
				let data = {
					ParentSellerId: this.dialog.data.ParentSellerId,
					Parent: this.dialog.data.Parent,
					SellerId: this.dialog.data.SellerId,
					RoleName: this.dialog.data.RoleName,
					RoleData: JSON.stringify(newroledata),
				}
				app.post('/admin/role/add', data, () => {
					this.dialog.show = false
					this.$message.success('操作成功')
					this.handleQuery(this.page)
				})
			}
		},
		getTreeData() {
			let setdisable = (node) => {
				for (let n in node) {
					if (typeof node[n] == 'object') {
						setdisable(node[n])
					} else {
						node[n] = 0
					}
				}
			}
			setdisable(this.dialog.superroledata)
			let setenable = (parent, node) => {
				for (let n in node) {
					if (typeof node[n] == 'object') {
						let p = parent + `.${n}`
						setenable(p, node[n])
					} else {
						if (node[n] == 1) {
							let p = parent.split('.')
							let pn = this.dialog.superroledata
							for (let j = 0; j < p.length; j++) {
								pn = pn[p[j]]
							}
							pn[n] = 1
						}
					}
				}
			}
			for (let n in this.dialog.parentroledata) {
				let parent = `${n}`
				setenable(parent, this.dialog.parentroledata[n])
			}
			let menu = []
			let submenu = (node, root) => {
				for (let n in root) {
					if (typeof root[n] == 'object') {
						let subnode = {
							path: node.path + '.' + n,
							name: n,
							children: [],
						}
						node.children.push(subnode)
						submenu(subnode, root[n])
					} else {
						let path = node.path + '.' + n
						let p = path.split('.')
						let pr = this.dialog.parentroledata
						for (let i = 0; i < p.length; i++) {
							pr = pr[p[i]]
						}
						if (pr == 1) {
							let subnode = {
								path: path,
								name: n,
								leaf: true,
							}
							node.children.push(subnode)
						}
					}
				}
			}
			for (let n in this.dialog.superroledata) {
				let node = {
					path: n,
					name: n,
					children: [],
				}
				menu.push(node)
				submenu(node, this.dialog.superroledata[n])
			}
			let selected = []
			let getselected = (parent, node) => {
				for (let n in node) {
					if (typeof node[n] == 'object') {
						let p = parent + `.${n}`
						getselected(p, node[n])
					} else {
						if (node[n] == 1) {
							selected.push(`${parent}.${n}`)
						}
					}
				}
			}
			for (let n in this.dialog.roledata) {
				let parent = `${n}`
				getselected(parent, this.dialog.roledata[n])
			}
			for (let i = 0; i < menu.length; i++) {
				if (!menu[i].children) continue
				for (let j = 0; j < menu[i].children.length; j++) {
					if (!menu[i].children[j].children) continue
					for (let k = 0; k < menu[i].children[j].children.length; k++) {
						if (!menu[i].children[j].children[k].children) continue
						if (menu[i].children[j].children[k].children.length == 0) {
							menu[i].children[j].children.splice(k, 1)
							k--
						}
					}
				}
			}
			for (let i = 0; i < menu.length; i++) {
				if (!menu[i].children) continue
				for (let j = 0; j < menu[i].children.length; j++) {
					if (!menu[i].children[j].children) continue
					if (menu[i].children[j].children.length == 0) {
						menu[i].children.splice(j, 1)
						j--
					}
				}
			}
			for (let i = 0; i < menu.length; i++) {
				if (!menu[i].children) continue
				if (menu[i].children.length == 0) {
					menu.splice(i, 1)
					i--
				}
			}
			return { menu, selected }
		},
	},
}
</script>
