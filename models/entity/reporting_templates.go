package entity

import (
	helperModel "github.com/GodeFvt/go-backend/helper/models"
	"github.com/gofrs/uuid"
)

type ReportingTemplate struct {
	ID         *uuid.UUID             `json:"id" db:"id"`
	DatasetID  string                 `json:"dataset_id" db:"dataset_id"`
	Name       string                 `json:"name" db:"name"`
	Columns    []*Columns             `json:"columns" db:"columns"`
	Positions  []*Position            `json:"positions" db:"positions"`
	ResourceID *string                `json:"resource_id" db:"resource_id"`
	CreatedAt  *helperModel.Timestamp `json:"created_at" db:"created_at"`
	UpdatedAt  *helperModel.Timestamp `json:"updated_at" db:"updated_at"`
}

type Position struct {
	TableName   string  `json:"table_name" db:"table_name"`
	ColumnsName string  `json:"column_name" db:"column_name"`
	X           float64 `json:"x" db:"x"`
	Y           float64 `json:"y" db:"y"`
}

func (u *ReportingTemplate) GenUUID() {
	id, _ := uuid.NewV4()
	u.ID = &id
}
