<template lang="pug">
	.upload-metadata
		form.meta-editor(@submit.prevent="submit")
			div
				input(id="title" v-model="title" placeholder="Book Title")
			div
				input(id="author" v-model="author" list="author_list" placeholder="Author")
			div
				textarea(id="description" v-model="description" placeholder="Description (optional)")
			div
				input(type="submit" :disabled="disableButton")

		datalist#author_list
			option(v-for="author in authorList" :value="author")

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
			authorList: [],
		}
	},
	async created() {
		try {
			const res = await fetch("http://localhost:8080/author")

			if (res.status === 200) {
				this.authorList = await res.json()
			}
		} catch (e) {
			// Its not too bad if we cannot get the list
		}
	},
	methods: {
		submit() {
			this.callback(this.file, {
				title: this.title,
				author: this.author,
				description: this.description,
			})
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
