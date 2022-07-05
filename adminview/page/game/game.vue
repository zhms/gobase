<template>
	<div class="container">
		<!-- 筛选 -->
		<div>
			<el-form :inline="true" :model="filters">
				<el-form-item label="运营商:" v-show="zong">
					<el-select v-model="filters.SellerId" placeholder="运营商" style="width: 130px">
						<el-option v-for="item in seller_nozong" :key="item.SellerId" :label="item.SellerName" :value="item.SellerId"> </el-option>
					</el-select>
				</el-form-item>
				<el-form-item label="游戏:">
					<el-select v-model="filters.GameId" placeholder="游戏" style="width: 130px">
						<el-option v-for="item in hgames_all" :key="item.id" :label="item.name" :value="item.id"> </el-option>
					</el-select>
				</el-form-item>
				<el-form-item>
					<el-button type="primary" icon="el-icon-refresh" class="mr10" @click="handleQuery">查询</el-button>
					<el-button type="primary" icon="el-icon-plus" class="mr10" v-show="auth2('增')" @click="handleAdd">添加</el-button>
				</el-form-item>
			</el-form>
		</div>
		<!-- 表 -->
		<div>
			<el-table :data="table_data" border max-height="620px" class="table" :cell-style="{ padding: '0' }">
				<el-table-column align="center" prop="SellerName" label="运营商" width="100"></el-table-column>
				<el-table-column align="center" prop="GameName" label="游戏名称" width="150"></el-table-column>
				<el-table-column align="center" prop="RoomName" label="房间名称" width="100"></el-table-column>
				<el-table-column align="center" prop="Address" label="投注地址" width="300"></el-table-column>
				<el-table-column align="center" label="状态" width="100">
					<template slot-scope="scope">
						<el-link :underline="false" type="danger" v-if="scope.row.State == 2">禁用</el-link>
						<el-link :underline="false" type="primary" v-if="scope.row.State == 1">启用</el-link>
					</template>
				</el-table-column>
				<el-table-column align="center" prop="Multipe" label="返奖赔率" width="100"></el-table-column>
				<el-table-column align="center" prop="MinTrx" label="最小TRX" width="100"></el-table-column>
				<el-table-column align="center" prop="MaxTrx" label="最大TRX" width="100"></el-table-column>
				<el-table-column align="center" prop="MinUsdt" label="最小USDT" width="100"></el-table-column>
				<el-table-column align="center" prop="MaxUsdt" label="最大USDT" width="100"></el-table-column>
				<el-table-column align="center" prop="MaxFeeRate" label="返还费率" width="100"></el-table-column>
				<el-table-column align="center" prop="AgentRate" label="返佣基数" width="100"></el-table-column>
				<el-table-column label="操作">
					<template slot-scope="scope">
						<el-button type="text" icon="el-icon-edit" v-show="auth2('改')" @click="handleModify(scope.$index)">编辑</el-button>
					</template>
				</el-table-column>
			</el-table>
		</div>
		<div class="pagination">
			<el-pagination style="margin-top: 5px" background layout="total, prev, pager, next, jumper" :hide-on-single-page="true" :total="total" @current-change="handleQuery" :page-size="pagesize"></el-pagination>
		</div>
		<!--对话框-->
		<div>
			<el-dialog :title="dialog.title" :visible.sync="dialog.show" width="610px" center>
				<el-form :inline="true" label-width="100px">
					<el-form-item label="运营商:" v-show="zong">
						<el-select v-model="dialog.data.SellerId" placeholder="运营商" style="width: 390px" :disabled="dialog.type == 'alter'">
							<el-option v-for="item in seller_nozong" :key="item.SellerId" :label="item.SellerName" :value="item.SellerId"> </el-option>
						</el-select>
					</el-form-item>
					<el-form-item label="游戏:">
						<el-select v-model="dialog.data.GameId" placeholder="游戏" style="width: 390px" :disabled="dialog.type == 'alter'">
							<el-option v-for="item in hgames" :key="item.id" :label="item.name" :value="item.id"> </el-option>
						</el-select>
					</el-form-item>
					<el-form-item label="房间:">
						<el-select v-model="dialog.data.RoomLevel" placeholder="房间" style="width: 390px" :disabled="dialog.type == 'alter'">
							<el-option v-for="item in hrooms" :key="item.id" :label="item.name" :value="item.id"> </el-option>
						</el-select>
					</el-form-item>
					<el-form-item label="地址:">
						<el-input v-model="dialog.data.Address" style="width: 300px" :disabled="true"></el-input>
						<el-button type="primary" icon="el-icon-refresh" class="mr10" style="margin-left: 15px" @click="handleFlushAddress" v-show="dialog.type == 'add'">刷新</el-button>
					</el-form-item>
					<el-form-item label="赔率:" v-show="dialog.data.GameId == 1 || dialog.data.GameId == 2 || dialog.data.GameId == 4">
						<el-input v-model="dialog.data.Multipe" style="width: 390px"></el-input>
					</el-form-item>
					<el-form-item label="庄赔率:" v-show="dialog.data.GameId == 5">
						<el-input v-model="dialog.data.mz" style="width: 140px"></el-input>
					</el-form-item>
					<el-form-item label="闲赔率:" v-show="dialog.data.GameId == 5">
						<el-input v-model="dialog.data.mx" style="width: 140px"></el-input>
					</el-form-item>
					<el-form-item label="和赔率:" v-show="dialog.data.GameId == 5">
						<el-input v-model="dialog.data.mh" style="width: 140px"></el-input>
					</el-form-item>
					<el-form-item label="牛一赔率:" v-show="dialog.data.GameId == 3">
						<el-input v-model="dialog.data.mn1" style="width: 140px"></el-input>
					</el-form-item>
					<el-form-item label="牛二赔率:" v-show="dialog.data.GameId == 3">
						<el-input v-model="dialog.data.mn2" style="width: 140px"></el-input>
					</el-form-item>
					<el-form-item label="牛三赔率:" v-show="dialog.data.GameId == 3">
						<el-input v-model="dialog.data.mn3" style="width: 140px"></el-input>
					</el-form-item>
					<el-form-item label="牛四赔率:" v-show="dialog.data.GameId == 3">
						<el-input v-model="dialog.data.mn4" style="width: 140px"></el-input>
					</el-form-item>
					<el-form-item label="牛五赔率:" v-show="dialog.data.GameId == 3">
						<el-input v-model="dialog.data.mn5" style="width: 140px"></el-input>
					</el-form-item>
					<el-form-item label="牛六赔率:" v-show="dialog.data.GameId == 3">
						<el-input v-model="dialog.data.mn6" style="width: 140px"></el-input>
					</el-form-item>
					<el-form-item label="牛七赔率:" v-show="dialog.data.GameId == 3">
						<el-input v-model="dialog.data.mn7" style="width: 140px"></el-input>
					</el-form-item>
					<el-form-item label="牛八赔率:" v-show="dialog.data.GameId == 3">
						<el-input v-model="dialog.data.mn8" style="width: 140px"></el-input>
					</el-form-item>
					<el-form-item label="牛九赔率:" v-show="dialog.data.GameId == 3">
						<el-input v-model="dialog.data.mn9" style="width: 140px"></el-input>
					</el-form-item>
					<el-form-item label="牛牛赔率:" v-show="dialog.data.GameId == 3">
						<el-input v-model="dialog.data.mn0" style="width: 140px"></el-input>
					</el-form-item>
					<el-form-item label="最小TRX:">
						<el-input v-model="dialog.data.MinTrx" style="width: 390px"></el-input>
					</el-form-item>
					<el-form-item label="最大TRX:">
						<el-input v-model="dialog.data.MaxTrx" style="width: 390px"></el-input>
					</el-form-item>
					<el-form-item label="最小USDT:">
						<el-input v-model="dialog.data.MinUsdt" style="width: 390px"></el-input>
					</el-form-item>
					<el-form-item label="最大USDT:">
						<el-input v-model="dialog.data.MaxUsdt" style="width: 390px"></el-input>
					</el-form-item>
					<el-form-item label="返回费率:">
						<el-input v-model="dialog.data.MaxFeeRate" style="width: 390px"></el-input>
					</el-form-item>
					<el-form-item label="返佣基数:">
						<el-input v-model="dialog.data.AgentRate" style="width: 390px"></el-input>
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
			filters: {
				GameId: null,
			},
		}
	},
	components: {},
	computed: {},
	created() {
		if (this.filters.SellerId <= 0) this.filters.SellerId = null
		this.handleQuery()
	},
	methods: {
		auth2(o) {
			return app.auth2('游戏管理', '游戏列表', o)
		},
		handleQuery(page) {
			this.page = page || 1
			if (typeof this.page != 'number') this.page = 1
			var data = {
				page: this.page,
				pagesize: this.pagesize,
				SellerId: parseInt(this.filters.SellerId || 0),
				GameId: this.filters.GameId || '',
			}
			app.post('/game/list', data, (result) => {
				for (let i = 0; i < result.data.data.length; i++) {
					for (let j = 0; j < this.hgames.length; j++) {
						if (result.data.data[i].GameId == this.hgames[j].id) {
							result.data.data[i].GameName = this.hgames[j].name
						}
					}
					for (let j = 0; j < this.hrooms.length; j++) {
						if (result.data.data[i].RoomLevel == this.hrooms[j].id) {
							result.data.data[i].RoomName = this.hrooms[j].name
						}
					}
					for (let j = 0; j < this.seller.length; j++) {
						if (this.seller[j].SellerId == result.data.data[i].SellerId) result.data.data[i].SellerName = this.seller[j].SellerName
					}
				}
				this.table_data = result.data.data
			})
		},
		handleFlushAddress() {
			if (this.dialog.data.SellerId <= 0) {
				this.$message.error('请先选择运营商')
				return
			}
			app.post('/game/address', { SellerId: this.dialog.data.SellerId }, (result) => {
				this.dialog.data.Address = result.data.Address
				this.dialog.data = JSON.parse(JSON.stringify(this.dialog.data))
			})
		},
		handleAdd() {
			this.dialog.type = 'add'
			this.dialog.title = '新增游戏'
			this.dialog.data = {
				SellerId: null,
				GameId: null,
				RoomLevel: null,
				MinTrx: null,
				MaxTrx: null,
				MinUsdt: null,
				MaxUsdt: null,
				mz: null,
				mx: null,
				mh: null,
				MultipeZhuang: null,
				mn0: null,
				mn1: null,
				mn2: null,
				mn3: null,
				mn4: null,
				mn5: null,
				mn6: null,
				mn7: null,
				mn8: null,
				mn9: null,
			}
			this.dialog.show = true
		},
		handleModify(index) {
			this.current_row = index
			this.dialog.type = 'alter'
			let data = this.table_data[index]
			if (data.GameId == 5) {
				let d = data.Multipe.split(',')
				data.mz = d[0]
				data.mx = d[1]
				data.mh = d[2]
			}
			if (data.GameId == 3) {
				let d = data.Multipe.split(',')
				data.mn1 = d[0]
				data.mn2 = d[1]
				data.mn3 = d[2]
				data.mn4 = d[3]
				data.mn5 = d[4]
				data.mn6 = d[5]
				data.mn7 = d[6]
				data.mn8 = d[7]
				data.mn9 = d[8]
				data.mn0 = d[9]
			}
			this.dialog.data = app.clone(this.table_data[index])
			this.dialog.show = true
		},
		handleDel(index) {},

		handleConfirm() {
			if (this.dialog.data.GameId == 5) this.dialog.data.Multipe = `${this.dialog.data.mz},${this.dialog.data.mx},${this.dialog.data.mh}`
			if (this.dialog.data.GameId == 3)
				this.dialog.data.Multipe = `${this.dialog.data.mn1},${this.dialog.data.mn2},${this.dialog.data.mn3},${this.dialog.data.mn4},${this.dialog.data.mn5},${this.dialog.data.mn6},${this.dialog.data.mn7},${this.dialog.data.mn8},${this.dialog.data.mn9},${this.dialog.data.mn0}`
			if (this.dialog.type == 'alter') {
				this.dialog.data.MinTrx = Number(this.dialog.data.MinTrx)
				this.dialog.data.MaxTrx = Number(this.dialog.data.MaxTrx)
				this.dialog.data.MinUsdt = Number(this.dialog.data.MinUsdt)
				this.dialog.data.MaxUsdt = Number(this.dialog.data.MaxUsdt)
				this.dialog.data.MaxFeeRate = Number(this.dialog.data.MaxFeeRate)
				this.dialog.data.AgentRate = Number(this.dialog.data.AgentRate)
				app.post('/game/modify', this.dialog.data, () => {
					this.handleQuery(this.page)
					this.dialog.show = false
				})
			}
			if (this.dialog.type == 'add') {
				this.dialog.data.MinTrx = Number(this.dialog.data.MinTrx)
				this.dialog.data.MaxTrx = Number(this.dialog.data.MaxTrx)
				this.dialog.data.MinUsdt = Number(this.dialog.data.MinUsdt)
				this.dialog.data.MaxUsdt = Number(this.dialog.data.MaxUsdt)
				this.dialog.data.MaxFeeRate = Number(this.dialog.data.MaxFeeRate)
				this.dialog.data.AgentRate = Number(this.dialog.data.AgentRate)
				app.post('/game/add', this.dialog.data, () => {
					this.handleQuery(this.page)
					this.dialog.show = false
				})
			}
		},
	},
}
</script>
