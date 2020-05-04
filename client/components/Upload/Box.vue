<template lang="pug">
	label
		.upload-box(
			ref="box"
			@drop.stop.prevent="dropFile"
			@dragover.stop.prevent="dragOver"
			@dragleave.stop.prevent="dragExit"
			:class="{'file-hovering': fileHovering}"
		)
			.text
				p #[b Click to choose files]
				p or drag them here
			input(type="file" multiple @input="dialogFile")
</template>

<script>
export default {
	data() {
		return {
			fileHovering: false,
		}
	},
	methods: {
		// Called when a file is let
		dropFile(e) {
			this.fileHovering = false
			const files = e.dataTransfer.files
			if (files) {
				for (let i = 0, len = files.length; i < len; i++) {
					// Not all files are assigned a mime. As an example Firefox does
					// not assign one for ".m4b"s. We will just have to assume the
					// user knows what they are doing if we cannot check the mime.
					const file = files[i]
					const type = file.type
					if (type && !type.includes("audio")) {
						this.$store.dispatch("upload", {
							files,
							err: `${file.name}: invalid file type "${type}"`,
						})
						return
					}
				}
				this.handleFiles(files)
			}
		},
		dragOver(e) {
			e.dataTransfer.dropEffect = "copy"
			this.fileHovering = true
		},
		dragExit(e) {
			this.fileHovering = false
		},
		async dialogFile(e) {
			await this.handleFiles(e.srcElement.files)
			e.srcElement.value = ""
		},
		async handleFiles(arr) {
			// A file input element has only a single FileList which is reused
			// when needed. To keep references to the files we need to create a
			// new array. Additionally this gives us access to array functions
			// not available with FileList.
			const files = []
			for (let i = 0, len = arr.length; i < len; i++) {
				const file = arr[i]
				files[i] = file
			}
			files.sort((a, b) => {
				if (a.name < b.name) {
					return -1
				} else if (a.name > b.name) {
					return 1
				} else {
					return 0
				}
			})

			await this.$store.dispatch("upload", { files })
			this.$store.dispatch("getAudiobooks") // refresh list
		},
	},
}
</script>

<style lang="less" scoped>
@import "../../globals.less";

.upload-box {
	.container();

	width: 100%;
	height: 100%;
	text-align: center;
	cursor: pointer;
	padding: 4em;
	transition: 250ms all ease;
}

label {
	cursor: pointer;
}

.file-hovering {
	background: lighten(@background, 25%);
}
</style>
