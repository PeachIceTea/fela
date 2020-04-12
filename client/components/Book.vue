<template lang="pug">
	.book
		.info
			h2 {{book.title}} by {{book.author}}
		.audiobooks
			h2 Audiobooks
			.list
				button.item(v-for="audiobook in audiobooks" @click="play(audiobook)") Play Audiobook \#{{audiobook.id}}
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
			this.$store.dispatch("player/play", { audiobook, book: this.book })
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
