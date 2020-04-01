<template lang="pug">
	.upload-container
		.file-list
			.file(v-for="file in files")
				span {{ file.name }} - Progress: {{ file.progress }}% - {{ file.status }}

				UploadMetadata(:callback="callback" :file="file" v-if="!file.bookInfoSubmitted")
				span(v-else) {{ file.metadataStatus }}

		label(for="fileupload")
			.upload(
				@drop="dropFile"
				@dragover.stop.prevent="dragOver"
			)
				p.text Drag file here or click to upload file.
			input(id="fileupload" type="file" multiple @input="dialogFile")
</template>

<script>
import UploadMetadata from "./UploadMetadata.vue"

export default {
	data() {
		return {
			files: [],
		}
	},
	methods: {
		dragOver(e) {
			e.dataTransfer.dropEffect = "copy"
		},
		dropFile(e) {
			// FIXME: do drag and drop
			const files = e.dataTransfer.files

			for (let i = 0, len = files.length; i < len; i++) {
				let file = files[i]
			}
		},
		async dialogFile(e) {
			const files = e.srcElement.files

			for (let i = 0, len = files.length; i < len; i++) {
				let file = files[i]

				//FIXME: Test for valid mime

				const fileInfo = {
					//TODO: Improve variable names
					//TODO: Represent the status text as numbers
					// Data for file upload
					name: file.name,
					progress: 0,
					fileID: 0,
					file,
					status: "pending",

					// Data for book info submission
					bookInfoSubmitted: false,
					bookInfoSubmissionSuccess: false,
					bookInfoSubmissionStatusText: "",
				}

				this.files.push(fileInfo)

				const form = new FormData()
				form.set("file", file)

				const xhr = new XMLHttpRequest()
				xhr.responseType = "json"

				xhr.upload.onprogress = function(e) {
					fileInfo.progress = Math.floor((e.loaded / e.total) * 100)
					fileInfo.status = "uploading"
				}

				xhr.onerror = function(e) {
					fileInfo.status = "error: could not reach server"
					fileInfo.bookInfoSubmitted = true
				}

				xhr.onload = function(e) {
					const id = xhr.response.file_id
					if (id === 0) {
						//FIXME: Deal with reupload in some smart way
						console.error("Reupload file")
						return
					}

					if (xhr.status === 200) {
						fileInfo.status = "done"
					} else {
						console.log(xhr.response)
						fileInfo.status = "error: " + xhr.response.error
						fileInfo.bookInfoSubmitted = true
					}
					fileInfo.fileID = id
				}

				xhr.open("POST", "http://localhost:8080/upload")
				xhr.send(form)
			}

			e.srcElement.value = ""
		},
		async callback(fileInfo, data) {
			fileInfo.bookInfoSubmitted = true
			data.file_id = fileInfo.fileID

			try {
				const res = await fetch("http://localhost:8080/book/create", {
					method: "POST",
					headers: {
						"Content-Type": "application/json",
					},
					body: JSON.stringify(data),
				})
			} catch (e) {
				fileInfo.bookInfoSubmissionSuccess = false
				fileInfo.bookInfoSubmissionStatusText =
					"Could not connect to server."
			}
		},
	},
	components: { UploadMetadata },
}
</script>

<style lang="stylus" scoped>
.upload
	width: 250px
	height: 150px
	border: 2px dashed #000
	position: relative
	cursor: pointer
	display: inline-block

	> .text
		margin: 0;
		position: absolute;
		top: 50%;
		left: 50%;
		margin-right: -50%;
		transform: translate(-50%, -50%)

.file
	margin: 1.25em
	border: 2px solid #000

#fileupload
	width: 1px
	position: absolute
	top: -1
</style>
