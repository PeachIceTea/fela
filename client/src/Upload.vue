<template lang="pug">
	.upload-container
		.file-list
			.file(v-for="file in files")
				span {{ file.name }} - Progress: {{ file.progress }}%
				.meta-editor
					input(id="title" placeholder="Book Title")
					input(id="author" placeholder="Author")
					input(id="series" placeholder="Series Name")
					input(id="series-number" placeholder="Series Number")


		label(for="fileupload")
			.upload(
				@drop="dropFile"
				@dragover.stop.prevent="dragOver"
			)
				p.text Drag file here or click to upload file.
			input(id="fileupload" type="file" multiple @input="dialogFile")
</template>

<script>
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
			// TODO: do drag and drop
			const files = e.dataTransfer.files

			for (let i = 0, len = files.length; i < len; i++) {
				let file = files[i]
			}
		},
		async dialogFile(e) {
			const files = e.srcElement.files

			for (let i = 0, len = files.length; i < len; i++) {
				let file = files[i]

				//TODO: Test for valid mime

				const status = {
					name: file.name,
					progress: 0,
					fileID: 0,
					file,
				}

				this.files.push(status)

				const form = new FormData()
				form.set("file", file)

				const xhr = new XMLHttpRequest()
				xhr.responseType = "json"

				xhr.upload.onprogress = function(e) {
					status.progress = Math.floor((e.loaded / e.total) * 100)
				}

				xhr.onerror = function(e) {
					//TODO: proper error handling
					console.error("error occured during upload", e)
				}

				xhr.onload = function(e) {
					const id = xhr.response.file_id
					if (id === 0) {
						//TODO: Deal with reupload in some smart way
						console.error("Reupload file")
						return
					}

					status.fileID = id
				}

				xhr.open("POST", "http://localhost:8080/upload")
				xhr.send(form)
			}

			e.srcElement.value = ""
		},
	},
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
