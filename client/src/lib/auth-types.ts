// Authentication types

export interface User {
  id: number;
  email: string;
  name?: string;
  phone?: string;
  role?: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  name?: string;
  phone?: string;
}

export interface AuthResponse {
  token: string;
  user?: User;
  message?: string;
}

export interface AuthError {
  error: string;
  message?: string;
}

