export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  display_name?: string;
  phone?: string;
  locale?: string;
  timezone?: string;
}

export interface RegisterResponse {
  success: boolean;
  message: string;
  data: {
    id: number;
    username: string;
    email: string;
    display_name?: string;
    phone?: string;
    locale: string;
    timezone: string;
    created_at: string;
  };
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  success: boolean;
  message: string;
  data: {
    id: number;
    username: string;
    email: string;
    display_name?: string;
    phone?: string;
    locale: string;
    timezone: string;
    token?: string;
  };
}

export interface ApiError {
  error: string;
  message: string;
}
