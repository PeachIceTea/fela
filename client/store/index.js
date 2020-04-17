import Vue from "vue"
import Vuex from "vuex"

Vue.use(Vuex)

import auth from "./auth"
import audiobook from "./audiobook"

export default new Vuex.Store({
	state: {
		notifications: [],
	},
	actions: {
		notifyError({ commit }, { msg }) {
			console.error(msg)
		},
	},
	mutations: {},
	modules: { auth, audiobook },
	strict: process.env.NODE_ENV !== "production",
})
