import {
  getJWTAuthenToken,
  removeJWTtAuthenToken,
} from "../helpers/authen-token";
import { configs } from "./config";

import axios from "axios";

export const userSerive = {
  login,
  logout,
  register,
  getProfile,
};

function login(username, password) {
  const body = {
    username: username,
    password: password,
  };

  return axios.post(`${configs.apiUrl}/auth`, body);
}

function register(username, password) {
  const params = new URLSearchParams();
  params.append("username", username);
  params.append("password", password);

  return axios.post(`${configs.apiUrl}/register`, { params: params });
}

async function getProfile() {
  const params = new URLSearchParams();
  params.append("token", getJWTAuthenToken());
  return await axios.get(`${configs.apiUrl}/profile`, { params: params });
}

function logout() {
  removeJWTtAuthenToken();
}
