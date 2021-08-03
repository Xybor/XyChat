<template>
  <div class="profile container mt-5">
    <h3>Profile</h3>
    <p>ID: {{ profile.id }}</p>
    <p>Username: {{ profile.username }}</p>
    <p>Role: {{ profile.role }}</p>
  </div>
</template>

<script>
import { userSerive } from "../services/userService";
export default {
  data() {
    return {
      profile: {},
    };
  },
  mounted() {
    var rs = userSerive.getProfile();
    rs.then((response) => {
      this.profile = response.data.data;
    });
  },
};
</script>
