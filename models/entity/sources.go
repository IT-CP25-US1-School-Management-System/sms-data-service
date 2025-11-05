package entity

import (
	helperModel "github.com/GodeFvt/go-backend/helper/models"
	"github.com/gofrs/uuid"
)

type Sources struct {
	ID           *uuid.UUID             `json:"id" db:"id"`
	Name         string                 `json:"name" db:"name"`
	Description  *string                `json:"description" db:"description"`
	Type         string                 `json:"type" db:"type"`
	IsActive     bool                   `json:"is_active" db:"is_active"`
	Sensitivity  string                 `json:"sensitivity" db:"sensitivity"`
	DBType       string                 `json:"db_type" db:"db_type"`
	Host         string                 `json:"host" db:"host"`
	Port         int                    `json:"port" db:"port"`
	Username     string                 `json:"username" db:"username"`
	Password     string                 `json:"password" db:"password"`
	DatabaseName string                 `json:"database_name" db:"database_name"`
	Params       *string                `json:"params" db:"params"`
	CreatedAt    *helperModel.Timestamp `json:"created_at" db:"created_at"`
	UpdatedAt    *helperModel.Timestamp `json:"updated_at" db:"updated_at"`
}

func (u *Sources) GenUUID() {
	id, _ := uuid.NewV4()
	u.ID = &id
}
