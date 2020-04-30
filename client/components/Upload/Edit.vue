<template lang="pug">
	.upload-edit
		form(@submit.prevent="submit")
			h2 Edit book \#{{ id }}
			.input-row
				label Title
					.input-container
						input(
							type="text"
							v-model="title"
							placeholder="Title"
						)
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
						span(v-if="!cover")
							| Click here to select a new cover
						span(v-else="cover") "{{ cover.name }}"
					input(
						type="file"
						accept="image/png, image/jpeg"
						@input="newCover"
					)
			.input-row
				input(type="submit" value="Update")

			.input-row
				button(@click.prevent="deleteAudiobook")
					span(v-if="!deleting") Delete Audiobook
					span(v-else) Are you sure?

			.message.err(v-show="err") Error: {{ err }}
			.message.success(v-show="success") Audiobook was updated.
		.current-cover
			h2 Current cover
			img(
				:src="coverURL(id)"
				@error="noImage"
				alt="Current cover"
			)
		datalist#authors
			option(v-for="author in authors" :value="author")
</template>

<script>
import { updateAudiobook, deleteAudiobook, coverURL } from "../../api"
import PlacholderCover from "../../placeholder-cover.jpg"

export default {
	data() {
		return {
			title: "",
			author: "",
			cover: null,
			err: "",
			success: false,
			deleting: 0,
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
				this.$store.state.auth.loggedIn.id,
			)
		},
		async deleteAudiobook() {
			if (!this.deleting) {
				this.deleting = setTimeout(() => (this.deleting = 0), 2000)
				return
			}

			clearTimeout(this.deleting)
			const res = await deleteAudiobook(this.id)
			if (res.err) {
				this.err = res.err
			}
			this.$store.dispatch("getAudiobooks")
			this.$router.push("/")
		},
		coverURL,
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

<style lang="less" scoped>
@import "../../globals.less";

.upload-edit {
	display: flex;
	justify-content: center;
	margin: 1em;
	.container();
}
form {
	flex: 1;
	display: inline-block;
	margin-right: 1em;
}

.current-cover {
	display: inline-block;
}

img {
	height: 500px;
	border-radius: @border-radius;
	.boxShadow();
}

h2 {
	font-size: 28px;
}

.input-row {
	margin-bottom: 1em;
}

input,
.file-input {
	border: 0;
	padding: @input-padding;
	width: 100%;
	max-width: 700px;
	border-radius: @border-radius;
	font-size: 20px;
}

input[type="text"],
.file-input {
	.boxShadow();
}

button {
	width: 100%;
	max-width: 700px;
}

.file-input {
	background: @offwhite;
	color: @black-text;
	cursor: pointer;
}

.message {
	.container();

	background: @offwhite;
	color: @black-text;
	padding-top: calc(1em - 4px);
	border-top: 8px solid transparent;
}

.err {
	border-top-color: @err;
}

.success {
	border-top-color: @success;
}

.bar-err {
	background: @err;
}

.bar-success {
	background: @success;
}
</style>
