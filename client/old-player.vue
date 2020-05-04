<template lang="pug">
	.player(v-show="audiobook.files")
		.progressbar(
			ref="progress"
			@click="handleProgressClick"
			@mouseenter="showProgressInfo"
			@mousemove="moveProgressInfo"
			@mouseleave="hideProgressInfo"
		)
			.played(:style="{width: `${currentTime/totalDuration*100}%`}")
		.content
			.col
				.cover
					img(
						:src="coverURL(audiobook.id)"
						@error="noImage"
					)
				.book-info(style="display: inline")
					p {{ audiobook.title }}
					p by {{ audiobook.author }}
			.col
				.control(@click="rewind")
					Rewind.control-symbol
				.control(@click="toggle")
					Pause.control-symbol(v-show="!paused")
					Play.control-symbol(v-show="paused")
				.control(@click="forward")
					FastForward.control-symbol
			.col
				| {{ currentTime | duration }} /
				|  {{ totalDuration | duration }}
		.progress-info(
			ref="progressInfo"
			v-show="progressInfo.show"
			:style="{top: `${progressInfo.y}px`,left: `${progressInfo.x}px`}"
		)
			| {{ progressInfo.value | duration }}
		audio(:src="fileUrl" ref="audio")
</template>

<script>
import { coverURL, audiobookURL } from "../api"
import Play from "./svg/Play.vue"
import Pause from "./svg/Pause.vue"
import FastForward from "./svg/FastForward.vue"
import Rewind from "./svg/Rewind.vue"
import PlacholderCover from "../placeholder-cover.jpg"
export default {
	data() {
		return {
			paused: true,
			time: 0,
			volume: 1,
			progressInfo: {
				show: false,
				value: 0,
				x: 0,
			},
			file: 0,
		}
	},
	computed: {
		audiobook() {
			return this.$store.state.audiobook.playing
		},
		currentFile() {
			if (this.audiobook.files) {
				return this.audiobook.files[this.file]
			} else {
				return false
			}
		},
		currentChapterIndex() {
			const chapters = this.currentFile.metadata.chapters
			for (let i = 0, len = chapters.length; i < len; i++) {
				let current = chapters[i],
					next = chapters[i + 1]
				if (!next) {
					return i
				}

				const startTime = parseFloat(current.start_time)
				const endTime = parseFloat(current.end_time)
				if (this.time >= startTime && this.time < endTime) {
					return i
				}
			}
		},
		// Returns the metadata for the current chapter. Only usable for single
		// file audiobooks with chapter information.
		currentChapter() {
			return this.currentFile.metadata.chapters[this.currentChapterIndex]
		},
		duration() {
			if (this.currentFile) {
				if (this.multipleFiles) {
					return this.currentFile.duration
				} else {
					return (
						parseFloat(this.currentChapter.end_time) -
						parseFloat(this.currentChapter.start_time)
					)
				}
			} else {
				return 0
			}
		},
		currentTime() {
			if (this.multipleFiles) {
				return this.time
			} else {
			}
		},
		fileUrl() {
			const audiobook = this.audiobook
			if (audiobook.files) {
				return audiobookURL(this.audiobook.id, this.currentFile.name)
			}
		},
		multipleFiles() {
			if (this.audiobook.files) {
				return this.audiobook.files.length > 1
			} else {
				return false
			}
		},
		chapterName() {
			if (this.currentFile) {
				if (this.multipleFiles) {
					// Test if we have a chapter object that includes a title
					// for the chapter. Its not guranteed to exists or to have
					// any tags or to have a title tag. What a mess.
					const chapterObj = this.currentFile.metadata.chapters[0]
					if (
						chapterObj &&
						chapterObj.tags &&
						chapterObj.tags.title
					) {
						return chapterObj.tags.title
					}

					// A hack to see if the files title tag actually is the
					// title of the chapter and not of the entire book. The
					// first and second book not having the same title does
					// obviously not gurantee that they contain the title of
					// the chapter.
					const firstTitle = this.audiobook.files[0].metadata.format
						.tags.title
					const nextTitle = this.audiobook.files[1].metadata.format
						.tags.title
					if (firstTitle && firstTitle !== nextTitle) {
						return this.currentFile.metadata.format.tags.title
					} else {
						return `Chapter ${this.chapter + 1}`
					}
				} else {
					let chapter

					// The only way with a single file audiobook to know the
					// chapters is to read them from the metadata. Sadly not
					// all formats (namely MP3) support chapter metadata. So,
					// it is not save to assume that we have any chapter
					// information.
					const chapterObj = this.currentFile.metadata.chapters[
						this.currentChapter
					]
					if (
						chapterObj &&
						chapterObj.tags &&
						chapterObj.tags.title
					) {
						return chapterObj.tags.title
					}
				}
			}
		},
	},
	methods: {
		rewind() {
			this.seek(this.time - 30)
		},
		toggle() {
			const el = this.$refs.audio
			if (el) {
				el.paused ? el.play() : el.pause()
			}
		},
		forward() {
			this.seek(this.time + 30)
		},
		increaseVolume() {
			this.setVolume(this.volume + 0.1)
		},
		decreaseVolume() {
			this.setVolume(this.volume - 0.1)
		},

		// Handles keyboard events for media control.
		keyHandler(e) {
			if (e.srcElement.tagName !== "INPUT") {
				switch (e.key) {
					case " ":
						e.preventDefault()
						return this.toggle()
					case "ArrowLeft":
						e.preventDefault()
						return this.rewind()
					case "ArrowRight":
						e.preventDefault()
						return this.forward()
					case "ArrowUp":
						e.preventDefault()
						return this.increaseVolume()
					case "ArrowDown":
						e.preventDefault()
						return this.decreaseVolume()
				}
			}
		},

		// Moves to a given position in the audio
		seek(time) {
			this.$refs.audio.time = time
		},

		// Helps
		handleProgressClick(e) {
			const bounds = this.$refs.progress.getBoundingClientRect()
			this.seek((this.duration * (e.clientX - bounds.x)) / bounds.width)
		},

		playbackEnded(e) {
			if (this.multipleFiles) {
				this.file++
			}
		},

		// Moves the timer while hovering over the progressbar
		moveProgressInfo(e) {
			const bounds = this.$refs.progress.getBoundingClientRect()
			this.progressInfo.value = (this.duration * e.clientX) / bounds.width
			this.progressInfo.x =
				e.clientX -
				this.$refs.progressInfo.getBoundingClientRect().width / 2
		},

		showProgressInfo() {
			this.progressInfo.show = true
		},

		hideProgressInfo() {
			this.progressInfo.show = false
		},
		noImage(e) {
			e.srcElement.src = PlacholderCover
		},
		coverURL,

		// Unimplemented volume controls
		// TODO: implement
		setVolume(vol) {
			this.$refs.audio.volume = vol
			this.volume = vol
		},
		handleVolumeClick(e) {
			const volume = this.$refs.volume
			const bounds = volume.getBoundingClientRect()
			const setVolume = e => {
				this.setVolume(
					Math.min(
						1,
						Math.max(0, (e.clientX - bounds.x) / bounds.width),
					),
				)
			}
			const removeListener = () => {
				document.removeEventListener("mousemove", setVolume)
				document.removeEventListener("mouseup", removeListener)
			}
			setVolume(e)
			document.addEventListener("mousemove", setVolume)
			document.addEventListener("mouseup", removeListener)
		},
	},
	created() {
		document.addEventListener("keydown", this.keyHandler)
	},
	mounted() {
		const el = this.$refs.audio
		el.addEventListener("timeupdate", e => {
			this.time = el.time
		})
		el.addEventListener("play", e => (this.paused = false))
		el.addEventListener("pause", e => (this.paused = true))
		el.addEventListener("ended", this.playbackEnded)
	},
	destroyed() {
		document.removeEventListener("keydown", this.keyHandler)
	},
	components: { Play, Pause, FastForward, Rewind },
}
</script>

