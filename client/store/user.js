import { getAllUsers } from "../api"

export default {
	state: {
		list: [],
	},
	actions: {
		async getAllUsers({ commit }) {
			const res = await getAllUsers()
			if (res.err) {
				return res
			}

			commit("setUserList", res.users)
		},
	},
	mutations: {
		setUserList(state, list) {
			state.list = list
		},
	},
	getters: {},
}
