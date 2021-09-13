import { getAllAudiobooks, getAudiobook, getUserProgress } from "../api"

export default {
	state: {
		playing: {},
		list: [],
		editing: {},
		userProgress: []
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
		async getEditingAudiobook({ commit }, id) {
			const { audiobook, err } = await getAudiobook(id)
			if (err) {
				return { err }
			}

			commit("setEditingAudiobook", audiobook)
		},
		async getUserProgress({commit}) {
			const {progress, err} = await getUserProgress()
			if (err) {
				return {err}
			}

			commit("setUserProgress", progress)
		},
		async playAudiobook({ commit }, id) {
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
		setEditingAudiobook(state, audiobook) {
			state.editing = audiobook
		},
		setUserProgress(state, progress) {
			state.userProgress = progress
		}
	},
	getters: {},
}
