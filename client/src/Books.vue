<template lang="pug">
	.books
		.book(v-for="book in books")
			div
				span
					router-link(:to="{ name: 'book', params: { id: book.id } }") {{ book.title }}
					|  by {{ book.author }}
			br

</template>

<script>
export default {
	data() {
		return {
			books: null,
			error: null,
		}
	},
	async beforeRouteEnter(from, to, next) {
		try {
			const res = await fetch("http://localhost:8080/book")

			if (res.status === 200) {
				const books = await res.json()
				next(vm => {
					vm.books = books
				})
			} else {
				//TODO: Error handling
			}
		} catch (e) {}
	},
}
</script>
