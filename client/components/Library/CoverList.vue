<template lang="pug">
	.cover-library
		.cover(
			v-for="audiobook in audiobooks"
			@click="play(audiobook)"
		)
			img(
				:src="coverURL(audiobook.id)"
				@error="noImage"
			)
</template>

<script>
export default {
	props: ["audiobooks", "play", "coverURL", "noImage"],
}
</script>

<style lang="less" scoped>
.cover-library {
	display: grid;
	grid-template-columns: repeat(3, 1fr);
}

// The ::before on .book and position absolute on img are there to force the
// covers into a 1:1 aspect ratio which most audiobooks covers have. This
// prevents odd shaped covers from distorting other covers in their row.
// Source: https://stackoverflow.com/a/20117454/13394124
.cover {
	cursor: pointer;
	position: relative;
	overflow: hidden;

	&::before {
		content: "";
		display: block;
		padding-top: 100%;
	}

	img {
		position: absolute;
		top: 0;
		height: 100%;
		width: 100%;
		transition: 500ms all ease;

		&:hover {
			/* filter: blur(5px) brightness(0.75); */
			transform: scale(1.03);
		}
	}
}

@media (min-width: 1400px) {
	.cover-library {
		grid-template-columns: repeat(4, 1fr);
	}
}

@media (max-width: 720px) {
	.cover-library {
		grid-template-columns: repeat(2, 1fr);
	}
}

@media (max-width: 550px) {
	.cover-library {
		display: grid;
		grid-template-columns: repeat(1, 1fr);
	}

	.library-header {
		font-size: 25px;
		margin: 24px 0px 10px 10px;
	}
}
</style>
