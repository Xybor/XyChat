import {
  getAccountInfo,
  getLoginStatus,
  setAccountInfo,
  setLogin,
  setLogout,
} from "../../helpers/localStorageManager";
import { userSerive } from "../../services/userService";
import router from "../../router";
import { capitalizeFirstLetter } from "../../helpers/stringUtils";

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
      // Check if login success
      if (response.data.meta.errno == 0) {
        let accountInfo = {
          username: username,
          id: 0,
        };
        setLogin();
        setAccountInfo(accountInfo.id, accountInfo.username);
        commit("loginSuccess", accountInfo);

        dispatch("alert/success", "Login successfully", { root: true });
        router.push("/profile");
      } else {
        commit("loginFailure");
        dispatch("alert/error", response.data.meta.error, { root: true });
      }
    }).catch((err) => {
      // Check if error have response
      if (err.response) {
        commit("loginFailure");
        dispatch(
          "alert/error",
          capitalizeFirstLetter(err.response.data.meta.error),
          { root: true }
        );
      }
    });
  },
  register({ dispatch, commit }, { username, password }) {
    commit("registerRequest", { username });
    var rs = userSerive.register(username, password);
    rs.then((response) => {
      if (response.data.meta.errno == 0) {
        commit("registerSuccess");
        dispatch("alert/success", "Register successfully", { root: true });
      }
    }).catch((err) => {
      if (err.response) {
        commit("registerFailure");
        dispatch(
          "alert/error",
          capitalizeFirstLetter(err.response.data.meta.error),
          { root: true }
        );
      }
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
      if (response.data.meta.errno == 0) {
        commit("updateAccountInfo", response.data.data);
        dispatch("alert/success", "Check session successfully", { root: true });
      } else {
        dispatch("alert/error", "Session is expired", { root: true });
        commit("logout");
        setLogout();
        router.push("/login");
      }
    }).catch((err) => {
      if (err.response) {
        dispatch("alert/error", "Something went wrong", { root: true });
        commit("logout");
        setLogout();
        router.push("/login");
      }
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
