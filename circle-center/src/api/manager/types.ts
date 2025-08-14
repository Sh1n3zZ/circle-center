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

export interface ProjectTokenModel {
  id: number;
  name: string;
  active: boolean;
  last_used_at: string;
  created_at: string;
}

export interface ListTokensResponse extends ApiEnvelope<ProjectTokenModel[]> {}

export interface CreateTokenRequest {
  name?: string;
}

export interface CreateTokenData {
  token_id: number;
  token: string;
}

export interface CreateTokenResponse extends ApiEnvelope<CreateTokenData> {}
export interface DeleteTokenResponse extends Partial<ApiEnvelope<Record<string, never>>> {}

// ==========================
// XML Import Types
// ==========================

export interface IconImportComponent {
  // Keep the same keys as backend JSON
  name: string;
  pkg: string;
  componentInfo: string;
  drawable: string;
}

export interface ParseXmlRequest {
  appfilter?: string;
  appmap?: string;
  theme?: string;
}

export interface ParseXmlResponse {
  status: "success" | "error";
  components?: IconImportComponent[];
  message?: string;
}

export interface ImportSummary {
  total: number;
  created: number;
  duplicates: number;
  errors: number;
  errorMsgs: string[];
}

export interface ConfirmImportRequest {
  projectId: number;
  components: IconImportComponent[];
}

export interface ConfirmImportResponse {
  status: "success" | "error";
  summary?: ImportSummary;
  message?: string;
}