import { apiCall } from "."

export default {
	namespaced: true,
	state: {
		books: [],
		book: {},
		audiobooks: [],
		authors: [],
	},
	getters: {},
	actions: {
		async getBooks({ commit }) {
			//TODO: Error handling
			const { books } = await (await apiCall("/book")).json()
			if (!books) return
			commit("setBooks", books)
		},
		async getBook({ commit }, id) {
			const { book, audiobooks } = await (
				await apiCall(`/book/${id}`)
			).json()
			commit("setBook", { book, audiobooks })
		},
		async getAuthors({ commit }) {
			const { authors } = await (await apiCall("/author")).json()
			commit("setAuthors", { authors })
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
		setAuthors(state, { authors }) {
			state.authors = authors
		},
	},
}
