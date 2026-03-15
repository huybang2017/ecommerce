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

@Injectable()
export class AuthInterceptor implements HttpInterceptor {
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

  private isAuthEndpoint(url?: string): boolean {
    if (!url) return false;
    return (
      url.includes("/auth/refresh") ||
      url.includes("/auth/login") ||
      url.includes("/auth/register") ||
      url.includes("/auth/logout")
    );
  }

  private handle401Error(
    req: HttpRequest<any>,
    next: HttpHandler,
  ): Observable<HttpEvent<any>> {
    const http = this.injector.get(HttpClient);
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
              }),
            );
          }
          // Refresh failed → Don't retry, throw error
          console.log("[AUTH] Refresh failed, aborting request");
          return throwError(() => new Error("Token refresh failed"));
        }),
      );
    }
    console.log("[AUTH] Access token expired, initiating refresh...");
    this.isRefreshing = true;
    this.refreshSubject.next(null);

    return http
      .post(
        `${env.api}/auth/refresh`,
        {},
        { withCredentials: true, headers: { "X-App-Role": "admin" } },
      )
      .pipe(
        switchMap(() => {
          console.log("[AUTH] ✅ Token refreshed successfully");
          this.refreshSubject.next(true);
          return next.handle(
            req.clone({
              withCredentials: true,
              setHeaders: { "X-App-Role": "admin" },
            }),
          );
        }),
        catchError((refreshError) => {
          console.error("[AUTH] ❌ Token refresh failed:", refreshError);
          this.refreshSubject.next(false);

          http
            .post(
              `${env.api}/auth/logout`,
              {},
              { withCredentials: true, headers: { "X-App-Role": "admin" } },
            )
            .subscribe({
              next: () => console.log("[AUTH] Server session cleared"),
              error: (err) => console.error("[AUTH] Logout error:", err),
              complete: () => {
                this.router.navigate(["/signin"], {
                  queryParams: { reason: "session_expired" },
                });
              },
            });

          return throwError(() => refreshError);
        }),
        finalize(() => {
          this.isRefreshing = false;
          console.log("[AUTH] Refresh flow completed, reset flag");
        }),
      );
  }
}
