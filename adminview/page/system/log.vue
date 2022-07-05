<template>
	<div class="container">
		<!-- 筛选 -->
		<div>
			<el-form :inline="true" :model="filters">
				<el-form-item label="管理员:">
					<el-input v-model="filters.Account" style="width: 150px" :clearable="true"></el-input>
				</el-form-item>
				<el-form-item label="操作:">
					<el-input v-model="filters.Opt" style="width: 150px" :clearable="true"></el-input>
				</el-form-item>
				<el-form-item label="运营商:" v-show="zong">
					<el-select v-model="filters.SellerId" placeholder="运营商" style="width: 130px" @change="handleSelectSeller">
						<el-option v-for="item in seller" :key="item.SellerId" :label="item.SellerName" :value="item.SellerId"> </el-option>
					</el-select>
				</el-form-item>
				<el-form-item>
					<el-button type="primary" icon="el-icon-refresh" v-on:click="handleQuery">查询</el-button>
				</el-form-item>
			</el-form>
		</div>
		<!-- 表 -->
		<div>
			<el-table :data="table_data" border class="table" max-height="700px" :cell-style="{ padding: '0' }">
				<el-table-column align="center" prop="Id" label="序号" width="80"></el-table-column>
				<el-table-column align="center" prop="Account" label="管理员" width="100"></el-table-column>
				<el-table-column align="center" prop="SellerName" label="运营商" width="150"></el-table-column>
				<el-table-column align="center" prop="Ip" label="ip" width="130"></el-table-column>
				<el-table-column align="center" prop="Opt" label="操作类型" width="200"></el-table-column>
				<el-table-column align="center" prop="CreateTime" label="时间" width="160"></el-table-column>
				<el-table-column label="内容">
					<template slot-scope="scope">
						<el-button type="text" icon="el-icon-document-copy" @click="handleCopy(scope.$index)">复制</el-button>
					</template>
				</el-table-column>
			</el-table>
			<div class="pagination">
				<el-pagination style="margin-top: 5px" background layout="total, prev, pager, next, jumper" :hide-on-single-page="true" :total="total" @current-change="handleQuery" :page-size="pagesize"></el-pagination>
			</div>
		</div>
	</div>
</template>
<script>
import { app } from '@/api/app.js'
import base from '@/api/base.js'
import '@/assets/css/k.css'
export default {
	extends: base,
	data() {
		return {
			filters: {
				Account: null,
				Opt: null,
			},
		}
	},
	created() {
		this.handleQuery(1)
	},
	methods: {
		handleSelectSeller() {
			app.setCurrentSeller(this.filters.SellerId)
		},
		handleCopy(index) {
			var oInput = document.createElement('input')
			oInput.value = this.table_data[index].Data
			document.body.appendChild(oInput)
			oInput.select()
			document.execCommand('Copy')
			oInput.remove()
			this.$message.success('复制成功')
		},
		handleQuery(page) {
			this.page = page || 1
			if (typeof this.page == 'object') this.page = 1
			var data = {
				page: this.page,
				pagesize: this.pagesize,
				Account: this.filters.Account || '',
				Opt: this.filters.Opt || '',
				SellerId: this.filters.SellerId || 0,
			}
			app.post('/admin/opt_log', data, (result) => {
				this.table_data = result.data.data
				this.total = result.data.total
				for (var i = 0; i < this.table_data.length; i++) {
					this.table_data[i].CreateTime = this.$moment(this.table_data[i].CreateTime).format('YYYY-MM-DD hh:mm:ss')
					for (let j = 0; j < this.seller.length; j++) {
						if (this.seller[j].SellerId == this.table_data[i].SellerId) this.table_data[i].SellerName = this.seller[j].SellerName
					}
				}
			})
		},
	},
}
</script>
