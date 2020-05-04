<template lang="pug">
	.edit-user
		form(@submit.prevent="submit")
			h1(v-if="!isNewUser") Edit User
			h1(v-else) New User
			.input-row
				label Name
					.input-container
						input(
							type="text"
							v-model="name"
							placeholder="Username"
						)
			.input-row
				label Password
					.input-container
						input(
							type="password"
							v-model="password"
							placeholder="Password"
						)
			.input-row
				label Role
					.input-container
						select(v-model="role")
							option(
								value="user"
								:selected="!user.role || user.role === 'user'"
							) User
							option(
								value="uploader"
								:selected="user.role === 'uploader'"
							) Uploader
							option(
								value="admin"
								:selected="user.role === 'admin'"
							) Admin
			.input-row
				input(type="submit" :value="isNewUser ? 'Create' : 'Update'")

			.input-row(v-if="!isNewUser")
				button(@click.prevent="deleteUser")
					span(v-if="!deleting") Delete User
					span(v-else) Are you sure?
			.message(v-show="err") Error: {{ err }}
		</template>

<script>
import { register, updateUser, deleteUser } from "../../api"

export default {
	data() {
		return {
			name: "",
			password: "",
			role: "",
			deleting: false,
			err: "",
		}
	},
	computed: {
		user() {
			return this.getUser(this.$route.params.id)
		},
		isNewUser() {
			return this.$route.params.id === "new"
		},
	},
	methods: {
		getUser(id) {
			const list = this.$store.state.user.list
			for (let i = 0, len = list.length; i < len; i++) {
				const user = list[i]
				if (user.id == id) {
					return user
				}
			}
			return {}
		},
		async submit() {
			if (this.isNewUser) {
				const res = await register(this.name, this.password, this.role)
				if (res.err) {
					this.err = res.err
					return
				}
			} else {
				const update = {}
				if (this.user.name !== this.name) update.name = this.name
				if (this.password) update.password = this.password
				if (this.user.role !== this.role) update.role = this.role

				const res = await updateUser(this.user.id, update)
				if (res.err) {
					this.err = res.err
					return
				}
			}
			this.$store.dispatch("getAllUsers")
			this.$router.push("/admin")
		},
		async deleteUser() {
			if (!this.deleting) {
				this.deleting = setTimeout(() => (this.deleting = 0), 2000)
				return
			}

			clearTimeout(this.deleting)
			const res = await deleteUser(this.user.id)
			console.log(res)
			if (res.err) {
				this.err = res.err
			}
			this.$store.dispatch("getAllUsers")
			this.$router.push("/admin")
		},
	},
	beforeRouteEnter: routeHandler,
	beforeRouteUpdate: routeHandler,
}

function routeHandler(to, from, next) {
	next(async vm => {
		await vm.$store.dispatch("getAllUsers")
		const user = vm.getUser(to.params.id)
		if (user) {
			vm.name = user.name
			vm.role = user.role
		}
	})
}
</script>

<style lang="less" scoped>
@import "../../globals.less";

.edit-user {
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
	border-top-color: @err;
}
</style>
