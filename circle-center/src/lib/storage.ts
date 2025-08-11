import secureLocalStorage from 'react-secure-storage';
import { TOKEN_REFRESH_THRESHOLD } from '../api/config';

// Type definitions for better TypeScript support
export type StorageValue = string | number | boolean | object | null;
export type StorageSetValue = string | number | boolean | object;

/**
 * Secure Storage Wrapper
 * 
 * This wrapper provides a secure interface for localStorage operations using react-secure-storage.
 * All data is encrypted using browser fingerprints, ensuring that data can only be decrypted
 * on the same browser where it was encrypted.
 * 
 * Supported data types: string, number, boolean, object
 */
export class SecureStorage {
  private static instance: SecureStorage;

  // Private constructor to implement singleton pattern
  private constructor() {}

  /**
   * Get singleton instance of SecureStorage
   */
  public static getInstance(): SecureStorage {
    if (!SecureStorage.instance) {
      SecureStorage.instance = new SecureStorage();
    }
    return SecureStorage.instance;
  }

  /**
   * Store a value securely in localStorage
   * @param key - The key to store the value under
   * @param value - The value to store (string, number, boolean, or object)
   * @returns Promise<boolean> - Success status
   */
  public async setItem(key: string, value: StorageSetValue): Promise<boolean> {
    try {
      if (key.trim() === '') {
        throw new Error('Storage key cannot be empty');
      }

      secureLocalStorage.setItem(key, value);
      return true;
    } catch (error) {
      console.error(`Failed to set item in secure storage:`, error);
      return false;
    }
  }

  /**
   * Retrieve a value from secure localStorage
   * @param key - The key to retrieve the value for
   * @returns The stored value or null if not found
   */
  public getItem<T = StorageValue>(key: string): T | null {
    try {
      if (key.trim() === '') {
        throw new Error('Storage key cannot be empty');
      }

      const value = secureLocalStorage.getItem(key);
      return value as T;
    } catch (error) {
      console.error(`Failed to get item from secure storage:`, error);
      return null;
    }
  }

  /**
   * Remove a specific item from secure localStorage
   * @param key - The key to remove
   * @returns Promise<boolean> - Success status
   */
  public async removeItem(key: string): Promise<boolean> {
    try {
      if (key.trim() === '') {
        throw new Error('Storage key cannot be empty');
      }

      secureLocalStorage.removeItem(key);
      return true;
    } catch (error) {
      console.error(`Failed to remove item from secure storage:`, error);
      return false;
    }
  }

  /**
   * Clear all items from secure localStorage
   * @returns Promise<boolean> - Success status
   */
  public async clear(): Promise<boolean> {
    try {
      secureLocalStorage.clear();
      return true;
    } catch (error) {
      console.error(`Failed to clear secure storage:`, error);
      return false;
    }
  }

  /**
   * Check if a key exists in secure localStorage
   * @param key - The key to check
   * @returns boolean - Whether the key exists
   */
  public hasItem(key: string): boolean {
    try {
      if (key.trim() === '') {
        return false;
      }

      const value = secureLocalStorage.getItem(key);
      return value !== null;
    } catch (error) {
      console.error(`Failed to check item existence in secure storage:`, error);
      return false;
    }
  }

  /**
   * Get all keys from secure localStorage
   * Note: This method accesses the underlying localStorage to get keys
   * @returns string[] - Array of all storage keys
   */
  public getAllKeys(): string[] {
    try {
      const keys: string[] = [];
      const prefix = import.meta.env.VITE_SECURE_LOCAL_STORAGE_PREFIX || 'secure__';
      
      for (let i = 0; i < localStorage.length; i++) {
        const key = localStorage.key(i);
        if (key && key.startsWith(prefix)) {
          // Remove the prefix to get the actual key
          keys.push(key.substring(prefix.length));
        }
      }
      
      return keys;
    } catch (error) {
      console.error(`Failed to get all keys from secure storage:`, error);
      return [];
    }
  }

  /**
   * Get the size of secure localStorage
   * @returns number - Number of items in storage
   */
  public getSize(): number {
    try {
      return this.getAllKeys().length;
    } catch (error) {
      console.error(`Failed to get secure storage size:`, error);
      return 0;
    }
  }
}

