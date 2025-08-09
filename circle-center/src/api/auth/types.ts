// Authentication related types

export interface RefreshTokenResponse {
  success: boolean;
  message: string;
  data: {
    token: string;
    expires_at: number;
  };
}

export interface LogoutResponse {
  success: boolean;
  message: string;
}

export interface UserProfileResponse {
  success: boolean;
  message: string;
  data: {
    user_id: number;
    username: string;
    email: string;
  };
}

export interface AuthError {
  error: string;
  message: string;
  code?: string;
}
