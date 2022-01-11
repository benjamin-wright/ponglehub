import { createRouter, createWebHistory } from "vue-router";
import Home from "../views/Home.vue";
import Game from "../views/Game.vue";

const routes = [
  {
    path: "/",
    name: "Home",
    component: Home,
    meta: {
      requiresAuth: true,
    },
  },
  {
    path: "/naughts-and-crosses",
    name: "naughts",
    component: Game,
    meta: {
      requiresAuth: true,
      game: "Naughts & Crosses",
    },
    props: {
      game: "naughts-and-crosses",
    },
  },
  {
    path: "/draughts",
    name: "draughts",
    component: Game,
    meta: {
      requiresAuth: true,
      game: "Draughts",
    },
    props: {
      game: "draughts",
    },
  },
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

export default router;
