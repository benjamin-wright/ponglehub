import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import Messenger from "messaging";

let messenger = new Messenger();

messenger.menuOptions(["start", "stop"], (event) =>
  console.info(`Got message from parent: ${event}`)
);

createApp(App).use(store).use(router).mount("#app");
