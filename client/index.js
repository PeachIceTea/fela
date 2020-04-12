import Vue from "vue"
import VueRouter from "vue-router"

import store from "./store"

import App from "./components/App.vue"
import Library from "./components/Library.vue"
import Book from "./components/Book.vue"
import Upload from "./components/Upload"
import NotFound from "./components/NotFound.vue"

Vue.use(VueRouter)

const fillZero = n => (n < 10 ? `0${n}` : n)
Vue.filter("formatDuration", t => {
	const h = Math.floor(t / (60 * 60))
	t -= h * 60 * 60

	const m = fillZero(Math.floor(t / 60))
	t -= m * 60

	const s = fillZero(Math.floor(t))
	return `${h}:${m}:${s}`
})

new Vue({
	el: document.getElementById("vue-el"),
	render: h => h(App),
	router: new VueRouter({
		mode: "history",
		routes: [
			{ path: "/", component: Library },
			{ path: "/book/:id", component: Book },
			{ path: "/upload", component: Upload },
			{ path: "*", component: NotFound },
		],
	}),
	store,
})
