<template lang="pug">
	.navbar
		router-link.nav-el(to="/" :class="{active: isActive('/')}") Library
		.subroute(:class="{active: isActive('/')}")
			input.nav-el(type="text" placeholder="Search ..." @input="search")
			.nav-el(@click="nextView") View: {{ viewName }}
			.nav-el(@click="nextOrder") Sorted: {{ orderName }}
		.nav-el.nav-break(v-if="uploader")
		router-link.nav-el(
			to="/audiobook/upload"
			v-if="uploader"
			:class="{active: isActive('/audiobook/upload')}"
		) Upload
		router-link.nav-el(
			to="/admin"
			v-if="admin"
			:class="{active: isActive('/admin')}") Admin
		.nav-el.nav-break(v-if="uploader")
		router-link.nav-el(
			to="/settings"
			:class="{active: isActive('/settings')}") Settings
		.logout.nav-el(@click="logout") Logout
</template>

<script>
export default {
	computed: {
		admin() {
			return this.$store.getters.isAdmin
		},
		uploader() {
			return this.$store.getters.isUploader
		},
		orderName() {
			switch (this.$store.state.ui.order) {
				case 0:
					return "Alphabetically (A-Z)"
				case 1:
					return "Alphabetically (Z-A)"
				case 2:
					return "Upload Date (Ascending)"
				case 3:
					return "Upload Date (Descending)"
				case 4:
					return "Default"
			}
		},
		viewName() {
			switch (this.$store.state.ui.view) {
				case 0:
					return "Cover"
				case 1:
					return "List"
			}
		},
	},
	methods: {
		logout() {
			this.$store.dispatch("logout")
			this.$router.push("/login")
		},
		isActive(link) {
			return this.$route.path === link
		},
		nextOrder() {
			this.$store.commit("nextOrder")
		},
		nextView() {
			this.$store.commit("nextView")
		},
		search(e) {
			this.$store.commit("setSearch", e.srcElement.value)
		},
	},
}
</script>

<style lang="less" scoped>
@import "../globals.less";

.navbar {
	background: @highlight;

	height: 100%;
	padding-top: 15px;
	user-select: none;
}

.nav-el {
	display: block;
	text-align: center;
	color: @white-text;
	text-decoration: none;
	font-size: 20px;
	padding: 15px 0;
	cursor: pointer;
	font-size: 25px;
	transition: 250ms all ease-out;
	z-index: 3;
	text-shadow: @text-shadow;

	&:visited {
		color: @white-text;
	}

	&:hover {
		background: darken(@highlight, 15%);
	}
}

.active {
	background: darken(@highlight, 15%);
}

.nav-break {
	background: transparent !important;
	border: 0;
	cursor: default;
}

.subroute {
	background: darken(@highlight, 10%);
	display: none;

	.nav-el {
		font-size: 18px;
	}

	&.active {
		position: relative;
		display: block;
	}

	input[type="text"] {
		padding: @input-padding !important;
		width: 100%;
		font-size: 20px;
		border: 0;
		border-radius: 0;
		background: darken(@highlight, 10%);
		cursor: text;

		&::placeholder {
			color: @white-text;
			opacity: 1;
		}
	}
}
</style>
