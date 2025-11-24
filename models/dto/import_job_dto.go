package dto

import (
	helperModel "github.com/GodeFvt/go-backend/helper/models"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/entity"
)

// CreateImportTemplateDTO for API 1: POST /v1/reporting/import/template
type CreateImportTemplateDTO struct {
	DatasetID string `json:"dataset_id" validate:"required"`
	Version   string `json:"version" validate:"required"`
	Format    string `json:"format" validate:"required,oneof=csv xlsx"`
}

// CreateImportJobDTO for API 2: POST /v1/reporting/import/job
type CreateImportJobDTO struct {
	DatasetID string `json:"dataset_id" validate:"required"`
	Version   string `json:"version" validate:"required"`
	Format    string `json:"format" validate:"required,oneof=csv xlsx"`
}

func (dto *CreateImportJobDTO) ToEntity() *entity.ImportJob {
	return &entity.ImportJob{
		DatasetID: dto.DatasetID,
		Version:   dto.Version,
		Format:    dto.Format,
	}
}

// ImportJobResponseDTO for API 3: GET /v1/reporting/import/job/:job_id
type ImportJobResponseDTO struct {
	Status       string                 `json:"status"`
	ErrorMessage string                 `json:"error_message,omitempty"`
	CreatedAt    *helperModel.Timestamp `json:"created_at"`
	CompletedAt  *helperModel.Timestamp `json:"completed_at,omitempty"`
}
