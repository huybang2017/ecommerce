import { Injectable } from "@angular/core";
import {
  HttpClient,
  HttpParams,
  HttpHeaders,
  HttpErrorResponse,
} from "@angular/common/http";
import { Observable, throwError } from "rxjs";
import { catchError, timeout as rxTimeout, retry } from "rxjs/operators";
import { environment as env } from "../../environments/environment";

export interface ApiOptions {
  params?: HttpParams | Record<string, any>;
  headers?: HttpHeaders | Record<string, string>;
  withCredentials?: boolean;
  timeoutMs?: number;
  retry?: number;
}
export interface ApiError {
  status?: number;
  message: string;
  details?: any;
  isNetworkError?: boolean;
}

@Injectable({ providedIn: "root" })
export class HttpApiService {
  private readonly base = env.api;
  private readonly defaultHeaders: Record<string, string> = {
    Accept: "application/json",
    "Content-Type": "application/json",
  };
  constructor(private http: HttpClient) {}
  private buildParams(
    params?: HttpParams | Record<string, any>,
  ): HttpParams | undefined {
    if (!params) return undefined;
    if (params instanceof HttpParams) return params;
    let httpParams = new HttpParams();
    Object.keys(params).forEach((key) => {
      const value = (params as Record<string, any>)[key];
      if (value === null || value === undefined) return;
      if (Array.isArray(value)) {
        value.forEach((v) => (httpParams = httpParams.append(key, v)));
      } else {
        httpParams = httpParams.set(key, String(value));
      }
    });
    return httpParams;
  }

  private buildHeaders(
    headers?: HttpHeaders | Record<string, string>,
  ): HttpHeaders {
    let httpHeaders = new HttpHeaders(this.defaultHeaders);
    if (!headers) return httpHeaders;
    if (headers instanceof HttpHeaders) {
      headers
        .keys()
        .forEach(
          (k) => (httpHeaders = httpHeaders.set(k, headers.get(k) ?? "")),
        );
      return httpHeaders;
    }
    Object.keys(headers).forEach(
      (k) => (httpHeaders = httpHeaders.set(k, headers[k])),
    );
    return httpHeaders;
  }
  private buildOptions(options?: ApiOptions) {
    return {
      params: this.buildParams(options?.params),
      headers: this.buildHeaders(options?.headers),
      withCredentials: options?.withCredentials ?? true,
    } as const;
  }
  private readonly DEFAULT_TIMEOUT_MS = 30000;
  private handleError(err: HttpErrorResponse) {
    const apiErr: ApiError = {
      status: err.status || undefined,
      message: err.error?.message || err.message || "Unknown error",
      details: err.error || err.message,
      isNetworkError: err.status === 0 || !err.status,
    };
    return throwError(() => apiErr);
  }
  private request<T>(
    method: string,
    path: string,
    body?: any,
    options?: ApiOptions,
  ): Observable<T> {
    const url = `${this.base}${path}`;
    const opts = this.buildOptions(options);
    let req$ = this.http.request<T>(method, url, { body, ...opts });
    const timeoutMs = options?.timeoutMs ?? this.DEFAULT_TIMEOUT_MS;
    req$ = req$.pipe(rxTimeout(timeoutMs));
    if (options?.retry && options.retry > 0) {
      req$ = req$.pipe(retry(options.retry));
    }
    return req$.pipe(
      catchError((err: HttpErrorResponse) => this.handleError(err)),
    );
  }
  get<T>(path: string, options?: ApiOptions): Observable<T> {
    return this.request<T>("GET", path, undefined, options);
  }
  post<T>(path: string, body?: any, options?: ApiOptions): Observable<T> {
    return this.request<T>("POST", path, body, options);
  }

  put<T>(path: string, body?: any, options?: ApiOptions): Observable<T> {
    return this.request<T>("PUT", path, body, options);
  }
  patch<T>(path: string, body?: any, options?: ApiOptions): Observable<T> {
    return this.request<T>("PATCH", path, body, options);
  }
  delete<T>(path: string, options?: ApiOptions): Observable<T> {
    return this.request<T>("DELETE", path, undefined, options);
  }
}
