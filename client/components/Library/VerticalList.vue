<template lang="pug">
	.vertical-list
		.container.item(v-for="audiobook in audiobooks")
			.small-cover(@click="play(audiobook)")
				img(
					:src="coverURL(audiobook.id)"
					@error="noImage"
				)
			.text
				.title {{ audiobook.title }}
				.by by
					span.author  {{ audiobook.author }}
			.buttons
				button(@click="edit(audiobook)" v-if="uploader") Edit
				button(@click="play(audiobook)") Play
</template>

<script>
export default {
	props: ["audiobooks", "play", "coverURL", "noImage", "uploader"],
	methods: {
		// Moves user to the edit page for a given audiobook.
		edit(audiobook) {
			this.$router.push(`/audiobook/${audiobook.id}/edit`)
		},
	},
}
</script>

<style lang="less" scoped>
@import "../../globals.less";

.vertical-list {
	margin: 1em;
}

.container {
	padding: 0.3em;
}

.list-library {
	padding: 1em;
}

.item {
	margin-bottom: 0.35em;
	display: flex;
	height: 4em;
}

.text {
	align-self: center;
	text-shadow: @text-shadow;
}

.title,
.author {
	font-size: 100%;
	font-weight: 600;
}

.small-cover {
	height: 100%;
	margin-right: 1em;
	img {
		max-height: 100%;
		max-width: 100%;
	}
}

.buttons {
	flex: 1;
	display: flex;
	justify-content: flex-end;
	align-items: center;
}

button {
	height: 100%;
	width: 4em;

	&:last-child {
		margin-left: 0.5em;
	}
}

@media (max-width: 480px) {
	.cover-library {
		display: grid;
		grid-template-columns: repeat(1, 1fr);
	}

	.buttons {
		display: none;
	}

	.author {
		font-size: 15px;
	}

	.text {
		height: 100%;
		overflow: scroll;
	}

	.item {
		height: 7em;
	}

	.small-cover {
		height: 100%;
		width: auto;
		img {
			height: 100%;
			width: auto;
			max-width: unset;
			max-height: unset;
		}
	}

	.library-header {
		font-size: 25px;
		margin: 24px 0px 10px 10px;
	}
}
</style>
