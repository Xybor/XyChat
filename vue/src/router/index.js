import { createRouter, createWebHistory } from "vue-router";
import { getJWTAuthenToken } from "../helpers/authen-token";
import Home from "../views/Home.vue";
// import Login from "../views/Login.vue";
// import Profile from "../views/Profile.vue";
// import Register from "../views/Register.vue";

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
    component: () => import(/* webpackChunkName: "" */ "../views/Login.vue"),
  },
  {
    path: "/register",
    name: "Register",
    component: () => import(/* webpackChunkName: "" */ "../views/Register.vue"),
  },
  {
    path: "/profile",
    name: "Profile",
    component: () => import(/* webpackChunkName: "" */ "../views/Profile.vue"),
  },
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

router.beforeEach((to, from, next) => {
  const publicPages = ["/", "/home", "/login", "/register"];
  const authenRequired = !publicPages.includes(to.path);
  const loggedIn = getJWTAuthenToken();

  if (authenRequired && !loggedIn) {
    return next("/login");
  }

  next();
});

export default router;
