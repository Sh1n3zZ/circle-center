package manager

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	accountsvc "circle-center/panel/account/svc"
	mutils "circle-center/panel/manager/utils"
	managerdb "circle-center/repository/sqlc/manager"
)

// ProjectService handles project-related business logic
// It coordinates auth, quotas, and DB access.
type ProjectService struct {
	queries    *managerdb.Queries
	authClient *accountsvc.AuthClient
}

// NewProjectService constructs a ProjectService instance
func NewProjectService(db *sql.DB, authClient *accountsvc.AuthClient) *ProjectService {
	return &ProjectService{
		queries:    managerdb.New(db),
		authClient: authClient,
	}
}

// CreateProjectRequest represents the payload for creating a project
// visibility: "private" | "public"
type CreateProjectRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=255"`
	Slug        *string `json:"slug,omitempty"`
	PackageName *string `json:"package_name,omitempty"`
	Visibility  *string `json:"visibility,omitempty"`
	Description *string `json:"description,omitempty"`
}

// UpdateProjectRequest represents the payload for updating a project
// All fields are optional; only provided ones will be updated
type UpdateProjectRequest struct {
	Name        *string `json:"name,omitempty"`
	Slug        *string `json:"slug,omitempty"`
	PackageName *string `json:"package_name,omitempty"`
	Visibility  *string `json:"visibility,omitempty"`
	Description *string `json:"description,omitempty"`
}

// AssignRoleRequest represents assigning a role to a user in a project
type AssignRoleRequest struct {
	TargetUserID uint64 `json:"target_user_id" binding:"required"`
	Role         string `json:"role" binding:"required"` // owner|admin|editor|viewer (owner change is restricted)
}

// CreateProjectResponse represents the response for project creation
type CreateProjectResponse struct {
	ID          uint64 `json:"id"`
	OwnerUserID uint64 `json:"owner_user_id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	PackageName string `json:"package_name,omitempty"`
	Visibility  string `json:"visibility"`
	Description string `json:"description,omitempty"`
	IconCount   uint32 `json:"icon_count"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CreateProject creates a new project for the authenticated user
func (s *ProjectService) CreateProject(ctx context.Context, token string, req *CreateProjectRequest) (*CreateProjectResponse, error) {
	if s.authClient == nil {
		return nil, fmt.Errorf("auth client not initialized")
	}

	claims, err := s.authClient.ValidateToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	ownerUserID := claims.UserID

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, fmt.Errorf("project name is required")
	}

	visibility := managerdb.ProjectsVisibilityPrivate
	if req.Visibility != nil && *req.Visibility != "" {
		switch strings.ToLower(strings.TrimSpace(*req.Visibility)) {
		case "private":
			visibility = managerdb.ProjectsVisibilityPrivate
		case "public":
			visibility = managerdb.ProjectsVisibilityPublic
		default:
			return nil, fmt.Errorf("invalid visibility: %s", *req.Visibility)
		}
	}

	slug := ""
	if req.Slug != nil && strings.TrimSpace(*req.Slug) != "" {
		slug = mutils.Slugify(*req.Slug)
	} else {
		slug = mutils.Slugify(name)
	}
	if slug == "" {
		return nil, fmt.Errorf("invalid slug")
	}

	quota, err := s.queries.CheckUserQuota(ctx, managerdb.CheckUserQuotaParams{
		OwnerUserID: ownerUserID,
		UserID:      ownerUserID,
	})
	if err == nil {
		if can, convErr := mutils.AsBool(quota.CanCreateProject); convErr == nil && !can {
			return nil, fmt.Errorf("project limit reached for user")
		}
	}

	var pkg sql.NullString
	if req.PackageName != nil {
		p := strings.TrimSpace(*req.PackageName)
		if p != "" {
			pkg = sql.NullString{String: p, Valid: true}
		}
	}

	var desc sql.NullString
	if req.Description != nil {
		d := strings.TrimSpace(*req.Description)
		if d != "" {
			desc = sql.NullString{String: d, Valid: true}
		}
	}

	result, err := s.queries.CreateProject(ctx, managerdb.CreateProjectParams{
		OwnerUserID: ownerUserID,
		Name:        name,
		Slug:        slug,
		PackageName: pkg,
		Visibility:  visibility,
		Description: desc,
	})
	if err != nil {
		e := strings.ToLower(err.Error())
		if strings.Contains(e, "duplicate") || strings.Contains(e, "unique") {
			return nil, fmt.Errorf("project slug already exists")
		}
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	projectID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get project id: %w", err)
	}

	_, _ = s.queries.CreateUserProjectRole(ctx, managerdb.CreateUserProjectRoleParams{
		UserID:    ownerUserID,
		ProjectID: uint64(projectID),
		Role:      managerdb.UserProjectRolesRoleOwner,
	})

	project, err := s.queries.GetProjectByID(ctx, uint64(projectID))
	if err != nil {
		return nil, fmt.Errorf("failed to load created project: %w", err)
	}

	resp := &CreateProjectResponse{
		ID:          project.ID,
		OwnerUserID: project.OwnerUserID,
		Name:        project.Name,
		Slug:        project.Slug,
		PackageName: mutils.NullString(project.PackageName),
		Visibility:  string(project.Visibility),
		Description: mutils.NullString(project.Description),
		IconCount:   project.IconCount,
		CreatedAt:   project.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   project.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
	}
	return resp, nil
}

