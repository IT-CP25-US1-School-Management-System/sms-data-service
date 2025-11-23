package entity

import (
	helperModel "github.com/GodeFvt/go-backend/helper/models"
	"github.com/gofrs/uuid"
)

type ExportJob struct {
	JobId           *uuid.UUID              `json:"job_id" db:"job_id"`
	DatasetId        string                 `json:"dataset_id" db:"dataset_id"`
	View             string                 `json:"view" db:"view"`
	Format           string                 `json:"format" db:"format"`
	Version          string                 `json:"version" db:"version"`
	DestinationUri   string                 `json:"destination_uri" db:"destination_uri"`
	Status           string                 `json:"status" db:"status"`
	CreatedAt       *helperModel.Timestamp  `json:"created_at" db:"created_at"`
	CompletedAt     *helperModel.Timestamp  `json:"completed_at" db:"completed_at"`
	ErrorMessage     string                 `json:"error_message" db:"error_message"`

	// Additional fields
	OriginalFilename string               `json:"original_filename,omitempty"`
	FileSize	   int64                  `json:"file_size,omitempty"`

}

func (u *ExportJob) GenUUID() {
	id, _ := uuid.NewV4()
	u.JobId = &id
}