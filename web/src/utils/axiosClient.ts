/* eslint-disable @typescript-eslint/no-explicit-any */
import { env } from "@/env/client";
import axios from "axios";
import { getSession } from "next-auth/react";

const MAX_RETRIES = 3;


export const Axios = axios.create({
  baseURL: env.NEXT_PUBLIC_API_BASE_URL,
  timeout: 1000, 
  headers: {
    "Content-Type": "application/json",
    Accept: "application/json",
  },
});

Axios.interceptors.request.use(
  async (config) => {
    const session = await getSession();
    const token = session && session.user ? session.user.access_token : "";

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
    const config = error.config;

    if (!config._retryCount) config._retryCount = 0;

    const status = error.response?.status;

    const shouldRetry = config._retryCount < MAX_RETRIES && (!error.response || status === 401);

    if (shouldRetry) {
      config._retryCount += 1;

      if (status === 401) {
        const newSession = await getSession();

        const newToken = newSession?.user?.access_token;

        if (newToken) {
          config.headers.Authorization = `Bearer ${newToken}`;
        } else {          
          if (typeof window !== 'undefined') {
            window.location.href = '/auth/login';
          }
        }
      }

      const delay = (ms: number) => new Promise((res) => setTimeout(res, ms));
      await delay(config._retryCount * 1000);

      return Axios(config);
    }

    const shouldRedirect = config._retryCount > MAX_RETRIES &&  status === 401;
    if(shouldRedirect){
      if (typeof window !== 'undefined') {
        window.location.href = '/auth/login';
      }
    }

    return Promise.reject(error);
  }
);
