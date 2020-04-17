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

export async function updateUser(id, { name, password, role }) {
	return await PUT(`/user/${id}`, { name, password, role })
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
export function upload(files, progressCallback) {
	return new Promise(resolve => {
		const form = new FormData()
		// An file input element has only a single FileList which is reused
		// when needed. To keep references to the files we need to create a
		// new array.
		const arr = []
		for (let i = 0, len = files.length; i < len; i++) {
			const file = files[i]
			arr[i] = file
			form.append("file", file)
		}
		arr.sort((a, b) => {
			if (a.name < b.name) {
				return -1
			} else if (a.name > b.name) {
				return 1
			} else {
				return 0
			}
		})

		const xhr = new XMLHttpRequest()
		xhr.responseType = "json"
		xhr.upload.onprogress = e => {
			progressCallback(Math.floor(e.loaded / e.total))
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
			xhr.open("POST", url("/audiobook/upload"))
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
	progressCallback,
) {
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
		progressCallback(Math.floor(e.loaded / e.total))
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
		xhr.open("PUT", url(`/audiobook/${id}`))
		xhr.send(form)
	} catch (e) {
		console.error(e)
		resolve({ err: "could not connect to server" })
	}
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
			Authorization: `Bearer ${store.state.auth.token}`,
			"Content-Type": contentType,
		},
	}
	return config
}

export async function api(route, config) {
	try {
		return await (await fetch(url(route), config)).json()
	} catch (e) {
		return { err: "cannot connect to server" }
	}
}

export function url(route) {
	return `http://localhost:8080/api/v1${route}`
}
