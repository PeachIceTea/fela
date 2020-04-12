<template lang="pug">
	.player
		.progress-container
			.progress-bar(@mousemove="hover" @mouseleave="hideInfo" ref="progressBar" @click="jump")
				.played(:style="{width: `${progress}%`}")
		.everything-else
			.info.col
				.book-info {{ book.title }} by {{ book.author }}
				.progress {{ formatTime(timestamp) }} / {{ formatTime(duration) }}
			.controls.col
				.play-btn(@click="toogle")
					Pause(v-show="!paused")
					Play(v-show="paused")
			.col
			audio(:src="`http://localhost:8080/files/${file.hash}`" ref="audio" autoplay)
		.hover-info(v-show="hoverInfo" :style="hoverStyle" ref="hoverInfo") {{ formatTime(duration * (hoverPercent / 100)) }}
</template>

<script>
import Play from "./svg/Play.vue"
import Pause from "./svg/Pause.vue"

export default {
	data() {
		return {
			timestamp: 0,
			duration: 0,
			paused: true,
			hoverInfo: false,
			hoverPercent: 0,
			hoverStyle: {
				left: "0px",
			},
		}
	},
	computed: {
		book() {
			return this.$store.state.player.book
		},
		file() {
			return this.$store.state.player.file
		},
		progress() {
			return (this.timestamp / this.duration) * 100
		},
	},
	created() {
		document.addEventListener("keydown", this.spaceHandler)
	},
	mounted() {
		const audio = this.$refs.audio
		audio.addEventListener("timeupdate", e => {
			this.timestamp = audio.currentTime
		})
		audio.addEventListener("durationchange", e => {
			this.duration = audio.duration
		})
		audio.addEventListener("play", e => (this.paused = false))
		audio.addEventListener("pause", e => (this.paused = true))
	},
	destroyed() {
		document.removeEventListener("keydown", this.spaceHandler)
	},
	methods: {
		toogle() {
			const audio = this.$refs.audio
			if (audio) audio.paused ? audio.play() : audio.pause()
		},
		spaceHandler(e) {
			if (e.key === " ") this.toogle()
		},
		formatTime(time) {
			const hours = Math.floor(time / (60 * 60))
			time -= hours * 60 * 60

			const minutes = this.fillZero(Math.floor(time / 60))
			time -= minutes * 60

			const seconds = this.fillZero(Math.floor(time))
			return `${hours}:${minutes}:${seconds}`
		},
		fillZero(time) {
			time = time.toString()
			return time.length < 2 ? `0${time}` : time
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
		jump(e) {
			console.log(this.duration * (e.clientX / screen.width))
			this.$refs.audio.currentTime =
				this.duration * (e.clientX / screen.width)
		},
	},
	components: { Play, Pause },
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

.play-btn
	display: inline
	transition: 100ms all ease
	cursor: pointer
	fill: #ddd

	&:hover
		fill: #fff

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
	transition: 100ms all ease

.hover-info
	position: absolute
	bottom: (progressBarHeight + playerHeight) + 0.25em
</style>
