<template lang="pug">
	.audiobook
		h2 Files
		p(v-for="file in files")
			a(:href="`http://localhost:8080/files/${file.hash}`") {{file.name}}
</template>

<script>
export default {
	computed: {
		files() {
			return this.$store.state.book.files
		},
	},
	beforeRouteEnter(to, from, next) {
		next(
			async vm =>
				await vm.$store.dispatch("book/getAudiobook", to.params.id),
		)
	},
	beforeRouteUpdate(to, from, next) {
		next(
			async vm =>
				await vm.$store.dispatch("book/getAudiobook", to.params.id),
		)
	},
}
</script>
