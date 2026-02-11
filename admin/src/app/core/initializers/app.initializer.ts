import { Injectable } from '@angular/core';
import { AuthService } from '../services/auth.service';
import { AuthStateService } from '../services/auth-state.service';
import { catchError, finalize } from 'rxjs/operators';
import { of } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class AppInitializer {
  constructor(
    private authService: AuthService,
    private authState: AuthStateService,
  ) {}

  /**
   * Initialize app by checking if user has a valid session
   * Called on app bootstrap
   */
  initialize() {
    return this.authService.checkSession().pipe(
      catchError(() => {
        this.authState.setAuthenticated(false);
        return of(null);
      }),
      finalize(() => {
        // Always mark checking as complete
        this.authState.setIsChecking(false);
      })
    );
  }
}
