<template lang="pug">
	.navbar
		router-link(to="/").nav-el Library
		router-link(to="/upload").nav-el Upload
		.nav-el.nav-break
		router-link(to="/404").nav-el Logout
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
</style>
