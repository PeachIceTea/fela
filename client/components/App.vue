<template lang="pug">
	.app
		h1 Fela
		p
			router-view
		p
			router-link(to="/") Return home
		Player
</template>

<script>
import Player from "./Player"

export default {
	components: { Player },
}
</script>
