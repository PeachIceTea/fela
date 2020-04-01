import Vue from "vue"
import VueRouter from "vue-router"

import App from "./App.vue"
import Home from "./Home.vue"
import Upload from "./Upload.vue"
import Books from "./Books.vue"
import Book from "./Book.vue"
import NotFound from "./NotFound.vue"

Vue.use(VueRouter)

new Vue({
	el: document.body,
	render: h => h(App),
	router: new VueRouter({
		mode: "history",
		routes: [
			{ path: "/", component: Home },
			{ path: "/upload", component: Upload },
			{ path: "/book", component: Books },
			{ path: "/book/:id", component: Book, name: "book" },

			{ path: "*", component: NotFound },
		],
	}),
})
