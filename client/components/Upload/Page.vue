<template lang="pug">
	.upload-page
		.half.upload-form
			Box
		.half.uploads
			List(v-show="!$route.params.id")
			router-view(v-show="$route.params.id")
</template>

<script>
import Box from "./Box.vue"
import List from "./List.vue"

export default {
	components: { Box, List },
}
</script>

<style lang="stylus" scoped>
@import "../../globals"

.upload-page
	display: flex
	padding: 1em

.half
	flex: 1

.uploads
	margin-left: 1em
</style>
