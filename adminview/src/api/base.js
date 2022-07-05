import { app } from '@/api/app.js'
export default {
	data() {
		return {
			filters: {
				SellerId: app.currentSeller(),
			},
			zong: app.zong(),
			seller: app.getSeller(),
			seller_noall: app.getSellerNoAll(),
			seller_nozong: app.getSellerNoZong(),
			page: 1,
			pagesize: 15,
			total: 0,
			table_data: null,
			options: {},
			current_row: null,
			dialog: {
				show: false,
				title: '',
				type: '',
				options: {},
				data: {},
			},
			hgames_all: [
				{
					id: 0,
					name: '全部游戏',
				},
				{
					id: 1,
					name: '哈希单双',
				},
				{
					id: 2,
					name: '哈希大小',
				},
				{
					id: 3,
					name: '哈希牛牛',
				},
				{
					id: 4,
					name: '幸运哈希',
				},
				{
					id: 5,
					name: '幸运庄闲',
				},
			],
			hgames: [
				{
					id: 1,
					name: '哈希单双',
				},
				{
					id: 2,
					name: '哈希大小',
				},
				{
					id: 3,
					name: '哈希牛牛',
				},
				{
					id: 4,
					name: '幸运哈希',
				},
				{
					id: 5,
					name: '幸运庄闲',
				},
			],
			hrooms: [
				{
					id: 1,
					name: '初级场',
				},
				{
					id: 2,
					name: '中级场',
				},
				{
					id: 3,
					name: '高级场',
				},
			],
		}
	},
}
