<template>
  <div class="chat">
    <div class="alert alert-info" role="alert">
      You don't have any chats yet
      <br />
      <button class="btn btn-primary mt-2" @click="onMatchButton">
        Match now!
      </button>
    </div>
  </div>
</template>

<script>
import { chatService } from "../services/chatService";
export default {
  setup() {
    var onMatchButton = () => {
      chatService.match();
    };

    return {
      onMatchButton,
    };
  },
};
</script>
