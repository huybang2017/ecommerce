// Route configuration for Next.js middleware
// Defines public routes, protected routes, and admin routes

export const PUBLIC_ROUTES = [
  "/",
  "/login",
  "/register",
  "/products",
  "/search",
  "/api/public",
];

export const AUTH_ROUTES = ["/login", "/register"];

export const PROTECTED_ROUTES = ["/cart", "/checkout", "/orders", "/profile"];

export const ADMIN_ROUTES = ["/admin", "/inventory", "/dashboard"];

// Helper function to check if path matches a pattern
export function isPathMatch(path: string, patterns: string[]): boolean {
  return patterns.some((pattern) => {
    if (pattern === path) return true;
    if (pattern.endsWith("*")) {
      const base = pattern.slice(0, -1);
      return path.startsWith(base);
    }
    return path.startsWith(pattern);
  });
}
