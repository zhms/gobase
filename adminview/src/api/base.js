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
		}
	},
}
