<template lang="pug">
	.admin
		.row.header
			.column Username
			.column Password
			.column Role
			.column Created
			.column Last Updated
			.column Save
			.column Delete
		.row.user(v-for="user in users")
			.column
				input(type="text" :value="user.name" placeholder="Username")
			.column
				input(type="password" placeholder="Password")
			.column
				select
					option(value="listener" :selected="user.role === 'listener'") Listener
					option(value="uploader" :selected="user.role === 'uploader'") Uploader
					option(value="admin" :selected="user.role === 'admin'") Admin
			.column {{ user.created_at }}
			.column {{ user.updated_at }}
			.column
				button Save
			.column
				button Delete
		.row
			button New User
</template>

<script>
export default {
	computed: {
		users() {
			return this.$store.state.user.list
		},
	},
	beforeRouteEnter(to, from, next) {
		next(async vm => {
			vm.$store.dispatch("getAllUsers")
		})
	},
}
</script>

<style lang="less" scoped>
@import "../globals.less";

.admin {
	display: grid;
	padding: 1em;
}

.row {
	display: grid;
	grid-template-columns: repeat(7, 1fr);
}

.column {
	padding: 0.25em;
}

.header {
	padding: 1em;
}

.user {
	.container();
	font-size: 20px;
	line-height: 5;
	padding: 0 1em;
}

input,
select {
	background: lighten(@background, 10%);
	color: @white-text;
	border: 0;
	border-bottom: 1px @offwhite solid;
}
</style>
