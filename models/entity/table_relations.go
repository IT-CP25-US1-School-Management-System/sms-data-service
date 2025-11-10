package entity

import (
	helperModel "github.com/GodeFvt/go-backend/helper/models"
	"github.com/gofrs/uuid"
)

type TableRelations struct {
	ID         *uuid.UUID             `json:"id" db:"id"`
	SourceID   *uuid.UUID             `json:"source_id" db:"source_id"`
	TableFrom  string                 `json:"table_from" db:"table_from"`
	ColumnFrom string                 `json:"column_from" db:"column_from"`
	TableTo    string                 `json:"table_to" db:"table_to"`
	ColumnTo   string                 `json:"column_to" db:"column_to"`
	CreatedAt  *helperModel.Timestamp `json:"created_at" db:"created_at"`
}

func (tr *TableRelations) GenUUID() {
	id, _ := uuid.NewV4()
	tr.ID = &id
}
