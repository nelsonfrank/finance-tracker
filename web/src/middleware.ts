import { getToken } from "next-auth/jwt";
import { NextRequest, NextResponse } from "next/server";

export async function middleware(request: NextRequest) {
    const token = await getToken({
      req: request,
      secret: process.env.NEXTAUTH_SECRET
    })
  
    const isLoginPage = request.nextUrl.pathname.startsWith('/auth')
  
    if (isLoginPage && token) {
      return NextResponse.redirect(new URL('/dashboard', request.url))
    }
  
    const isProtectedRoute = request.nextUrl.pathname.startsWith('/dashboard')

    if (isProtectedRoute && !token) {
      const loginUrl = new URL('/auth/login', request.url)

      loginUrl.searchParams.set('callbackUrl', request.nextUrl.pathname)

      return NextResponse.redirect(loginUrl)
    }
    return NextResponse.next();
}

export const config = {
  matcher: ["/dashboard/:path*", "/auth/:path*"],
};
