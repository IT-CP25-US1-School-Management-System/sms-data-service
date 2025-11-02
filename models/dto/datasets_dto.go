package dto

import "github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/entity"

type UpsertDatasetsDTO struct {
	ID          string   `json:"id" validate:"required"`
	Name        string   `json:"name" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Domain      string   `json:"domain" validate:"required"`
	Owner       string   `json:"owner" validate:"required"`
	Sensitivity string   `json:"sensitivity" validate:"required,oneof=public internal confidential restricted"`
	HasPii      bool     `json:"has_pii" validate:"required"`
	Tags        []string `json:"tags" validate:"required,dive,required"`
}

func (dto *UpsertDatasetsDTO) UpsertDatasetsDTOToEntity() *entity.Datasets {
	return &entity.Datasets{
		ID:          dto.ID,
		Name:        dto.Name,
		Description: dto.Description,
		Domain:      dto.Domain,
		Owner:       dto.Owner,
		Sensitivity: dto.Sensitivity,
		HasPii:      dto.HasPii,
		Tags:        dto.Tags,
	}
}
