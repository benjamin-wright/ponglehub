<template>
  <div class="nav">
    <router-link to="/">Ponglehub</router-link>
    <p>{{ game }}</p>
    <a v-if="loggedIn" v-on:click="this.$store.commit('logOut')">Log out</a>
  </div>
</template>

<script>
export default {
  name: "NavBar",
  computed: {
    game: function () {
      return this.$route.meta.game;
    },
    loggedIn: function () {
      return this.$store.state.loggedIn;
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.nav {
  padding: 1em;
  background: #2c3e50;
  color: #dbeeff;
  text-align: left;
  display: flex;
  flex-direction: row;
  justify-content: space-between;
}

a {
  font-weight: bold;
  text-decoration: none;
  text-transform: uppercase;
  cursor: pointer;
}

a.router-link-exact-active {
  color: #dbeeff;
}

a:visited {
  color: #dbeeff;
}

p {
  margin: 0;
}
</style>
