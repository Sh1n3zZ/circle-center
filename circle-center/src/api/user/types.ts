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
    email_sent: boolean;
    email_error?: string;
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
    avatar_url?: string;
    token?: string;
    expires_at?: number;
  };
}

export interface ResendVerificationEmailRequest {
  email: string;
}

export interface ResendVerificationEmailResponse {
  success: boolean;
  message: string;
  data: {
    email_sent: boolean;
    email_error?: string;
  };
}

export interface VerifyEmailRequest {
  token: string;
  email: string;
}

export interface VerifyEmailResponse {
  success: boolean;
  message: string;
  data: {
    success: boolean;
    message: string;
  };
}

export interface UploadAvatarResponse {
  success: boolean;
  message: string;
  data: {
    path: string;
    url: string;
  };
}

export interface ApiError {
  error: string;
  message: string;
  code?: string;
  email?: string;
}
