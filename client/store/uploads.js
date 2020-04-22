import { uploadAudiobook } from "../api"

export default {
	state: [],
	actions: {
		async upload({ commit }, upload) {
			upload.progress = 0
			upload.err = ""
			commit("addUpload", upload)
			if (upload.err) {
				return
			}

			commit("startUpload", upload)
			const res = await uploadAudiobook(upload.files, progress => {
				commit("updateProgress", { upload, progress })
			})
			if (res.err) {
				commit("assignError", { upload, err: res.err })
				return res
			}
			commit("assignAudiobookID", { upload, id: res.audiobook_id })
			return res
		},
	},
	mutations: {
		addUpload(state, upload) {
			state.push(upload)
		},
		startUpload(state, upload) {
			upload.progress = 0
		},
		updateProgress(state, { upload, progress }) {
			upload.progress = progress
		},
		assignError(state, { upload, err }) {
			upload.err = err
		},
		assignAudiobookID(state, { upload, id }) {
			upload.id = id
		},
	},
	getters: {},
}
