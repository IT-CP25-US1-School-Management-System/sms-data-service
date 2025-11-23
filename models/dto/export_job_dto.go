package dto

import (
	helperModel "github.com/GodeFvt/go-backend/helper/models"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/entity"
)

type ExportJobDTO struct {
	DatasetId string `json:"dataset_id" validate:"required"`
	View      string `json:"view" validate:"required"`
	Format    string `json:"format" validate:"required,oneof=csv pdf"`
	Version   string `json:"version" validate:"required"`
}

func (dto *ExportJobDTO) ExportJobDTOToEntity() *entity.ExportJob {
	return &entity.ExportJob{
		DatasetId: dto.DatasetId,
		View:      dto.View,
		Format:    dto.Format,
		Version:   dto.Version,
	}
}

type ExportJobResponseDTO struct {
	Status         string                 `json:"status"`
	DatasetId      string                 `json:"dataset_id"`
	View           string                 `json:"view"`
	Format         string                 `json:"format"`
	Version        string                 `json:"version"`
	DestinationUri string                 `json:"destination_uri" db:"destination_uri"`
	CreatedAt      *helperModel.Timestamp `json:"created_at"`
	CompletedAt    *helperModel.Timestamp `json:"completed_at"`
	ErrorMessage   string                 `json:"error_message"`
}
