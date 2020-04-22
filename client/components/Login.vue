<template lang="pug">
	.login
		form(@submit.prevent="submit")
			input.input(type="text" placeholder="Username" v-model="name" required)
			input.input(
				type="password"
				placeholder="Password"
				v-model="password"
				required
			)
			input.button(type="submit" value="Login" ref="btn")
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

.login
	display: flex
	height: 100%
	width: 100%
	flex-direction: column
	align-items: center

form
	display: flex
	flex-direction: column
	align-items: center
	margin-top: 10vh
	box-shadow: box-shadow
	width: form-width

input:first-child
	border-radius: border-radius border-radius 0 0

input:last-child
	border-radius: 0 0 border-radius border-radius
</style>