// ListProjects returns current user's projects with pagination
func (s *ProjectService) ListProjects(ctx context.Context, token string, limit, offset int32) ([]*CreateProjectResponse, error) {
	if s.authClient == nil {
		return nil, fmt.Errorf("auth client not initialized")
	}
	claims, err := s.authClient.ValidateToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	rows, err := s.queries.ListProjectsByOwner(ctx, managerdb.ListProjectsByOwnerParams{
		OwnerUserID: claims.UserID,
		Limit:       limit,
		Offset:      offset,
	})
	if err != nil {
		return nil, err
	}
	list := make([]*CreateProjectResponse, 0, len(rows))
	for _, p := range rows {
		list = append(list, &CreateProjectResponse{
			ID:          p.ID,
			OwnerUserID: p.OwnerUserID,
			Name:        p.Name,
			Slug:        p.Slug,
			PackageName: mutils.NullString(p.PackageName),
			Visibility:  string(p.Visibility),
			Description: mutils.NullString(p.Description),
			IconCount:   p.IconCount,
			CreatedAt:   p.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   p.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	return list, nil
}

// GetProject returns a single project by id if owned by current user
func (s *ProjectService) GetProject(ctx context.Context, token string, projectID uint64) (*CreateProjectResponse, error) {
	if s.authClient == nil {
		return nil, fmt.Errorf("auth client not initialized")
	}
	claims, err := s.authClient.ValidateToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	p, err := s.queries.GetProjectByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("project not found")
	}
	if p.OwnerUserID != claims.UserID {
		return nil, fmt.Errorf("forbidden")
	}

	resp := &CreateProjectResponse{
		ID:          p.ID,
		OwnerUserID: p.OwnerUserID,
		Name:        p.Name,
		Slug:        p.Slug,
		PackageName: mutils.NullString(p.PackageName),
		Visibility:  string(p.Visibility),
		Description: mutils.NullString(p.Description),
		IconCount:   p.IconCount,
		CreatedAt:   p.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   p.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
	}
	return resp, nil
}

// UpdateProject updates editable fields for an owner's project
func (s *ProjectService) UpdateProject(ctx context.Context, token string, projectID uint64, req *UpdateProjectRequest) (*CreateProjectResponse, error) {
	if s.authClient == nil {
		return nil, fmt.Errorf("auth client not initialized")
	}
	claims, err := s.authClient.ValidateToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	project, err := s.queries.GetProjectByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("project not found")
	}
	if project.OwnerUserID != claims.UserID {
		return nil, fmt.Errorf("forbidden")
	}

	// resolve fields (use existing when not provided)
	name := project.Name
	if req.Name != nil {
		name = strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, fmt.Errorf("project name is required")
		}
	}

	slug := project.Slug
	if req.Slug != nil {
		val := strings.TrimSpace(*req.Slug)
		if val == "" {
			return nil, fmt.Errorf("invalid slug")
		}
		slug = mutils.Slugify(val)
	}

	visibility := project.Visibility
	if req.Visibility != nil && *req.Visibility != "" {
		switch strings.ToLower(strings.TrimSpace(*req.Visibility)) {
		case "private":
			visibility = managerdb.ProjectsVisibilityPrivate
		case "public":
			visibility = managerdb.ProjectsVisibilityPublic
		default:
			return nil, fmt.Errorf("invalid visibility: %s", *req.Visibility)
		}
	}

	var pkg sql.NullString = project.PackageName
	if req.PackageName != nil {
		p := strings.TrimSpace(*req.PackageName)
		if p == "" {
			pkg = sql.NullString{}
		} else {
			pkg = sql.NullString{String: p, Valid: true}
		}
	}

	var desc sql.NullString = project.Description
	if req.Description != nil {
		d := strings.TrimSpace(*req.Description)
		if d == "" {
			desc = sql.NullString{}
		} else {
			desc = sql.NullString{String: d, Valid: true}
		}
	}

	err = s.queries.UpdateProject(ctx, managerdb.UpdateProjectParams{
		Name:        name,
		Slug:        slug,
		PackageName: pkg,
		Visibility:  visibility,
		Description: desc,
		ID:          projectID,
		OwnerUserID: project.OwnerUserID,
	})
	if err != nil {
		e := strings.ToLower(err.Error())
		if strings.Contains(e, "duplicate") || strings.Contains(e, "unique") {
			return nil, fmt.Errorf("project slug already exists")
		}
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	updated, err := s.queries.GetProjectByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to load updated project: %w", err)
	}

	resp := &CreateProjectResponse{
		ID:          updated.ID,
		OwnerUserID: updated.OwnerUserID,
		Name:        updated.Name,
		Slug:        updated.Slug,
		PackageName: mutils.NullString(updated.PackageName),
		Visibility:  string(updated.Visibility),
		Description: mutils.NullString(updated.Description),
		IconCount:   updated.IconCount,
		CreatedAt:   updated.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   updated.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
	}
	return resp, nil
}

