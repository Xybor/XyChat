import axios from "axios";
import { configs } from "./config";

const axiosInstance = axios.create({
  baseURL: configs.apiUrl,
  timeout: 2000,
  withCredentials: true,
  credentials: "include",
});

export { axiosInstance };
