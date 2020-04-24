export default {
	state: {
		order: 0,
		view: 0,
		search: "",
	},
	actions: {},
	mutations: {
		nextOrder(state) {
			state.order++
			state.order %= 4
		},
		nextView(state) {
			state.view++
			state.view %= 2
		},
		setSearch(state, search) {
			state.search = search
		},
	},
	getters: {},
}
