import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";

const landingPageUrl = "http://localhost:3000";
// const draughtsUrl = "http://localhost:3001";
const authUrl = "http://localhost:4000";

router.beforeEach((to, from, next) => {
  if (store.state.loggedIn) {
    return next();
  }

  if (process.env.NODE_ENV == "development") {
    console.log("dev mode: bypassing login");
    store.commit("logIn");
    return next();
  }

  window.location = `${authUrl}/auth/login?redirect=${landingPageUrl}${to.path}`;
});

createApp(App).use(store).use(router).mount("#app");

