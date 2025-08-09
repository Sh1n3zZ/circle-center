export const BASE_URL = (import.meta.env.VITE_API_BASE_URL as string) || "http://localhost:8080"; // fallback

export const API_PREFIX = (import.meta.env.VITE_API_PREFIX as string) || "/v1";

// token refresh configuration
export const TOKEN_REFRESH_THRESHOLD = ((): number => {
  const envValue = import.meta.env.VITE_TOKEN_REFRESH_THRESHOLD;
  if (envValue) {
    const parsed = parseInt(envValue, 10);
    return isNaN(parsed) ? 3 * 24 * 60 * 60 * 1000 : parsed;
  }
  return 3 * 24 * 60 * 60 * 1000; // 3 days in milliseconds as default
})();
