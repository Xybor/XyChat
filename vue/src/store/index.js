import Vuex from "vuex";
import { account } from "./modules/account";
import { alert } from "./modules/alert";

export const store = new Vuex.Store({
  modules: {
    account,
    alert,
  },
});
