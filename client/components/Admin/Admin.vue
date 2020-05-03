<template lang="pug">
	.admin
		.row.header
			.column Username
			.column Role
			.column Created
			.column Last Updated
			.column Edit
		.row.user(v-for="user in users")
			.column {{ user.name | title}}
			.column {{ user.role | title }}
			.column {{ user.created_at }}
			.column {{ user.updated_at }}
			.column
				button(@click="edit(user)") Edit
		.row
			button(@click="newUser") New User
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
			await vm.$store.dispatch("getAllUsers")
		})
	},
	methods: {
		edit(user) {
			this.$router.push(`/user/${user.id}/edit`)
		},
		newUser(user) {
			this.$router.push(`/user/new/edit`)
		},
	},
}
</script>

<style lang="less" scoped>
@import "../../globals.less";

.admin {
	display: grid;
	padding: 1em;
}

.row {
	display: grid;
	grid-template-columns: repeat(5, 1fr);
	margin-bottom: 1em;
}

.column {
	padding: 0.25em;
}

.user {
	.container();
	font-size: 18px;
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
