import { apiPath, apiPOST } from "."

export default {
	namespaced: true,
	state: {
		uploads: [],
	},
	getters: {},
	actions: {
		upload({ state, commit }, files) {
			const index = state.uploads.length
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

			const tmp = {
				files: arr,
				chapterized: undefined,
				progress: 0,
			}
			commit("pushUploads", tmp)

			const xhr = new XMLHttpRequest()
			xhr.responseType = "json"
			xhr.upload.onprogress = e => {
				commit("updateProgress", {
					index,
					progress: Math.floor((e.loaded / e.total) * 100),
				})
			}
			xhr.onerror = e => {
				//TODO: Error handling
				throw e
			}
			xhr.onload = e => {
				if (xhr.status !== 200) {
					//TODO: Error handling
					throw e
				}

				const list = xhr.response.file_ids
				for (let i = 0, len = list.length; i < len; i++) {
					const { name, id } = list[i]
					for (let i = 0; i < len; i++) {
						const file = arr[i]
						if (file.name === name) {
							file.id = id
							break
						}
					}
				}
			}
			try {
				xhr.open("POST", apiPath("/file/upload"))
				xhr.send(form)
			} catch (e) {
				//TODO: Error handling
				throw e
			}
		},
		async submitBook({ state, commit }, { title, author, index, book }) {
			const upload = state.uploads[index]
			try {
				let res = await apiPOST("/book/new", { title, author })
				if (res.status !== 200) {
					//TODO: Error handling
					console.error(res)
					throw "POST /book/new failed"
				}
				const bookID = (await res.json()).book_id
				console.log(bookID)

				res = await apiPOST("/audiobook/new", { book_id: bookID })
				if (res.status !== 200) {
					//TODO: Error handling
					console.error(res)
					throw "POST /audiobook/new failed"
				}
				const audiobookID = (await res.json()).audiobook_id

				const tmp = []
				if (book) {
					tmp.push({ file_id: book.id, chapter: -1 })
				} else {
					for (let i = 0, len = upload.files.length; i < len; i++) {
						const file = upload.files[i]
						tmp.push({ file_id: file.id, chapter: i })
					}
				}

				const assignments = {
					audiobook_id: audiobookID,
					assignments: tmp,
				}
				res = await apiPOST("/file/assign", assignments)
				if (res.status !== 200) {
					//TODO: Error handling
					console.error(res)
					throw "POST /file/assign failed"
				}

				console.log(assignments)

				commit("markDone", { index, book })
			} catch (e) {
				//TODO: Error handling
				throw e
			}
		},
	},
	mutations: {
		pushUploads(state, files) {
			state.uploads.push(files)
		},
		setChapterized(state, { index, chapterized }) {
			state.uploads[index].chapterized = chapterized
		},
		updateProgress(state, { index, progress }) {
			state.uploads[index].progress = progress
		},
		markDone(state, { index, book }) {
			if (book) {
				book.done = true
			} else {
				state.uploads[index].done = true
			}
		},
	},
}
