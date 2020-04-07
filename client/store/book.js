import { apiCall } from "."

export default {
	namespaced: true,
	state: {
		books: [],
		book: null,
		audiobooks: [],
	},
	getters: {},
	actions: {
		async getBooks({ commit }) {
			//TODO: Error handling
			const { books } = await (await apiCall("/book")).json()
			commit("setBooks", books)
		},
		async getBook({ commit }, id) {
			const { book, audiobooks } = await (
				await apiCall(`/book/${id}`)
			).json()
			commit("setBook", { book, audiobooks })
		},
	},
	mutations: {
		setBooks(state, books) {
			state.books = books
		},
		setBook(state, { book, audiobooks }) {
			state.book = book
			state.audiobooks = audiobooks
		},
	},
}
