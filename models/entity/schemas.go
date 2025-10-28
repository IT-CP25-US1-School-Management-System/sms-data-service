package entity

import (
	helperModel "github.com/GodeFvt/go-backend/helper/models"
	"github.com/gofrs/uuid"
)

type Schemas struct {
	ID            int64                  `json:"id" db:"id"`
	SourceID     *uuid.UUID             `json:"sources_id" db:"sources_id"`
	Schema        string                 `json:"schema" db:"schema"`
	Discovered_at *helperModel.Timestamp `json:"discovered_at" db:"discovered_at"`
}
