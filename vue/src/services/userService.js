import { axiosInstance } from "../helpers/getAxios";

export const userSerive = {
  login,
  register,
  getProfile,
  changePassword,
};

function login(username, password) {
  const params = new URLSearchParams();
  params.append("username", username);
  params.append("password", password);
  return axiosInstance.post("/auth", params);
}

function register(username, password) {
  const params = new URLSearchParams();
  params.append("username", username);
  params.append("password", password);

  return axiosInstance.post("/register", params);
}

function getProfile() {
  return axiosInstance.get("/profile");
}

function changePassword(userId, oldpassword, newpassword) {
  const params = new URLSearchParams();
  params.append("oldpassword", oldpassword);
  params.append("newpassword", newpassword);

  return axiosInstance.put(`/users/${userId}/password`, params);
}
