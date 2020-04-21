<template lang="pug">
	.upload-status
		.upload(v-for="upload, i in uploads")
			.file-names Files:
				.name(v-for="file in upload.files") {{ file.name }}
			.progress Progress: {{ upload.progress * 100 }}%
			.err(v-if="upload.err") Error: {{ upload.err }}
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
