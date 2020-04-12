<template lang="pug">
	.navbar
		router-link.nav-el(to="/") Library
		router-link.nav-el(to="/upload") Upload
		.nav-el.nav-break
		router-link.nav-el(to="/settings") Settings
		router-link.nav-el(to="/404") Logout
</template>

<style lang="stylus" scoped>
.navbar
	height: 100%
	display: flex
	flex-direction: column
	background: #040404
	padding: 1.25em

.nav-el
	display: flex
	justify-content: center
	align-items: center
	height: 3em
	cursor: pointer
	border-radius: 0.25em
	text-decoration: none
	outline: 0
	color: #fff

	&:visited
		color: #fff

	&:hover
		background: #333

.nav-break
	background: transparent !important
	cursor: default
</style>
