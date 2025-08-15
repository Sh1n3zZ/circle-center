package utils

import (
	"database/sql"
	"encoding/json"
)

// ConvertNullStringToRawMessage converts sql.NullString to *json.RawMessage
func ConvertNullStringToRawMessage(ns sql.NullString) *json.RawMessage {
	if !ns.Valid {
		return nil
	}
	raw := json.RawMessage(ns.String)
	return &raw
}

// ConvertRawMessageToNullString converts *json.RawMessage to sql.NullString
func ConvertRawMessageToNullString(raw *json.RawMessage) sql.NullString {
	if raw == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: string(*raw), Valid: true}
}
