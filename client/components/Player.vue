<template lang="pug">
	.player
		.progress-container
			.progress-bar(@mousemove="hover" @mouseleave="hideInfo" ref="progressBar" @click="progressClick")
				.played(:style="{width: `${progress}%`}")
		.everything-else
			.info.col
				.book-info {{ book.title }} by {{ book.author }}
				.progress {{ timestamp | formatDuration }} / {{ duration | formatDuration }}
			.controls.col
				.control(@click="rewind")
					Rewind.control-symbol
				.control(@click="toogle")
					Pause.control-symbol(v-show="!paused")
					Play.control-symbol(v-show="paused")
				.control(@click="forward")
					FastForward.control-symbol
			.col
		.hover-info(v-show="hoverInfo" :style="hoverStyle" ref="hoverInfo") {{ duration * (hoverPercent / 100) | formatDuration }}
		audio(:src="fileURL" ref="audio" autoplay)
</template>

<script>
import { mapState } from "vuex"

import Play from "./svg/Play.vue"
import Pause from "./svg/Pause.vue"
import FastForward from "./svg/FastForward.vue"
import Rewind from "./svg/Rewind.vue"

export default {
	data() {
		return {
			currentTime: 0,
			currentFileIndex: 0,
			fileDuration: 0,
			paused: true,
			hoverInfo: false,
			hoverPercent: 0,
			hoverStyle: {
				left: "0px",
			},
		}
	},
	computed: Object.assign(
		{
			progress() {
				return (this.timestamp / this.duration) * 100
			},
			duration() {
				if (this.files.length > 1) {
					return this.audiobook.duration
				} else {
					return this.fileDuration
				}
			},
			timestamp() {
				const len = this.files.length
				if (len > 1) {
					let time = this.currentTime
					for (let i = 0; i < this.currentFileIndex; i++) {
						time += this.files[i].duration
					}
					return time
				} else {
					return this.currentTime
				}
			},
			file() {
				if (this.files[0]) {
					return this.files[this.currentFileIndex]
				} else {
					return {}
				}
			},
			fileURL() {
				return this.file.hash
					? `http://localhost:8080/files/${this.file.hash}`
					: ""
			},
		},
		mapState("player", ["book", "audiobook", "files"]),
	),
	created() {
		document.addEventListener("keydown", this.keyHandler)
	},
	mounted() {
		const audio = this.$refs.audio
		audio.addEventListener("timeupdate", e => {
			this.currentTime = audio.currentTime
		})
		audio.addEventListener("durationchange", e => {
			this.fileDuration = audio.duration
		})
		audio.addEventListener("ended", e => {
			if (this.files.length > 1) {
				this.currentFileIndex++
			}
		})
		audio.addEventListener("play", e => (this.paused = false))
		audio.addEventListener("pause", e => (this.paused = true))
	},
	destroyed() {
		document.removeEventListener("keydown", this.keyHandler)
	},
	methods: {
		toogle() {
			const audio = this.$refs.audio
			if (audio) audio.paused ? audio.play() : audio.pause()
		},
		keyHandler(e) {
			switch (e.key) {
				case " ":
					return this.toogle()
				case "ArrowRight":
					return this.forward()
				case "ArrowLeft":
					return this.rewind()
			}
		},
		hover(e) {
			this.hoverInfo = true
			this.hoverPercent = (e.clientX / screen.width) * 100

			const width = this.$refs.hoverInfo.offsetWidth
			const max = screen.width - width
			const left = Math.min(Math.max(e.clientX - width / 2, 0), max)
			this.hoverStyle.left = `${left}px`
		},
		hideInfo() {
			this.hoverInfo = false
		},
		progressClick(e) {
			this.seek(this.duration * (e.clientX / screen.width))
		},
		seek(to) {
			let fileIndex = 0
			for (let i = 0, len = this.files.length, t = 0; i < len; i++) {
				const file = this.files[i]
				if (to < t + file.duration) {
					to -= t
					fileIndex = i
					break
				}

				t += file.duration
			}

			if (this.currentFileIndex !== fileIndex) {
				this.currentFileIndex = fileIndex
				const tmp = () => {
					this.$refs.audio.currentTime = to
					this.$refs.audio.removeEventListener("loadeddata", tmp)
				}
				this.$refs.audio.addEventListener("loadeddata", tmp)
			} else {
				this.$refs.audio.currentTime = to
			}
		},
		rewind() {
			//TODO: Allow user to customize jump
			this.seek(this.timestamp - 30)
		},
		forward() {
			this.seek(this.timestamp + 30)
		},
	},
	components: { Play, Pause, FastForward, Rewind },
}
</script>

<style lang="stylus" scoped>
playerHeight = 5em
progressBarHeight = 1.5em

.player
	display: flex
	background: #282828
	color: #fff
	width: 100%
	height: playerHeight
	flex-direction: column

.everything-else
	display: flex
	padding: 1em
	flex: 1

.col
	flex: 1

.controls
	display: flex
	justify-content: center
	align-items: center

.control
	margin-left: 2em
	display: inline
	transition: 100ms all ease
	cursor: pointer
	fill: #ddd

	&:hover
		fill: #fff

.control-symbol
	height: 100%
	width: 2em

.progress-container
	width: 100%
	height: progressBarHeight

.progress-bar
	height: 0.5em
	background: lighten(#282828, 10%)
	cursor: pointer
	transition: 100ms all ease
	margin: -0.5em

	&:hover
		height: progressBarHeight
		margin-top: - progressBarHeight

.played
	height: 100%
	width: 0%
	background: #fff
	transition: 250ms all ease

.hover-info
	position: absolute
	bottom: (progressBarHeight + playerHeight) + 0.25em
</style>
