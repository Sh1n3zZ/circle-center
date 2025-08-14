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

// RequestService handles icon request related logic
type RequestService struct {
	db      *sql.DB
	queries *managerdb.Queries
}

// NewRequestService constructs a RequestService instance
func NewRequestService(db *sql.DB) *RequestService {
	return &RequestService{db: db, queries: managerdb.New(db)}
}

// RequestManagerResponse is the API response expected by the client
type RequestManagerResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// appsPayload represents the incoming JSON with app components
type appsPayload struct {
	Components []appComponent `json:"components"`
}

type appComponent struct {
	Name          string `json:"name"`
	Pkg           string `json:"pkg"`
	ComponentInfo string `json:"componentInfo"`
	Drawable      string `json:"drawable"`
}

// ProcessUpload validates token, stores the request and its items in a transaction
// archivePath should be a persisted path to the uploaded zip file; can be empty
func (s *RequestService) ProcessUpload(ctx context.Context, token string, appsJSON string, archivePath string, userID *uint64) (*RequestManagerResponse, error) {
	if strings.TrimSpace(token) == "" {
		return &RequestManagerResponse{Status: "error", Message: "missing token"}, fmt.Errorf("missing token")
	}

	// Validate API key by token hash
	key, err := s.queries.GetProjectAPIKeyByHash(ctx, mutils.HashSHA256Hex(token))
	if err != nil {
		return &RequestManagerResponse{Status: "error", Message: "invalid token"}, fmt.Errorf("invalid token")
	}
	_ = s.queries.UpdateAPIKeyLastUsed(ctx, key.ID)

	// Validate apps JSON and parse components
	var payload appsPayload
	if err := json.Unmarshal([]byte(appsJSON), &payload); err != nil {
		return &RequestManagerResponse{Status: "error", Message: "invalid apps json"}, fmt.Errorf("invalid apps json: %w", err)
	}
	if len(payload.Components) == 0 {
		return &RequestManagerResponse{Status: "error", Message: "no components"}, fmt.Errorf("no components")
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return &RequestManagerResponse{Status: "error", Message: "internal error"}, err
	}
	qtx := s.queries.WithTx(tx)

	// Insert icon_requests row
	var reqBy sql.NullInt64
	if userID != nil {
		reqBy = sql.NullInt64{Int64: int64(*userID), Valid: true}
	}
	var archive sql.NullString
	if strings.TrimSpace(archivePath) != "" {
		archive = sql.NullString{String: archivePath, Valid: true}
	}
	res, err := qtx.CreateIconRequest(ctx, managerdb.CreateIconRequestParams{
		ProjectID:         key.ProjectID,
		RequestedByUserID: reqBy,
		Source:            managerdb.IconRequestsSourceApi,
		AppsJson:          json.RawMessage(appsJSON),
		ArchivePath:       archive,
	})
	if err != nil {
		_ = tx.Rollback()
		return &RequestManagerResponse{Status: "error", Message: "failed to create request"}, fmt.Errorf("create icon request: %w", err)
	}
	requestID, err := res.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return &RequestManagerResponse{Status: "error", Message: "failed to create request"}, fmt.Errorf("request id: %w", err)
	}

	// Insert items; skip duplicates in the same batch by component_info
	seen := make(map[string]struct{})
	for _, c := range payload.Components {
		comp := strings.TrimSpace(c.ComponentInfo)
		if comp == "" {
			continue
		}
		if _, ok := seen[comp]; ok {
			continue
		}
		seen[comp] = struct{}{}
		_, err := qtx.CreateRequestItem(ctx, managerdb.CreateRequestItemParams{
			RequestID:     uint64(requestID),
			ProjectID:     key.ProjectID,
			Name:          strings.TrimSpace(c.Name),
			Pkg:           strings.TrimSpace(c.Pkg),
			ComponentInfo: comp,
			Drawable:      strings.TrimSpace(c.Drawable),
		})
		if err != nil {
			// best-effort: on unique constraint within same request, continue
			if e := strings.ToLower(err.Error()); strings.Contains(e, "duplicate") || strings.Contains(e, "unique") {
				continue
			}
			_ = tx.Rollback()
			return &RequestManagerResponse{Status: "error", Message: "failed to save items"}, fmt.Errorf("create request item: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return &RequestManagerResponse{Status: "error", Message: "internal error"}, err
	}

	return &RequestManagerResponse{Status: "success", Message: "request received"}, nil
}

// no-op helpers can be added here if needed
