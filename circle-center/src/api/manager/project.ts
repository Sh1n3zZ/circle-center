import { get, post, put, request } from '../client';
import type {
  AssignRoleRequest,
  AssignRoleResponse,
  CreateProjectRequest,
  CreateProjectResponse,
  DeleteProjectResponse,
  ListProjectRolesResponse,
  ListProjectsResponse,
  UpdateProjectRequest,
  UpdateProjectResponse,
} from './types';

// Manager project APIs
export const projectApi = {
  // List projects of current user (owner)
  async listMyProjects(limit = 50, offset = 0): Promise<ListProjectsResponse> {
    // Assuming backend will provide an endpoint in the future like /manager/projects/mine
    // For now call a general endpoint; backend currently not implemented. Stub to /manager/projects?limit=&offset=
    const res = await get<ListProjectsResponse>(
      `/manager/projects?limit=${limit}&offset=${offset}`
    );
    return res.data;
  },

  // Create a project
  async createProject(
    data: CreateProjectRequest
  ): Promise<CreateProjectResponse> {
    const res = await post<CreateProjectResponse>('/manager/projects', data);
    return res.data;
  },

  // Update a project
  async updateProject(
    id: number,
    data: UpdateProjectRequest
  ): Promise<UpdateProjectResponse> {
    const res = await put<UpdateProjectResponse>(
      `/manager/projects/${id}`,
      data
    );
    return res.data;
  },

  // Delete a project
  async deleteProject(id: number): Promise<DeleteProjectResponse> {
    const res = await request<DeleteProjectResponse>({
      url: `/manager/projects/${id}`,
      method: 'DELETE',
    });
    return res.data;
  },

  // Assign a role to a user under a project (owner only)
  async assignRole(
    projectId: number,
    data: AssignRoleRequest
  ): Promise<AssignRoleResponse> {
    const res = await post<AssignRoleResponse>(
      `/manager/projects/${projectId}/roles`,
      data
    );
    return res.data;
  },

  // List project roles/collaborators
  async listProjectRoles(projectId: number): Promise<ListProjectRolesResponse> {
    const res = await get<ListProjectRolesResponse>(
      `/manager/projects/${projectId}/roles`
    );
    return res.data;
  },

  // Remove a collaborator from a project (owner only)
  async removeCollaborator(
    projectId: number,
    userId: number
  ): Promise<{ success: boolean }> {
    const res = await request<{ success: boolean }>({
      url: `/manager/projects/${projectId}/roles/${userId}`,
      method: 'DELETE',
    });
    return res.data;
  },
};
