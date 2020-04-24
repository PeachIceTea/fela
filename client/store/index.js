import Vue from "vue"
import Vuex from "vuex"

Vue.use(Vuex)

import auth from "./auth"
import audiobook from "./audiobook"
import uploads from "./uploads"
import user from "./user"
import ui from "./ui"

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
	modules: { auth, audiobook, uploads, user, ui },
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
