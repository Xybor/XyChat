import { createApp } from "vue";
import App from "./App.vue";
import { Vuex } from "vuex";
import router from "./router";
import { store } from "./store/index";

import "bootstrap/dist/css/bootstrap.min.css";
import "bootstrap";
import "../assets/css/main.css";

import Toast from "vue-toastification";
// Import the CSS or use your own!
import "vue-toastification/dist/index.css";

const options = {
  transition: "Vue-Toastification__bounce",
  maxToasts: 20,
  newestOnTop: true,
};

createApp(App).use(store).use(router).use(Toast, options).mount("#app");
