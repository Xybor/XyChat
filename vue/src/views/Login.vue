<template>
  <div class="login container mt-5">
    <h3>Login</h3>
    <div class="row align-items-center">
      <div class="input-group col-6">
        <span class="input-group-text" id="username-addon">@</span>
        <input
          type="text"
          class="form-control"
          placeholder="Username"
          aria-label="Username"
          aria-describedby="username-addon"
          maxlength="16"
          v-model="username"
        />
        <span class="input-group-text" id="password-addon">**</span>
        <input
          type="password"
          class="form-control"
          placeholder="Password"
          aria-label="Password"
          aria-describedby="password-addon"
          v-model="password"
        />
        <button
          class="btn btn-outline-success"
          type="submit"
          @click="handleSubmit"
        >
          Login
        </button>
      </div>
    </div>
  </div>
</template>

<script>
// @ is an alias to /src
//import HelloWorld from "@/components/HelloWorld.vue";
import { mapState, mapActions } from "vuex";

export default {
  data() {
    return {
      username: "",
      password: "",
    };
  },
  computed: {
    ...mapState("account", ["status"]),
  },
  methods: {
    ...mapActions("account", ["login", "logout"]),
    handleSubmit() {
      this.login({username: this.username, password: this.password});
    },
  },
};
</script>
