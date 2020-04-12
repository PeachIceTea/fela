import { apiCall } from "."

const defaultState = {
	book: {},
	audiobook: {},
	files: [],
}

export default {
	namespaced: true,
	state: Object.assign({}, defaultState),
	getters: {},
	actions: {
		async play({ commit, dispatch }, { audiobook, book }) {
			commit("clear")
			commit("setNew", { audiobook, book })

			const { files } = await (
				await apiCall(`/audiobook/${audiobook.id}`)
			).json()
			commit("setFiles", { files })
		},
	},
	mutations: {
		clear(state) {
			Object.assign(state, defaultState)
		},
		setNew(state, { book, audiobook }) {
			state.book = book
			state.audiobook = audiobook
		},
		setFiles(state, { files }) {
			state.files = files
		},
	},
}
