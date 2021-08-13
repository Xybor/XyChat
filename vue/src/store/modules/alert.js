const state = {
  type: null,
  message: null,
  triggerAlert: false,
};

const actions = {
  success({ commit }, message) {
    commit("success", message);
  },
  error({ commit }, message) {
    commit("error", message);
  },
};

const mutations = {
  success(state, message) {
    state.type = "success";
    state.message = message;
    state.triggerAlert = !state.triggerAlert;
  },
  error(state, message) {
    state.type = "error";
    state.message = message;
    state.triggerAlert = !state.triggerAlert;
  },
};

export const alert = {
  namespaced: true,
  state,
  actions,
  mutations,
};
