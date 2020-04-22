<template lang="pug">
	.upload-status
		.upload(v-for="upload, i in uploads")
			.progress-bar(:style="{width: `${upload.progress * 100}%`}")
			.content
				.file-names.collapsed {{ upload.files.length }} file
					span(v-if="upload.files.length !== 1") s
					| :&#32;
					span(v-for="file, i in upload.files") "{{ file.name }}"
						span(v-if="i + 1 !== upload.files.length") ,&#32;
				.message.err(v-if="upload.err") Error: {{ upload.err }}
				.message.success(v-if="!upload.err && upload.progress === 1")
					| Upload finished
</template>

<script>
export default {
	computed: {
		uploads() {
			const uploads = this.$store.state.uploads
			const reversed = []
			for (let i = uploads.length; i > 0; i--) {
				reversed.push(uploads[i - 1])
			}
			return reversed
		},
	},
}
</script>

<style lang="stylus" scoped>
@import "../../globals.styl"

.upload
	border: 1px offwhite solid
	margin-bottom: 1em

.progress-bar
	height: 0.5em
	background: offwhite

.content
	padding: 1em

.collapsed
	overflow: hidden
	white-space: nowrap
	text-overflow: ellipsis
	width: 100%
	position:relative
	cursor: pointer
</style>
