import { createStore } from "vuex";
import Events from "../events";

const events = new Events("http://localhost:4000");

export default createStore({
  state: {
    loggedIn: false,
    user: "",
    menuOptions: {
      options: [],
      handler: null,
    },
  },
  mutations: {
    logIn: function (state, user) {
      state.loggedIn = true;
      state.user = user;
    },
    logOut: function (state) {
      state.loggedIn = false;
      state.user = "";
    },
    menuOptions: function (state, options) {
      state.menuOptions = options;
    },
  },
  actions: {
    logIn(context) {
      return events
        .getUserData()
        .then((user) => context.commit("logIn", user.name));
    },
    logOut(context) {
      return events.logOut().then(() => context.commit("logOut"));
    },
  },
  modules: {},
});
