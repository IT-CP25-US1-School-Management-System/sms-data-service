package entity

import (
	helperModel "github.com/GodeFvt/go-backend/helper/models"
	"github.com/gofrs/uuid"
)

type ImportJob struct {
	JobID        *uuid.UUID             `json:"job_id" db:"job_id"`
	DatasetID    string                 `json:"dataset_id" db:"dataset_id"`
	Version      string                 `json:"version" db:"version"`
	Format       string                 `json:"format" db:"format"`
	Status       string                 `json:"status" db:"status"`
	ErrorMessage string                 `json:"error_message" db:"error_message"`
	CreatedAt    *helperModel.Timestamp `json:"created_at" db:"created_at"`
	CompletedAt  *helperModel.Timestamp `json:"completed_at" db:"completed_at"`
	ResourceID   string                 `json:"resource_id" db:"resource_id"`
}

func (i *ImportJob) GenUUID() {
	id, _ := uuid.NewV4()
	i.JobID = &id
}
