<template lang="pug">
	.app
		.menu-btn(@click="toggleNav")
			Menu(v-show="!showNav")
			Cancel(v-show="showNav")
		nav(:class="{show: showNav}")
			Navbar
		main
			router-view
		footer
			Player
</template>

<script>
import Menu from "./svg/Menu.vue"
import Cancel from "./svg/Cancel"

import Navbar from "./Navbar.vue"
import Player from "./Player.vue"

export default {
	data() {
		return {
			showNav: false,
		}
	},
	methods: {
		toggleNav() {
			this.showNav = !this.showNav
		},
	},
	watch: {
		$route() {
			this.showNav = false
		},
	},
	components: { Navbar, Player, Menu, Cancel },
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

.menu-btn {
	display: none;
}

@media (max-width: 480px) {
	nav {
		display: none;

		&.show {
			display: block !important;
			position: absolute;
			width: 100vw;
			height: 100vh;
			z-index: 10;
			box-shadow: 100px 0px 200px 10px #000;
		}
	}

	.menu-btn {
		display: flex;
		position: absolute;
		z-index: 100;
		top: 10px;
		right: 10px;
		fill: @offwhite;
		font-size: 2em;
		width: 1.5em;
		background: @highlight;
		justify-items: center;
		box-shadow: @box-shadow;
		border-radius: @border-radius;
		cursor: pointer;

		svg {
			width: 100%;
			height: 1.5em;
		}
	}

	main {
		-webkit-overflow-scrolling: touch;
	}
}
</style>
