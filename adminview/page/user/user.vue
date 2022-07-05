<template>
	<div class="container">
		<!-- 筛选 -->
		<div>
			<el-form :inline="true" :model="filters">
				<el-form-item label="账号:">
					<el-input v-model="filters.Account" style="width: 120px" :clearable="true" placeholder="账号"></el-input>
				</el-form-item>
				<el-form-item label="Id:">
					<el-input v-model="filters.UserId" style="width: 120px" :clearable="true" placeholder="玩家id"></el-input>
				</el-form-item>
				<el-form-item>
					<el-button type="primary" icon="el-icon-refresh" class="mr10" @click="handleQuery">查询</el-button>
				</el-form-item>
			</el-form>
		</div>
		<!--表-->
		<div>
			<el-table :data="table_data" border max-height="620px" class="table" :cell-style="{ padding: '0' }">
				<el-table-column align="center" prop="Id" label="id" width="80"></el-table-column>
				<el-table-column align="center" prop="UserId" label="玩家id" width="100"></el-table-column>
				<el-table-column align="center" prop="Account" label="账号" width="220"></el-table-column>
				<el-table-column align="center" label="昵称" width="130">
					<template slot-scope="scope">
						<span style="cursor: pointer; color: rgb(64, 158, 255)" @click="handleLook(scope.$index)">{{ table_data[scope.$index].NickName }}</span>
					</template>
				</el-table-column>
				<el-table-column align="center" prop="Score" label="金币" width="120"></el-table-column>
				<el-table-column align="center" prop="BankScore" label="银行金币" width="120"></el-table-column>
				<el-table-column align="center" prop="Agent" label="代理" width="100"></el-table-column>
				<el-table-column align="center" prop="RegIp" label="注册Ip" width="120"></el-table-column>
				<el-table-column align="center" prop="RegOs" label="平台" width="100"></el-table-column>
				<el-table-column align="center" prop="RegisterTime" label="注册时间" width="160"></el-table-column>
			</el-table>
		</div>
		<div class="pagination">
			<el-pagination style="margin-top: 5px" background layout="total, prev, pager, next, jumper" :hide-on-single-page="true" :total="total" @current-change="handleQuery" :page-size="pagesize"></el-pagination>
		</div>
		<!--对话框-->
		<div>
			<el-dialog :title="dialog_title" :visible.sync="dialog" width="1400px" top="2vh">
				<el-tabs v-model="activeName" type="card" center width="500px" @tab-click="handleSelectTab">
					<el-tab-pane label="基础信息" name="0">
						<div style="height: 70vh; overflow: auto">
							<el-form>
								<el-form :inline="true" label-width="130px" style="margin-left: 30px">
									<el-form-item label=" "></el-form-item>
								</el-form>
								<el-form :inline="true" label-width="130px" style="margin-left: 50px">
									<el-form-item label="玩家ID:" style="width: 200px; margin-top: -15px">{{ dialog_data.UserId }}</el-form-item>
									<el-form-item label="代理:" style="width: 200px; margin-top: -15px">{{ dialog_data.Agent }}</el-form-item>
									<el-form-item label="注册来源:" style="width: 200px; margin-top: -15px">{{ dialog_data.RegOs }}</el-form-item>
								</el-form>
								<el-form :inline="true" label-width="130px" style="margin-left: 50px">
									<el-form-item label="运营商:" style="width: 200px; margin-top: -15px">{{ dialog_data.SellerId }}</el-form-item>
									<el-form-item label="状态:" v-bind:class="dialog_data.IsDisabled ? 'red' : ''" style="width: 200px; margin-top: -15px">{{ dialog_data.IsDisabled == 1 ? '禁用' : '正常' }}</el-form-item>
									<el-form-item label="登录次数:" style="width: 300px; margin-top: -15px">{{ dialog_data.LoginCount }}</el-form-item>
								</el-form>
								<el-form :inline="true" label-width="130px" style="margin-left: 50px">
									<el-form-item label="账号:" style="width: 250px; margin-top: -15px">{{ dialog_data.Account }}</el-form-item>
									<el-form-item label="测试:" v-bind:class="dialog_data.IsTester ? 'red' : ''" style="width: 200px; margin-top: -15px">{{ dialog_data.IsTester == 0 ? '否' : '是' }}</el-form-item>
								</el-form>
								<el-form :inline="true" label-width="130px" style="margin-left: 50px">
									<el-form-item label="昵称:" style="width: 200px; margin-top: -15px">{{ dialog_data.NickName }}</el-form-item>
									<el-form-item label="超管:" v-bind:class="dialog_data.IsAdmin ? 'red' : ''" style="width: 200px; margin-top: -15px">{{ dialog_data.IsAdmin == 0 ? '否' : '是' }}</el-form-item>
								</el-form>
								<el-form :inline="true" label-width="130px" style="margin-left: 50px">
									<el-form-item label="金币:" style="width: 200px; margin-top: -15px">{{ dialog_data.Score }}</el-form-item>
								</el-form>
								<el-form :inline="true" label-width="130px" style="margin-left: 50px">
									<el-form-item label="银行金币:" style="width: 200px; margin-top: -15px">{{ dialog_data.BankScore }}</el-form-item>
								</el-form>
								<el-form :inline="true" label-width="130px" style="margin-left: 50px">
									<el-form-item label="备注:" style="width: 1000px; margin-top: -15px"> {{ dialog_data.Note }} </el-form-item>
								</el-form>
								<el-form :inline="true" style="margin-left: 100px">
									<el-form-item>
										<el-button type="primary" @click="handleShowNoteDialog()" v-show="auth('修改备注')">修改备注</el-button>
										<el-button type="primary" @click="handleModifyField('IsTester')" v-show="auth('设置取消测试')">{{ this.dialog_data.IsTester == 0 ? '设置测试' : '取消测试' }}</el-button>
										<el-button
											type="primary"
											v-show="auth('赠送金币')"
											@click="
												add_score = null
												add_score_reason = null
												dialog_score = true
											"
											>赠送金币</el-button
										>
										<el-button type="primary" @click="handleModifyField('IsDisabled')" v-show="auth('冻结解冻')">{{ dialog_data.IsDisabled == 0 ? '冻结' : '解冻' }}</el-button>
										<el-button type="primary" @click="handleModifyField('IsAdmin')" v-show="auth('设置取消超管')">{{ dialog_data.IsAdmin == 0 ? '设置超管' : '取消超管' }}</el-button>
										<el-button type="primary" v-show="false">绑定代理</el-button>
									</el-form-item>
								</el-form>
							</el-form>
						</div>
					</el-tab-pane>
					<el-tab-pane label="登陆注册" name="1">
						<div style="height: 70vh; overflow: auto">
							<el-form>
								<el-form :inline="true" label-width="130px" style="margin-left: 30px">
									<el-form-item label=" " style="margin-left: 165px"></el-form-item>
									<el-form-item label=" " style="margin-left: 265px"></el-form-item>
								</el-form>
								<el-form :inline="true" label-width="130px" style="margin-left: 50px">
									<el-form-item label="登录IP:" style="width: 400px; margin-top: -15px">{{ dialog_data.LastLoginIp }}</el-form-item>
									<el-form-item label="注册IP:" style="width: 400px; margin-top: -15px">{{ dialog_data.RegIp }}</el-form-item>
								</el-form>
								<el-form :inline="true" label-width="130px" style="margin-left: 50px">
									<el-form-item label="登录设备:" style="width: 400px; margin-top: -15px">{{ dialog_data.LastLoginOs }}</el-form-item>
									<el-form-item label="注册设备:" style="width: 400px; margin-top: -15px">{{ dialog_data.RegOs }}</el-form-item>
								</el-form>
								<el-form :inline="true" label-width="130px" style="margin-left: 50px">
									<el-form-item label="登录时间:" style="width: 400px; margin-top: -15px">{{ dialog_data.LastLoginTime }}</el-form-item>
									<el-form-item label="注册时间:" style="width: 400px; margin-top: -15px">{{ dialog_data.RegTime }}</el-form-item>
								</el-form>
								<el-form :inline="true" label-width="130px" style="margin-left: 50px">
									<el-form-item label="登录机器:" style="width: 400px; margin-top: -15px">{{ dialog_data.LastLoginMachineId }}</el-form-item>
									<el-form-item label="注册机器:" style="width: 400px; margin-top: -15px">{{ dialog_data.RegMachineId }}</el-form-item>
								</el-form>
							</el-form>
						</div>
					</el-tab-pane>
					<el-tab-pane label="充值记录" name="4">充值记录</el-tab-pane>
					<el-tab-pane label="兑换记录" name="5">兑换记录</el-tab-pane>
					<el-tab-pane label="财富信息" name="6">财富信息</el-tab-pane>
					<el-tab-pane label="游戏信息" name="7">游戏信息</el-tab-pane>
				</el-tabs>
				<!-- 修改备注对话框 -->
				<el-dialog width="415px" center title="修改备注" :visible.sync="dialog_note" append-to-body>
					<el-input type="area" v-model="dialog_data_copy.Note"></el-input>
					<div slot="footer" class="dialog-footer">
						<el-button type="primary" @click="handleModifyField('Note')">确定</el-button>
					</div>
				</el-dialog>
				<!-- 赠送对话框 -->
				<el-dialog width="415px" center title="赠送金币" :visible.sync="dialog_score" append-to-body>
					<el-form :inline="true">
						<el-form-item label="金币数量:" label-width="110px">
							<el-input v-model="add_score"></el-input>
						</el-form-item>
						<el-form-item label="赠送原因:" label-width="110px">
							<el-input v-model="add_score_reason"></el-input>
						</el-form-item>
					</el-form>
					<div slot="footer" class="dialog-footer">
						<el-button type="primary" @click="handleAddScore()">确定</el-button>
					</div>
				</el-dialog>
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
			filters: {
				Account: null,
				UserId: null,
			},
			activeName: '0',
			dialog: false,
			dialog_title: '',
			dialog_data: {},
			current_row: null,
			dialog_note: false,
			dialog_data_copy: {},
			dialog_score: false,
			add_score: null,
			add_score_reason: null,
			current_tab: null,
		}
	},
	components: {},
	computed: {},
	created() {
		this.handleQuery()
	},
	methods: {
		auth(o) {
			return app.auth2('玩家管理', '账号管理', o)
		},
		tagChange() {
			this.dialog_data_copy = app.clone(this.dialog_data)
		},
		handleQuery(page) {
			this.page = page || 1
			if (typeof this.page != 'number') this.page = 1
			var data = {
				page: this.page,
				pagesize: this.pagesize,
				SellerId: app.getInfo().SellerId,
				Account: this.filters.Account,
				UserId: parseInt(this.filters.UserId),
			}
			app.post('/user/list', data, (result) => {
				this.table_data = result.data.data
				this.total = result.data.total
				for (var i = 0; i < this.table_data.length; i++) {
					for (let j = 0; j < this.seller.length; j++) {
						if (this.seller[j].SellerId == this.table_data[i].SellerId) this.table_data[i].SellerName = this.seller[j].SellerName
					}
					this.table_data[i].RegisterTime = this.$moment(this.table_data[i].RegisterTime).format('YYYY-MM-DD hh:mm:ss')
					// this.table_data[i].RegTime = this.$moment(this.table_data[i].RegTime).format('YYYY-MM-DD hh:mm:ss')
				}
				console.log(this.seller)
				console.log(this.table_data)
			})
		},
		handleLook(index) {
			this.current_row = index
			this.dialog_data = app.clone(this.table_data[index])
			this.dialog_data_copy = app.clone(this.dialog_data)
			this.activeName = '0'
			this.current_tab = null
			this.dialog = true
		},
		handleShowNoteDialog() {
			this.dialog_data_copy = app.clone(this.dialog_data)
			this.dialog_note = true
		},
		handleModifyField(field) {
			var data = {
				UserId: this.dialog_data.UserId,
				field: field,
				value: this.dialog_data_copy[field],
			}
			if (field == 'IsDisabled' || field == 'IsAdmin' || field == 'IsTester') {
				if (data.value) {
					data.value = 0
				} else {
					data.value = 1
				}
			}
			app.post('/user/account/modify', data, () => {
				this.dialog_data[field] = data.value
				this.table_data[this.current_row][field] = data.value
				this.dialog_data_copy = app.clone(this.dialog_data)
				this.dialog_note = false
				this.$message.success('操作成功')
			})
		},
		handleAddScore() {
			var data = {
				UserId: this.dialog_data.UserId,
				Score: this.add_score,
				Reason: this.add_score_reason,
			}
			app.post('/user/addscore', data, (result) => {
				this.table_data[this.current_row].Score = result.Score
				this.table_data[this.current_row].BankScore = result.BankScore
				this.dialog_data = app.clone(this.table_data[this.current_row])
				this.dialog_data_copy = app.clone(this.dialog_data)
				this.dialog_score = false
				this.$message.success('操作成功')
			})
		},
		handleSelectTab(tab) {
			if (this.current_tab != this.activeName) {
				this.current_tab = this.activeName
			}
		},
		handleCopy(field, index) {
			var input = document.getElementById('kcopy')
			if (!input) {
				input = document.createElement('input')
				input.setAttribute('id', 'kcopy')
				input.setAttribute('readonly', 'readonly')
				document.body.appendChild(input)
			}
			input.setAttribute('value', this.gamerecord_data[index][field])
			input.select()
			input.setSelectionRange(0, 9999)
			document.execCommand('Copy')
			this.$message.success('操作成功')
		},
	},
}
</script>
