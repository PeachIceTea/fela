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
		// Runs when the form is submitted. Passes the values of the two inputs
		// to the API and redirects to the main page if the API does not return
		// an error. It also starts a small animation to indicate to the user
		// that the login is happening.
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

		// Animates three dots on the login button to indicate to the user that
		// their request is being processed.
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

<style lang="less" scoped>
@import "../globals.less";

.login {
	display: flex;
	align-items: center;
	flex-direction: column;
}

form {
	width: 250px;
	margin-top: 100px;

	.boxShadow();
}

input {
	font-size: 20px;
	width: 100%;
	border: 0;
	outline: 0;
	border-radius: 0px;
	padding: 11px 10px 9px;
	color: @black-text;

	&:first-child {
		border-radius: @border-radius @border-radius 0 0;
	}
}

input[type="submit"] {
	border-radius: 0 0 @border-radius @border-radius;
	background: @highlight;
	color: @white-text;
	box-shadow: none;
}

.err {
	.container();

	border-top: 5px solid @err;
	margin-top: 1em;
	width: 100%;
}
</style>
