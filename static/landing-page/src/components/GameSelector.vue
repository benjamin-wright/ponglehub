<template>
  <div class="list">
    <router-link to="/naughts-and-crosses">
      <img src="../assets/naughts-and-crosses.png" width="128" height="128" />
    </router-link>
    <router-link to="/draughts">
      <img src="../assets/draughts.png" width="128" height="128" />
    </router-link>
  </div>
</template>

<script>
export default {
  name: "GameSelector",
  props: {},
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

.list {
  display: flex;
}

.list a {
  margin: 2em;
}

</style>
