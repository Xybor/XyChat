const USER_INFO_KEY = "UserInfo";
const LOGIN_STATUS_KEY = "LoginStatus";

export function getAccountInfo() {
  return JSON.parse(localStorage.getItem(USER_INFO_KEY));
}

export function setAccountInfo(id, username) {
  localStorage.setItem(
    USER_INFO_KEY,
    JSON.stringify({ id: id, username: username })
  );
}

export function unsetAccountInfo() {
  localStorage.removeItem(USER_INFO_KEY);
}

export function getLoginStatus() {
  return localStorage.getItem(LOGIN_STATUS_KEY) == "true";
}

export function setLogin() {
  localStorage.setItem(LOGIN_STATUS_KEY, "true");
}

export function setLogout() {
  localStorage.removeItem(LOGIN_STATUS_KEY, "false");
}

export function unsetLoginStatus() {
  localStorage.removeItem(LOGIN_STATUS_KEY);
}
