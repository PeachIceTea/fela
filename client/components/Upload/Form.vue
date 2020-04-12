<template lang="pug">
	.upload-form
		.upload-list
			Upload(v-for="upload, i  in uploads" :upload="upload" :index="i")

		label(for="fileupload")
			.upload-box(
				@drop.stop.prevent="dropFile"
				@dragover.stop.prevent="dragOver"
			)
				p.text Drag file here or click to upload file.
			input(id="fileupload" type="file" multiple @input="dialogFile")
</template>

<script>
import Upload from "./Upload.vue"

export default {
	methods: {
		dropFile(e) {
			const files = e.dataTransfer.files
			this.$store.dispatch("upload/upload", files)
			window.files = files
		},
		dragOver(e) {
			e.dataTransfer.dropEffect = "copy"
		},
		dialogFile(e) {
			const files = e.srcElement.files
			this.$store.dispatch("upload/upload", files)
			e.srcElement.value = ""
		},
	},
	computed: {
		uploads() {
			return this.$store.state.upload.uploads
		},
	},
	created() {
		console.log("hi")
		this.$store.dispatch("book/getAuthors")
	},
	components: { Upload },
}
</script>

<style lang="stylus" scoped>
.upload-box
	width: 250px
	height: 150px
	border: 2px dashed #fff
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
#fileupload
	width: 1px
	position: absolute
	top: -1
</style>
