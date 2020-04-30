import { login, testToken } from "../api"

export default {
	state: {
		token: null,
		loggedIn: {},
	},
	actions: {
		async login({ commit }, { name, password }) {
			const res = await login(name, password)
			if (res.err) {
				return res
			} else {
				commit("setloggedIn", res.token)
				localStorage.setItem("token", res.token)
				localStorage.setItem("user", JSON.stringify(res.user))
				return res
			}
		},
		async restoreToken({ commit }) {
			const token = localStorage.getItem("token")
			const user = JSON.parse(localStorage.getItem("user"))
			if (token && user) {
				commit("setToken", token)
				const res = await testToken()
				if (res.err) {
					commit("setToken", null)
					localStorage.removeItem("token")
					localStorage.removeItem("user")
					return
				}
				commit("setloggedIn", user)
			}
		},
		logout({ commit }) {
			localStorage.removeItem("token")
			localStorage.removeItem("user")
			commit("setToken", null)
			commit("setloggedIn", {})
		},
	},
	mutations: {
		setToken(state, token) {
			state.token = token
		},
		setloggedIn(state, user) {
			state.loggedIn = user
		},
	},
	getters: {
		isUploader(state) {
			const role = state.loggedIn.role
			return role === "uploader" || role === "admin"
		},
		isAdmin(state) {
			return state.loggedIn.role === "admin"
		},
	},
}
