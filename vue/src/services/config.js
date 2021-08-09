let currentHost = window.location.host;
if (process.env.NODE_ENV == "development") {
  currentHost = `${window.location.hostname}:1999`;
}

export const configs = {
  apiUrl: `http://${currentHost}/api/v1`,
  wsUrl: `ws://${currentHost}/ws/v1`,
};
