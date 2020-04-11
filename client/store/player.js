import { apiCall } from "."

export default {
	namespaced: true,
	state: {
		files: [],
		file: {},
		fileIndex: -1,
	},
	getters: {},
	actions: {
		async getFiles({ commit }, id) {
			const { files } = await (await apiCall(`/audiobook/${id}`)).json()
			commit("setFiles", { files })
			commit("setCurrentFile", { index: 0 })
		},
	},
	mutations: {
		setFiles(state, { files }) {
			state.files = files
		},
		setCurrentFile(state, { index }) {
			state.file = state.files[index]
			state.fileIndex = index
		},
	},
}
