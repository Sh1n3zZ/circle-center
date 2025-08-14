package manager

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"

	mutils "circle-center/panel/manager/utils"
	managerdb "circle-center/repository/sqlc/manager"
)

// XMLIOService handles parsing XML files into icon components and persisting them.
type XMLIOService struct {
	db      *sql.DB
	queries *managerdb.Queries
}

// NewXMLIOService constructs a new XMLIOService
func NewXMLIOService(db *sql.DB) *XMLIOService {
	return &XMLIOService{db: db, queries: managerdb.New(db)}
}

// ParseXMLInputs parses optional XML strings and merges them into a consolidated list.
// Any of appfilterXML, appmapXML, themeXML can be empty. At least one must be non-empty.
func (s *XMLIOService) ParseXMLInputs(appfilterXML, appmapXML, themeXML string) ([]mutils.IconRequestComponent, error) {
	if strings.TrimSpace(appfilterXML) == "" && strings.TrimSpace(appmapXML) == "" && strings.TrimSpace(themeXML) == "" {
		return nil, errors.New("no xml provided")
	}
	var af, am, th []mutils.IconRequestComponent
	var err error
	if strings.TrimSpace(appfilterXML) != "" {
		if af, err = mutils.ParseAppFilterXML(appfilterXML); err != nil {
			return nil, err
		}
	}
	if strings.TrimSpace(appmapXML) != "" {
		if am, err = mutils.ParseAppMapXML(appmapXML); err != nil {
			return nil, err
		}
	}
	if strings.TrimSpace(themeXML) != "" {
		if th, err = mutils.ParseThemeResourcesXML(themeXML); err != nil {
			return nil, err
		}
	}
	merged := mutils.MergeParsedComponents(af, am, th)
	// Ensure fields are trimmed and derive missing fields where possible
	for i := range merged {
		merged[i].Name = strings.TrimSpace(merged[i].Name)
		merged[i].Drawable = strings.TrimSpace(merged[i].Drawable)
		merged[i].ComponentInfo = strings.TrimSpace(merged[i].ComponentInfo)
		merged[i].Package = strings.TrimSpace(merged[i].Package)
		if merged[i].Package == "" {
			merged[i].Package = mutils.InferPackageFromComponent(merged[i].ComponentInfo)
		}
		if merged[i].Name == "" {
			// Fallback to drawable as name
			merged[i].Name = merged[i].Drawable
		}
	}
	return merged, nil
}

// ImportSummary returns the result of persisting parsed components
type ImportSummary struct {
	Total      int      `json:"total"`
	Created    int      `json:"created"`
	Duplicates int      `json:"duplicates"`
	Errors     int      `json:"errors"`
	ErrorMsgs  []string `json:"errorMsgs"`
}

// SaveIcons persists components into the icons table for the given project.
// It skips duplicates on (project_id, component_info) and counts them.
func (s *XMLIOService) SaveIcons(ctx context.Context, projectID uint64, components []mutils.IconRequestComponent) (*ImportSummary, error) {
	summary := &ImportSummary{Total: len(components)}
	for _, c := range components {
		name := strings.TrimSpace(c.Name)
		if name == "" {
			name = strings.TrimSpace(c.Drawable)
		}
		pkg := strings.TrimSpace(c.Package)
		if pkg == "" {
			pkg = mutils.InferPackageFromComponent(c.ComponentInfo)
		}
		comp := strings.TrimSpace(c.ComponentInfo)
		drawable := strings.TrimSpace(c.Drawable)
		if comp == "" || drawable == "" || pkg == "" || name == "" {
			summary.Errors++
			summary.ErrorMsgs = append(summary.ErrorMsgs, "missing required fields")
			continue
		}

		_, err := s.queries.CreateIcon(ctx, managerdb.CreateIconParams{
			ProjectID:     projectID,
			Name:          name,
			Pkg:           pkg,
			ComponentInfo: comp,
			Drawable:      drawable,
			Status:        managerdb.IconsStatusPending,
			Metadata:      json.RawMessage(nil),
		})
		if err != nil {
			// best-effort: treat unique constraint as duplicate
			e := strings.ToLower(err.Error())
			if strings.Contains(e, "duplicate") || strings.Contains(e, "unique") {
				summary.Duplicates++
				continue
			}
			summary.Errors++
			summary.ErrorMsgs = append(summary.ErrorMsgs, err.Error())
			continue
		}
		summary.Created++
	}
	return summary, nil
}
