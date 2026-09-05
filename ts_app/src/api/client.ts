import axios from "axios";
import { getToken } from "./authStorage";

export const api = axios.create({
  baseURL: "/api",
  headers: {
    "Content-Type": "application/json",
  },
});

api.interceptors.request.use((config) => {
  const token = getToken();

  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }

  return config;
});

export const healthApi = axios.create({
  baseURL: "/",
  headers: {
    "Content-Type": "application/json",
  },
});
