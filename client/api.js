import store from "./store"

// User
export async function getAllUsers() {
	return await GET("/user")
}

export async function getUser(id) {
	return await GET(`/user/${id}`)
}

export async function register(name, password, role) {
	return await POST("/user/register", { name, password, role })
}

export async function login(name, password) {
	return await POST("/user/login", { name, password })
}

export async function updateUser(id, { name, password, role, confirmation }) {
	return await PUT(`/user/${id}`, { name, password, role, confirmation })
}

export async function deleteUser(id) {
	return await DELETE(`/user/${id}`)
}

export async function testToken() {
	return await GET("/token-test")
}

// Audiobook
export async function getAllAudiobooks() {
	return await GET("/audiobook")
}

export async function getAudiobook(id) {
	return await GET(`/audiobook/${id}`)
}

export async function getAudiobookFiles(id) {
	return await GET(`/audiobook/${id}/files`)
}

export async function deleteAudiobook(id) {
	return await DELETE(`/audiobook/${id}`)
}

// Complex
// Handles upload
export function uploadAudiobook(files, progressCallback = () => {}) {
	return new Promise(resolve => {
		const form = new FormData()
		for (let i = 0, len = files.length; i < len; i++) {
			form.append("file", files[i])
		}

		const xhr = new XMLHttpRequest()
		xhr.responseType = "json"
		xhr.upload.onprogress = e => {
			progressCallback(e.loaded / e.total)
		}
		xhr.onerror = e => {
			console.error(e)
			resolve({ err: "could not connect to server" })
		}
		xhr.onload = e => {
			if (xhr.status !== 200) {
				resolve({ err: xhr.response.err })
			} else {
				resolve({ audiobookID: xhr.response.audiobook_id })
			}
		}
		try {
			xhr.open("POST", apiURL("/audiobook/upload"))
			xhr.setRequestHeader("Authorization", getAuthHeader())
			xhr.send(form)
		} catch (e) {
			console.error(e)
			resolve({ err: "could not connect to server" })
		}
	})
}

export function updateAudiobook(
	id,
	{ title, author, cover },
	progressCallback = () => {},
) {
	return new Promise(resolve => {
		const form = new FormData()
		if (title) {
			form.set("title", title)
		}
		if (author) {
			form.set("author", author)
		}
		if (cover) {
			form.set("cover", cover)
		}

		const xhr = new XMLHttpRequest()
		xhr.responseType = "json"
		xhr.upload.onprogress = e => {
			progressCallback(Math.floor((e.loaded / e.total) * 100))
		}
		xhr.onerror = e => {
			console.error(e)
			resolve({ err: "could not connect to server" })
		}
		xhr.onload = e => {
			if (xhr.status !== 200) {
				resolve({ err: xhr.response.err })
			} else {
				resolve({})
			}
		}
		try {
			xhr.open("PUT", apiURL(`/audiobook/${id}`))
			xhr.setRequestHeader("Authorization", getAuthHeader())
			xhr.send(form)
		} catch (e) {
			console.error(e)
			resolve({ err: "could not connect to server" })
		}
	})
}

// Helper
export function GET(route) {
	return api(route, fetchConfig({ method: "GET" }))
}

export function POST(route, body) {
	return api(
		route,
		fetchConfig({
			method: "POST",
			contentType: "application/json",
			body: JSON.stringify(body),
		}),
	)
}

export function PUT(route, body) {
	return api(
		route,
		fetchConfig({
			method: "PUT",
			contentType: "application/json",
			body: JSON.stringify(body),
		}),
	)
}

export function DELETE(route) {
	return api(route, fetchConfig({ method: "DELETE" }))
}

export function fetchConfig({ method, body, contentType }) {
	const config = {
		method,
		body,
		headers: {
			Authorization: getAuthHeader(),
			"Content-Type": contentType,
		},
	}
	return config
}

export const getAuthHeader = () => `Bearer ${store.state.auth.token}`

export async function api(route, config) {
	try {
		return await (await fetch(apiURL(route), config)).json()
	} catch (e) {
		return { err: "cannot connect to server" }
	}
}

export function apiURL(route) {
	return `/api/v1${route}`
}

// Fileurl
export function audiobookURL(id, filename) {
	const token = encodeURIComponent(store.state.auth.token)
	name = encodeURIComponent(name)
	return `/files/audio/${id}/${filename}?auth=${token}`
}

export function coverURL(id) {
	const token = encodeURIComponent(store.state.auth.token)
	return `/files/cover/${id}.jpg?auth=${token}&r=${Math.random()
		.toString(36)
		.replace(/[^a-z]+/g, "")
		.substr(0, 5)}`
}
