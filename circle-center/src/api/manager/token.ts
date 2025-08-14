import { get, post, request } from "../client";
import type {
  ListTokensResponse,
  CreateTokenRequest,
  CreateTokenResponse,
  DeleteTokenResponse,
} from "./types";

export const tokenApi = {
  async list(projectId: number): Promise<ListTokensResponse> {
    const res = await get<ListTokensResponse>(`/manager/projects/${projectId}/tokens`);
    return res.data;
  },

  async create(projectId: number, body: CreateTokenRequest): Promise<CreateTokenResponse> {
    const res = await post<CreateTokenResponse>(`/manager/projects/${projectId}/tokens`, body);
    return res.data;
  },

  async delete(projectId: number, tokenId: number): Promise<DeleteTokenResponse> {
    const res = await request<DeleteTokenResponse>({ url: `/manager/projects/${projectId}/tokens/${tokenId}`, method: "DELETE" });
    return res.data;
  },

};


