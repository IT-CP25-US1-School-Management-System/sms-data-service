package entity

import (
	helperModel "github.com/GodeFvt/go-backend/helper/models"
	"github.com/gofrs/uuid"
)

type ReportingTemplateExportJob struct {
	JobID               *uuid.UUID             `json:"job_id" db:"job_id"`
	ReportingTemplateID *uuid.UUID             `json:"reporting_template_id" db:"reporting_template_id"`
	ResourceID          string                 `json:"resource_id" db:"resource_id"`
	Status              string                 `json:"status" db:"status"`
	CreatedAt           *helperModel.Timestamp `json:"created_at" db:"created_at"`
	CompletedAt         *helperModel.Timestamp `json:"completed_at" db:"completed_at"`
	ErrorMessage        string                 `json:"error_message" db:"error_message"`

	// Additional fields
	OriginalFilename string `json:"original_filename,omitempty"`
	FileSize         int64  `json:"file_size,omitempty"`
}

func (u *ReportingTemplateExportJob) GenUUID() {
	id, _ := uuid.NewV4()
	u.JobID = &id
}
