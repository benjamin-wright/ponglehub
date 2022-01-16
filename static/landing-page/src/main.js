import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";

const landingPageUrl = "http://localhost:3000";
const authUrl = "http://localhost:4000";

router.beforeEach((to, from, next) => {
  if (store.state.loggedIn) {
    return next();
  }

  store
    .dispatch("logIn")
    .then(() => next())
    .catch(
      () =>
        (window.location = `${authUrl}/auth/login?redirect=${landingPageUrl}${to.path}`)
    );
});

createApp(App).use(store).use(router).mount("#app");
