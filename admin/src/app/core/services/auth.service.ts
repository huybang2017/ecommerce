import { Injectable } from "@angular/core";
import { HttpClient } from "@angular/common/http";
import { Router } from "@angular/router";
import { Observable, BehaviorSubject, throwError, of } from "rxjs";
import { tap, map, catchError } from "rxjs/operators";
import { environment as env } from "../../environments/environment";
import { LoginRequest, LoginResponse } from "../../shared/models/auth.model";
import { AuthStateService } from "./auth-state.service";

/**
 * ============================================================================
 * AUTH SERVICE - Production-Grade Implementation
 * ============================================================================
 *
 * Features:
 * 1. Reactive auth state management (BehaviorSubject)
 * 2. Token refresh
 * 3. Session validation
 * 4. Complete state cleanup on logout
 * 5. User profile caching
 *
 * Design Pattern: Singleton (providedIn: 'root') + Observer (BehaviorSubject)
 *
 * Why BehaviorSubject?
 * - Always has a current value (last emitted)
 * - New subscribers immediately receive the current value
 * - Perfect for reactive state management
 * ============================================================================
 */

interface User {
  id: number;
  email: string;
  username: string;
  full_name: string;
  role: string;
}

@Injectable({ providedIn: "root" })
export class AuthService {
  // ✅ Use BehaviorSubject for reactive state management
  // Components can subscribe to currentUser$ and automatically update
  private currentUserSubject = new BehaviorSubject<User | null>(null);
  public currentUser$ = this.currentUserSubject.asObservable();

  constructor(
    private http: HttpClient,
    private authState: AuthStateService,
    private router: Router
  ) {
    // Initialize auth state on service creation
    this.initializeAuthState();
  }

  /**
   * Check if user is authenticated on app initialization
   * This runs once when the service is created
   */
  private initializeAuthState(): void {
    this.checkSession().subscribe({
      next: (isValid) => {
        if (isValid) {
          console.log('[AUTH SERVICE] User authenticated on init');
        } else {
          console.log('[AUTH SERVICE] User not authenticated on init');
        }
      },
      error: () => {
        console.log('[AUTH SERVICE] Session check failed on init');
        this.currentUserSubject.next(null);
        this.authState.setAuthenticated(false);
      },
    });
  }

  /**
   * Login with email and password
   * Sets httpOnly cookies (admin_access_token, admin_session_id, admin_refresh_token)
   */
  login(payload: LoginRequest): Observable<LoginResponse> {
    return this.http
      .post<LoginResponse>(`${env.api}/auth/login`, payload, {
        withCredentials: true,
        headers: { "X-App-Role": "admin" },
      })
      .pipe(
        tap((response) => {
          // Update state
          if (response.user) {
            this.currentUserSubject.next(response.user as any);
          }
          this.authState.setAuthenticated(true);
          console.log("[AUTH SERVICE] Login successful:", response.user?.email);
        }),
        catchError((error) => {
          console.error("[AUTH SERVICE] Login failed:", error);
          this.authState.setAuthenticated(false);
          this.currentUserSubject.next(null);
          return throwError(() => error);
        })
      );
  }

  /**
   * Check if current session is valid
   * Returns true if backend responds with 200 OK
   * Used by auth guard to validate session before allowing navigation
   */
  checkSession(): Observable<boolean> {
    return this.http
      .get<{ data: User }>(`${env.api}/users/profile`, {
        withCredentials: true,
        headers: { "X-App-Role": "admin" },
      })
      .pipe(
        map((response) => {
          // Update current user
          this.currentUserSubject.next(response.data);
          this.authState.setAuthenticated(true);
          console.log("[AUTH SERVICE] Session valid:", response.data.email);
          return true;
        }),
        catchError((error) => {
          console.log("[AUTH SERVICE] Session invalid:", error.status);
          this.currentUserSubject.next(null);
          this.authState.setAuthenticated(false);
          return of(false);
        })
      );
  }

  /**
   * Get current user profile
   * Makes API call to fetch fresh data
   * Used by role guards to check user role
   */
  getCurrentUser(): Observable<User | null> {
    // If we already have user in memory, return it
    const currentUser = this.currentUserSubject.value;
    if (currentUser) {
      console.log("[AUTH SERVICE] Returning cached user:", currentUser.email);
      return of(currentUser);
    }

    // Otherwise, fetch from server
    console.log("[AUTH SERVICE] Fetching user from server");
    return this.http
      .get<{ data: User }>(`${env.api}/users/profile`, {
        withCredentials: true,
        headers: { "X-App-Role": "admin" },
      })
      .pipe(
        map((response) => {
          this.currentUserSubject.next(response.data);
          return response.data;
        }),
        catchError((error) => {
          console.error("[AUTH SERVICE] Failed to fetch user:", error);
          this.currentUserSubject.next(null);
          return of(null);
        })
      );
  }

  /**
   * Refresh access token using session_id or refresh_token cookie
   * Called automatically by interceptor when access token expires
   */
  refreshToken(): Observable<void> {
    console.log("[AUTH SERVICE] Refreshing access token");
    return this.http
      .post<void>(
        `${env.api}/auth/refresh`,
        {},
        {
          withCredentials: true,
          headers: { "X-App-Role": "admin" },
        }
      )
      .pipe(
        tap(() => {
          console.log("[AUTH SERVICE] ✅ Token refreshed successfully");
        }),
        catchError((error) => {
          console.error("[AUTH SERVICE] ❌ Token refresh failed:", error);
          return throwError(() => error);
        })
      );
  }

  /**
   * Logout - Complete State Cleanup
   *
   * Steps:
   * 1. Call backend logout endpoint (clears cookies)
   * 2. Clear client-side state (BehaviorSubject, AuthStateService)
   * 3. Clear localStorage/sessionStorage
   * 4. Navigate to signin page
   *
   * CRITICAL: Even if logout API fails, still clear client state
   */
  logout(): Observable<void> {
    console.log("[AUTH SERVICE] Logging out");

    return this.http
      .post<void>(
        `${env.api}/auth/logout`,
        {},
        {
          withCredentials: true,
          headers: { "X-App-Role": "admin" },
        }
      )
      .pipe(
        tap(() => {
          console.log("[AUTH SERVICE] ✅ Logout successful");
          this.clearClientState();
        }),
        catchError((error) => {
          console.error("[AUTH SERVICE] ❌ Logout failed:", error);
          // ✅ CRITICAL: Clear client state even if API fails
          this.clearClientState();
          return throwError(() => error);
        })
      );
  }

  /**
   * Clear all client-side state
   * Called after logout (success or failure)
   */
  private clearClientState(): void {
    console.log("[AUTH SERVICE] Clearing client state");

    // ✅ 1. Clear reactive state
    this.currentUserSubject.next(null);
    this.authState.setAuthenticated(false);

    // ✅ 2. Clear localStorage (non-sensitive data only)
    // DON'T store tokens in localStorage!
    localStorage.removeItem("user_preferences");
    localStorage.removeItem("theme");
    localStorage.removeItem("language");

    // ✅ 3. Clear sessionStorage
    sessionStorage.clear();

    // ✅ 4. Navigate to signin
    this.router.navigate(["/signin"], {
      queryParams: { reason: "logged_out" },
    });
  }

  /**
   * Check if user is authenticated (synchronous)
   * Returns current auth state from BehaviorSubject
   */
  isAuthenticated(): boolean {
    return this.currentUserSubject.value !== null;
  }

  /**
   * Get current user (synchronous)
   * Returns current user from BehaviorSubject
   */
  getCurrentUserSync(): User | null {
    return this.currentUserSubject.value;
  }
}
