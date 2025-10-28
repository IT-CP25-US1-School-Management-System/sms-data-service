package entity

import (
	helperModel "github.com/GodeFvt/go-backend/helper/models"
	"github.com/gofrs/uuid"
)

type Sources struct {
	ID            *uuid.UUID             `json:"id" db:"id"`
	Name          string                 `json:"name" db:"name"`
	Type          string                 `json:"type" db:"type"`
	ConnectionRef string                 `json:"connection_ref" db:"connection_ref"`
	Sensitivity   string                 `json:"sensitivity" db:"sensitivity"`
	Config        string                 `json:"config" db:"config"`
	CreatedAt     *helperModel.Timestamp `json:"created_at" db:"created_at"`
}