package entity

import helperModel "github.com/GodeFvt/go-backend/helper/models"

type Datasets struct {
	ID          string                 `json:"id" db:"id"`
	Name        string                 `json:"name" db:"name"`
	Description string                 `json:"description" db:"description"`
	Domain      string                 `json:"domain" db:"domain"`
	Owner       string                 `json:"owner" db:"owner"`
	Sensitivity string                 `json:"sensitivity" db:"sensitivity"`
	HasPii      bool                   `json:"has_pii" db:"has_pii"`
	Tags        []string               `json:"tags" db:"tags"`
	CreatedAt   *helperModel.Timestamp `json:"created_at" db:"created_at"`
	UpdatedAt   *helperModel.Timestamp `json:"updated_at" db:"updated_at"`
}
