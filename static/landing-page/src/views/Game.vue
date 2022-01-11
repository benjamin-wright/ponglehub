<template>
  <iframe src="http://localhost:8081" />
</template>

<script>
// @ is an alias to /src
export default {
  name: "Game",
  props: {
    game: String,
  },
};
</script>

<style scoped>
iframe {
  height: calc(100% - 3.8em);
  width: 100%;
}
</style>
