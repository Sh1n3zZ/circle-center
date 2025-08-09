import { post } from "../client";
import type { 
  RegisterRequest, 
  RegisterResponse, 
  LoginRequest, 
  LoginResponse,
  ResendVerificationEmailRequest,
  ResendVerificationEmailResponse
} from "./types";

export const userApi = {
  /**
   * Register a new user
   */
  register: async (data: RegisterRequest): Promise<RegisterResponse> => {
    const response = await post<RegisterResponse>("/account/register", data);
    return response.data;
  },

  /**
   * Login user
   */
  login: async (data: LoginRequest): Promise<LoginResponse> => {
    const response = await post<LoginResponse>("/account/login", data);
    return response.data;
  },

  /**
   * Resend verification email
   */
  resendVerificationEmail: async (data: ResendVerificationEmailRequest): Promise<ResendVerificationEmailResponse> => {
    const response = await post<ResendVerificationEmailResponse>("/account/resend-verification", data);
    return response.data;
  },
};
