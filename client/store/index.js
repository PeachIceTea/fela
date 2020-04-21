import Vue from "vue"
import Vuex from "vuex"

Vue.use(Vuex)

import auth from "./auth"
import audiobook from "./audiobook"
import user from "./user"
import uploads from "./uploads"

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
	modules: { auth, audiobook, user, uploads },
	getters: {
		authors(state) {
			const list = state.audiobook.list
			const authors = []
			for (let i = 0, len = list.length; i < len; i++) {
				const book = list[i]
				if (book.author && !authors.includes(book.author)) {
					authors.push(book.author)
				}
			}
			return authors
		},
	},
	strict: process.env.NODE_ENV !== "production",
})
