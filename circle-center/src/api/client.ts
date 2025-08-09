import axios, { type AxiosRequestConfig, type AxiosResponse } from "axios";
import { BASE_URL, API_PREFIX } from "./config";
import { authStorage } from "../lib/storage";

// Create a single axios instance to be used throughout the application.
const client = axios.create({
  baseURL: `${BASE_URL}${API_PREFIX}`,
  timeout: 10000, // 10 seconds timeout
  headers: {
    "Content-Type": "application/json",
  },
});

// Token refresh flag to prevent multiple concurrent refresh requests
let isRefreshing = false;
let failedQueue: Array<{
  resolve: (value?: any) => void;
  reject: (reason?: any) => void;
}> = [];

const processQueue = (error: any, token: string | null = null) => {
  failedQueue.forEach(({ resolve, reject }) => {
    if (error) {
      reject(error);
    } else {
      resolve(token);
    }
  });
  
  failedQueue = [];
};

// Request interceptor for attaching auth token and refresh
client.interceptors.request.use(
  async (config) => {
    const skipAuthEndpoints = ['/account/login', '/account/register', '/account/resend-verification'];
    const isAuthEndpoint = skipAuthEndpoints.some(endpoint => config.url?.includes(endpoint));
    
    if (isAuthEndpoint) {
      return config;
    }

    const authHeader = authStorage.getAuthHeader();
    
    if (authHeader) {
      config.headers["Authorization"] = authHeader;
    }

    if (authStorage.shouldRefreshToken() && !isRefreshing && !config.url?.includes('/account/refresh')) {
      try {
        isRefreshing = true;
        
        const { authApi } = await import('./auth/auth');
        const refreshResponse = await authApi.refreshToken();
        
        await authStorage.setAuthToken(refreshResponse.data.token, refreshResponse.data.expires_at);
        
        config.headers["Authorization"] = `Bearer ${refreshResponse.data.token}`;
        
        processQueue(null, refreshResponse.data.token);
      } catch (refreshError) {
        await authStorage.clearAuth();
        processQueue(refreshError, null);
        
        if (typeof window !== 'undefined') {
          window.location.href = '/login';
        }
        
        return Promise.reject(refreshError);
      } finally {
        isRefreshing = false;
      }
    }

    return config;
  },
  (error) => Promise.reject(error),
);

// Response interceptor for handling global errors, refreshing tokens, etc.
client.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    if (error.response?.status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        }).then(token => {
          originalRequest.headers['Authorization'] = `Bearer ${token}`;
          return client(originalRequest);
        }).catch(err => {
          return Promise.reject(err);
        });
      }

      originalRequest._retry = true;
      isRefreshing = true;

      try {
        const { authApi } = await import('./auth/auth');
        const refreshResponse = await authApi.refreshToken();
        
        await authStorage.setAuthToken(refreshResponse.data.token, refreshResponse.data.expires_at);
        
        originalRequest.headers['Authorization'] = `Bearer ${refreshResponse.data.token}`;
        
        processQueue(null, refreshResponse.data.token);
        
        return client(originalRequest);
      } catch (refreshError) {
        await authStorage.clearAuth();
        processQueue(refreshError, null);
        
        if (typeof window !== 'undefined') {
          window.location.href = '/login';
        }
        
        return Promise.reject(refreshError);
      } finally {
        isRefreshing = false;
      }
    }

    return Promise.reject(error);
  },
);

// -----------------------------
// Generic request wrapper
// -----------------------------

export interface RequestOptions<T = any> extends AxiosRequestConfig<T> {
  /**
   * Mark the request data to be sent as FormData (multipart/form-data).
   * If true and `data` is a plain object, it will be converted automatically.
   */
  asFormData?: boolean;
}

/**
 * Unified request function that allows per-call control over content type while
 * keeping a single axios instance under the hood.
 */
export function request<T = any>(options: RequestOptions<T>): Promise<AxiosResponse<T>> {
  const { asFormData, data, headers, transformRequest, ...rest } = options;

  let finalHeaders = headers ?? {};
  let finalData: unknown = data;
  let finalTransformRequest = transformRequest;

  if (asFormData) {
    if (data && typeof data === "object" && !(data instanceof FormData)) {
      const formData = new FormData();
      Object.entries(data as Record<string, unknown>).forEach(([key, value]) => {
        if (Array.isArray(value)) {
          value.forEach((v) => formData.append(key, v as any));
        } else if (value !== undefined && value !== null) {
          formData.append(key, value as any);
        }
      });
      finalData = formData;
    }

    finalHeaders = {
      ...finalHeaders,
      "Content-Type": "multipart/form-data",
    };

    finalTransformRequest = (reqData: unknown) => reqData;
  }

  return client.request<T>({
    data: finalData,
    headers: finalHeaders,
    transformRequest: finalTransformRequest,
    ...rest,
  });
}

// Optional helper shortcuts
export const get = <T = any>(url: string, config?: RequestOptions) =>
  request<T>({ url, method: "GET", ...config });

export const post = <T = any>(url: string, data?: any, config?: RequestOptions) =>
  request<T>({ url, method: "POST", data, ...config });

// Authentication helper functions
export const authHelpers = {
  /**
   * Store authentication data after successful login
   */
  async storeAuthData(token: string, expiresAt: number, userInfo?: any): Promise<boolean> {
    const tokenStored = await authStorage.setAuthToken(token, expiresAt);
    
    if (userInfo) {
      const userInfoStored = await authStorage.setUserInfo(userInfo);
      return tokenStored && userInfoStored;
    }
    
    return tokenStored;
  },

  /**
   * Clear all authentication data (for logout)
   */
  async clearAuthData(): Promise<boolean> {
    return await authStorage.clearAuth();
  },

  /**
   * Check if user is currently authenticated
   */
  isAuthenticated(): boolean {
    return authStorage.isAuthenticated();
  },

  /**
   * Get current user info
   */
  getCurrentUser<T = any>(): T | null {
    return authStorage.getUserInfo<T>();
  },

  /**
   * Check if token needs refresh
   */
  shouldRefreshToken(): boolean {
    return authStorage.shouldRefreshToken();
  }
};

export default client;
