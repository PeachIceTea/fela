<template lang="pug">
	.upload
		.header
			.upload-no Upload \#{{index+1}}
			.upload-progress Progress {{upload.progress}}%
		.content
			.question.chapterized(
				v-if="upload.files.length !== 1 && upload.chapterized === undefined"
			)
				button(@click="setChapterized(false)") Multiple Books
				button(@click="setChapterized(true)") Single Book
			div(v-else)
				BookForm(v-if="upload.files.length === 1 || upload.chapterized" :callback="submitCallback" :upload="upload")
				BookForm(v-else v-for="book in upload.files" :book="book" :callback="submitCallback" :upload="upload")

</template>

<script>
import BookForm from "./BookForm"

export default {
	props: {
		upload: Object,
		index: Number,
	},
	methods: {
		setChapterized(chapterized) {
			this.$store.commit("upload/setChapterized", {
				index: this.index,
				chapterized,
			})
		},
		submitCallback(title, author, book) {
			if (this.upload.chapterized) {
				this.$store.dispatch("upload/submitBook", {
					title,
					author,
					index: this.index,
				})
			} else {
				this.$store.dispatch("upload/submitBook", {
					title,
					author,
					index: this.index,
					book,
				})
			}
		},
	},
	components: { BookForm },
}
</script>
