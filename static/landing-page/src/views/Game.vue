<template>
  <iframe :src="url" />
</template>

<script>
// @ is an alias to /src
export default {
  name: "Game",
  props: {
    game: String,
  },
  computed: {
    url: function () {
      if (this.game == "draughts") return "http://localhost:3001";
      if (this.game == "naughts-and-crosses") return "http://localhost:3002";
      return "";
    },
  },
};
</script>

<style scoped>
iframe {
  height: calc(100% - 3.8em);
  width: 100%;
}
</style>
