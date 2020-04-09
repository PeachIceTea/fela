<template lang="pug">
	.book-form(v-if="!formDone")
		.file-name(v-if="book") {{book.name}}
		form.book-form(@submit.prevent="submit")
			label
				p Title
				p
					input(type="text" v-model="title")
			label
				p Author
				p
					input(type="text" v-model="author" list="authors")
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
					input(type="submit" value="Submit" :disabled="disabledButton")
			datalist#authors
				option(v-for="author in authors" :value="author")
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
		formDone() {
			if (this.book) {
				return this.book.done
			} else {
				return this.upload.done
			}
		},
		disabledButton() {
			return !(this.title && this.author)
		},
		authors() {
			return this.$store.state.book.authors
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
		upload: Object,
	},
}
</script>
