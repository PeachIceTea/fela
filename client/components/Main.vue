<template lang="pug">
	.app
		.not-player
			nav
				Navbar
			main
				router-view
		footer
			Player
</template>

<script>
import Navbar from "./Navbar.vue"
import Player from "./Player.vue"

export default {
	components: { Navbar, Player },
}
</script>

<style lang="stylus" scoped>
@import "../globals.styl"

.app
	height: 100vh
	width: 100vw
	display: flex
	flex-direction: column

.not-player
	height: 100%
	display: flex
	flex-direction: row
	flex-grow: 1
	min-height: 0

nav
	flex: 1

main
	overflow-y: scroll
	scrollbar-color: highlight darken(background, 50%)
	flex: 10
	background: background

footer
	height: 100%
	flex: 100
	border: 1px solid #000
</style>
