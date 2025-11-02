package entity

import (
	helperModel "github.com/GodeFvt/go-backend/helper/models"
	"github.com/gofrs/uuid"
)

type Tables struct {
	ID        int64                  `json:"id" db:"id"`
	SourceID  *uuid.UUID             `json:"sources_id" db:"sources_id"`
	Schema    string                 `json:"schema" db:"schema"`
	TableName string                 `json:"table_name" db:"table_name"`
	CreatedAt *helperModel.Timestamp `json:"created_at" db:"created_at"`
}
