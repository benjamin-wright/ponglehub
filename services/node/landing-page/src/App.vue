<template>
  <div id="app">
    <p>
      <router-link to="/foo">Foo Component</router-link>
      <router-link to="/bar">Bar Component</router-link>
      <button v-on:click="logout">logout</button>
    </p>

    <p>hi there</p>
    <router-view></router-view>

  </div>
</template>

<script>

import axios from 'axios'

axios.defaults.withCredentials = true

export default {
  name: 'App',
  methods: {
    logout: function(event) {
      console.log("hi!");
      axios
        .post('https://auth.ponglehub.co.uk/api/logout')
        .then(response => {
          this.result = response.data
        })
        .catch(error => { console.log("Failed: ", error.message); this.result = error; })
        .finally(() => this.done = true)
    }
  }
}

</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
