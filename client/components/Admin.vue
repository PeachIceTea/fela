<template lang="pug">
	.admin
		.user(v-for="user in users")
			.name Name: {{user.name | title }}
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
	padding: 1em;
}

.user {
	.container();
}
</style>