<style lang="less" scoped>
@import "../globals.less";

.player {
	width: 100%;
	background: darken(@background, 5%);
}

.progressbar {
	height: 0.5em;
	width: 100%;
	background: lighten(@background, 5%);
	transition: 250ms all ease;
	cursor: pointer;

	&:hover {
		height: 1.25em;
	}
}

.content {
	display: flex;
	flex-direction: row;
	padding: 1em;
}

.col {
	flex: 1;
	display: flex;
	flex-direction: row;
	align-items: center;
	justify-content: center;
}

.middle {
	flex-direction: column;
	align-items: unset;
	justify-content: unset;
}

.book-info {
	flex: 1;
	font-size: 20px;
}

.playback-info {
	flex: 1;
}

p {
	margin: 0.5em 0;
}

.cover {
	max-height: 100px;
	margin-right: 1em;

	&:last-child {
		justify-content: flex-end;
	}
}

img {
	max-height: inherit;
	max-width: inherit;
	border-radius: @border-radius;

	.boxShadow();
}

.controls {
	display: flex;
}

.timer {
	text-align: center;
	margin-bottom: 1em;
	font-size: 20px;
}

.control {
	flex: 1;
	display: flex;
	transition: 100ms all ease;
	cursor: pointer;
	fill: @black-text;
	justify-content: center;
	font-size: 20px;

	&:hover {
		fill: @white-text;
	}
}

.control-symbol {
	height: 100%;
	width: 2em;
	cursor: pointer;
}

.progress-info {
	position: absolute;
	bottom: 10em;
	text-shadow: 2px 2px 3px rgba(0, 0, 0, 1);
}

.played {
	height: 100%;
	background: @offwhite;
}
</style>
