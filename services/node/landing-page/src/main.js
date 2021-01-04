import Vue from 'vue'
import VueRouter from 'vue-router'
import App from './App.vue'
import HelloWorld from './components/HelloWorld.vue'

Vue.config.productionTip = false
Vue.use(VueRouter)

const authGuard = (to, from) => {
  console.log("Not logged in or out: ", from.path, " -> ", to.path)
  window.location.replace(`https://auth.ponglehub.co.uk/login?redirect=https://games.ponglehub.co.uk/#${to.fullPath}`)
}

const routes = [
  { path: '/foo', component: HelloWorld, props: { msg: "Foo component!" }, beforeEnter: authGuard },
  { path: '/bar', component: HelloWorld, props: { msg: "Bar component!"}, beforeEnter: authGuard },
  { path: '*', component: HelloWorld, props: { msg: "Missing component!"} },
]

const router = new VueRouter({
  routes
})

new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
