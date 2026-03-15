import { Injectable } from "@angular/core";
import { HttpClient } from "@angular/common/http";
import { Router } from "@angular/router";
import { Observable, BehaviorSubject, throwError, of } from "rxjs";
import { tap, map, catchError } from "rxjs/operators";
import { environment as env } from "../../environments/environment";
import { LoginRequest, LoginResponse } from "../../shared/models/auth.model";
import { AuthStateService } from "./auth-state.service";

interface User {
  id: number;
  email: string;
  username: string;
  full_name: string;
  role: string;
}

@Injectable({ providedIn: "root" })
export class AuthService {
  private currentUserSubject = new BehaviorSubject<User | null>(null);
  public currentUser$ = this.currentUserSubject.asObservable();

  constructor(
    private http: HttpClient,
    private authState: AuthStateService,
    private router: Router,
  ) {
    this.initializeAuthState();
  }

  private initializeAuthState(): void {
    this.checkSession().subscribe({
      next: (isValid) => {
        if (isValid) {
          console.log("[AUTH SERVICE] User authenticated on init");
        } else {
          console.log("[AUTH SERVICE] User not authenticated on init");
        }
      },
      error: () => {
        console.log("[AUTH SERVICE] Session check failed on init");
        this.currentUserSubject.next(null);
        this.authState.setAuthenticated(false);
      },
    });
  }

  login(payload: LoginRequest): Observable<LoginResponse> {
    return this.http
      .post<LoginResponse>(`${env.api}/auth/login`, payload, {
        withCredentials: true,
        headers: { "X-App-Role": "admin" },
      })
      .pipe(
        tap((response) => {
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
        }),
      );
  }

  checkSession(): Observable<boolean> {
    return this.http
      .get<{ data: User }>(`${env.api}/users/profile`, {
        withCredentials: true,
        headers: { "X-App-Role": "admin" },
      })
      .pipe(
        map((response) => {
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
        }),
      );
  }

  getCurrentUser(): Observable<User | null> {
    const currentUser = this.currentUserSubject.value;
    if (currentUser) {
      console.log("[AUTH SERVICE] Returning cached user:", currentUser.email);
      return of(currentUser);
    }
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
        }),
      );
  }

  refreshToken(): Observable<void> {
    console.log("[AUTH SERVICE] Refreshing access token");
    return this.http
      .post<void>(
        `${env.api}/auth/refresh`,
        {},
        {
          withCredentials: true,
          headers: { "X-App-Role": "admin" },
        },
      )
      .pipe(
        tap(() => {
          console.log("[AUTH SERVICE] ✅ Token refreshed successfully");
        }),
        catchError((error) => {
          console.error("[AUTH SERVICE] ❌ Token refresh failed:", error);
          return throwError(() => error);
        }),
      );
  }

  logout(): Observable<void> {
    console.log("[AUTH SERVICE] Logging out");
    return this.http
      .post<void>(
        `${env.api}/auth/logout`,
        {},
        {
          withCredentials: true,
          headers: { "X-App-Role": "admin" },
        },
      )
      .pipe(
        tap(() => {
          console.log("[AUTH SERVICE] ✅ Logout successful");
          this.clearClientState();
        }),
        catchError((error) => {
          console.error("[AUTH SERVICE] ❌ Logout failed:", error);
          this.clearClientState();
          return throwError(() => error);
        }),
      );
  }

  private clearClientState(): void {
    console.log("[AUTH SERVICE] Clearing client state");
    this.currentUserSubject.next(null);
    this.authState.setAuthenticated(false);
    localStorage.removeItem("user_preferences");
    localStorage.removeItem("theme");
    localStorage.removeItem("language");
    sessionStorage.clear();
    this.router.navigate(["/signin"], {
      queryParams: { reason: "logged_out" },
    });
  }

  isAuthenticated(): boolean {
    return this.currentUserSubject.value !== null;
  }

  getCurrentUserSync(): User | null {
    return this.currentUserSubject.value;
  }
}
