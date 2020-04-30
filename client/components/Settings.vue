<template lang="pug">
	.settings
		h1 Settings
		form(@submit.prevent="updatePassword")
			h2 Change password
			.input-row
				label Current password
					.input-container
						input(
							type="password"
							v-model="password.current"
							placeholder="Current password"
						)
			.input-row
				label New password
					.input-container
						input(
							type="password"
							v-model="password.new"
							placeholder="New password"
						)
			.input-row
				label Confirm new password
					.input-container
						input(
							type="password"
							v-model="password.confirm"
							placeholder="New password confirmation"
						)
			.input-row
				input(type="submit" value="Update Password")
			.message.err(v-show="this.password.err") {{ this.password.err }}
			.message.success(
				v-show="this.password.success"
			) {{ this.password.success }}
</template>

<script>
import { updateUser } from "../api"

const initialPasswordState = {
	current: "",
	new: "",
	confirm: "",
	err: "",
	success: "",
}

export default {
	data() {
		return {
			password: { ...initialPasswordState },
		}
	},
	methods: {
		async updatePassword() {
			this.password.err = ""
			if (!this.password.current || !this.password.new) return

			if (this.password.new !== this.password.confirm) {
				this.password.err = "passwords have to match"
				return
			}

			const res = await updateUser(this.$store.state.auth.loggedIn.id, {
				confirmation: this.password.current,
				password: this.password.new,
			})
			if (res.err) {
				this.password.err = res.err
				return
			}
			this.password = {
				...initialPasswordState,
				success: "password updated",
			}
		},
	},
}
</script>

<style lang="less" scoped>
@import "../globals.less";

.settings {
	margin: 1em;

	.container();
}

.input-row {
	margin-bottom: 1em;
}

input,
select,
button {
	border: 0;
	padding: @input-padding;
	width: 100%;
	max-width: 700px;
	border-radius: @border-radius;
	font-size: 20px;
}

.message {
	.container();

	background: @offwhite;
	color: @black-text;
	padding-top: calc(1em - 4px);
	border-top: 8px solid transparent;
	max-width: 700px;
}

.err {
	border-top-color: @err;
}

.success {
	border-top-color: @success;
}
</style>
