<template>
  <iframe ref="game-frame" :src="url" />
</template>

<script>
import Messenger from "messaging";

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
  created: function () {
    this.messenger = new Messenger();
    this.messenger.listenMenuOptions((options) => {
      this.$store.commit("menuOptions", {
        options,
        handler: (event) => {
          this.messenger.selectMenuOption(this.$refs["game-frame"], event);
        },
      });
    });
  },
  unmounted: function () {
    this.messenger.stop();
  },
};
</script>

<style scoped>
iframe {
  height: calc(100% - 3.8em);
  width: 100%;
  border: none;
}
</style>
