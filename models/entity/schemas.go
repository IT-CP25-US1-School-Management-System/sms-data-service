package entity

import (
	helperModel "github.com/GodeFvt/go-backend/helper/models"
	"github.com/gofrs/uuid"
)

type Schemas struct {
	ID        *uuid.UUID             `json:"id" db:"id"`
	SourceID  *uuid.UUID             `json:"sources_id" db:"sources_id"`
	Schema    string                 `json:"schema" db:"schema"`
	CreatedAt *helperModel.Timestamp `json:"created_at" db:"created_at"`
}

func (u *Schemas) GenUUID() {
	id, _ := uuid.NewV4()
	u.ID = &id
}
