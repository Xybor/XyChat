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

async function login(username, password) {
  const params = new URLSearchParams();
  params.append("username", username);
  params.append("password", password);

  return await axios.get(`${configs.apiUrl}/auth`, { params: params });
}

function register(username, password) {
  const params = new URLSearchParams();
  params.append("username", username);
  params.append("password", password);

  return axios.get(`${configs.apiUrl}/register`, { params: params });
}

async function getProfile() {
  const params = new URLSearchParams();
  params.append("token", getJWTAuthenToken());
  return await axios.get(`${configs.apiUrl}/profile`, { params: params });
}

function logout() {
  removeJWTtAuthenToken();
}
