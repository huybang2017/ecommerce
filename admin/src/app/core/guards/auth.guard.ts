import { inject } from '@angular/core';
import { Router, CanActivateFn, ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';
import { AuthService } from '../services/auth.service';
import { map, catchError, of } from 'rxjs';

export const authGuard: CanActivateFn = (
  route: ActivatedRouteSnapshot,
  state: RouterStateSnapshot
) => {
  const authService = inject(AuthService);
  const router = inject(Router);
  console.log('[AUTH GUARD] Checking authentication for:', state.url);
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

export const adminRoleGuard: CanActivateFn = (
  route: ActivatedRouteSnapshot,
  state: RouterStateSnapshot
) => {
  const authService = inject(AuthService);
  const router = inject(Router);

  console.log('[ROLE GUARD] Checking admin role for:', state.url);

  return authService.getCurrentUser().pipe(
    map((user) => {
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