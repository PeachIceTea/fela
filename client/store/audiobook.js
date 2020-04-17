import { getAllAudiobooks, getAudiobook } from "../api"

export default {
	state: {
		playing: {},
		list: [],
	},
	actions: {
		async getAudiobooks({ commit }) {
			const { audiobooks, err } = await getAllAudiobooks()
			if (err) {
				return { err }
			}

			commit("setAudiobookList", audiobooks)
			return {}
		},
		async getAudiobook({ commit }, id) {
			const { audiobook, err } = await getAudiobook(id)
			if (err) {
				return { err }
			}

			commit("setPlayingAudiobook", audiobook)
		},
	},
	mutations: {
		setAudiobookList(state, list) {
			state.list = list
		},
		setPlayingAudiobook(state, audiobook) {
			state.playing = audiobook
		},
	},
	getters: {},
}
