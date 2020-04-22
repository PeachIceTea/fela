<template lang="pug">
	.upload-page
		.half.upload-form
			Status
			Box
		.half.uploads
			List(v-show="!$route.params.id")
			router-view(v-show="$route.params.id")
</template>

<script>
import Box from "./Box.vue"
import List from "./List.vue"
import Status from "./Status"

export default {
	components: { Box, List, Status },
}
</script>

<style lang="stylus" scoped>
@import "../../globals"

.upload-page
	display: flex
	padding: 1em

.half
	flex: 1
	width: 50%

.uploads
	margin-left: 1em
</style>
