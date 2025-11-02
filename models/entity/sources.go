package entity

import (
	helperModel "github.com/GodeFvt/go-backend/helper/models"
	"github.com/gofrs/uuid"
)

type Sources struct {
	ID          *uuid.UUID             `json:"id" db:"id"`
	Name        string                 `json:"name" db:"name"`
	Description string                 `json:"description" db:"description"`
	Type        string                 `json:"type" db:"type"`
	CreatedAt   *helperModel.Timestamp `json:"created_at" db:"created_at"`
}
