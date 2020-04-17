import Vue from "vue"
import VueRouter from "vue-router"

Vue.use(VueRouter)
import store from "./store"

import App from "./components/App.vue"
import Login from "./components/Login.vue"
import Main from "./components/Main.vue"
import Library from "./components/Library"
import NotFound from "./components/NotFound"

const fillZero = n => (n < 10 ? `0${n}` : n)
Vue.filter("formatDuration", t => {
	const h = Math.floor(t / (60 * 60))
	t -= h * 60 * 60

	const m = Math.floor(t / 60)
	t -= m * 60

	const s = Math.floor(t)
	return `${h}:${fillZero(m)}:${fillZero(s)}`
})

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
					{ path: "*", component: NotFound },
				],
			},
		],
	}),
	store,
	components: { App },
})
