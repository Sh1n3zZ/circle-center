import axios, { type AxiosRequestConfig, type AxiosResponse } from "axios";
import { BASE_URL } from "./config";

// Create a single axios instance to be used throughout the application.
const client = axios.create({
  baseURL: BASE_URL,
  timeout: 10000, // 10 seconds timeout
  headers: {
    "Content-Type": "application/json",
  },
});

// Request interceptor for attaching auth token, logging, etc.
client.interceptors.request.use(
  (config) => {
    // Example: attach auth token if exists.
    const token = localStorage.getItem("token");
    if (token) {
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-ignore
      config.headers["Authorization"] = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error),
);

// Response interceptor for handling global errors, refreshing tokens, etc.
client.interceptors.response.use(
  (response) => response,
  (error) => {
    // Here you can handle global errors (e.g., 401 Unauthorized)
    // console.error("API error", error);
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
    // Convert plain object data to FormData if needed.
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

    // Ensure the correct content type header. Let the browser set the boundary.
    finalHeaders = {
      ...finalHeaders,
      "Content-Type": "multipart/form-data",
    };

    // When sending FormData, avoid default JSON stringify transform.
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

export default client;
