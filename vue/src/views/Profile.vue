<template>
  <div class="profile container mt-5">
    <div class="row">
      <div class="card col-lg-4 no-radius">
        <div class="card-body text-center mt-3">
          <img
            class="rounded-circle"
            src="https://via.placeholder.com/256"
            alt="avatar"
          />
          <div class="mt-3">
            <span><b>Huy Gay</b></span>
            <br />
            <span class="text-muted">@{{ profile.username }}</span>
          </div>
        </div>
      </div>
      <div class="card col-lg-8 no-radius">
        <div class="card-body">
          <h5>Profile settings</h5>
          <div class="row">
            <div class="col-6">
              <div class="form-group">
                <label for="firstname">First name</label>
                <input type="text" class="form-control" id="firstname" />
              </div>
            </div>

            <div class="col-6">
              <div class="form-group">
                <label for="lastname">Last name</label>
                <input type="text" class="form-control" id="lastname" />
              </div>
            </div>
          </div>

          <div class="form-group mt-2">
            <label for="exampleFormControlInput1">Role</label>
            <select class="form-select" aria-label="Default select example">
              <option value="1" selected>One</option>
              <option value="2">Two</option>
              <option value="3">Three</option>
            </select>
          </div>

          <button class="btn btn-outline-success mt-3" @click="changeProfile">
            Change profile
          </button>

          <h5 class="mt-5">Change password</h5>
          <div class="form-group">
            <label for="oldPassword">Old password</label>
            <input
              type="email"
              class="form-control"
              id="oldPassword"
              placeholder="********"
              v-model="password.oldPassword"
            />
          </div>

          <div class="form-group mt-2">
            <label for="newPassword">New password</label>
            <input
              type="email"
              class="form-control"
              id="newPassword"
              placeholder="********"
              v-model="password.newPassword"
            />
          </div>

          <div class="form-group mt-2">
            <label for="verifyPassword">Verify new password</label>
            <input
              type="email"
              class="form-control"
              id="verifyPassword"
              placeholder="********"
              v-model="password.verifyPassword"
            />
          </div>

          <button class="btn btn-outline-success mt-3" @click="changePassword">
            Change password
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { userSerive } from "../services/userService";
import { ref, onMounted } from "vue";
export default {
  setup() {
    const profile = ref({});
    const password = ref({
      oldPassword: "",
      newPassword: "",
      verifyPassword: "",
    });

    const getProfile = () => {
      const rs = userSerive.getProfile();
      rs.then((response) => {
        console.log(response.data.data);
        profile.value = response.data.data;
      });
    };

    onMounted(() => {
      getProfile();
    });

    var changeProfile = () => {};
    var changePassword = () => {};

    return {
      profile,
      password,
    };
  },
};
</script>

<style scoped>
label {
  color: gray;
}
</style>
