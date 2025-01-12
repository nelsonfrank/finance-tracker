import { cookies } from "next/headers";
import { type NextRequest } from "next/server";
import * as API from "@/data/backend/api";

export async function POST(req: NextRequest) {
  const payload = await req.json();

  const cookieStore = await cookies();

  const response = await API.loginAPI(payload);
  const data = response.data;

  const expirationDate = new Date();
  expirationDate.setMinutes(expirationDate.getMinutes() + 5);

  cookieStore.set("access_token", data.access_token, {
    path: "/",
    expires: expirationDate,
    httpOnly: false,
    secure: false, // Set to true in production
    sameSite: "lax",
  });

  cookieStore.set("refresh_token", data.refresh_token, {
    path: "/",
    expires: expirationDate,
    httpOnly: true,
    secure: true, // Set to true in production
    sameSite: "lax",
  });
  return Response.json(data.user);
}
