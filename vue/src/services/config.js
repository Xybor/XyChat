const currentHost = window.location.host;

export const configs = {
  apiUrl: `http://${currentHost}:1999/api/v1`,
  wsUrl: `ws://${currentHost}:1999/ws/v1`,
};
