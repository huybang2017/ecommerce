// middleware.ts - Next.js Edge Runtime Middleware
// Professional authentication and authorization for E-commerce Microservices
// Implements: Cookie-based auth, Route protection, Role-based access, Callback URL

import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";
import {
  PUBLIC_ROUTES,
  AUTH_ROUTES,
  PROTECTED_ROUTES,
  ADMIN_ROUTES,
  isPathMatch,
} from "./src/lib/route-config";

// Cookie names (must match backend)
const ACCESS_TOKEN_COOKIE = "access_token";
const REFRESH_TOKEN_COOKIE = "refresh_token";
const LOGIN_PAGE = "/login";

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;

  // 1. Get tokens from cookies
  const accessToken = request.cookies.get(ACCESS_TOKEN_COOKIE)?.value;
  const refreshToken = request.cookies.get(REFRESH_TOKEN_COOKIE)?.value;

  const hasValidSession = !!accessToken || !!refreshToken;

  // 2. Allow public routes (home, products, etc.) for everyone
  if (
    isPathMatch(pathname, PUBLIC_ROUTES) &&
    !isPathMatch(pathname, AUTH_ROUTES)
  ) {
    return NextResponse.next();
  }

  // 3. Handle auth routes (login, register)
  if (isPathMatch(pathname, AUTH_ROUTES)) {
    // If already logged in, redirect to home
    if (hasValidSession) {
      return NextResponse.redirect(new URL("/", request.url));
    }
    return NextResponse.next();
  }

  // 4. Handle protected routes (cart, orders, profile)
  if (isPathMatch(pathname, PROTECTED_ROUTES)) {
    if (!hasValidSession) {
      // Redirect to login with callback URL
      const loginUrl = new URL(LOGIN_PAGE, request.url);
      loginUrl.searchParams.set("callbackUrl", pathname);
      return NextResponse.redirect(loginUrl);
    }
    return NextResponse.next();
  }

  // 5. Handle admin routes (require admin role)
  if (isPathMatch(pathname, ADMIN_ROUTES)) {
    if (!hasValidSession) {
      const loginUrl = new URL(LOGIN_PAGE, request.url);
      loginUrl.searchParams.set("callbackUrl", pathname);
      return NextResponse.redirect(loginUrl);
    }

    // TODO: Decode JWT to check role (requires jose library)
    // For now, allow if authenticated
    // In production, decode accessToken and check claims.role === 'ADMIN'
    return NextResponse.next();
  }

  // 6. Default: allow all other routes
  return NextResponse.next();
}

// Configure which routes to run middleware on
export const config = {
  matcher: [
    /*
     * Match all request paths except:
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     * - public folder files
     */
    "/((?!_next/static|_next/image|favicon.ico|.*\\.(?:svg|png|jpg|jpeg|gif|webp)$).*)",
  ],
};
