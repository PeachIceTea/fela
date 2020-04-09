import Vue from "vue"
import VueRouter from "vue-router"

import store from "./store"

import App from "./components/App.vue"
import Home from "./components/Home.vue"
import Books from "./components/Books.vue"
import Book from "./components/Book"
import Audiobook from "./components/Audiobook"
import Upload from "./components/Upload"
import NotFound from "./components/NotFound.vue"

Vue.use(VueRouter)

new Vue({
	el: document.body,
	render: h => h(App),
	router: new VueRouter({
		mode: "history",
		routes: [
			{ path: "/", component: Home },

			{ path: "/book", component: Books },
			{ path: "/book/:id", component: Book, name: "book" },

			{ path: "/audiobook/:id", component: Audiobook, name: "audiobook" },

			{ path: "/upload", component: Upload },
			{ path: "*", component: NotFound },
		],
	}),
	store,
})
