// Types for manager project APIs
// All comments are in English for maintainability

export type ProjectVisibility = 'private' | 'public';

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

export interface DeleteProjectResponse
  extends ApiEnvelope<{ id: number } | Record<string, never>> {}

export interface AssignRoleRequest {
  target_user_id: number;
  role: 'admin' | 'editor' | 'viewer'; // changing owner is not supported here
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
export interface DeleteTokenResponse
  extends Partial<ApiEnvelope<Record<string, never>>> {}

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
  status: 'success' | 'error';
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
  status: 'success' | 'error';
  summary?: ImportSummary;
  message?: string;
}

// ==========================
// Icon Management Types
// ==========================

export interface IconModel {
  id: number;
  projectId: number; // Keep as projectId to match backend JSON field
  name: string;
  pkg: string;
  componentInfo: string;
  drawable: string;
  status: 'pending' | 'in_progress' | 'published' | 'rejected';
  metadata?: Record<string, any>;
  createdAt: string;
  updatedAt: string;
}

export interface ListIconsParams {
  limit?: number;
  offset?: number;
  status?: string;
  package?: string;
  search?: string;
}

export interface ListIconsResponse
  extends ApiEnvelope<{
    icons: IconModel[];
    total: number;
    totalPages: number;
    currentPage: number;
    limit: number;
    offset: number;
  }> {}

export interface GetIconResponse extends ApiEnvelope<IconModel> {}

export interface CreateIconRequest {
  name: string;
  pkg: string;
  componentInfo: string;
  drawable: string;
  status?: string;
  metadata?: Record<string, any>;
}

export interface CreateIconResponse extends ApiEnvelope<IconModel> {}

export interface UpdateIconRequest {
  name?: string;
  pkg?: string;
  componentInfo?: string;
  drawable?: string;
  status?: string;
  metadata?: Record<string, any>;
}

export interface UpdateIconResponse extends ApiEnvelope<IconModel> {}

export interface DeleteIconResponse extends ApiEnvelope<{ message: string }> {}

export interface IconStats {
  total_icons: number;
  pending_count: number;
  in_progress_count: number;
  published_count: number;
  rejected_count: number;
}

export interface GetIconStatsResponse extends ApiEnvelope<IconStats> {}

// ==========================
// Icon Upload & Retrieval Types
// ==========================

export interface UploadIconRequest {
  component_info: string;
  file: File;
}

export interface UploadIconResponse
  extends ApiEnvelope<{
    path: string;
    content_type: string;
  }> {}

export interface GetIconImageResponse
  extends ApiEnvelope<{
    data: ArrayBuffer;
    content_type: string;
  }> {}
