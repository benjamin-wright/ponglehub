<template>
  <div id="app">
    <p>
      <router-link to="/foo">Foo Component</router-link>
      <router-link to="/bar">Bar Component</router-link>
    </p>
    <router-view></router-view>
  </div>
</template>

<script>
import HelloWorld from './components/HelloWorld.vue'
import VueRouter from 'vue-router'

const routes = [
  { path: '/foo', component: HelloWorld, props: { msg: "Foo component!"} },
  { path: '/bar', component: HelloWorld, props: { msg: "Bar component!"}  }
]

const router = new VueRouter({
  routes
})

export default {
  name: 'App',
  components: {
    // eslint-disable-next-line vue/no-unused-components
    router
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
