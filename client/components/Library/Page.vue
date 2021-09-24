<template lang="pug">
	.library
		.search(v-if="search")
			.cover-container(v-show="coverViewActive")
				CoverList(:audiobooks="audiobooks" :play="play" :coverURL="coverURL" :noImage="noImage")
			.list-library(v-show="!coverViewActive")
				VerticalList(:audiobooks="audiobooks" :play="play" :coverURL="coverURL" :noImage="noImage" :uploader="uploader")
		div(v-if="!search")
			.cover-container(v-show="coverViewActive")
				.library-header Started Books
				CoverList(:audiobooks="audiobooks.started" :play="play" :coverURL="coverURL" :noImage="noImage")
				.library-header Unread Books
				CoverList(:audiobooks="audiobooks.unread" :play="play" :coverURL="coverURL" :noImage="noImage")
			.list-library(v-show="!coverViewActive")
				.library-header Started Books
				VerticalList(:audiobooks="audiobooks.started" :play="play" :coverURL="coverURL" :noImage="noImage" :uploader="uploader")
				.library-header Unread Books
				VerticalList(:audiobooks="audiobooks.unread" :play="play" :coverURL="coverURL" :noImage="noImage" :uploader="uploader")
</template>

<script>
import Fuse from "fuse.js"
import { coverURL } from "../../api"
import PlacholderCover from "../../images/placeholder-cover.jpg"

import CoverList from "./CoverList.vue"
import VerticalList from "./VerticalList.vue"

export default {
	computed: {
		// Retrieves the audiobook list from the store. Either sorts it using
		// ui.order or uses ui.search to filter the list. When searching the
		// order is always best to worst matches.
		audiobooks() {
			let arr = this.$store.state.audiobook.list.concat() // concat to create new array

			if (this.search) {
				return new Fuse(arr, {
					keys: ["title", "author"],
				})
					.search(this.search)
					.map(e => {
						return e.item
					})
			}

			const [started, unread] = [[], []]
			arr.forEach(audiobook => {
				if (
					this.userProgress.find(
						progress => progress.audiobook === audiobook.id,
					)
				) {
					started.push(audiobook)
				} else {
					unread.push(audiobook)
				}
			})

			switch (this.$store.state.ui.order) {
				case 0:
					started.sort(sort.titleDesc)
					unread.sort(sort.titleDesc)
					break
				case 1:
					started.sort(sort.titleAsc)
					unread.sort(sort.titleAsc)
					break
				case 2:
					started.sort(sort.uploadAsc)
					unread.sort(sort.uploadAsc)
					break
				case 3:
					started.sort(sort.titleDesc)
					unread.sort(sort.titleDesc)
					break
				case 4:
					started.sort(sort.defaultStarted(this.userProgress))
					unread.sort(sort.titleDesc)
					break
			}
			return { started, unread }
		},

		search() {
			return this.$store.state.ui.search
		},
		uploader() {
			return this.$store.getters.isUploader
		},

		coverViewActive() {
			return this.$store.state.ui.view === 0
		},

		userProgress() {
			return this.$store.state.audiobook.userProgress
		},
		showUnread() {
			return this.$store.state.ui.showUnread
		},
	},
	methods: {
		// Moves the audiobook to the actively played on in the store.
		play(audiobook) {
			this.$store.dispatch("playAudiobook", audiobook.id)
		},

		// Called when an cover fails to load, this happens when no cover has
		// been uploaded. Replaces the source of the image with the placeholder.
		noImage(e) {
			e.srcElement.src = PlacholderCover
		},

		// Makes coverURL api helper usable in the Vuejs template
		coverURL,
	},
	created() {
		this.$store.dispatch("getAudiobooks")
		this.$store.dispatch("getUserProgress")
	},
	components: { CoverList, VerticalList },
}

const sort = {
	defaultStarted(userProgress) {
		return (a, b) => {
			const recentA = userProgress.find(tmp => tmp.audiobook === a.id)
			const recentB = userProgress.find(tmp => tmp.audiobook === b.id)
			const aDate = recentA.updated_at
				? new Date(recentA.updated_at)
				: new Date(recentA.created_at)
			const bDate = recentB.updated_at
				? new Date(recentB.updated_at)
				: new Date(recentB.created_at)
			if (aDate < bDate) return 1
			if (aDate > bDate) return -1
			return 0
		}
	},
	titleDesc(a, b) {
		if (a.title < b.title) return -1
		if (a.title > b.title) return 1
		return 0
	},
	titleAsc(a, b) {
		if (a.title < b.title) return 1
		if (a.title > b.title) return -1
		return 0
	},
	uploadAsc(a, b) {
		const aDate = new Date(a.created_at)
		const bDate = new Date(b.created_at)
		if (aDate < bDate) return -1
		if (aDate > bDate) return 1
		return 0
	},
	uploadDesc(a, b) {
		const aDate = new Date(a.created_at)
		const bDate = new Date(b.created_at)
		if (aDate < bDate) return 1
		if (aDate > bDate) return -1
		return 0
	},
}
</script>

<style lang="less" scoped>
.library-header {
	font-size: 34px;
	font-weight: 600;
	margin: 0.4em 0 0.3em 0.5em;
}
</style>
