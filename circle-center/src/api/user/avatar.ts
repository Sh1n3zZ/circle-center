import { post } from "../client";
import { BASE_URL, API_PREFIX } from "../config";
import type { UploadAvatarResponse } from "./types";

export const avatarApi = {
  /**
   * Upload user avatar
   * @param file - The image file to upload
   * @returns Promise with upload result
   */
  uploadAvatar: async (file: File): Promise<UploadAvatarResponse['data']> => {
    const formData = new FormData();
    formData.append("file", file);

    const response = await post<UploadAvatarResponse>("/account/avatar", formData, {
      asFormData: true,
    });

    return response.data.data;
  },

  /**
   * Get avatar URL for display
   * @param path - Avatar path from upload response
   * @param size - Optional size parameter for resizing
   * @param quality - Optional quality parameter (1-100)
   * @returns Complete avatar URL
   */
  getAvatarUrl: (path: string, size?: number, quality?: number): string => {
    const baseUrl = `${BASE_URL}${API_PREFIX}/account/avatar`;
    
    const params = new URLSearchParams();
    if (size) params.append("size", size.toString());
    if (quality) params.append("quality", quality.toString());
    
    const queryString = params.toString();
    return `${baseUrl}/${path}${queryString ? `?${queryString}` : ""}`;
  },
};
