import Vue from "vue"
import VueRouter from "vue-router"
Vue.use(VueRouter)

import store from "./store"

import App from "./components/App.vue"
import Login from "./components/Login.vue"

import Main from "./components/Main.vue"
import Library from "./components/Library"
import Upload from "./components/Upload"
import Edit from "./components/Upload/Edit.vue"

import Admin from "./components/Admin"
import EditUser from "./components/Admin/Edit.vue"

import Settings from "./components/Settings.vue"

import NotFound from "./components/NotFound"

// Formats a number to a string representing duration, e.g. 88762.131156 to
// 24:39:22". Registered as the "duration" filter to use in Vuejs templates.
function formatDuration(duration) {
	const fillZero = n => (n < 10 ? `0${n}` : n)
	const h = Math.floor(duration / (60 * 60))
	duration -= h * 60 * 60

	const m = Math.floor(duration / 60)
	duration -= m * 60

	const s = Math.floor(duration)
	return `${h}:${fillZero(m)}:${fillZero(s)}`
}
Vue.filter("duration", formatDuration)

// Returns string with the first character capitalized and all other in
// lowercase. Registered as the "title" filter to use in Vuejs templates.
function fromatTitle(str) {
	return `${str[0].toUpperCase()}${str.slice(1).toLowerCase()}`
}
Vue.filter("title", fromatTitle)

// Guards prevent navigation to a protected route by non authorized users. These
// proctetions are implemented on server side as well and for the client simply
// serve to improve the expierence of the user.

// Route requires user to be not logged in.
async function nonAuthGuard(to, from, next) {
	await store.dispatch("restoreToken")
	if (store.state.auth.token) {
		next("/")
	} else {
		next()
	}
}

// Route requires user to be logged in.
async function authGuard(to, from, next) {
	await store.dispatch("restoreToken")
	if (store.state.auth.token) {
		next()
	} else {
		next("/login")
	}
}

// Route requries "uploader" or "admin" role.
function uploaderGuard(to, from, next) {
	if (store.getters.isUploader) {
		next()
	} else {
		next("/")
	}
}

// Route requires "admin" role.
function adminGuard(to, from, next) {
	if (store.getters.isAdmin) {
		next()
	} else {
		next("/")
	}
}

async function main() {
	store.dispatch("restorePlaybackRate")
	// Create the Vue instance
	new Vue({
		el: document.body,
		render: h => h("App"),
		router: new VueRouter({
			mode: "history",
			routes: [
				{
					path: "/login",
					component: Login,
					beforeEnter: nonAuthGuard,
				},
				{
					path: "/",
					component: Main,
					beforeEnter: authGuard,
					children: [
						{ path: "/", component: Library },
						{
							path: "/audiobook/upload",
							component: Upload,
							beforeEnter: uploaderGuard,
						},
						{
							path: "/audiobook/:id/edit",
							component: Edit,
							beforeEnter: uploaderGuard,
						},
						{
							path: "/admin",
							component: Admin,
							beforeEnter: adminGuard,
						},
						{
							path: "/user/:id/edit",
							component: EditUser,
							beforeEnter: adminGuard,
						},
						{ path: "/settings", component: Settings },
						{ path: "*", component: NotFound },
					],
				},
			],
		}),
		store,
		components: { App },
	})
}

main()
