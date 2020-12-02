import Vue from 'vue'
import VueRouter from 'vue-router'
import App from './App.vue'
import HelloWorld from './components/HelloWorld.vue'

Vue.config.productionTip = false
Vue.use(VueRouter)

const routes = [
  { path: '/foo', component: HelloWorld, props: { msg: "Foo component!"} },
  { path: '/bar', component: HelloWorld, props: { msg: "Bar component!"} },
  { path: '*', component: HelloWorld, props: { msg: "Missing component!"} },
]

const router = new VueRouter({
  routes
})

new Vue({
  router,
  render: h => h(App),
}).$mount('#app')
