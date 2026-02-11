import { Injectable } from "@angular/core";
import { BehaviorSubject, Observable } from "rxjs";

@Injectable({ providedIn: "root" })
export class AuthStateService {
  private isAuthenticatedSubject = new BehaviorSubject<boolean>(false);
  private isCheckingSubject = new BehaviorSubject<boolean>(true);

  isAuthenticated$ = this.isAuthenticatedSubject.asObservable();
  isChecking$ = this.isCheckingSubject.asObservable();

  setAuthenticated(value: boolean): void {
    this.isAuthenticatedSubject.next(value);
  }

  setIsChecking(value: boolean): void {
    this.isCheckingSubject.next(value);
  }

  getIsAuthenticated(): boolean {
    return this.isAuthenticatedSubject.value;
  }
}
