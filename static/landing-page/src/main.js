import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";

const url = "http://localhost:8080";

router.beforeEach((to, from, next) => {
  if (store.state.loggedIn) {
    next();
  } else {
    window.location = `http://localhost:3000/auth/login?redirect=${url}${to.path}`;
  }
});

createApp(App).use(store).use(router).mount("#app");
