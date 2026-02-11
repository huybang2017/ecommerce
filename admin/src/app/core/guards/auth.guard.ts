import { inject } from '@angular/core';
import { Router, CanActivateFn, ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';
import { AuthService } from '../services/auth.service';
import { map, catchError, of } from 'rxjs';

/**
 * ============================================================================
 * AUTH GUARD - Functional Style (Angular 17+)
 * ============================================================================
 *
 * Protects routes by checking if user has valid session.
 *
 * Usage in routes:
 * ```typescript
 * {
 *   path: 'dashboard',
 *   component: DashboardComponent,
 *   canActivate: [authGuard]
 * }
 * ```
 *
 * How it works:
 * 1. Calls AuthService.checkSession() to verify session
 * 2. If valid → Allow navigation (return true)
 * 3. If invalid → Redirect to /signin with returnUrl (return false)
 *
 * Why not trust client-side flags?
 * - Client-side flags can be manipulated
 * - Server-side validation is the source of truth
 * - checkSession() makes an API call to verify token validity
 * ============================================================================
 */
export const authGuard: CanActivateFn = (
  route: ActivatedRouteSnapshot,
  state: RouterStateSnapshot
) => {
  const authService = inject(AuthService);
  const router = inject(Router);

  console.log('[AUTH GUARD] Checking authentication for:', state.url);

  // Call backend to verify session is still valid
  // This is better than trusting client-side flags
  return authService.checkSession().pipe(
    map((isValid) => {
      if (isValid) {
        console.log('[AUTH GUARD] ✅ Session valid, allowing access');
        return true;
      }

      console.warn('[AUTH GUARD] ❌ Session invalid, redirecting to login');
      router.navigate(['/signin'], {
        queryParams: { returnUrl: state.url }, // Save intended destination
      });
      return false;
    }),
    catchError((error) => {
      console.error('[AUTH GUARD] Session check failed:', error);
      router.navigate(['/signin'], {
        queryParams: { returnUrl: state.url },
      });
      return of(false);
    })
  );
};

/**
 * ============================================================================
 * ADMIN ROLE GUARD - For admin-only routes
 * ============================================================================
 *
 * Usage:
 * ```typescript
 * {
 *   path: 'users',
 *   component: UsersComponent,
 *   canActivate: [authGuard, adminRoleGuard] // ← Both guards
 * }
 * ```
 *
 * Why chain guards?
 * - authGuard: Checks if user is authenticated
 * - adminRoleGuard: Checks if user has admin role
 * - Both must pass for access to be granted
 *
 * NOTE: This guard assumes user is already authenticated (authGuard passed)
 * ============================================================================
 */
export const adminRoleGuard: CanActivateFn = (
  route: ActivatedRouteSnapshot,
  state: RouterStateSnapshot
) => {
  const authService = inject(AuthService);
  const router = inject(Router);

  console.log('[ROLE GUARD] Checking admin role for:', state.url);

  return authService.getCurrentUser().pipe(
    map((user) => {
      // Check if user has ADMIN role
      if (user?.role === 'ADMIN') {
        console.log('[ROLE GUARD] ✅ User is admin, allowing access');
        return true;
      }

      console.warn('[ROLE GUARD] ❌ User is not admin, denying access');
      router.navigate(['/access-denied'], {
        queryParams: { reason: 'insufficient_permissions' },
      });
      return false;
    }),
    catchError((error) => {
      console.error('[ROLE GUARD] Role check failed:', error);
      router.navigate(['/signin']);
      return of(false);
    })
  );
};

/**
 * ============================================================================
 * SELLER ROLE GUARD - For seller-only routes
 * ============================================================================
 *
 * Usage:
 * ```typescript
 * {
 *   path: 'products/manage',
 *   component: ManageProductsComponent,
 *   canActivate: [authGuard, sellerRoleGuard]
 * }
 * ```
 * ============================================================================
 */
export const sellerRoleGuard: CanActivateFn = (
  route: ActivatedRouteSnapshot,
  state: RouterStateSnapshot
) => {
  const authService = inject(AuthService);
  const router = inject(Router);

  console.log('[ROLE GUARD] Checking seller role for:', state.url);

  return authService.getCurrentUser().pipe(
    map((user) => {
      if (user?.role === 'SELLER') {
        console.log('[ROLE GUARD] ✅ User is seller, allowing access');
        return true;
      }

      console.warn('[ROLE GUARD] ❌ User is not seller, denying access');
      router.navigate(['/access-denied']);
      return false;
    }),
    catchError(() => {
      router.navigate(['/signin']);
      return of(false);
    })
  );
};

/**
 * ============================================================================
 * MULTI-ROLE GUARD - For routes accessible by multiple roles
 * ============================================================================
 *
 * Usage:
 * ```typescript
 * {
 *   path: 'dashboard',
 *   component: DashboardComponent,
 *   canActivate: [authGuard, multiRoleGuard],
 *   data: { allowedRoles: ['ADMIN', 'SELLER'] }
 * }
 * ```
 * ============================================================================
 */
export const multiRoleGuard: CanActivateFn = (
  route: ActivatedRouteSnapshot,
  state: RouterStateSnapshot
) => {
  const authService = inject(AuthService);
  const router = inject(Router);

  // Get allowed roles from route data
  const allowedRoles = route.data['allowedRoles'] as string[] | undefined;

  if (!allowedRoles || allowedRoles.length === 0) {
    console.error('[MULTI-ROLE GUARD] No allowedRoles specified in route data');
    return of(true); // Allow by default if no roles specified
  }

  console.log('[MULTI-ROLE GUARD] Checking roles:', allowedRoles, 'for:', state.url);

  return authService.getCurrentUser().pipe(
    map((user) => {
      if (user && allowedRoles.includes(user.role)) {
        console.log('[MULTI-ROLE GUARD] ✅ User has allowed role:', user.role);
        return true;
      }

      console.warn('[MULTI-ROLE GUARD] ❌ User role not allowed:', user?.role);
      router.navigate(['/access-denied']);
      return false;
    }),
    catchError(() => {
      router.navigate(['/signin']);
      return of(false);
    })
  );
};

/**
 * ============================================================================
 * USAGE EXAMPLES
 * ============================================================================
 *
 * // app.routes.ts
 *
 * import { Routes } from '@angular/router';
 * import { authGuard, adminRoleGuard, multiRoleGuard } from './core/guards/auth.guard';
 *
 * export const routes: Routes = [
 *   // Public routes (no guard)
 *   { path: '', component: HomeComponent },
 *   { path: 'signin', component: SigninComponent },
 *
 *   // Protected routes (auth required)
 *   {
 *     path: 'dashboard',
 *     component: DashboardComponent,
 *     canActivate: [authGuard]
 *   },
 *
 *   // Admin-only routes
 *   {
 *     path: 'users',
 *     component: UsersComponent,
 *     canActivate: [authGuard, adminRoleGuard]
 *   },
 *
 *   // Multi-role routes
 *   {
 *     path: 'analytics',
 *     component: AnalyticsComponent,
 *     canActivate: [authGuard, multiRoleGuard],
 *     data: { allowedRoles: ['ADMIN', 'SELLER'] }
 *   },
 *
 *   // Fallback
 *   { path: 'access-denied', component: AccessDeniedComponent },
 *   { path: '**', redirectTo: '' }
 * ];
 *
 * ============================================================================
 */
