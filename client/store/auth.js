import { login, testToken } from "../api"

export default {
	state: {
		token: null,
	},
	actions: {
		async login({ commit }, { name, password }) {
			const res = await login(name, password)
			if (res.err) {
				return res
			} else {
				commit("setToken", res.token)
				localStorage.setItem("token", res.token)
				return res
			}
		},
		async restoreToken({ commit }) {
			const token = localStorage.getItem("token")
			if (token) {
				commit("setToken", token)
				const res = await testToken()
				if (res.err) {
					commit("setToken", null)
					localStorage.removeItem("token")
				}
			}
		},
		logout({ commit }) {
			localStorage.removeItem("token")
			commit("setToken", null)
		},
	},
	mutations: {
		setToken(state, token) {
			state.token = token
		},
	},
	getters: {},
}
