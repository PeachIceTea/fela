<template lang="pug">
	.upload-edit
		form(@submit.prevent="submit")
			h2 Edit book \#{{ id }}
			.input-row
				label Title
					.input-container
						input(type="text" v-model="title" placeholder="Title")
			.input-row
				label Author
					.input-container
						input(
							type="text"
							v-model="author"
							placeholder="Author"
							list="authors"
						)
			.input-row
				label New cover
					.file-input
						span(v-if="!cover") Click here to select a new cover
						span(v-else="cover") "{{ cover.name }}"
					.input-container
						input(
							type="file"
							accept="image/png, image/jpeg"
							@input="newCover"
						)
			.input-row
				input(type="submit" value="Update")

			.message.error(v-show="err") Error: {{ err }}
			.message.success(v-show="success") Audiobook was updated.
		.current-cover
			h2 Current cover
			img(
				:src="`http://localhost:8080/files/cover/${id}.jpg`"
				@error="noImage"
				alt="Current cover"
			)
		datalist#authors
			option(v-for="author in authors" :value="author")
</template>

<script>
import { updateAudiobook } from "../../api"
import PlacholderCover from "../../placeholder-cover.jpg"

export default {
	data() {
		return {
			title: "",
			author: "",
			cover: null,
			err: "",
			success: false,
		}
	},
	computed: {
		id() {
			return this.$route.params.id
		},
		authors() {
			return this.$store.getters.authors
		},
	},
	methods: {
		noImage(e) {
			e.srcElement.src = PlacholderCover
		},
		newCover(e) {
			this.cover = e.srcElement.files[0]
			e.srcElement.value = ""
		},
		async submit() {
			this.err = ""
			this.success = false
			const res = await updateAudiobook(this.id, {
				title: this.title,
				author: this.author,
				cover: this.cover,
			})

			if (res.err) {
				return (this.err = res.err)
			}

			this.success = true

			this.$store.dispatch("getAudiobooks")
			this.$store.dispatch(
				"getUserUploads",
				this.$store.state.auth.loggedInUser.id,
			)
		},
	},
	beforeRouteEnter: routeHandler,
	beforeRouteUpdate: routeHandler,
}

function routeHandler(to, from, next) {
	next(async vm => {
		await vm.$store.dispatch("getEditingAudiobook", to.params.id)
		const a = vm.$store.state.audiobook.editing
		a.title && (vm.title = a.title)
		a.author && (vm.author = a.author)
	})
}
</script>

<style lang="stylus" scoped>
@import "../../globals"
borderRadius = 3px
form-width = 249px
font-size = 20px
padding = input-padding

.upload-edit
	display: flex
	border: 1px solid offwhite
	padding: 1em

form
	flex: 1
	margin-right: 1em
	box-shadow: box-shadow

h2
	margin: 0 0 0.75em

.input-row
	margin-bottom: 1em

img
	flex: 1
	max-width: 100%
	object-fit: contain
	background: url("../../placeholder-cover.jpg")

label
	font-size: 20px
	font-weight: 500

input, .file-input, .message
	display: block
	width: 100%
	border: 0
	border-radius: 3px
	color: black-text
	padding: padding
	outline: 0
	background: offwhite
	font-size: font-size

.file-input
	cursor: pointer

	> span
		cursor: inherit

.message
	margin: 1em 0
	font-size: font-size

.error
	border-top: 5px red solid

.success
	border-top: 5px green solid

input[type="submit"]
	background: highlight
	color: white-text
	cursor: pointer
	text-shadow: text-shadow

	&::-moz-focus-inner
		border: 0

input[type="file"]
	width: 1px
	height: 1px
	position: absolute
	left: -1px
	outline: none
	padding: 0
</style>
