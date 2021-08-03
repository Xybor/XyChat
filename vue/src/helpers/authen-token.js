export function getJWTAuthenToken() {
  return localStorage.getItem("jwt_token");
}

export function setJWTtAuthenToken(token) {
  return localStorage.setItem("jwt_token", token);
}

export function removeJWTtAuthenToken() {
  return localStorage.removeItem("jwt_token");
}

export function getUsername() {
  return localStorage.getItem("username");
}

export function setUsername(username) {
  return localStorage.setItem("username", username);
}

export function removeUsername() {
  return localStorage.removeItem("username");
}
