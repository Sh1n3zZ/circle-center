import { post, request } from '../client';
import type { GetIconImageResponse, UploadIconResponse } from './types';

// ==========================
// Icon Upload & Retrieval API
// ==========================

/**
 * Upload an icon file for a specific project and component
 * @param projectId - The project ID
 * @param componentInfo - The component info to match with icon record
 * @param file - The icon file to upload
 */
export async function uploadIcon(
  projectId: number,
  componentInfo: string,
  file: File
): Promise<UploadIconResponse> {
  const formData = new FormData();
  formData.append('component_info', componentInfo);
  formData.append('file', file);

  const response = await post<UploadIconResponse>(
    `/manager/icons/${projectId}/upload`,
    formData,
    {
      asFormData: true,
    }
  );

  return response.data;
}

/**
 * Get icon image by relative path (requires authentication)
 * @param relPath - The relative path from uploads root (e.g., "icons/123/my_icon.png")
 */
export async function getIcon(relPath: string): Promise<GetIconImageResponse> {
  try {
    const response = await request<ArrayBuffer>({
      url: `/manager/icons/${relPath}`,
      method: 'GET',
      responseType: 'arraybuffer',
    });

    return {
      success: true,
      message: 'Icon retrieved successfully',
      data: {
        data: response.data,
        content_type:
          response.headers['content-type'] || 'application/octet-stream',
      },
    };
  } catch (error: any) {
    if (error.response?.data?.code === 'ICON_FILE_NOT_FOUND') {
      throw new Error('ICON_FILE_NOT_FOUND');
    }
    throw error;
  }
}

/**
 * Get icon URL for display (creates a blob URL from the icon data)
 * @param relPath - The relative path from uploads root
 */
export async function getIconUrl(relPath: string): Promise<string> {
  try {
    const response = await getIcon(relPath);
    const blob = new Blob([response.data.data], {
      type: response.data.content_type,
    });
    return URL.createObjectURL(blob);
  } catch (error) {
    console.error('Failed to get icon URL:', error);
    throw error;
  }
}

/**
 * Revoke a blob URL to free memory
 * @param url - The blob URL to revoke
 */
export function revokeIconUrl(url: string): void {
  if (url.startsWith('blob:')) {
    URL.revokeObjectURL(url);
  }
}
