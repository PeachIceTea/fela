<template lang="pug">
	.library
		.book(
			v-for="audiobook in audiobooks"
			@click="play(audiobook)"
		)
			img(
				:src="`http://localhost:8080/files/cover/${audiobook.id}.jpg`"
				@error="noImage(audiobook)"
			)
</template>

<script>
export default {
	computed: {
		audiobooks() {
			return this.$store.state.audiobook.list
		},
	},
	methods: {
		play(audiobook) {
			this.$store.dispatch("getAudiobook", audiobook.id)
		},
		noImage(b) {
			console.log(b)
			b.noCover = true
		},
	},
	created() {
		this.$store.dispatch("getAudiobooks")
	},
}
</script>

<style lang="stylus" scoped>
@import "../globals"

.library
	display: grid
	grid-template-columns: repeat(4, 1fr)

.book
	cursor: pointer
	position: relative
	width: 100%

	img
		z-index: 2
		width: 100%
		height: 100%
		transition: 500ms all ease

		&:hover
			filter: blur(5px) brightness(0.75)

.text-container
	z-index: 1
	position: absolute
	top: 0
	left: 0
	height: 100%
	width: 100%
	pointer-events: auto

.img-text
	padding: 10px
</style>
