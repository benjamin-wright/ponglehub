<template>
  <div class="nav">
    <router-link to="/">Ponglehub</router-link>
    <p>{{ game }}</p>
    <template v-if="options && options.length > 0">
      <NavMenu
        :items="options"
        @select:option="this.$store.state.menuOptions.handler($event)"
        @select:logout="
          this.$store.commit('logOut');
          this.$router.push('/');
        "
      />
    </template>
    <template v-if="!options || options.length == 0">
      <a v-if="loggedIn" v-on:click="this.$store.commit('logOut')">Log out</a>
    </template>
  </div>
</template>

<script>
import Messenger from "messaging";
import NavMenu from "@/components/NavMenu";

export default {
  name: "NavBar",
  components: {
    NavMenu,
  },
  created: function () {
    this.messenger = new Messenger();
  },
  unmounted: function () {
    this.messenger.stop();
  },
  computed: {
    game: function () {
      return this.$route.meta.game;
    },
    loggedIn: function () {
      return this.$store.state.loggedIn;
    },
    options: function () {
      return this.$store.state.menuOptions.options;
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.nav {
  position: relative;
  padding: 1em;
  background: #2c3e50;
  color: #dbeeff;
  text-align: left;
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  transform-style: preserve-3d;
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
