import { createRouter, createWebHistory } from "vue-router";
import { getLoginStatus } from "../helpers/localStorageManager";
import Home from "../views/Home.vue";
import Login from "../views/Login.vue";
import Register from "../views/Register.vue";
import Profile from "../views/Profile.vue";
import Chat from "../views/Chat.vue";
import ErrorPage from "../views/404.vue";

const routes = [
  {
    path: "/",
    name: "Home",
    component: Home,
  },
  {
    path: "/login",
    name: "Login",
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: Login,
  },
  {
    path: "/register",
    name: "Register",
    component: Register,
  },
  {
    path: "/profile",
    name: "Profile",
    component: Profile,
  },
  {
    path: "/chat",
    name: "Chat",
    component: Chat,
  },
  {
    path: "/:pathMatch(.*)*",
    component: ErrorPage,
  },
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

router.beforeEach((to, from, next) => {
  const publicPages = ["/", "/home", "/login", "/register"];
  const authenRequired = !publicPages.includes(to.path);
  const loggedIn = getLoginStatus();

  if (authenRequired && !loggedIn) {
    return next("/login");
  }

  next();
});

export default router;
