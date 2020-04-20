<template lang="pug">
	body(@dragover.stop.prevent="dragOver" @drop.stop.prevent="drop")
		router-view
</template>

<script>
export default {
	methods: {
		dragOver(e) {
			e.dataTransfer.dropEffect = "none"
		},
		drop(e) {},
	},
}
</script>

<style lang="stylus">
@import "../globals.styl"

*
	box-sizing: border-box

div, span, p
	cursor: default
	user-select: none

html, body
	margin: 0
	height: 100%
	overflow: hidden
	background: background
	color: white-text
</style>
