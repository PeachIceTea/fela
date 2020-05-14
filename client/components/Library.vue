<template lang="pug">
	.library
		.cover-library(v-show="coverViewActive")
			.cover(
				v-for="audiobook in audiobooks"
				@click="play(audiobook)"
			)
				img(
					:src="coverURL(audiobook.id)"
					@error="noImage"
				)
		.list-library(v-show="!coverViewActive")
			.container.item(v-for="audiobook in audiobooks")
				.small-cover
					img(
						:src="coverURL(audiobook.id)"
						@error="noImage"
					)
				.text
					.title {{ audiobook.title }}
					.by by
						span.author  {{ audiobook.author }}
				.buttons
					button(@click="edit(audiobook)" v-if="uploader") Edit
					button(@click="play(audiobook)") Play
</template>

<script>
import Fuse from "fuse.js"
import { coverURL } from "../api"
import PlacholderCover from "../placeholder-cover.jpg"

export default {
	computed: {
		// Retrieves the audiobook list from the store. Either sorts it using
		// ui.order or uses ui.search to filter the list. When searching the
		// order is always best to worst matches.
		audiobooks() {
			const arr = this.$store.state.audiobook.list.concat()
			if (this.search) {
				return new Fuse(arr, {
					keys: ["title", "author"],
				})
					.search(this.search)
					.map(e => {
						return e.item
					})
			} else {
				switch (this.$store.state.ui.order) {
					case 0:
						return arr.sort((a, b) => {
							if (a.title < b.title) return -1
							if (a.title > b.title) return 1
							return 0
						})
					case 1:
						return arr.sort((a, b) => {
							if (a.title < b.title) return 1
							if (a.title > b.title) return -1
							return 0
						})
					case 2:
						return arr.sort((a, b) => {
							const aDate = new Date(a.created_at)
							const bDate = new Date(b.created_at)
							if (aDate < bDate) return -1
							if (aDate > bDate) return 1
							return 0
						})
					case 3:
						return arr.sort((a, b) => {
							const aDate = new Date(a.created_at)
							const bDate = new Date(b.created_at)
							if (aDate < bDate) return 1
							if (aDate > bDate) return -1
							return 0
						})
				}
			}
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
	},
	methods: {
		// Moves the audiobook to the actively played on in the store.
		play(audiobook) {
			this.$store.dispatch("playAudiobook", audiobook.id)
		},

		// Moves user to the edit page for a given audiobook.
		edit(audiobook) {
			this.$router.push(`/audiobook/${audiobook.id}/edit`)
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
	},
}
</script>

<style lang="less" scoped>
@import "../globals.less";

.list-library {
	padding: 1em;
}

.item {
	margin-bottom: 1em;
	display: flex;
	height: 15em;
}

.text {
	align-self: center;
	text-shadow: @text-shadow;
}

.title {
	font-size: 25px;
	font-weight: 600;
}

.small-cover {
	height: 100%;
	margin-right: 1em;
	img {
		max-height: 100%;
		max-width: 100%;
	}
}

.buttons {
	flex: 1;
	display: flex;
	justify-content: flex-end;
	align-items: center;

	button {
		height: 5em;
		width: 8em;

		&:last-child {
			margin-left: 1em;
		}
	}
}

.cover-library {
	display: grid;
	grid-template-columns: repeat(3, 1fr);
}

// The ::before on .book and position absolute on img are there to force the
// covers into a 1:! aspect ratio in which most audiobooks covers have. This
// prevents odd shaped covers from distorting other covers in their row.
// Source: https://stackoverflow.com/a/20117454/13394124
.cover {
	cursor: pointer;
	position: relative;
	overflow: hidden;

	&::before {
		content: "";
		display: block;
		padding-top: 100%;
	}

	img {
		position: absolute;
		top: 0;
		height: 100%;
		width: 100%;
		transition: 500ms all ease;

		&:hover {
			/* filter: blur(5px) brightness(0.75); */
			transform: scale(1.03);
		}
	}
}

@media (min-width: 1400px) {
	.cover-library {
		display: grid;
		grid-template-columns: repeat(4, 1fr);
	}
}
</style>
