<template lang="pug">
	.login
		.msg(v-show="err") {{ err }}
		form(@submit.prevent="submit")
			input(type="text" placeholder="Username" v-model="name" required)
			input(
				type="password"
				placeholder="Password"
				v-model="password"
				required
			)
			input(type="submit" value="Login" ref="btn")
</template>

<script>
import Vue from "vue"
export default {
	data() {
		return {
			name: "",
			password: "",
			err: undefined,
			animate: false,
			frame: 0,
		}
	},
	methods: {
		async submit() {
			this.animate = true
			this.animation()

			const res = await this.$store.dispatch("login", {
				name: this.name,
				password: this.password,
			})
			if (!res.err) {
				this.$router.push("/")
			} else {
				this.err = res.err
			}

			this.animate = false
		},
		animation() {
			if (this.animate) {
				let text = "Logging in ."
				if (this.frame === 1) {
					text += "."
				} else if (this.frame === 2) {
					text += ".."
				}

				this.$refs.btn.value = text
				this.frame++
				this.frame %= 3
				setTimeout(this.animation, 500)
			} else {
				this.$refs.btn.value = "Login"
			}
		},
	},
}
</script>

<style lang="stylus" scoped>
@import "../globals.styl"
borderRadius = 3px
form-width = 249px
font-size = 20px
padding = 11px 10px 9px

.login
	display: flex
	height: 100%
	width: 100%
	flex-direction: column
	align-items: center
	font-size: font-size

form
	display: flex
	flex-direction: column
	align-items: center
	margin-top: 10vh
	box-shadow: 1px 1px 5px 0px rgba(0,0,0,0.25)
	width: form-width

input
	display: block
	width: 100%
	border: 0
	color: black-text
	padding: padding
	outline: 0
	background: offwhite
	font-size: font-size

input[type="submit"]
	background: highlight
	color: white-text
	cursor: pointer
	transition: 100ms ease all

	&::-moz-focus-inner
		border: 0

input:first-child
	border-radius: borderRadius borderRadius 0 0

input:last-child
	border-radius: 0 0 borderRadius borderRadius
</style>
