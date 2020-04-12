<template lang="pug">
	.book
		.info
			h2 {{book.title}} by {{book.author}}
		.audiobooks
			h2 Audiobooks
			table
				tr
					th ID
					th Codec
					th Duration
				tr(v-for="audiobook in audiobooks" @click="play(audiobook)")
					th {{audiobook.id}}
					th {{audiobook.codec}}
					th {{audiobook.duration | formatDuration}}
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

<style lang="stylus" scoped>
table
	border-spacing: 0
	width: 100%

th
	text-align: left

tr:not(:first-child)&:hover
	background: #333
	color: white
	cursor: pointer
</style>
