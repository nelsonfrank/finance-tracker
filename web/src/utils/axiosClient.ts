/* eslint-disable @typescript-eslint/no-explicit-any */
import { env } from "@/env/client";
import axios from "axios";
import Cookies from "js-cookie";
import { getCookie } from "./cookiesUtils";

let isRefreshing = false;
let failedQueue: any[] = [];

const processQueue = (error: any, token: string | null = null) => {
  failedQueue.forEach((prom) => {
    if (token) {
      prom.resolve(token);
    } else {
      prom.reject(error);
    }
  });
  failedQueue = [];
};

export const Axios = axios.create({
  baseURL: env.NEXT_PUBLIC_API_BASE_URL,
  timeout: 10000, // Request timeout in milliseconds
  headers: {
    "Content-Type": "application/json",
    Accept: "application/json",
  },
});

Axios.interceptors.request.use(
  async (config) => {
    const token = await getCookie("access_token");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    // Handle request errors
    return Promise.reject(error);
  }
);
Axios.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    if (error.response?.status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        // Queue the failed request until the refresh process completes
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        }).then((token) => {
          originalRequest.headers.Authorization = `Bearer ${token}`;
          return Axios(originalRequest);
        });
      }

      originalRequest._retry = true;
      isRefreshing = true;

      try {
        // Call refresh-token endpoint
        const refreshResponse = await Axios.post(
          "/v1/auth/refresh-token",
          {},
          { withCredentials: true } // Send cookies to backend
        );

        const newAccessToken = refreshResponse.data.access_token;
        // Update token in cookies
        Cookies.set("access_token", newAccessToken);

        // Retry failed requests
        processQueue(null, newAccessToken);

        // Retry the original request
        originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;
        return Axios(originalRequest);
      } catch (refreshError) {
        console.log({ refreshError });
        // Handle token refresh failure
        processQueue(refreshError, null);

        if (typeof window !== "undefined") {
          window.location.href = "/auth/login"; // Redirect to login
        }

        return Promise.reject(refreshError);
      } finally {
        isRefreshing = false;
      }
    }

    return Promise.reject(error);
  }
);
