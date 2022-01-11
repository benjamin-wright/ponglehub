import { createStore } from "vuex";

export default createStore({
  state: {
    loggedIn: false,
  },
  mutations: {
    logIn: function (state) {
      state.loggedIn = true;
    },
    logOut: function (state) {
      state.loggedIn = true;
    },
  },
  actions: {},
  modules: {},
});
