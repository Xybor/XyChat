var currentHost = window.location.host;
var schema = window.location.protocol;

if (process.env.NODE_ENV == "development") {
  currentHost = `${window.location.hostname}:1999`;
}

export const configs = {
  apiUrl: `${schema}//${currentHost}/api/v1`,
  wsUrl: `ws://${currentHost}/ws/v1`,
};
