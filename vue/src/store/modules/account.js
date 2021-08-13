import {
  getAccountInfo,
  getLoginStatus,
  setAccountInfo,
  setLogin,
  setLogout,
} from "../../helpers/localStorageManager";
import { userSerive } from "../../services/userService";
import router from "../../router";

const isLogin = getLoginStatus();
const state = isLogin
  ? {
      isLoggedIn: true,
      accountInfo: getAccountInfo(),
    }
  : {
      isLoggedIn: false,
      accountInfo: {
        username: null,
        id: 0,
      },
    };

const actions = {
  login({ dispatch, commit }, { username, password }) {
    commit("loginRequest", { token: null, username: username });
    var rs = userSerive.login(username, password);
    rs.then((response) => {
      if (response.data.meta.errno == 0) {
        // Login Success
        setLogin();
        setAccountInfo(0, username);

        commit("loginSuccess", {
          id: 0,
          username: username,
        });
        dispatch("alert/success", "Login successfully", { root: true });
        router.push("/profile");
      } else {
        // Login fail
        commit("loginFailure");
        dispatch("alert/error", "Login Fail", { root: true });
      }
    }).catch((err) => {
      // Login error
      commit("loginFailure");
      dispatch("alert/error", "Something went wrong", { root: true });
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
    setLogout();
    dispatch("alert/success", "Logout success", { root: true });
  },
  checkToken({ dispatch, commit }) {
    const rs = userSerive.getProfile();
    rs.then((response) => {
      console.log(response);
      if (response.data.meta.errno == 0) {
        commit("updateAccountInfo", {
          id: response.data.data.id,
        });
        dispatch("alert/success", "Check session successfully", { root: true });
      } else {
        dispatch("alert/error", "Session is expired", { root: true });
        commit("logout");
        setLogout();
        router.push("/login");
      }
    }).catch((err) => {
      dispatch("alert/error", "Something went wrong", { root: true });
      commit("logout");
      setLogout();
      router.push("/login");
    });
  },
};

const mutations = {
  loginRequest(state, data) {
    state.isLoggedIn = true;
    state.accountInfo = {
      username: data.username,
      id: data.id,
    };
  },
  loginSuccess(state, data) {
    state.isLoggedIn = true;
    state.accountInfo = {
      ...state.accountInfo,
      ...data,
    };
  },
  loginFailure(state) {
    state.isLoggedIn = false;
    state.accountInfo = null;
  },
  logout(state) {
    state.isLoggedIn = false;
    state.accountInfo = null;
  },
  registerRequest(state) {
    state.isLoggedIn = false;
    state.accountInfo = null;
  },
  registerSuccess(state) {
    state.isLoggedIn = false;
    state.accountInfo = null;
  },
  registerFailure(state) {
    state.isLoggedIn = false;
    state.accountInfo = null;
  },
  updateAccountInfo(state, accountInfo) {
    state.account = {
      ...state.account,
      ...accountInfo,
    };
  },
};

export const account = {
  namespaced: true,
  state,
  actions,
  mutations,
};
