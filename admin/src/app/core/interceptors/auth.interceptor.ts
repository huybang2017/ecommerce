import { Injectable, Injector } from "@angular/core";
import {
  HttpInterceptor,
  HttpRequest,
  HttpHandler,
  HttpEvent,
  HttpErrorResponse,
  HttpClient,
} from "@angular/common/http";
import { Observable, BehaviorSubject, throwError } from "rxjs";
import { catchError, switchMap, filter, take, finalize } from "rxjs/operators";
import { environment as env } from "../../environments/environment";
import { Router } from "@angular/router";

/**
 * ============================================================================
 * AUTH INTERCEPTOR - Production-Grade Implementation
 * ============================================================================
 *
 * Features:
 * 1. Automatic X-App-Role header injection (admin)
 * 2. Smart 401 handling with token refresh
 * 3. Race condition prevention using BehaviorSubject
 * 4. Request retry after successful refresh
 * 5. Automatic logout on refresh failure
 *
 * Design Pattern: Singleton + Observer (RxJS)
 *
 * Race Condition Handling:
 * - Multiple requests hit 401 simultaneously
 * - Only ONE refresh token request is made (isRefreshing flag)
 * - Other requests subscribe to refreshSubject and WAIT
 * - When refresh succeeds, all waiting requests retry
 * - BehaviorSubject ensures late subscribers get the last emitted value
 * ============================================================================
 */
@Injectable()
export class AuthInterceptor implements HttpInterceptor {
  // ✅ Use BehaviorSubject instead of Subject
  // BehaviorSubject replays the last value to new subscribers
  // This prevents race conditions where late subscribers miss the emission
  private isRefreshing = false;
  private refreshSubject = new BehaviorSubject<boolean | null>(null);

  constructor(
    private injector: Injector,
    private router: Router,
  ) {}

  intercept(
    req: HttpRequest<any>,
    next: HttpHandler,
  ): Observable<HttpEvent<any>> {
    // Clone request: send cookies + identify this app as "admin"
    const authReq = req.clone({
      withCredentials: true,
      setHeaders: { "X-App-Role": "admin" },
    });

    return next.handle(authReq).pipe(
      catchError((error: HttpErrorResponse) => {
        if (error.status === 401 && !this.isAuthEndpoint(authReq.url)) {
          return this.handle401Error(authReq, next);
        }
        return throwError(() => error);
      }),
    );
  }

  /**
   * Check if URL is an auth endpoint (login, register, refresh, logout)
   * These endpoints should NOT trigger refresh token logic
   */
  private isAuthEndpoint(url?: string): boolean {
    if (!url) return false;
    return (
      url.includes("/auth/refresh") ||
      url.includes("/auth/login") ||
      url.includes("/auth/register") ||
      url.includes("/auth/logout")
    );
  }

  /**
   * ========================================================================
   * HANDLE 401 ERROR - Smart Retry with Race Condition Prevention
   * ========================================================================
   *
   * Scenario:
   * - Request A, B, C all hit 401 at the same time (token expired)
   * - WITHOUT this logic: 3 separate refresh calls → Race condition
   * - WITH this logic: Only 1 refresh call, others wait and retry
   *
   * Implementation:
   * 1. First request sets isRefreshing = true
   * 2. Other requests subscribe to refreshSubject and WAIT
   * 3. When refresh succeeds, refreshSubject emits true
   * 4. All waiting requests wake up and retry
   * 5. If refresh fails, all requests get rejected + logout
   *
   * BehaviorSubject vs Subject:
   * - BehaviorSubject: Replays last value to late subscribers ✅
   * - Subject: Late subscribers miss the emission ❌
   * ========================================================================
   */
  private handle401Error(
    req: HttpRequest<any>,
    next: HttpHandler,
  ): Observable<HttpEvent<any>> {
    const http = this.injector.get(HttpClient);

    // ────────────────────────────────────────────────────────────────────
    // CASE 1: Refresh already in progress → WAIT
    // ────────────────────────────────────────────────────────────────────
    if (this.isRefreshing) {
      console.log("[AUTH] Refresh in progress, queueing request...");

      return this.refreshSubject.pipe(
        filter((success) => success !== null), // ✅ Wait for non-null value
        take(1), // Only take first emission
        switchMap((success) => {
          if (success) {
            console.log("[AUTH] Refresh succeeded, retrying request");
            return next.handle(
              req.clone({
                withCredentials: true,
                setHeaders: { "X-App-Role": "admin" },
              })
            );
          }
          // Refresh failed → Don't retry, throw error
          console.log("[AUTH] Refresh failed, aborting request");
          return throwError(() => new Error("Token refresh failed"));
        })
      );
    }

    // ────────────────────────────────────────────────────────────────────
    // CASE 2: First 401 → Initiate refresh
    // ────────────────────────────────────────────────────────────────────
    console.log("[AUTH] Access token expired, initiating refresh...");
    this.isRefreshing = true;
    this.refreshSubject.next(null); // ✅ Reset to null (queued requests will wait)

    return http
      .post(
        `${env.api}/auth/refresh`,
        {},
        { withCredentials: true, headers: { "X-App-Role": "admin" } }
      )
      .pipe(
        switchMap(() => {
          console.log("[AUTH] ✅ Token refreshed successfully");
          this.refreshSubject.next(true); // ✅ Notify queued requests: SUCCESS

          // Retry original request with fresh token
          return next.handle(
            req.clone({
              withCredentials: true,
              setHeaders: { "X-App-Role": "admin" },
            })
          );
        }),
        catchError((refreshError) => {
          console.error("[AUTH] ❌ Token refresh failed:", refreshError);
          this.refreshSubject.next(false); // ✅ Notify queued requests: FAILURE

          // ✅ Call logout endpoint to clear server-side session
          http
            .post(
              `${env.api}/auth/logout`,
              {},
              { withCredentials: true, headers: { "X-App-Role": "admin" } }
            )
            .subscribe({
              next: () => console.log("[AUTH] Server session cleared"),
              error: (err) => console.error("[AUTH] Logout error:", err),
              complete: () => {
                // Navigate to signin with reason
                this.router.navigate(["/signin"], {
                  queryParams: { reason: "session_expired" },
                });
              },
            });

          return throwError(() => refreshError);
        }),
        finalize(() => {
          // ✅ CRITICAL: Always reset refresh state when done
          // This ensures next 401 can trigger a new refresh
          this.isRefreshing = false;
          console.log("[AUTH] Refresh flow completed, reset flag");
        })
      );
  }
}
