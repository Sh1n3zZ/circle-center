// Types for manager project APIs
// All comments are in English for maintainability

export type ProjectVisibility = "private" | "public";

export interface ProjectModel {
  id: number;
  owner_user_id: number;
  name: string;
  slug: string;
  package_name?: string;
  visibility: ProjectVisibility;
  description?: string;
  icon_count: number;
  created_at: string;
  updated_at: string;
}

export interface ApiEnvelope<T> {
  success: boolean;
  message: string;
  data: T;
}

export interface ListProjectsResponse extends ApiEnvelope<ProjectModel[]> {}

export interface CreateProjectRequest {
  name: string;
  slug?: string;
  package_name?: string;
  visibility?: ProjectVisibility;
  description?: string;
}

export interface CreateProjectResponse extends ApiEnvelope<ProjectModel> {}

export interface UpdateProjectRequest {
  name?: string;
  slug?: string;
  package_name?: string;
  visibility?: ProjectVisibility;
  description?: string;
}

export interface UpdateProjectResponse extends ApiEnvelope<ProjectModel> {}

export interface DeleteProjectResponse extends ApiEnvelope<{ id: number } | Record<string, never>> {}

export interface AssignRoleRequest {
  target_user_id: number;
  role: "admin" | "editor" | "viewer"; // changing owner is not supported here
}

export interface AssignRoleResponse extends ApiEnvelope<{ ok: boolean }> {}
