import { post, get } from "../client";
import type { 
  RefreshTokenResponse, 
  LogoutResponse, 
  UserProfileResponse 
} from "./types";

export const authApi = {
  /**
   * Refresh JWT token
   */
  refreshToken: async (): Promise<RefreshTokenResponse> => {
    const response = await post<RefreshTokenResponse>("/account/refresh");
    return response.data;
  },

  /**
   * Logout user (revoke current token)
   */
  logout: async (): Promise<LogoutResponse> => {
    const response = await post<LogoutResponse>("/account/logout");
    return response.data;
  },

  /**
   * Get current user profile
   */
  getUserProfile: async (): Promise<UserProfileResponse> => {
    const response = await get<UserProfileResponse>("/account/profile");
    return response.data;
  },
};
