<template lang="pug">
	.book
		div
			p {{ book.title }} by {{ book.author }}
			p {{ book.description }}

		.file(v-for="file in files")
			a(:href="'http://localhost:8080/files/' + file.hash") {{ file.name }}
</template>

<script>
export default {
	data() {
		return {
			book: {},
			files: [],
			error: null,
		}
	},
	async beforeRouteEnter(to, from, next) {
		try {
			const res = await fetch(
				`http://localhost:8080/book/${to.params.id}`,
			)

			if (res.status === 200) {
				const body = await res.json()
				next(vm => {
					vm.book = body.book
					vm.files = body.files
				})
			} else {
				//TODO: Error handling
				console.error(res)
			}
		} catch (e) {
			console.error(e)
		}
	},
}
</script>
