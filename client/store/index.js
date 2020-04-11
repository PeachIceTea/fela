import Vue from "vue"
import Vuex from "vuex"

import upload from "./upload"
import book from "./book"
import player from "./player"

Vue.use(Vuex)

export function apiPath(path) {
	return `http://localhost:8080/api/v1${path}`
}

export function apiCall(path, options) {
	return fetch(apiPath(path), options)
}

export function apiPOST(path, body, options = { headers: {} }) {
	options.headers["Content-Type"] = "application/json"
	options.method = "POST"

	return apiCall(path, Object.assign(options, { body: JSON.stringify(body) }))
}

export default new Vuex.Store({
	modules: { upload, book, player },
	strict: true,
})
