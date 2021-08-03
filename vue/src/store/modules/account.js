import {
  getJWTAuthenToken,
  setJWTtAuthenToken,
  getUsername,
  setUsername,
} from "../../helpers/authen-token";
import { userSerive } from "../../services/userService";
import router from "../../router";

const savedToken = getJWTAuthenToken();
const savedUsername = getUsername();
const state = savedToken
  ? {
      status: { loggedIn: true },
      user: { token: savedToken, username: savedUsername },
    }
  : { status: { loggedIn: false }, user: null };
const actions = {
  login({ dispatch, commit }, { username, password }) {
    commit("loginRequest", { token: null, username: username });
    var rs = userSerive.login(username, password);
    rs.then((response) => {
      const data = response.data.data;
      if (data.token) {
        const revicedToken = data.token;
        const revicedUsername = data.username;

        setJWTtAuthenToken(revicedToken);
        setUsername(revicedUsername);
        router.push("/profile");
        commit("loginSuccess", {
          token: revicedToken,
          username: revicedUsername,
        });
      } else {
        commit("loginFailure");
        dispatch("alert/error", "Login Fail", { root: true });
      }
    }).catch((err) => {
      commit("loginFailure");
      dispatch("alert/error", "Login Fail", { root: true });
    });
  },
  register({ dispatch, commit }, { username, password }) {
    commit("registerRequest", { username });
    var rs = userSerive.register(username, password);
    rs.then((response) => {
      commit("registerSuccess");
      dispatch("alert/success", "Register success fully", { root: true });
    }).catch((err) => {
      commit("registerFailure");
      dispatch("alert/error", "Register fail", { root: true });
    });
  },
  logout({ dispatch, commit }) {
    commit("logout");
    router.push("/");
    userSerive.logout();
  },
};

const mutations = {
  loginRequest(state, user) {
    state.status = { loggedIn: true };
    state.user = user;
  },
  loginSuccess(state, user) {
    state.status = { loggedIn: true };
    state.user = user;
  },
  loginFailure(state) {
    state.status = { loggedIn: false };
    state.user = null;
  },
  logout(state) {
    state.status = { loggedIn: false };
    state.user = null;
  },
  registerRequest(state, user) {
    state.status = { registering: true };
  },
  registerSuccess(state, user) {
    state.status = { registering: false };
  },
  registerFailure(state, error) {
    state.status = { registering: false };
  },
};

export const account = {
  namespaced: true,
  state,
  actions,
  mutations,
};
