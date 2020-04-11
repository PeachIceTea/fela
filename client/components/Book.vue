<template lang="pug">
	.book
		p {{book.title}} by {{book.author}}
		br
		p Audiobooks
		p(v-for="audiobook, i in audiobooks")
			button(@click="play(audiobook)") Play Audiobook
</template>

<script>
export default {
	computed: {
		book() {
			return this.$store.state.book.book
		},
		audiobooks() {
			return this.$store.state.book.audiobooks
		},
	},
	methods: {
		play(audiobook) {
			this.$store.dispatch("player/getFiles", audiobook.id)
		},
	},
	beforeRouteEnter(to, from, next) {
		next(async vm => await vm.$store.dispatch("book/getBook", to.params.id))
	},
	beforeRouteUpdate(to, from, next) {
		next(async vm => await vm.$store.dispatch("book/getBook", to.params.id))
	},
}
</script>
