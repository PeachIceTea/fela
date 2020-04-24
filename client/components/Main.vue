<template lang="pug">
	.app
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

<style lang="less" scoped>
@import "../globals.less";

@scrollbar-background: darken(@background, 5%);

.app {
	min-height: 100%;
	width: 100%;
	height: 100%;
	position: relative;
	display: grid;
	grid-template-rows: 1fr auto;
	grid-template-columns: auto 1fr;
	grid-template-areas:
		"nav-bar main-content"
		"currently-playing currently-playing";
}

nav {
	grid-area: nav-bar;
	width: 350px;
	box-shadow: 2px 0px 10px 0px rgba(0, 0, 0, 0.25);
	z-index: 2;
}

main {
	grid-area: main-content;
	width: 100%;
	height: 100%;
	position: relative;
	overflow-y: scroll;
	overflow-x: hidden;
	min-height: 0;
	scrollbar-color: @highlight @scrollbar-background;
}

main::-webkit-scrollbar {
	background: @scrollbar-background;
}

main::-webkit-scrollbar-track {
	background: @scrollbar-background;
}

main::-webkit-scrollbar-thumb {
	background: @highlight;
}

footer {
	grid-area: currently-playing;
	width: 100%;
	z-index: 4;
	box-shadow: 0px -2px 10px 0px rgba(0, 0, 0, 0.25);

	.child {
		display: flex;
		height: auto;
		-webkit-box-orient: vertical;

		-webkit-box-direction: normal;
	}
}
</style>
