import Cookies from "js-cookie";

export async function getCookie(key: string) {
  let token;
  if (typeof window !== "undefined") {
    token = Cookies.get(key);
  } else {
    // Access token from cookies (server-side)
    // eslint-disable-next-line @typescript-eslint/no-require-imports
    const { cookies } = require("next/headers");
    const serverCookies = await cookies();
    token = serverCookies.get(key)?.value;
  }
  return token;
}
