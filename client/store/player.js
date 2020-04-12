import { apiCall } from "."

const defaultState = {
	book: {},
	audiobook: {},
	files: [],
	file: {},
	fileIndex: -1,
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
			commit("setCurrentFile", { index: 0 })
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
		setCurrentFile(state, { index }) {
			state.file = state.files[index]
			state.fileIndex = index
		},
	},
}
