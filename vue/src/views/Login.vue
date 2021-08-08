<template>
  <div class="login container mt-5">
    <h3 class="text-center">Login</h3>
    <div class="row justify-content-md-center">
      <div class="form-group col-6">
        <label for="username">Username</label>
        <input
          type="text"
          class="form-control"
          id="username"
          placeholder="Username"
          v-model="username"
        />
      </div>
    </div>
    <div class="row justify-content-md-center mt-2">
      <div class="form-group col-6">
        <label for="password">Password</label>
        <input
          type="password"
          class="form-control"
          id="password"
          placeholder="******"
          v-model="password"
          v-on:keyup.enter="onEnterLogin"
        />
      </div>
    </div>

    <div class="row justify-content-md-center mt-2">
      <div class="form-group col-6">
        <router-link to="/forgotPassword">Forgot password</router-link>
        <br />
        <router-link to="/register">Create new account</router-link>
      </div>
    </div>

    <div class="text-center mt-3">
      <button
        class="btn btn-outline-success text-center"
        type="submit"
        @click="handleSubmit"
      >
        Login
      </button>
    </div>
  </div>
</template>

<script>
import { ref } from "vue";
// @ is an alias to /src
//import HelloWorld from "@/components/HelloWorld.vue";
import { useStore } from "vuex";

export default {
  setup() {
    // Declear variable
    const username = ref("");
    const password = ref("");
    const store = useStore();

    // Methods
    const handleSubmit = () => {
      if (username.value == "" || password.value == "") {
        store.dispatch("alert/error", "Username or password can not empty");
        return false;
      }

      store.dispatch("account/login", {
        username: username.value,
        password: password.value,
      });
    };

    const onEnterLogin = () => {
      handleSubmit();
    };

    return {
      username,
      password,
      handleSubmit,
      onEnterLogin,
    };
  },
};
</script>
