import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";

const landingPageUrl = "http://localhost:3000";
// const draughtsUrl = "http://localhost:3001";
const authUrl = "http://localhost:3002";

router.beforeEach((to, from, next) => {
  if (process.env.NODE_ENV == "development") {
    store.commit("logIn");
    return next();
  }

  if (store.state.loggedIn) {
    return next();
  }

  window.location = `${authUrl}/auth/login?redirect=${landingPageUrl}${to.path}`;
});

createApp(App).use(store).use(router).mount("#app");

