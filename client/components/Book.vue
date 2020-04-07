<template lang="pug">
	.book
		p {{book.title}} by {{book.author}}
		br
		p Audiobooks
		p(v-for="audiobook, i in audiobooks") Audiobook {{audiobook.id}}
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
	beforeRouteEnter(to, from, next) {
		next(async vm => await vm.$store.dispatch("book/getBook", to.params.id))
	},
	beforeRouteUpdate(to, from, next) {
		next(async vm => await vm.$store.dispatch("book/getBook", to.params.id))
	},
}
</script>
