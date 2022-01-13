import { createStore } from "vuex";

export default createStore({
  state: {
    loggedIn: false,
    menuOptions: {
      options: [],
      handler: null,
    },
  },
  mutations: {
    logIn: function (state) {
      state.loggedIn = true;
    },
    logOut: function (state) {
      state.loggedIn = false;
    },
    menuOptions: function (state, options) {
      state.menuOptions = options;
    },
  },
  actions: {},
  modules: {},
});
