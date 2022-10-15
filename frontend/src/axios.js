import axios from "axios";

export const backendAxios = axios.create(
  process.env.NODE_ENV === "development"
    ? { baseURL: "http://127.0.0.1:8000" }
    : {}
);
