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
				return audiobookURL(
					this.audiobook.id,
					this.audiobook.files[0].name,
				)
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
			const bounds = this.$refs.progress.getBoundingClientRect()
			this.progressInfo.value =
				(this.totalDuration * e.clientX) / bounds.width
			this.progressInfo.x =
				e.clientX -
				this.$refs.progressInfo.getBoundingClientRect().width / 2
		},
		noImage(e) {
			e.srcElement.src = PlacholderCover
		},
		coverURL,
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
