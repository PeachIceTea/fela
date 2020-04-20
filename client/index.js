import Vue from "vue"
import VueRouter from "vue-router"

Vue.use(VueRouter)
import store from "./store"

import App from "./components/App.vue"
import Login from "./components/Login.vue"

import Main from "./components/Main.vue"
import Library from "./components/Library"

import Upload from "./components/Upload"
import UploadEdit from "./components/Upload/Edit.vue"

import NotFound from "./components/NotFound"

// Formats a number to a string representing duration,
// e.g. 88762.131156 to "24:39:22".
function formatDuration(duration) {
	const fillZero = n => (n < 10 ? `0${n}` : n)
	const h = Math.floor(duration / (60 * 60))
	duration -= h * 60 * 60

	const m = Math.floor(duration / 60)
	duration -= m * 60

	const s = Math.floor(duration)
	return `${h}:${fillZero(m)}:${fillZero(s)}`
}
Vue.filter("formatDuration", formatDuration)

async function main() {
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
					beforeEnter: async (to, from, next) => {
						await store.dispatch("restoreToken")
						if (store.state.auth.token) {
							next("/")
						} else {
							next()
						}
					},
				},
				{
					path: "/",
					component: Main,
					beforeEnter: async (to, from, next) => {
						await store.dispatch("restoreToken")
						if (store.state.auth.token) {
							next()
						} else {
							next("/login")
						}
					},
					children: [
						{ path: "/", component: Library },
						{
							path: "/upload",
							component: Upload,
							children: [
								{ path: "edit/:id", component: UploadEdit },
							],
						},
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
