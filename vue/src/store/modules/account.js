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
      isLoggedIn: true,
      token: savedToken,
      username: savedUsername,
      registering: false,
    }
  : {
      isLoggedIn: false,
      token: null,
      username: null,
      registering: false,
    };
const actions = {
  login({ dispatch, commit }, { username, password }) {
    commit("loginRequest", { token: null, username: username });
    var rs = userSerive.login(username, password);
    rs.then((response) => {
      const data = response.data.data;
      console.log(data);
      if (data.token) {
        const revicedToken = data.token;

        setJWTtAuthenToken(revicedToken);
        setUsername(username);
        commit("loginSuccess", {
          token: revicedToken,
          username: username,
        });
        dispatch("alert/success", "Login successfully", { root: true });
        router.push("/profile");
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
      dispatch("alert/success", "Register successfully", { root: true });
    }).catch((err) => {
      commit("registerFailure");
      dispatch("alert/error", "Register fail", { root: true });
    });
  },
  logout({ dispatch, commit }) {
    commit("logout");
    router.push("/");

    userSerive.logout();
    dispatch("alert/success", "Logout success", { root: true });
  },
};

const mutations = {
  loginRequest(state, data) {
    state.isLoggedIn = true;
    state.username = data.username;
    state.token = null;
  },
  loginSuccess(state, data) {
    state.isLoggedIn = true;
    state.username = data.username;
    state.token = data.token;
  },
  loginFailure(state) {
    state.isLoggedIn = false;
    state.username = null;
    state.token = null;
  },
  logout(state) {
    state.isLoggedIn = false;
    state.username = null;
    state.token = null;
  },
  registerRequest(state) {
    state.registering = true;
  },
  registerSuccess(state) {
    state.registering = true;
  },
  registerFailure(state) {
    state.registering = false;
  },
};

export const account = {
  namespaced: true,
  state,
  actions,
  mutations,
};
