<template lang="pug">
	.navbar
		router-link.nav-el(to="/") Library
		router-link.nav-el(to="/upload") Upload
		.nav-el.nav-break
		router-link.nav-el(to="/settings") Settings
		span.nav-el(@click="logout") Logout
</template>

<script>
export default {
	methods: {
		logout() {
			this.$store.dispatch("logout")
			this.$router.push("/login")
		},
	},
}
</script>

<style lang="stylus" scoped>
@import "../globals.styl"

.navbar
	height: 100%
	display: flex
	flex-direction: column
	background: highlight
	padding: 1.25em 0
	text-shadow: text-shadow

.nav-el
	display: flex
	justify-content: center
	align-items: center
	height: 75px
	cursor: pointer
	text-decoration: none
	outline: 0
	color: white-text
	font-size: 20px
	transition: 250ms all ease
	margin-top: 10px

	&:visited
		color: white-text

	&:hover
		background: darken(highlight, 15%)

.nav-break
	background: transparent !important
	border: 0
	cursor: default
</style>
