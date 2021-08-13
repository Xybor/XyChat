var currentHost = window.location.host;
var apiSchema = "https:";
var wsSchema = "ws:";

if (process.env.NODE_ENV == "development") {
  currentHost = `${window.location.hostname}:1999`;
  apiSchema = "http:";
}

export const configs = {
  apiUrl: `${apiSchema}//${currentHost}/api/v1`,
  wsUrl: `${wsSchema}://${currentHost}/ws/v1`,
};

console.log(configs);
