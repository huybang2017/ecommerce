/**
 * ============================================================================
 * NEXT.JS AUTH MIDDLEWARE - Edge Runtime Compatible
 * ============================================================================
 *
 * Purpose: Protect authenticated routes from unauthenticated access.
 *
 * Edge Runtime Limitations:
 * - Cannot decrypt httpOnly cookies (they're encrypted by browser)
 * - Cannot validate JWT tokens (no crypto libraries)
 * - Cannot access database
 *
 * Solution:
 * 1. Backend sets a secondary cookie "is_logged_in" (NOT httpOnly)
 * 2. Middleware reads this flag for preliminary check (fast, edge-compatible)
 * 3. API Gateway does real authentication (JWT validation)
 *
 * Flow:
 * - User logs in → Backend sets BOTH:
 *   - buyer_access_token (httpOnly, secure) ← Real token
 *   - is_logged_in=true (plain cookie) ← Middleware flag
 *
 * - User visits /profile:
 *   - Middleware checks is_logged_in → true → Allow
 *   - Page fetches data → API validates buyer_access_token → 200 OK
 *
 * - If JWT expires:
 *   - Middleware checks is_logged_in → true → Allow (⚠️ still valid)
 *   - Page fetches data → API returns 401
 *   - Client refreshes token OR redirects to login
 *
 * Security Note:
 * - is_logged_in is NOT secure (readable by JavaScript)
 * - But it only contains a boolean flag, not sensitive data
 * - Real authentication still uses httpOnly cookies
 * - This is the industry-standard approach for edge middleware
 * ============================================================================
 */

import { NextRequest, NextResponse } from "next/server";

// ────────────────────────────────────────────────────────────────────────────
// CONFIGURATION
// ────────────────────────────────────────────────────────────────────────────

const PUBLIC_ROUTES = [
  "/",
  "/login",
  "/register",
  "/about",
  "/contact",
  "/products", // Public product listing
  "/product", // Product detail pages (dynamic)
];

const PROTECTED_ROUTES = [
  "/profile",
  "/cart",
  "/checkout",
  "/orders",
  "/wishlist",
  "/settings",
];

const AUTH_ROUTES = ["/login", "/register"];

// ────────────────────────────────────────────────────────────────────────────
// MIDDLEWARE LOGIC
// ────────────────────────────────────────────────────────────────────────────

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;

  console.log("[MIDDLEWARE] Request:", pathname);

  // ── 1. Read auth state from plain cookie ──────────────────────────────────
  // This cookie is set by backend after login (NOT httpOnly)
  // We can read it because it's a plain cookie (not encrypted)
  const isLoggedIn = request.cookies.get("is_logged_in")?.value === "true";

  console.log("[MIDDLEWARE] Is logged in:", isLoggedIn);

  // ── 2. Check if route is protected ────────────────────────────────────────
  const isProtectedRoute = PROTECTED_ROUTES.some((route) =>
    pathname.startsWith(route)
  );

  const isAuthRoute = AUTH_ROUTES.some((route) => pathname.startsWith(route));

  const isPublicRoute = PUBLIC_ROUTES.some((route) =>
    pathname === route || pathname.startsWith(route)
  );

  // ── 3. Redirect logic ─────────────────────────────────────────────────────

  // CASE 1: User is logged in but trying to access login/register
  // → Redirect to home
  if (isLoggedIn && isAuthRoute) {
    console.log("[MIDDLEWARE] Already logged in, redirecting from auth page");
    return NextResponse.redirect(new URL("/", request.url));
  }

  // CASE 2: User is NOT logged in but trying to access protected route
  // → Redirect to login with returnUrl
  if (!isLoggedIn && isProtectedRoute) {
    console.log("[MIDDLEWARE] Not logged in, redirecting to login");
    const loginUrl = new URL("/login", request.url);
    loginUrl.searchParams.set("returnUrl", pathname); // Save intended destination
    return NextResponse.redirect(loginUrl);
  }

  // ── 4. Allow request ───────────────────────────────────────────────────────
  console.log("[MIDDLEWARE] Request allowed");

  // Optional: Add custom headers for debugging
  const response = NextResponse.next();
  response.headers.set("x-middleware-executed", "true");
  response.headers.set("x-is-logged-in", isLoggedIn.toString());

  return response;
}

// ────────────────────────────────────────────────────────────────────────────
// MATCHER CONFIGURATION
// ────────────────────────────────────────────────────────────────────────────

/**
 * Specify which routes trigger this middleware
 * - Include all protected routes
 * - Include auth routes (for logged-in redirect)
 * - Exclude static files, API routes, _next internals
 */
export const config = {
  matcher: [
    /*
     * Match all request paths except:
     * - _next/static (static files)
     * - _next/image (image optimization)
     * - favicon.ico (favicon file)
     * - public files (images, fonts, etc.)
     * - api routes (handled separately)
     */
    "/((?!_next/static|_next/image|favicon.ico|api|.*\\.(?:svg|png|jpg|jpeg|gif|webp|ico|woff|woff2|ttf|eot)$).*)",
  ],
};

/**
 * ============================================================================
 * BACKEND REQUIREMENT
 * ============================================================================
 *
 * Backend MUST set is_logged_in cookie after login!
 *
 * File: identity-service/internal/handler/cookie_helper.go
 *
 * func setAuthCookies(c *gin.Context, sessionID, accessToken, refreshToken string) {
 *   prefix := cookiePrefix(c)
 *
 *   // Set httpOnly cookies (secure)
 *   c.SetCookie(prefix+"session_id", sessionID, 604800, "/", "", false, true)
 *   c.SetCookie(prefix+"access_token", accessToken, 900, "/", "", false, true)
 *   c.SetCookie(prefix+"refresh_token", refreshToken, 604800, "/", "", false, true)
 *
 *   // ✅ NEW: Set plain cookie for Next.js middleware
 *   c.SetCookie("is_logged_in", "true", 604800, "/", "", false, false)
 *   //                                              ↑      ↑
 *   //                                          secure  httpOnly=false
 * }
 *
 * func clearAuthCookies(c *gin.Context) {
 *   prefix := cookiePrefix(c)
 *
 *   c.SetCookie(prefix+"session_id", "", -1, "/", "", false, true)
 *   c.SetCookie(prefix+"access_token", "", -1, "/", "", false, true)
 *   c.SetCookie(prefix+"refresh_token", "", -1, "/", "", false, true)
 *
 *   // ✅ NEW: Clear is_logged_in flag
 *   c.SetCookie("is_logged_in", "", -1, "/", "", false, false)
 * }
 *
 * ============================================================================
 */
