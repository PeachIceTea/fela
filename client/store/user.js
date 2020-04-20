import { getUser } from "../api"

export default {
	state: {
		uploads: [],
	},
	actions: {
		async getUserUploads({ commit }, id) {
			const res = await getUser(id)
			if (res.err) {
				return res
			} else {
				commit("setUploads", res.user.uploads)
			}
		},
	},
	mutations: {
		setUploads(state, uploads) {
			state.uploads = uploads
		},
	},
	getters: {},
}
