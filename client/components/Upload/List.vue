<template lang="pug">
	.upload-table
		.selector
			.tab-switcher(
				:class="{active: activeTab === 'user'}"
				@click="switchTab('user')") Your uploads
			.tab-switcher(
				:class="{active: activeTab === 'all'}"
				@click="switchTab('all')") All uploads
		table
			tr
				th ID
				th Title
				th Author
				th
					input.search(
						type="text"
						placeholder="Search"
						v-model="search")
			tr(v-for="upload in tableUploads")
				td {{ upload.id }}
				td {{ upload.title }}
				td {{ upload.author }}

				td.edit
					router-link.btn(:to="`/upload/edit/${upload.id}`") Edit
					span.btn(@click="showDeleteDialog(upload)") Delete

		.delete-container(v-if="potentialDelete")
			.delete-dialog
				.text
					p Are you sure you want to delete the audiobook
					p
						b.
							"{{ potentialDelete.title }} by
							{{ potentialDelete.author }}"?
					p This action is irrreversible!
				.options
					button(@click="deleteUpload") Yes, I am sure.
					button(@click="cancelDelete") Cancel
</template>

<script>
export default {
	data() {
		return {
			activeTab: "user", // "user" or "all"
			search: "",
			potentialDelete: null,
		}
	},
	computed: {
		tableUploads() {
			const search = this.search.toLowerCase()
			let list =
				this.activeTab === "user" ? this.userUploads : this.allUploads

			// Very, very basic filtering
			if (this.search) {
				list = list.filter(el => {
					return `${el.id} ${el.title} ${el.author}`
						.toLowerCase()
						.includes(search)
				})
			}

			return list
		},
		userUploads() {
			return this.$store.state.user.uploads
		},
		allUploads() {
			return this.$store.state.audiobook.list
		},
	},
	created() {
		this.$store.dispatch("getAudiobooks")
		this.$store.dispatch(
			"getUserUploads",
			this.$store.state.auth.loggedInUser.id,
		)
	},
	methods: {
		switchTab(to) {
			this.activeTab = to
		},
		showDeleteDialog(upload) {
			this.potentialDelete = upload
		},
		cancelDelete() {
			this.potentialDelete = null
		},
		deleteUpload() {},
	},
}
</script>

<style lang="stylus" scoped>
@import "../../globals"

border-width = 1px

.selector
	display: flex

.tab-switcher
	flex: 1
	padding: 0.25em
	font-size: 2em
	cursor: pointer
	border: border-width offwhite solid
	margin-bottom: - border-width

	&:not(.active)
		color: black-text

	&.active
		border-bottom: border-width background solid

	&:last-child
		margin-left: -1px

table
	width: 100%
	border-collapse: collapse
	border: border-width offwhite solid
	border-top: 0

	tr:last-of-type .edit
		border: 0

th
	text-align: left


td, th
	border-bottom: border-width solid offwhite
	padding: 8px 5px 3px

th:last-child
	width: 0

td:last-child
	display: flex

.search
	margin-right: 1em
	display: inline

.btn
	flex: 1
	color: offwhite
	font-weight: 600
	cursor: pointer
	text-decoration: underline

.delete-container
	position: absolute
	height: 100vh
	width: 100vw
	top: 0
	left: 0
	background: rgba(0,0,0,0.8)
	display: flex
	align-items: center
	justify-content: center

.delete-dialog
	background: background
	width: 750px
	text-align: center
	padding: 2em
	border-radius: 4px
	box-shadow: box-shadow

	p:first-child
		margin-top: 0

	.options
		display: flex

	button
		cursor: pointer
		flex: 1
		margin: 0 4em
		background: highlight
		color: white-text
		outline: 0
		border: 0
		box-shadow: box-shadow
		font-size: 20px
		border-radius: 3px
		text-shadow: text-shadow
		padding: input-padding

		&::-moz-focus-inner
			border: 0

	.text
		font-size: 20px
</style>
