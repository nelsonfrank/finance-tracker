import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

const PROTECTED_ROUTES = ["/dashboard", "/profile", "/settings"];

export function middleware(request: NextRequest) {
  const accessToken = request.cookies.get("access_token")?.value;

  if (
    PROTECTED_ROUTES.some((route) => request.nextUrl.pathname.startsWith(route))
  ) {
    // Validate the token
    if (!accessToken) {
      // Redirect to login if token is missing or invalid
      return NextResponse.redirect(new URL("/auth/login", request.url));
    }
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/((?!.*\\..*|_next).*)", "/", "/(api|trpc)(.*)"],
};
