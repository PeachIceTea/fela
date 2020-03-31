<template lang="pug">
	.upload-metadata
		form.meta-editor(@submit.prevent="submit")
			div
				input(id="title" v-model="title" placeholder="Book Title")
			div
				input(id="author" v-model="author" placeholder="Author")
			div
				textarea(id="description" v-model="description" placeholder="Description (optional)")
			div
				input(type="submit" :disabled="disableButton")

</template>

<script>
//TODO: Validation
//TODO: Fetch author list to give suggestions

export default {
	data() {
		return {
			title: "",
			author: "",
			description: "",
		}
	},
	methods: {
		submit() {
			this.callback(this.file, Object.assign({}, this.$data))
		},
	},
	computed: {
		disableButton() {
			return !this.file.fileID
		},
	},
	props: {
		file: Object,
		callback: Function,
	},
}
</script>
