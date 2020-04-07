<template lang="pug">
	.book
		p(v-for="book in books")
			router-link(:to="`/book/${book.id}`") {{book.title}}
			|  by {{book.author}}
</template>

<script>
export default {
	computed: {
		books() {
			return this.$store.state.book.books
		},
	},
	beforeRouteEnter(to, from, next) {
		next(async vm => await vm.$store.dispatch("book/getBooks"))
	},
}
</script>
