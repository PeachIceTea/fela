<template lang="pug">
	.player(v-show="audiobook.files")
		.progressbar(
			ref="progressBar"
			@click="handleSeekClick"
			@mouseenter="showSeekInfo"
			@mousemove="moveSeekInfo"
			@mouseleave="hideSeekInfo"
		)
			.played(:style="{width: `${chapterTime / chapterDuration * 100}%`}")
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
				.chapter(@click="toggleChapterList") {{ chapterName }}
			.col.middle
				.timer
					| {{ chapterTime | duration }} /
					| {{ chapterDuration | duration }}
				.controls
					.control(@click="previousChapter")
						Previous.control-symbol
					.control(@click="rewind")
						Rewind.control-symbol
					.control(@click="toggle")
						Pause.control-symbol(v-show="!paused")
						Play.control-symbol(v-show="paused")
					.control(@click="forward")
						FastForward.control-symbol
					.control(@click="nextChapter")
						Next.control-symbol
			.col
				.volume(@click="handleVolumeClick" ref="volume")
					.volume-slider(:style="{width:`${volume*100}%`}")

		.seek-info(
			ref="seekInfo"
			v-show="seekInfo.show"
			:style="{top: `${seekInfo.y}px`,left: `${seekInfo.x}px`}"
		)
			| {{ seekInfo.value | duration }}
		audio(:src="fileURL" autoplay ref="audio")
</template>

<script>
import { coverURL, audiobookURL } from "../api"
import PlacholderCover from "../placeholder-cover.jpg"

import Play from "./svg/Play.vue"
import Pause from "./svg/Pause.vue"
import FastForward from "./svg/FastForward.vue"
import Rewind from "./svg/Rewind.vue"
import Next from "./svg/Next"
import Previous from "./svg/Previous"

