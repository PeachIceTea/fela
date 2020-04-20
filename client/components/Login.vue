<template lang="pug">
	.login
		form(@submit.prevent="submit")
			input(type="text" placeholder="Username" v-model="name" required)
			input(
				type="password"
				placeholder="Password"
				v-model="password"
				required
			)
			input(type="submit" value="Login" ref="btn")
		.err(v-show="err") Error: {{ err }}
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
				try {
					await this.$router.push("/")
				} catch (e) {
					await this.$router.push("/")
				}
			} else {
				this.err = res.err
			}

			this.animate = false
		},
		animation() {
			const btn = this.$refs.btn
			if (this.animate) {
				let text = "Logging in ."
				if (this.frame === 1) {
					text += "."
				} else if (this.frame === 2) {
					text += ".."
				}

				btn.value = text
				this.frame++
				this.frame %= 3
				setTimeout(this.animation, 500)
			} else {
				if (btn) {
					btn.value = "Login"
				}
			}
		},
	},
}
</script>

<style lang="stylus" scoped>
@import "../globals.styl"

border-radius = 3px
form-width = 249px
font-size = 20px

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
	box-shadow: box-shadow
	width: form-width

input
	display: block
	width: 100%
	border: 0
	color: black-text
	padding: input-padding
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

	&:active
		background: darken(highlight, 10%)

input:first-child
	border-radius: border-radius border-radius 0 0

input:last-child
	border-radius: 0 0 border-radius border-radius

.err
	margin-top: 1em
	width: form-width
	font-size: font-size
	background: offwhite
	color: black-text
	border-radius: border-radius
	border-top: 5px red solid
	padding: 0.1em
</style>
