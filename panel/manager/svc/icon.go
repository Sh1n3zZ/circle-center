package manager

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	mutils "circle-center/panel/manager/utils"
	managerdb "circle-center/repository/sqlc/manager"
)

// IconService handles icon-related business logic
type IconService struct {
	db      *sql.DB
	queries *managerdb.Queries
}

// NewIconService constructs a new IconService
func NewIconService(db *sql.DB) *IconService {
	return &IconService{db: db, queries: managerdb.New(db)}
}

// IconModel represents an icon in the API response
type IconModel struct {
	ID            uint64           `json:"id"`
	ProjectID     uint64           `json:"projectId"`
	Name          string           `json:"name"`
	Package       string           `json:"pkg"`
	ComponentInfo string           `json:"componentInfo"`
	Drawable      string           `json:"drawable"`
	Status        string           `json:"status"`
	Metadata      *json.RawMessage `json:"metadata,omitempty"`
	CreatedAt     string           `json:"createdAt"`
	UpdatedAt     string           `json:"updatedAt"`
}

// ListIconsParams represents parameters for listing icons
type ListIconsParams struct {
	ProjectID uint64 `json:"projectId"`
	Status    string `json:"status,omitempty"`
	Package   string `json:"package,omitempty"`
	Search    string `json:"search,omitempty"`
	Limit     int32  `json:"limit"`
	Offset    int32  `json:"offset"`
}

// CreateIconRequest represents the request to create an icon
type CreateIconRequest struct {
	Name          string           `json:"name" binding:"required"`
	Package       string           `json:"pkg" binding:"required"`
	ComponentInfo string           `json:"componentInfo" binding:"required"`
	Drawable      string           `json:"drawable" binding:"required"`
	Status        string           `json:"status,omitempty"`
	Metadata      *json.RawMessage `json:"metadata,omitempty"`
}

// UpdateIconRequest represents the request to update an icon
type UpdateIconRequest struct {
	Name          string           `json:"name,omitempty"`
	Package       string           `json:"pkg,omitempty"`
	ComponentInfo string           `json:"componentInfo,omitempty"`
	Drawable      string           `json:"drawable,omitempty"`
	Status        string           `json:"status,omitempty"`
	Metadata      *json.RawMessage `json:"metadata,omitempty"`
}

