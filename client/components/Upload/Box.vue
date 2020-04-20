<template lang="pug">
	label
		.upload-box(
			ref="box"
			@drop.stop.prevent="dropFile"
			@dragover.stop.prevent="dragOver"
			@dragleave.stop.prevent="dragExit"
			:class="{dragover: fileHovering}"
		)
			.text
				p #[b Click to choose files]
				p or drag them here
			input(type="file" multiple @input="dialogFile")
</template>

<script>
export default {
	data() {
		return {
			fileHovering: false,
		}
	},
	methods: {
		dropFile(e) {
			this.fileHovering = false
			console.log(e)
		},
		dragOver(e) {
			e.dataTransfer.dropEffect = "copy"
			this.fileHovering = true
		},
		dragExit(e) {
			this.fileHovering = false
		},
		dialogFile(e) {},
	},
}
</script>

<style lang="stylus" scoped>
@import "../../globals"

.upload-box
	width: 100%
	outline: 2px dashed offwhite
	border-radius: 3px
	background: lighten(background, 2.5%)
	cursor: pointer
	display: inline-block
	padding: 5em
	text-align: center

	&.dragover
		background: lighten(background, 10%)

.text
	cursor: inherit

p
	cursor: inherit

input[type="file"]
	width: 1px
	height: 1px
	position: absolute
	left: -1px
	outline: none
</style>