// DeleteProject deletes a project owned by the authenticated user
func (s *ProjectService) DeleteProject(ctx context.Context, token string, projectID uint64) error {
	if s.authClient == nil {
		return fmt.Errorf("auth client not initialized")
	}
	claims, err := s.authClient.ValidateToken(ctx, token)
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	project, err := s.queries.GetProjectByID(ctx, projectID)
	if err != nil {
		return fmt.Errorf("project not found")
	}
	if project.OwnerUserID != claims.UserID {
		return fmt.Errorf("forbidden")
	}

	return s.queries.DeleteProject(ctx, managerdb.DeleteProjectParams{
		ID:          projectID,
		OwnerUserID: project.OwnerUserID,
	})
}

// AssignProjectRole creates or updates a collaborator role; owner-only operation
func (s *ProjectService) AssignProjectRole(ctx context.Context, token string, projectID uint64, req *AssignRoleRequest) error {
	if s.authClient == nil {
		return fmt.Errorf("auth client not initialized")
	}
	claims, err := s.authClient.ValidateToken(ctx, token)
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	project, err := s.queries.GetProjectByID(ctx, projectID)
	if err != nil {
		return fmt.Errorf("project not found")
	}
	if project.OwnerUserID != claims.UserID {
		return fmt.Errorf("forbidden")
	}

	roleStr := strings.ToLower(strings.TrimSpace(req.Role))
	var role managerdb.UserProjectRolesRole
	switch roleStr {
	case "admin":
		role = managerdb.UserProjectRolesRoleAdmin
	case "editor":
		role = managerdb.UserProjectRolesRoleEditor
	case "viewer":
		role = managerdb.UserProjectRolesRoleViewer
	case "owner":
		return fmt.Errorf("changing owner is not supported")
	default:
		return fmt.Errorf("invalid role: %s", req.Role)
	}

	// if exists update, else create
	_, err = s.queries.GetUserProjectRole(ctx, managerdb.GetUserProjectRoleParams{
		UserID:    req.TargetUserID,
		ProjectID: projectID,
	})
	if err == nil {
		return s.queries.UpdateUserProjectRole(ctx, managerdb.UpdateUserProjectRoleParams{
			Role:      role,
			UserID:    req.TargetUserID,
			ProjectID: projectID,
		})
	}

	_, err = s.queries.CreateUserProjectRole(ctx, managerdb.CreateUserProjectRoleParams{
		UserID:    req.TargetUserID,
		ProjectID: projectID,
		Role:      role,
	})
	return err
}