// ListIcons retrieves icons for a project with optional filtering
func (s *IconService) ListIcons(ctx context.Context, params ListIconsParams) ([]IconModel, int64, error) {
	var icons []managerdb.Icon
	var total int64
	var err error

	// Get total count first
	if params.Search != "" {
		// For search, we need to count separately since SearchIcons doesn't return count
		total, err = s.queries.CountProjectIcons(ctx, params.ProjectID)
		if err != nil {
			return nil, 0, err
		}
	} else if params.Status != "" {
		total, err = s.queries.CountIconsByStatus(ctx, managerdb.CountIconsByStatusParams{
			ProjectID: params.ProjectID,
			Status:    managerdb.IconsStatus(params.Status),
		})
		if err != nil {
			return nil, 0, err
		}
	} else {
		total, err = s.queries.CountProjectIcons(ctx, params.ProjectID)
		if err != nil {
			return nil, 0, err
		}
	}

	// Get icons based on filters
	if params.Search != "" {
		searchTerm := "%" + params.Search + "%"
		icons, err = s.queries.SearchIcons(ctx, managerdb.SearchIconsParams{
			ProjectID:     params.ProjectID,
			Name:          searchTerm,
			Pkg:           searchTerm,
			ComponentInfo: searchTerm,
			Name_2:        searchTerm,
			Pkg_2:         searchTerm,
			Limit:         params.Limit,
			Offset:        params.Offset,
		})
	} else if params.Status != "" {
		icons, err = s.queries.ListIconsByStatus(ctx, managerdb.ListIconsByStatusParams{
			ProjectID: params.ProjectID,
			Status:    managerdb.IconsStatus(params.Status),
			Limit:     params.Limit,
			Offset:    params.Offset,
		})
	} else if params.Package != "" {
		icons, err = s.queries.ListIconsByPackage(ctx, managerdb.ListIconsByPackageParams{
			ProjectID: params.ProjectID,
			Pkg:       params.Package,
		})
		// For package filter, we don't have pagination, so we need to handle it manually
		if err == nil && len(icons) > int(params.Offset) {
			end := int(params.Offset) + int(params.Limit)
			if end > len(icons) {
				end = len(icons)
			}
			icons = icons[params.Offset:end]
		}
	} else {
		icons, err = s.queries.ListProjectIcons(ctx, managerdb.ListProjectIconsParams{
			ProjectID: params.ProjectID,
			Limit:     params.Limit,
			Offset:    params.Offset,
		})
	}

	if err != nil {
		return nil, 0, err
	}

	// Convert to API model
	result := make([]IconModel, len(icons))
	for i, icon := range icons {
		result[i] = IconModel{
			ID:            icon.ID,
			ProjectID:     icon.ProjectID,
			Name:          icon.Name,
			Package:       icon.Pkg,
			ComponentInfo: icon.ComponentInfo,
			Drawable:      icon.Drawable,
			Status:        string(icon.Status),
			Metadata:      mutils.ConvertNullStringToRawMessage(icon.Metadata),
			CreatedAt:     icon.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:     icon.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return result, total, nil
}

// GetIcon retrieves a single icon by ID
func (s *IconService) GetIcon(ctx context.Context, projectID, iconID uint64) (*IconModel, error) {
	icon, err := s.queries.GetIconByID(ctx, iconID)
	if err != nil {
		return nil, err
	}

	// Verify the icon belongs to the project
	if icon.ProjectID != projectID {
		return nil, fmt.Errorf("icon not found in project")
	}

	return &IconModel{
		ID:            icon.ID,
		ProjectID:     icon.ProjectID,
		Name:          icon.Name,
		Package:       icon.Pkg,
		ComponentInfo: icon.ComponentInfo,
		Drawable:      icon.Drawable,
		Status:        string(icon.Status),
		Metadata:      mutils.ConvertNullStringToRawMessage(icon.Metadata),
		CreatedAt:     icon.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     icon.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// CreateIcon creates a new icon
func (s *IconService) CreateIcon(ctx context.Context, projectID uint64, req CreateIconRequest) (*IconModel, error) {
	// Validate status
	status := managerdb.IconsStatusPending
	if req.Status != "" {
		status = managerdb.IconsStatus(req.Status)
	}

	// Check for duplicate component_info in the same project
	_, err := s.queries.GetIconByComponent(ctx, managerdb.GetIconByComponentParams{
		ProjectID:     projectID,
		ComponentInfo: req.ComponentInfo,
	})
	if err == nil {
		// Icon already exists
		return nil, fmt.Errorf("icon with component info %s already exists", req.ComponentInfo)
	}

	// Create the icon
	result, err := s.queries.CreateIcon(ctx, managerdb.CreateIconParams{
		ProjectID:     projectID,
		Name:          strings.TrimSpace(req.Name),
		Pkg:           strings.TrimSpace(req.Package),
		ComponentInfo: strings.TrimSpace(req.ComponentInfo),
		Drawable:      strings.TrimSpace(req.Drawable),
		Status:        status,
		Metadata:      mutils.ConvertRawMessageToNullString(req.Metadata),
	})
	if err != nil {
		return nil, err
	}

	// Get the created icon
	iconID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return s.GetIcon(ctx, projectID, uint64(iconID))
}

// UpdateIcon updates an existing icon
func (s *IconService) UpdateIcon(ctx context.Context, projectID, iconID uint64, req UpdateIconRequest) (*IconModel, error) {
	// Get existing icon to verify it exists and belongs to the project
	existing, err := s.queries.GetIconByID(ctx, iconID)
	if err != nil {
		return nil, err
	}
	if existing.ProjectID != projectID {
		return nil, fmt.Errorf("icon not found in project")
	}

	// Prepare update fields
	name := existing.Name
	if req.Name != "" {
		name = strings.TrimSpace(req.Name)
	}

	pkg := existing.Pkg
	if req.Package != "" {
		pkg = strings.TrimSpace(req.Package)
	}

	componentInfo := existing.ComponentInfo
	if req.ComponentInfo != "" {
		componentInfo = strings.TrimSpace(req.ComponentInfo)
		// Check for duplicate component_info (excluding current icon)
		if componentInfo != existing.ComponentInfo {
			_, err := s.queries.GetIconByComponent(ctx, managerdb.GetIconByComponentParams{
				ProjectID:     projectID,
				ComponentInfo: componentInfo,
			})
			if err == nil {
				return nil, fmt.Errorf("icon with component info %s already exists", componentInfo)
			}
		}
	}

	drawable := existing.Drawable
	if req.Drawable != "" {
		drawable = strings.TrimSpace(req.Drawable)
	}

	status := existing.Status
	if req.Status != "" {
		status = managerdb.IconsStatus(req.Status)
	}

	metadata := existing.Metadata
	if req.Metadata != nil {
		metadata = mutils.ConvertRawMessageToNullString(req.Metadata)
	}

	// Update the icon
	err = s.queries.UpdateIcon(ctx, managerdb.UpdateIconParams{
		Name:          name,
		Pkg:           pkg,
		ComponentInfo: componentInfo,
		Drawable:      drawable,
		Status:        status,
		Metadata:      metadata,
		ID:            iconID,
		ProjectID:     projectID,
	})
	if err != nil {
		return nil, err
	}

	return s.GetIcon(ctx, projectID, iconID)
}

// DeleteIcon deletes an icon
func (s *IconService) DeleteIcon(ctx context.Context, projectID, iconID uint64) error {
	// Verify the icon exists and belongs to the project
	existing, err := s.queries.GetIconByID(ctx, iconID)
	if err != nil {
		return err
	}
	if existing.ProjectID != projectID {
		return fmt.Errorf("icon not found in project")
	}

	return s.queries.DeleteIcon(ctx, managerdb.DeleteIconParams{
		ID:        iconID,
		ProjectID: projectID,
	})
}

// GetIconStats retrieves statistics for icons in a project
func (s *IconService) GetIconStats(ctx context.Context, projectID uint64) (managerdb.GetIconStatsRow, error) {
	return s.queries.GetIconStats(ctx, projectID)
}
