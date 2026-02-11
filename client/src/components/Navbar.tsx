'use client';

/**
 * ============================================================================
 * NAVBAR COMPONENT - Hydration-Safe Authentication State
 * ============================================================================
 *
 * Problem:
 * - Server renders: "Not logged in" (no cookies on server)
 * - Client hydrates: "Logged in" (cookies available in browser)
 * - Result: HYDRATION MISMATCH ERROR
 *
 * Solution:
 * - Start with loading state (isLoggedIn = null)
 * - Fetch auth state ONLY on client-side (useEffect)
 * - Server and client both render "Loading..." initially
 * - After hydration, client updates to show correct state
 * - No mismatch because both start with same UI
 *
 * Alternative (Faster):
 * - Read is_logged_in cookie immediately
 * - No API call needed
 * - Trade-off: Might show stale state if cookie is outdated
 * ============================================================================
 */

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { api } from '@/lib/axios-client';

interface User {
  id: number;
  email: string;
  username: string;
  full_name: string;
  role: string;
}

export default function Navbar() {
  // ✅ Start with null (loading state)
  // This prevents hydration mismatch
  const [isLoggedIn, setIsLoggedIn] = useState<boolean | null>(null);
  const [user, setUser] = useState<User | null>(null);

  useEffect(() => {
    // ✅ Fetch auth state ONLY on client-side (after hydration)
    checkAuthStatus();
  }, []);

  /**
   * Check if user is authenticated by calling /users/profile
   * If successful → User is logged in
   * If 401 → User is not logged in
   */
  async function checkAuthStatus() {
    try {
      const { data } = await api.get('/api/v1/users/profile');
      setUser(data.data);
      setIsLoggedIn(true);
      console.log('[NAVBAR] User authenticated:', data.data.email);
    } catch (error: any) {
      console.log('[NAVBAR] User not authenticated:', error.response?.status);
      setIsLoggedIn(false);
      setUser(null);
    }
  }

  /**
   * Logout handler
   * Calls backend logout endpoint → Clears cookies → Redirects
   */
  async function handleLogout() {
    try {
      await api.post('/api/v1/auth/logout');
      console.log('[NAVBAR] Logout successful');
      setIsLoggedIn(false);
      setUser(null);
      window.location.href = '/';
    } catch (error) {
      console.error('[NAVBAR] Logout failed:', error);
      // Force redirect even if API fails
      window.location.href = '/';
    }
  }

  // ✅ Show loading state until auth status is determined
  // This ensures server and client render the same initial UI
  if (isLoggedIn === null) {
    return (
      <nav className="navbar">
        <div className="navbar-container">
          <Link href="/" className="navbar-logo">
            E-Commerce
          </Link>
          <div className="navbar-menu">
            <Link href="/products">Products</Link>
            <span className="navbar-loading">Loading...</span>
          </div>
        </div>
      </nav>
    );
  }

  // ✅ Now safe to render based on auth state (no hydration mismatch)
  return (
    <nav className="navbar">
      <div className="navbar-container">
        <Link href="/" className="navbar-logo">
          E-Commerce
        </Link>

        <div className="navbar-menu">
          <Link href="/products">Products</Link>

          {isLoggedIn ? (
            // Logged in menu
            <>
              <Link href="/cart">Cart</Link>
              <Link href="/orders">Orders</Link>
              <Link href="/profile" className="navbar-user">
                {user?.email}
              </Link>
              <button onClick={handleLogout} className="navbar-logout-btn">
                Logout
              </button>
            </>
          ) : (
            // Logged out menu
            <>
              <Link href="/login">Login</Link>
              <Link href="/register">Register</Link>
            </>
          )}
        </div>
      </div>
    </nav>
  );
}

/**
 * ============================================================================
 * ALTERNATIVE: Cookie-Based Check (Faster, No API Call)
 * ============================================================================
 *
 * Use this if you want faster initial render (no loading state)
 * Trade-off: Might show wrong state if is_logged_in cookie is stale
 *
 * ```typescript
 * import Cookies from 'js-cookie';
 *
 * export default function NavbarFast() {
 *   const [isLoggedIn, setIsLoggedIn] = useState(() => {
 *     return Cookies.get('is_logged_in') === 'true';
 *   });
 *
 *   // No hydration mismatch because:
 *   // - Server: Cookies not available → defaults to false
 *   // - Client: Cookies available → reads immediately
 *   // - Both render "Login" initially, then client updates to "Profile"
 *
 *   return (
 *     <nav>
 *       {isLoggedIn ? (
 *         <a href="/profile">Profile</a>
 *       ) : (
 *         <a href="/login">Login</a>
 *       )}
 *     </nav>
 *   );
 * }
 * ```
 *
 * Installation: npm install js-cookie @types/js-cookie
 * ============================================================================
 */

/**
 * ============================================================================
 * USAGE IN LAYOUT
 * ============================================================================
 *
 * // app/layout.tsx
 *
 * import Navbar from '@/components/Navbar';
 *
 * export default function RootLayout({ children }) {
 *   return (
 *     <html lang="en">
 *       <body>
 *         <Navbar />
 *         {children}
 *       </body>
 *     </html>
 *   );
 * }
 *
 * ============================================================================
 */

/**
 * ============================================================================
 * STYLING (Optional)
 * ============================================================================
 *
 * Add to your globals.css or create navbar.module.css:
 *
 * .navbar {
 *   background-color: #fff;
 *   box-shadow: 0 2px 4px rgba(0,0,0,0.1);
 *   padding: 1rem 0;
 * }
 *
 * .navbar-container {
 *   max-width: 1200px;
 *   margin: 0 auto;
 *   display: flex;
 *   justify-content: space-between;
 *   align-items: center;
 *   padding: 0 1rem;
 * }
 *
 * .navbar-logo {
 *   font-size: 1.5rem;
 *   font-weight: bold;
 *   color: #333;
 *   text-decoration: none;
 * }
 *
 * .navbar-menu {
 *   display: flex;
 *   gap: 1.5rem;
 *   align-items: center;
 * }
 *
 * .navbar-menu a {
 *   color: #666;
 *   text-decoration: none;
 *   transition: color 0.2s;
 * }
 *
 * .navbar-menu a:hover {
 *   color: #333;
 * }
 *
 * .navbar-user {
 *   font-weight: 500;
 *   color: #333;
 * }
 *
 * .navbar-logout-btn {
 *   background: none;
 *   border: 1px solid #ddd;
 *   padding: 0.5rem 1rem;
 *   border-radius: 4px;
 *   cursor: pointer;
 *   transition: all 0.2s;
 * }
 *
 * .navbar-logout-btn:hover {
 *   background-color: #f5f5f5;
 * }
 *
 * .navbar-loading {
 *   color: #999;
 * }
 *
 * ============================================================================
 */
