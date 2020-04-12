<template lang="pug">
	.app
		.everything-else
			nav.col.nav
				Navbar
			main.col.main
				router-view
		footer.play
			Player
</template>

<script>
import Player from "./Player.vue"
import Navbar from "./Navbar.vue"

export default {
	components: { Player, Navbar },
}
</script>

<style lang="stylus" scoped>
.app
	height: 100vh
	width: 100vw
	display: flex
	flex-direction: column
	overflow: hidden

.everything-else
	height: 100%
	display: flex
	flex-direction: row

.col
	height: 100%
	flex: 1
	border: 1px solid #000

.main
	flex: 10

.play
	margin-top: auto
	border: 1px solid #000
</style>
