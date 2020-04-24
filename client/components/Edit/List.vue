<template lang="pug">
	.upload-table
		.header
			.selector
				select(v-model="activeList")
					option(value="user") My Uploads
					option(value="all") All Uploads
			input.search(
				type="text"
				placeholder="Search"
				v-model="search")
		.table-container

			div(v-for="upload in tableUploads" :key="upload.id")
				.title {{ upload.title }}
				.author by {{ upload.author }}
				.edit
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
						button.button(@click="deleteUpload") Yes, I am sure.
						button.button(@click="cancelDelete") Cancel
</template>

<script>
import { deleteAudiobook } from "../../api"

export default {
	data() {
		return {
			activeList: "user", // "user" or "all"
			search: "",
			potentialDelete: null,
			deleteErr: "",
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
			const userID = this.$store.state.auth.loggedInUser.id
			return this.$store.state.audiobook.list.filter(
				a => a.uploader === userID,
			)
		},
		allUploads() {
			return this.$store.state.audiobook.list
		},
	},
	created() {
		this.$store.dispatch("getAudiobooks")
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
		async deleteUpload() {
			const res = await deleteAudiobook(this.potentialDelete.id)
			if (res.err) {
				this.deleteErr = res.err
			}

			this.$store.dispatch("getAudiobooks")
			this.potentialDelete = null
		},
	},
}
</script>

<style lang="less" scoped>
@import "../../globals.less";
</style>