export default {
	data() {
		return {
			fileIndex: 0,
			volume: 1,
			time: 0,
			paused: false,
			seekInfo: {
				show: false,
				value: 0,
				x: 0,
			},
			chapterList: false,
		}
	},
	computed: {
		// Get the audiobook to be played from the store.
		audiobook() {
			return this.$store.state.audiobook.playing
		},

		// Current file
		file() {
			if (this.audiobook.files) {
				return this.audiobook.files[this.fileIndex]
			}
		},

		fileURL() {
			if (this.file) {
				return audiobookURL(this.audiobook.id, this.file.name)
			} else {
				return ""
			}
		},

		// Returns array of chapters. The objects differ between chapterized
		// and single file audiobooks
		chapters() {
			if (this.file) {
				if (this.chapterized) {
					return this.audiobook.files
				} else {
					return this.audiobook.files[0].metadata.chapters
				}
			}
		},

		// Checks if we have multiple files or not.
		chapterized() {
			if (this.file) {
				return this.audiobook.files.length > 1
			} else {
				return false
			}
		},

		// Current chapter number
		chapter() {
			if (this.file) {
				if (this.chapterized) {
					return this.fileIndex
				} else {
					for (let i = 0, len = this.chapters.length; i < len; i++) {
						let current = this.chapters[i],
							next = this.chapters[i + 1]
						if (!next) {
							return i
						}

						const startTime = parseFloat(current.start_time)
						const endTime = parseFloat(current.end_time)
						if (this.time >= startTime && this.time < endTime) {
							return i
						}
					}
				}
			}
		},

		// Returns the chapter name of the c
		chapterName() {
			if (this.file) {
				return this.getChapterName(this.chapter)
			}
		},

		// Duration of the current chapter
		chapterDuration() {
			if (this.file) {
				if (this.chapterized) {
					return this.file.duration
				} else {
					const chapter = this.chapters[this.chapter]
					return (
						parseFloat(chapter.end_time) -
						parseFloat(chapter.start_time)
					)
				}
			}
		},

		// Position in current chapter
		chapterTime() {
			if (this.file) {
				if (this.chapterized) {
					return this.time
				} else {
					console
					return this.time - parseFloat(this.chapterObj.start_time)
				}
			}
		},

		// Current chapter object, only useful for a single file audiobooks that
		// contains chapter information.
		chapterObj() {
			return this.chapters[this.chapter]
		},

		// Total runtime of the audiobook
		totalDuration() {
			if (this.file) {
				return this.audiobook.files.reduce(
					(sum, a) => sum + a.duration,
					0,
				)
			} else {
				return 0
			}
		},
	},
	methods: {
		// Media controls
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
		toggle() {
			const el = this.$refs.audio
			if (el) {
				el.paused ? el.play() : el.pause()
			}
		},
		rewind() {
			if (!this.chapter) this.$refs.audio.currentTime = this.time - 30
		},
		forward() {
			this.$refs.audio.currentTime = this.time + 30
		},
		previousChapter() {
			if (this.chapterized) {
				this.fileIndex--
			} else {
				const nextChapter = this.chapters[this.chapter - 1]
				if (nextChapter) {
					this.$refs.audio.currentTime = parseFloat(
						nextChapter.start_time,
					)
				}
			}
		},
		nextChapter() {
			if (this.chapterized) {
				this.fileIndex++
			} else {
				const nextChapter = this.chapters[this.chapter + 1]
				if (nextChapter) {
					this.$refs.audio.currentTime = parseFloat(
						nextChapter.start_time,
					)
				}
			}
		},
		decreaseVolume() {
			this.setVolume(this.volume - 0.1)
		},
		increaseVolume() {
			this.setVolume(this.volume + 0.1)
		},
		setVolume(vol) {
			if (vol < 0) vol = 0
			if (vol > 1) vol = 1
			this.$refs.audio.volume = vol
			this.volume = vol
		},
		handleVolumeClick(e) {
			const bounds = this.$refs.volume.getBoundingClientRect()
			this.setVolume((e.clientX - bounds.x) / bounds.width)
		},

		// Seek info functions
		// Handles click events on the progressbar
		handleSeekClick(e) {
			const bounds = this.$refs.progressBar.getBoundingClientRect()
			this.seek(
				(this.chapterDuration * (e.clientX - bounds.x)) / bounds.width,
			)
		},
		seek(to) {
			this.$refs.audio.currentTime = to
		},

		// Handles mouse movement across the progress bar. Updates position and
		// content of the seek info.
		moveSeekInfo(e) {
			const bounds = this.$refs.progressBar.getBoundingClientRect()
			this.seekInfo.value =
				(this.chapterDuration * e.clientX) / bounds.width
			this.seekInfo.x =
				e.clientX -
				this.$refs.seekInfo.getBoundingClientRect().width / 2
		},
		showSeekInfo() {
			this.seekInfo.show = true
		},
		hideSeekInfo() {
			this.seekInfo.show = false
		},

		// Starts playback of the next file when the current one finishes
		playbackEnded() {
			if (
				this.chapterized &&
				this.audiobook.files.length > this.fileIndex + 1
			) {
				this.fileIndex++
			}
		},

		// Tries to find the best name for a chapter
		getChapterName(i) {
			if (this.chapterized) {
				// Test if we have a chapter object that includes a title
				// for the chapter. Its not guranteed to exists or to have
				// any tags or to have a title tag. What a mess.
				const chapterObj = this.chapters[i].metadata.chapters
				if (chapterObj && chapterObj.tags && chapterObj.tags.title) {
					return chapterObj.tags.title
				}

				// A hack to see if the file's title tag actually is the
				// title of the chapter and not of the entire book. The
				// first and second book not having the same title does
				// obviously not gurantee that they contain the title of
				// the chapter, but its the best I can come up with for now.
				const firstTitle = this.audiobook.files[0].metadata.format.tags
					.title
				const nextTitle = this.audiobook.files[1].metadata.format.tags
					.title
				if (firstTitle && firstTitle !== nextTitle) {
					return this.chapters[i].metadata.format.tags.title
				} else {
					return `Chapter ${this.chapter + 1}`
				}
			} else {
				// The only way with a single file audiobook to know the
				// chapters is to read them from the metadata. Sadly not
				// all formats (namely MP3) support chapter metadata. So,
				// it is not save to assume that we have any chapter
				// information.
				const chapterObj = this.chapters[i]
				if (chapterObj && chapterObj.tags && chapterObj.tags.title) {
					return chapterObj.tags.title
				}
			}
		},

		// Chapterlist
		toggleChapterList() {
			if (!this.chapterList) {
			}
			this.chapterList = !this.chapterList
		},

		// Makes the coverURL and fileURL api functions accessible within the
		// template.
		coverURL,

		// Replaces the cover image with a placeholder if the cover cannot be
		// fetched.
		noImage(e) {
			e.srcElement.src = PlacholderCover
		},
	},

	// Registers keyboard handler
	created() {
		document.addEventListener("keydown", this.keyHandler)
	},

	// Registers handlers for various events fired by the audio element.
	mounted() {
		const el = this.$refs.audio
		el.addEventListener("timeupdate", e => {
			this.time = el.currentTime
		})
		el.addEventListener("play", e => (this.paused = false))
		el.addEventListener("pause", e => (this.paused = true))
		el.addEventListener("ended", this.playbackEnded)
	},

	// Removes the keyboard handler, this has to be done for the keyboard
	// handler but not the other handlers because its registered to the entire
	// document which is not destoryed and thus the handler isn't cleaned up.
	destroyed() {
		document.removeEventListener("keydown", this.keyHandler)
	},
	components: { Play, Pause, FastForward, Rewind, Next, Previous },
}
</script>

<style lang="less" scoped>
@import "../globals.less";

.player {
	width: 100%;
	background: darken(@background, 5%);
	user-select: none;
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

.played {
	height: 100%;
	background: @offwhite;
	transition: 275ms all linear;
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

.seek-info {
	position: absolute;
	bottom: 10em;
	text-shadow: 2px 2px 3px rgba(0, 0, 0, 1);
}

.chapter {
	flex: 1;
}

.volume {
	flex: 1;
	height: 1em;
	background: lighten(@background, 5%);
	border-radius: 6px;
	cursor: pointer;
}

.volume-slider {
	height: 100%;
	background: @offwhite;
	border-radius: 6px;
	width: 0;
	transition: 250ms all ease;
}
</style>
