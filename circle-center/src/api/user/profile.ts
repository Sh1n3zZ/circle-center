import { get, put } from "../client";
import type { 
  GetUserProfileResponse, 
  UpdateUserProfileRequest, 
  UpdateUserProfileResponse 
} from "./types";

export const profileApi = {
  /**
   * Get current user's profile information
   */
  getUserProfile: async (): Promise<GetUserProfileResponse> => {
    const response = await get<GetUserProfileResponse>("/account/profile");
    return response.data;
  },

  /**
   * Update current user's profile information
   */
  updateUserProfile: async (data: UpdateUserProfileRequest): Promise<UpdateUserProfileResponse> => {
    const response = await put<UpdateUserProfileResponse>("/account/profile", data);
    return response.data;
  },
};
