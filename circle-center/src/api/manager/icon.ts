import { get, post, put, request } from "../client";
import type {
  ListIconsParams,
  ListIconsResponse,
  GetIconResponse,
  CreateIconRequest,
  CreateIconResponse,
  UpdateIconRequest,
  UpdateIconResponse,
  DeleteIconResponse,
  GetIconStatsResponse,
} from "./types";

export const iconApi = {
  // List icons with optional filtering and pagination
  async list(projectId: number, params: ListIconsParams = {}): Promise<ListIconsResponse> {
    const searchParams = new URLSearchParams();
    if (params.limit) searchParams.append("limit", params.limit.toString());
    if (params.offset) searchParams.append("offset", params.offset.toString());
    if (params.status) searchParams.append("status", params.status);
    if (params.package) searchParams.append("package", params.package);
    if (params.search) searchParams.append("search", params.search);

    const queryString = searchParams.toString();
    const url = `/manager/projects/${projectId}/icons${queryString ? `?${queryString}` : ""}`;
    const res = await get<ListIconsResponse>(url);
    return res.data;
  },

  // Get a single icon by ID
  async get(projectId: number, iconId: number): Promise<GetIconResponse> {
    const res = await get<GetIconResponse>(`/manager/projects/${projectId}/icons/${iconId}`);
    return res.data;
  },

  // Create a new icon
  async create(projectId: number, data: CreateIconRequest): Promise<CreateIconResponse> {
    const res = await post<CreateIconResponse>(`/manager/projects/${projectId}/icons`, data);
    return res.data;
  },

  // Update an existing icon
  async update(projectId: number, iconId: number, data: UpdateIconRequest): Promise<UpdateIconResponse> {
    const res = await put<UpdateIconResponse>(`/manager/projects/${projectId}/icons/${iconId}`, data);
    return res.data;
  },

  // Delete an icon
  async delete(projectId: number, iconId: number): Promise<DeleteIconResponse> {
    const res = await request<DeleteIconResponse>({
      url: `/manager/projects/${projectId}/icons/${iconId}`,
      method: "DELETE",
    });
    return res.data;
  },

  // Get icon statistics for a project
  async getStats(projectId: number): Promise<GetIconStatsResponse> {
    const res = await get<GetIconStatsResponse>(`/manager/projects/${projectId}/icons/stats`);
    return res.data;
  },
};
