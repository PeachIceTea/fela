<template lang="pug">
	.player
		button.play-btn(@click="play" :class="{playing: playing}") Play
		audio(:src="`http://localhost:8080/files/${file.hash}`" ref="audio")
</template>

<script>
export default {
	data() {
		return {}
	}
	computed: {
		file() {
			return this.$store.state.player.file
		},
		playing() {
			return this.$refs.audio && !this.$refs.audio.paused
		},
	},
	methods: {
		play() {
			this.$refs.audio.play()
		},
	},
}
</script>

<style lang="stylus" scoped>
.player
	box-sizing: border-box
	position: absolute
	display: inline-block
	bottom: 0
	background: #000
	color: #fff
	width: 100%
	height: 8em
	padding: 1em;

.play-btn {
  border: 0;
  background: transparent;
  box-sizing: border-box;
  width: 0;
  height: 2em;

  border-color: transparent transparent transparent #ddd;
  transition: 100ms all ease;
  cursor: pointer;

  // play state
  border-style: solid;
  border-width: 3em 0 3em 5em;

  & .paused {
    border-style: double;
    border-width: 0px 0 0px 60px;
  }

  &:hover {
    border-color: transparent transparent transparent #fff;
  }
}
</style>
