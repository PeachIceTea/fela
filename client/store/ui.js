export default {
	state: {
		order: 4,
		view: 0,
		search: "",
		playbackRate: 1.0,
		volume: 1.0,
	},
	actions: {
		restorePlaybackRate({ commit }) {
			const stored = localStorage.getItem("playback-rate")
			if (stored) {
				const parsed = parseFloat(stored)
				if (parsed) {
					commit("playbackRate", parsed)
				}
			}
		},
		setPlaybackRate({ commit }, rate) {
			localStorage.setItem("playback-rate", rate)
			commit("playbackRate", rate)
		},
	},
	mutations: {
		nextOrder(state) {
			state.order++
			state.order %= 5
		},
		nextView(state) {
			state.view++
			state.view %= 2
		},
		setSearch(state, search) {
			state.search = search
		},
		playbackRate(state, rate) {
			console.log(rate)
			state.playbackRate = rate
		},
	},
	getters: {},
}
