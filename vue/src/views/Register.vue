<template>
  <div class="login container mt-5">
    <h3 class="text-center">Register</h3>

    <div class="row justify-content-md-center">
      <div class="col-3">
        <div class="form-group">
          <label for="firstname">First name</label>
          <input type="text" class="form-control" id="firstname" />
        </div>
      </div>

      <div class="col-3">
        <div class="form-group">
          <label for="lastname">Last name</label>
          <input type="text" class="form-control" id="lastname" />
        </div>
      </div>
    </div>
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
      <div class="form-group mt-2 col-6">
        <label for="exampleFormControlInput1">Gender</label>
        <select class="form-select" aria-label="Default select example">
          <option value="1" selected>Male</option>
          <option value="2">Female</option>
          <option value="3">Other</option>
        </select>
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
        />
      </div>
    </div>
    <div class="row justify-content-md-center mt-2">
      <div class="form-group col-6">
        <label for="password">Verify Password</label>
        <input
          type="password"
          class="form-control"
          id="password"
          placeholder="******"
          v-model="password"
        />
      </div>
    </div>

    <div class="row justify-content-md-center mt-2">
      <div class="form-group col-6">
        Switch to <router-link to="/login">Login</router-link>
      </div>
    </div>

    <div class="text-center mt-3">
      <button
        class="btn btn-outline-success text-center"
        type="submit"
        @click="handleSubmit"
      >
        Create new account
      </button>
    </div>
  </div>
</template>

<script>
import { ref } from "vue";
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

      store.dispatch("account/register", {
        username: username.value,
        password: password.value,
      });
    };

    return {
      username,
      password,
      handleSubmit,
    };
  },
};
</script>
