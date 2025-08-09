import { post } from "../client";
import type { VerifyEmailRequest, VerifyEmailResponse } from "./types";

export const verificationApi = {
  /**
   * Verify user email with token and email address
   */
  verifyEmail: async (data: VerifyEmailRequest): Promise<VerifyEmailResponse> => {
    const response = await post<VerifyEmailResponse>("/account/verify", data);
    return response.data;
  },
};
