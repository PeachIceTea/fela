<template lang="pug">
	.upload-status
		.container.upload(v-for="upload, i in uploads")
			.progress-bar(
				:style="{width: `${upload.progress * 100}%`}"
				:class="{'bar-err': upload.err, 'bar-success': !upload.err && upload.progress === 1}"
			)
			.content
				details.message.file-names
					summary Uploading {{ upload.files.length }}
						span(v-if="upload.files.length === 1")  file
						span(v-else)  files
					span(v-for="file, i in upload.files") {{ file.name }}
							span(v-if="i + 1 !== upload.files.length") ,&#32;
				.message(v-show="!upload.err && upload.progress !== 1")
					| Progress {{ (upload.progress * 100).toFixed(2) }}%
				.message(v-if="upload.err") Error: {{ upload.err }}
				.message(v-if="!upload.err && upload.progress === 1")
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

<style lang="less" scoped>
@import "../../globals.less";

.upload {
	padding: 0;
	margin-bottom: 1em;
}

.progress-bar {
	height: 1em;
	background: @offwhite;
	border-radius: @border-radius @border-radius 0 0;
}

.content {
	padding: 1em;
}

.file-names {
	margin: 0.5em 0 1em;
}

.message {
	width: 100% !important;
	max-width: none;
}

.bar-err {
	background: @err;
}

.bar-success {
	background: @success;
}

summary {
	margin: 0.5em 0;
	cursor: pointer;
	user-select: none;
}
</style>
