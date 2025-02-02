package model

import (
	"encoding/json"
	"time"
)

type Model struct {
	ID        string          `json:"id"`
	Metadata  json.RawMessage `json:"metadata,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
