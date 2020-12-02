<template>
  <div id="app">
    <p>
      <router-link to="/foo">Foo Component</router-link>
      <router-link to="/bar">Bar Component</router-link>
    </p>

    <p>hi there</p>
    <router-view></router-view>

  </div>
</template>

<script>

export default {
  name: 'App'
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
