export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  message?: string;
  // extend this to include user/profile data if backend returns it
  user?: any;
}
