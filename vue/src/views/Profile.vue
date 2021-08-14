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
            <span
              ><b>{{ profile.fname }} {{ profile.lname }}</b></span
            >
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
                <input
                  type="text"
                  class="form-control"
                  id="firstname"
                  v-model="profile.fname"
                />
              </div>
            </div>

            <div class="col-6">
              <div class="form-group">
                <label for="lastname">Last name</label>
                <input
                  type="text"
                  class="form-control"
                  id="lastname"
                  v-model="profile.lname"
                />
              </div>
            </div>
          </div>

          <div class="form-group mt-2">
            <label for="exampleFormControlInput1">Gender</label>
            <select class="form-select" v-model="profile.gender">
              <option value="0"></option>
              <option value="male">Male</option>
              <option value="female">Female</option>
              <option value="other">Other</option>
            </select>
          </div>

          <button class="btn btn-outline-success mt-3" @click="changeProfile">
            Change profile
          </button>

          <h5 class="mt-5">Change password</h5>
          <div class="form-group">
            <label for="oldPassword">Old password</label>
            <input
              type="password"
              class="form-control"
              id="oldPassword"
              placeholder="********"
              v-model="password.oldPassword"
            />
          </div>

          <div class="form-group mt-2">
            <label for="newPassword">New password</label>
            <input
              type="password"
              class="form-control"
              id="newPassword"
              placeholder="********"
              v-model="password.newPassword"
            />
          </div>

          <div class="form-group mt-2">
            <label for="verifyPassword">Verify new password</label>
            <input
              type="password"
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
import { useStore } from "vuex";
import { capitalizeFirstLetter } from "../helpers/stringUtils";

export default {
  setup() {
    const store = useStore();
    const profile = ref({
      fname: "fname",
      lname: "lname",
      gender: "0",
      ...store.state.account.accountInfo,
    });
    const password = ref({
      oldPassword: "",
      newPassword: "",
      verifyPassword: "",
    });

    const getProfile = () => {
      const rs = userSerive.getProfile();
      rs.then((response) => {
        profile.value = {
          ...profile.value,
          ...response.data.data,
        };
      }).catch((err) => {});
    };

    onMounted(() => {
      getProfile();
    });

    var changeProfile = () => {};
    var changePassword = () => {
      if (password.value.newPassword != password.value.verifyPassword) {
        store.dispatch("alert/error", "Passwords don't match");
      } else {
        const rs = userSerive.changePassword(
          profile.value.id,
          password.value.oldPassword,
          password.value.newPassword
        );
        rs.then((response) => {
          if (response.data.meta.errno == 0) {
            store.dispatch("alert/success", "Password has been changed");
            password.value = {
              oldPassword: "",
              newPassword: "",
              verifyPassword: "",
            };
          } else {
            store.dispatch("alert/error", response.data.meta.error);
          }
        }).catch((err) => {
          if (err.response) {
            store.dispatch(
              "alert/error",
              capitalizeFirstLetter(err.response.data.meta.error)
            );
          }
        });
      }
    };

    return {
      profile,
      password,
      changePassword,
    };
  },
};
</script>

<style scoped>
label {
  color: gray;
}
</style>
