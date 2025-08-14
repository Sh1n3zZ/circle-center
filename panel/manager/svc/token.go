package manager

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	mutils "circle-center/panel/manager/utils"
	managerdb "circle-center/repository/sqlc/manager"
)

// TokenService provides project API token management
type TokenService struct {
	queries *managerdb.Queries
	db      *sql.DB
}

// NewTokenService builds a new TokenService
func NewTokenService(db *sql.DB) *TokenService {
	return &TokenService{queries: managerdb.New(db), db: db}
}

// TokenInfo is a safe projection of a project's API key without exposing hash
type TokenInfo struct {
	ID         uint64 `json:"id"`
	Name       string `json:"name"`
	Active     bool   `json:"active"`
	LastUsedAt string `json:"last_used_at"`
	CreatedAt  string `json:"created_at"`
}

// CreateToken issues a new API token for a project, returns plaintext token and record id
func (s *TokenService) CreateToken(ctx context.Context, projectID uint64, name string) (string, uint64, error) {
	label := strings.TrimSpace(name)
	if label == "" {
		label = "default"
	}
	token, err := mutils.GenerateSecureToken(32)
	if err != nil {
		return "", 0, fmt.Errorf("generate token: %w", err)
	}
	hash := mutils.HashSHA256Hex(token)
	res, err := s.queries.CreateProjectAPIKey(ctx, managerdb.CreateProjectAPIKeyParams{
		ProjectID: projectID,
		Name:      label,
		TokenHash: hash,
	})
	if err != nil {
		return "", 0, err
	}
	id, _ := res.LastInsertId()
	return token, uint64(id), nil
}

// DeleteToken removes a token by id for a project
func (s *TokenService) DeleteToken(ctx context.Context, projectID uint64, tokenID uint64) error {
	return s.queries.DeleteAPIKey(ctx, managerdb.DeleteAPIKeyParams{ID: tokenID, ProjectID: projectID})
}

// ListTokens returns all API keys for the project (safe fields only)
func (s *TokenService) ListTokens(ctx context.Context, projectID uint64) ([]*TokenInfo, error) {
	rows, err := s.queries.ListProjectAPIKeys(ctx, projectID)
	if err != nil {
		return nil, err
	}
	out := make([]*TokenInfo, 0, len(rows))
	for _, r := range rows {
		lastUsed := ""
		if r.LastUsedAt.Valid {
			lastUsed = r.LastUsedAt.Time.UTC().Format("2006-01-02T15:04:05Z07:00")
		}
		out = append(out, &TokenInfo{
			ID:         r.ID,
			Name:       r.Name,
			Active:     r.Active,
			LastUsedAt: lastUsed,
			CreatedAt:  r.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	return out, nil
}
