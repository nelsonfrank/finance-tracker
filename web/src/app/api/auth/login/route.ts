import { cookies } from "next/headers";
import { type NextRequest } from "next/server";
import * as API from "@/data/backend/api";
import redis from "redis";

export async function POST(req: NextRequest) {
  const payload = await req.json();

  const redisClient = redis.createClient();
  const response = await API.loginAPI(payload);
  const data = response.data;

  // store access token in redis
  storeAuthTokenInServer(redisClient, data.access_token);

  //store access token in browser
  storeAuthTokenInBrowser(data);

  return Response.json(data.user);
}

const storeAuthTokenInServer = async (
  redisClient: any,
  accessToken: string
) => {
  try {
    // Connect to Redis
    await redisClient.connect();

    // Store Access Token with expiration
    await redisClient.set("accessToken", accessToken, {
      EX: 15 * 60, // Expiration in 15 mins
    });
    // Store Refresh Token with expiration
    await redisClient.set("accessToken", accessToken, {
      EX: 15 * 60 * 60 * 24, // Expiration in 15 days
    });

    console.log("Access token stored successfully!");
  } catch (err) {
    console.error("Error storing access token in Redis:", err);
  } finally {
    // Close Redis connection
    await redisClient.disconnect();
  }
};

const storeAuthTokenInBrowser = async (data: API.loginResponse) => {
  const cookieStore = await cookies();

  const expirationDate = new Date();
  const refreshTokenExpirationDate = new Date();

  // Add 5 seconds to `expirationDate`
  expirationDate.setDate(refreshTokenExpirationDate.getDate() + 1);

  // Add 7 days to `refreshTokenExpirationDate`
  refreshTokenExpirationDate.setDate(refreshTokenExpirationDate.getDate() + 7);
  cookieStore.set("access_token", data.access_token, {
    path: "/",
    expires: expirationDate,
    httpOnly: false,
    secure: false, // Set to true in production
    sameSite: "lax",
  });

  cookieStore.set("refresh_token", data.refresh_token, {
    path: "/",
    expires: refreshTokenExpirationDate,
    httpOnly: true,
    secure: true, // Set to true in production
    sameSite: "lax",
  });
};
