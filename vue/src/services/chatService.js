import axios from "axios";
import { configs } from "../helpers/config";

export const chatService = {
  match,
};

function match() {
  const ws = new WebSocket(`${configs.wsUrl}/match`);
  ws.onopen = (e) => {
    console.log(e.data);
  };

  ws.onmessage = (e) => {
    console.log(e.data);
  };
}
