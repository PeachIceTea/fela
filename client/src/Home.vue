<template lang="pug">
	.home
		h1 Home, sweet home.
		.link-list
			router-link(to="/upload") Upload a new file
			br
			router-link(to="/book") Read a book
</template>
