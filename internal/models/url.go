package models

import (
	"github.com/google/uuid"
)

type URL struct {
	UserID   uuid.UUID `db:"user_id" json:"userID,omitempty"`
	Short    string    `db:"short_url" json:"short"`
	Long     string    `db:"long_url" json:"long"`
	IsActive bool      `db:"is_active" json:"is_active,omitempty"`
}
