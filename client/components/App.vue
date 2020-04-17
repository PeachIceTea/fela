<template lang="pug">
	body
		router-view
</template>

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
