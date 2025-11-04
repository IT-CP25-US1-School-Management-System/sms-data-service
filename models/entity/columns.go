package entity

import (
	helperModel "github.com/GodeFvt/go-backend/helper/models"
	"github.com/gofrs/uuid"
)

type Columns struct {
	ID              *uuid.UUID             `json:"id" db:"id"`
	SourceID        *uuid.UUID             `json:"sources_id" db:"sources_id"`
	Schema          string                 `json:"schema" db:"schema"`
	TableName       string                 `json:"table_name" db:"table_name"`
	ColumnsName     string                 `json:"column_name" db:"column_name"`
	DataType        string                 `json:"data_type" db:"data_type"`
	IsNullable      bool                   `json:"is_nullable" db:"is_nullable"`
	ColumnDefault   *string                `json:"column_default" db:"column_default"`
	OrdinalPosition *int                   `json:"ordinal_position" db:"ordinal_position"`
	CreatedAt       *helperModel.Timestamp `json:"created_at" db:"created_at"`
}

func (u *Columns) GenUUID() {
	id, _ := uuid.NewV4()
	u.ID = &id
}
