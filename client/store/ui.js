export default {
	state: {
		order: 4,
		view: 0,
		search: "",
		showUnread: false,
	},
	actions: {},
	mutations: {
		nextOrder(state) {
			state.order++
			state.order %= 5
		},
		nextView(state) {
			state.view++
			state.view %= 2
		},
		setSearch(state, search) {
			state.search = search
		},
		toggleShowUnread(state) {
			state.showUnread = !state.showUnread
		}
	},
	getters: {},
}
