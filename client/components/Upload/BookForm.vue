<template lang="pug">
	.book-form
		.file-name(v-if="book") {{book.name}}
		form.book-form(@submit.prevent="submit")
			label
				p Title
				p
					input(type="text" v-model="title")
			label
				p Author
				p
					input(type="text" v-model="author")
			//label
				p Series
				p
					input(type="text" v-model="series")
			// label
				p No
				p
					input(type="text" v-model="no")
			label
				p
					input(type="submit" value="Submit" :disabled="disabled")
</template>

<script>
export default {
	data() {
		return {
			title: "",
			author: "",
			series: "",
			no: "",
		}
	},
	computed: {
		disabled() {
			return !(this.title && this.author)
		},
	},
	methods: {
		submit() {
			this.callback(this.title, this.author, this.book)
		},
	},
	props: {
		book: File,
		callback: Function,
	},
}
</script>
