// internal/domain/models/base_model.go

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UUIDModel representa uma estrutura base que inclui um UUID como ID primário e timestamps.
type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type HTTPError struct {
	Message string `json:"message"`
}
