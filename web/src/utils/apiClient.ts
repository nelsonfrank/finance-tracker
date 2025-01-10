import { env } from "@/env/server";
import { getCookie } from "./cookiesUtils";
// import { getCookie } from "cookies-next";

export const BASE_URL = "http://localhost:8080";
export async function fetcher<T>(
  url: string,
  options?: RequestInit
): Promise<T> {
  const accessToken = getCookie("access_token");

  const requestUrl = formatRequestUrl(url);

  const response = await fetch(`${requestUrl}`, {
    headers: {
      Authorization: `Bearer ${accessToken}`,
    },
    ...options,
  });

  if (!response.ok) {
    throw new Error(`Error: ${response.status} ${response.statusText}`);
  }

  return response.json();
}

function formatRequestUrl(urlPath: string) {
  let formattedUrl: string;
  if (typeof window !== "undefined") {
    formattedUrl = urlPath;
  } else {
    formattedUrl = env.APP_BASE_URL + urlPath;
  }

  return formattedUrl;
}
