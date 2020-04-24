import { uploadAudiobook } from "../api"

export default {
	state: [],
	actions: {
		async upload({ commit }, { files, err }) {
			if (err) {
				commit("addUpload", { files, err })
				return
			}

			const upload = { files, err: "", progress: 0 }
			commit("addUpload", upload)
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
