<template lang="pug">
	.player(v-show="audiobook.files")
		.progress-bar(
			ref="progress"
			@click="handleProgressClick"
			@mouseenter="showProgressInfo"
			@mousemove="moveProgressInfo"
			@mouseleave="hideProgressInfo"
		)
			.played(:style="{width: `${currentTime/totalDuration*100}%`}")
		.everything-else
			.info.col
				.book-info.col(style="display: inline")
					p {{ audiobook.title }}
					p {{ audiobook.author }}
				.playback-info.col
					| {{ currentTime | formatDuration }} /
					|  {{ totalDuration | formatDuration }}
			.playback.col
				.control(@click="rewind")
					Rewind.control-symbol
				.control(@click="toggle")
					Pause.control-symbol(v-show="!paused")
					Play.control-symbol(v-show="paused")
				.control(@click="forward")
					FastForward.control-symbol
			.knobs.col
				.speed.col
				.volume.col(@mousedown="handleVolumeClick" ref="volume")
					.volume-bar
						.filled(:style="{width: `${volume*100}%`}")
		.progress-info(
			ref="progressInfo"
			v-show="progressInfo.show"
			:style="{top: `${progressInfo.y}px`,left: `${progressInfo.x}px`}"
		)
			| {{ progressInfo.value | formatDuration }}
		audio(:src="fileUrl" ref="audio")
</template>

<script>
import Play from "./svg/Play.vue"
import Pause from "./svg/Pause.vue"
import FastForward from "./svg/FastForward.vue"
import Rewind from "./svg/Rewind.vue"

export default {
	data() {
		return {
			paused: true,
			currentTime: 0,
			volume: 1,
			progressInfo: {
				show: false,
				value: 0,
				x: 0,
			},
		}
	},
	computed: {
		audiobook() {
			return this.$store.state.audiobook.playing
		},
		totalDuration() {
			if (this.audiobook.files) {
				return this.audiobook.files.reduce(
					(sum, a) => sum + a.duration,
					0,
				)
			} else {
				return 0
			}
		},
		fileUrl() {
			const audiobook = this.audiobook
			if (audiobook.files) {
				return `http://localhost:8080/files/audio/${audiobook.id}/${audiobook.files[0].name}`
			}
		},
	},
	methods: {
		rewind() {
			this.seek(this.currentTime - 30)
		},
		toggle() {
			const el = this.$refs.audio
			if (el) {
				el.paused ? el.play() : el.pause()
			}
		},
		forward() {
			this.seek(this.currentTime + 30)
		},
		increaseVolume() {
			this.setVolume(this.volume + 0.1)
		},
		decreaseVolume() {
			this.setVolume(this.volume - 0.1)
		},
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
		seek(time) {
			this.$refs.audio.currentTime = time
		},
		handleProgressClick(e) {
			const bounds = this.$refs.progress.getBoundingClientRect()
			this.seek(
				(this.totalDuration * (e.clientX - bounds.x)) / bounds.width,
			)
		},
		showProgressInfo() {
			this.progressInfo.show = true
		},
		hideProgressInfo() {
			this.progressInfo.show = false
		},
		moveProgressInfo(e) {
			this.progressInfo.value =
				(this.totalDuration * e.clientX) / screen.width
			this.progressInfo.x =
				e.clientX -
				this.$refs.progressInfo.getBoundingClientRect().width / 2
		},
	},
	created() {
		document.addEventListener("keydown", this.keyHandler)
	},
	mounted() {
		const el = this.$refs.audio
		el.addEventListener("timeupdate", e => {
			this.currentTime = el.currentTime
		})
		el.addEventListener("play", e => (this.paused = false))
		el.addEventListener("pause", e => (this.paused = true))
	},
	destroyed() {
		document.removeEventListener("keydown", this.keyHandler)
	},
	components: { Play, Pause, FastForward, Rewind },
}
</script>

<style lang="stylus" scoped>
@import "../globals"
player-height = 5.5em
progress-bar-height = 1.5em

.player
	display: flex
	width: 100%
	height: player-height
	flex-direction: column
	background: player-background
	text-shadow: text-shadow

.everything-else
	display: flex
	padding: 1em
	flex: 1
	height: 100%

.playback
	justify-content: center

.col
	display: flex
	align-items: center
	flex: 1

.info
	p
		margin: 5px 0 0 5px

.control
	margin-left: 2em
	display: inline
	transition: 100ms all ease
	cursor: pointer
	fill: black-text

	&:hover
		fill: white-text


.control-symbol
	height: 100%
	width: 2em
	cursor: pointer

.knobs
	justify-content: right

.volume-bar
	width: 100%
	height: progress-bar-height - 0.5em
	border-radius: 5px
	background: darken(black-text, 40%)
	cursor: pointer
	box-shadow: box-shadow

	.filled
		border-radius: 5px
		height: 100%
		width: 50%
		background: offwhite
		transition: 100ms all ease
		cursor: inherit
		box-shadow: inherit

.progress-bar
	z-index: 2
	height: 0.5em
	background: darken(black-text, 30%)
	cursor: pointer
	transition: 100ms all ease
	margin: -0.5em
	box-shadow: 0px -5px 15px 0px rgba(0,0,0,0.25)

	&:hover
		height: progress-bar-height
		margin-top: - progress-bar-height

	.played
		height: 100%
		width: 50%
		background: offwhite
		transition: 250ms all ease
		cursor: inherit

.progress-info
	position: absolute
	background: background
	padding: 0.5em
	bottom: player-height + progress-bar-height + 0.25em
</style>
