package dto

import (
	helperModel "github.com/GodeFvt/go-backend/helper/models"
)

type ReportingExportJobResponseDTO struct {
	Status           string                 `json:"status"`
	Url              string                 `json:"url"`
	OriginalFilename string                 `json:"original_filename"`
	FileSize         int64                  `json:"file_size"`
	ContentType      string                 `json:"content_type"`
	CreatedAt        *helperModel.Timestamp `json:"created_at"`
	CompletedAt      *helperModel.Timestamp `json:"completed_at"`
	ErrorMessage     string                 `json:"error_message,omitempty"`
}
