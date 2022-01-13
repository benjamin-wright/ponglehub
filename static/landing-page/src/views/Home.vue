<template>
  <div class="home">
    <h1>Welcome to PongleHub Games!</h1>
    <PongleFace />
    <a class="arrow">&#9660;</a>
  </div>
  <div class="selector">
    <GameSelector />
  </div>
</template>

<script>
// @ is an alias to /src
import PongleFace from "@/components/PongleFace.vue";
import GameSelector from "@/components/GameSelector.vue";

export default {
  name: "Home",
  components: {
    PongleFace,
    GameSelector,
  },
  mounted: function () {
    this.$store.commit("menuOptions", []);
  },
};
</script>

<style scoped>
h1 {
  margin-top: 2em;
}

.home {
  height: calc(100% - 3em);
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: center;
}

.selector {
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

.arrow {
  font-size: 2em;
  background: #2c3e50;
  color: #dbeeff;
  border-radius: 50%;
  width: 2em;
  height: 2em;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  margin-bottom: 1em;
}
</style>
