<template lang="pug">
	body(@dragover.stop.prevent="dragOver" @drop.stop.prevent="drop")
		router-view
</template>

<script>
export default {
	methods: {
		// These functions are registered to the body in order to prevent the
		// default behavior of browsers opening files when they are dragged into
		// the window. This makes uploading smoother as the user wont be
		// prompted by their browser to download the file they intended to
		// upload.
		dragOver(e) {
			e.dataTransfer.dropEffect = "none"
		},
		drop(e) {},
	},
}
</script>
