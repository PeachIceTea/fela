<template lang="pug">
	.library
		table
			tr
				th Title
				th Author
			tr(v-for="book in books" @click="nav(book)")
				td {{book.title}}
				td {{book.author}}
</template>

<script>
export default {
	data() {
		return {
			loading: false,
		}
	},
	computed: {
		books() {
			return this.$store.state.book.books
		},
	},
	methods: {
		async getBooks() {
			this.loading = true
			await this.$store.dispatch("book/getBooks")
			this.loading = false
		},
		nav(book) {
			this.$router.push(`/book/${book.id}`)
		},
	},
	created() {
		this.getBooks()
	},
	watch: {
		$route: "getBooks",
	},
}
</script>

<style lang="stylus" scoped>
.library
	padding: 0.5em

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
