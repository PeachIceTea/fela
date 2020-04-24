<template lang="pug">
	.upload-page
		Status
		Box
</template>

<script>
import Box from "./Box.vue"
import Status from "./Status.vue"

export default {
	components: { Box, Status },
}
</script>

<style lang="less" scoped>
@import "../../globals.less";

.upload-page {
	padding: 1em;
}
</style>
