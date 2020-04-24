<template lang="pug">
	body(@dragover.stop.prevent="dragOver" @drop.stop.prevent="drop")
		router-view
</template>

<script>
export default {
	methods: {
		dragOver(e) {
			e.dataTransfer.dropEffect = "none"
		},
		drop(e) {},
	},
}
</script>

<style lang="less" scoped>
@import "../globals.less";
</style>
