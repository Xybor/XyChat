<template>
  <div>
    <Navbar />
    <router-view />
  </div>
</template>

<script>
import Navbar from "@/components/Navbar.vue";
import { useStore } from "vuex";
import { useToast } from "vue-toastification";
import { onMounted } from "vue";

export default {
  components: {
    Navbar,
  },
  setup() {
    const store = useStore();
    const toast = useToast();

    document.title = "XyChat";

    const notify = (msg, type) => {
      toast(msg, {
        position: "bottom-right",
        timeout: 2000,
        closeOnClick: true,
        pauseOnFocusLoss: false,
        pauseOnHover: true,
        draggable: true,
        draggablePercent: 0.6,
        showCloseButtonOnHover: false,
        hideProgressBar: false,
        closeButton: "button",
        icon: true,
        rtl: false,
        type: type,
      });
    };

    onMounted(() => {
      store.watch(
        (state) => state.alert,
        (alertState) => {
          notify(alertState.message, alertState.type);
        },
        {
          deep: true,
        }
      );
    });
  },
};
</script>

<style>
#nav a.router-link-exact-active {
  color: #42b983;
  font-weight: bold;
}
</style>