// Authentication-specific storage methods
export class AuthStorage {
  private storage: SecureStorage;
  
  // Storage keys for authentication data - only token related
  private static readonly AUTH_TOKEN_KEY = 'auth_token';
  private static readonly AUTH_EXPIRES_AT_KEY = 'auth_expires_at';
  private static readonly REFRESH_TOKEN_KEY = 'refresh_token';

  constructor() {
    this.storage = SecureStorage.getInstance();
  }

  /**
   * Store authentication token
   * @param token - JWT token
   * @param expiresAt - Token expiration timestamp
   */
  public async setAuthToken(token: string, expiresAt?: number): Promise<boolean> {
    const success = await this.storage.setItem(AuthStorage.AUTH_TOKEN_KEY, token);
    
    if (success && expiresAt) {
      await this.storage.setItem(AuthStorage.AUTH_EXPIRES_AT_KEY, expiresAt);
    }
    
    return success;
  }

  /**
   * Get authentication token
   * @returns string | null - The stored token or null if not found
   */
  public getAuthToken(): string | null {
    return this.storage.getItem<string>(AuthStorage.AUTH_TOKEN_KEY);
  }

  /**
   * Get token expiration timestamp
   * @returns number | null - The expiration timestamp or null if not found
   */
  public getTokenExpiresAt(): number | null {
    return this.storage.getItem<number>(AuthStorage.AUTH_EXPIRES_AT_KEY);
  }

  /**
   * Check if token is expired
   * @returns boolean - Whether the token is expired
   */
  public isTokenExpired(): boolean {
    const expiresAt = this.getTokenExpiresAt();
    if (!expiresAt) {
      return true;
    }
    
    return Date.now() > expiresAt * 1000; // Convert to milliseconds
  }

  /**
   * Check if token needs refresh (based on configurable threshold, default 3 days)
   * @returns boolean - Whether the token needs refresh
   */
  public shouldRefreshToken(): boolean {
    const expiresAt = this.getTokenExpiresAt();
    if (!expiresAt) {
      return true;
    }
    
    const thresholdFromNow = Date.now() + TOKEN_REFRESH_THRESHOLD;
    return thresholdFromNow > expiresAt * 1000;
  }



  /**
   * Store refresh token (if using refresh token strategy)
   * @param refreshToken - Refresh token
   */
  public async setRefreshToken(refreshToken: string): Promise<boolean> {
    return await this.storage.setItem(AuthStorage.REFRESH_TOKEN_KEY, refreshToken);
  }

  /**
   * Get refresh token
   * @returns string | null - The stored refresh token or null if not found
   */
  public getRefreshToken(): string | null {
    return this.storage.getItem<string>(AuthStorage.REFRESH_TOKEN_KEY);
  }

  /**
   * Check if user is authenticated
   * @returns boolean - Whether user has valid authentication
   */
  public isAuthenticated(): boolean {
    const token = this.getAuthToken();
    if (!token) {
      return false;
    }
    
    return !this.isTokenExpired();
  }

  /**
   * Clear all authentication data
   * @returns Promise<boolean> - Success status
   */
  public async clearAuth(): Promise<boolean> {
    const results = await Promise.all([
      this.storage.removeItem(AuthStorage.AUTH_TOKEN_KEY),
      this.storage.removeItem(AuthStorage.AUTH_EXPIRES_AT_KEY),
      this.storage.removeItem(AuthStorage.REFRESH_TOKEN_KEY),
    ]);
    
    return results.every(result => result);
  }

  /**
   * Get authentication header value for API requests
   * @returns string | null - Bearer token header value or null if not authenticated
   */
  public getAuthHeader(): string | null {
    const token = this.getAuthToken();
    if (!token || this.isTokenExpired()) {
      return null;
    }
    
    return `Bearer ${token}`;
  }
}

// Export singleton instances for easy use
export const secureStorage = SecureStorage.getInstance();
export const authStorage = new AuthStorage();

// Export default as secureStorage for backward compatibility
export default secureStorage;
