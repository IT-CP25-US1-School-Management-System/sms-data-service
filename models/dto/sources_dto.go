package dto

import (
	helperModel "github.com/GodeFvt/go-backend/helper/models"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/entity"
	"github.com/gofrs/uuid"
)

type SourcesDTO struct {
	Name         string  `json:"name" validate:"required"`
	Description  *string `json:"description" validate:"omitempty"`
	Type         string  `json:"type" validate:"required,oneof=jdbc"`
	IsActive     bool    `json:"is_active" validate:"required"`
	Sensitivity  string  `json:"sensitivity" validate:"required,oneof=public internal confidential restricted"`
	DBType       string  `json:"db_type" validate:"required,oneof=postgres mysql"`
	Host         string  `json:"host" validate:"required"`
	Port         int     `json:"port"  validate:"required,min=1"`
	Username     string  `json:"username" validate:"required"`
	Password     string  `json:"password" validate:"required"`
	DatabaseName string  `json:"database_name" validate:"required"`
	Params       *string `json:"params" validate:"omitempty"`
}

func (dto *SourcesDTO) SourcesDTOToEntity() *entity.Sources {
	return &entity.Sources{
		Name:         dto.Name,
		Description:  dto.Description,
		Type:         dto.Type,
		IsActive:     dto.IsActive,
		Sensitivity:  dto.Sensitivity,
		DBType:       dto.DBType,
		Host:         dto.Host,
		Port:         dto.Port,
		Username:     dto.Username,
		Password:     dto.Password,
		DatabaseName: dto.DatabaseName,
		Params:       dto.Params,
	}
}

type UpdateSourcesDTO struct {
	Name             *string `json:"name" validate:"omitempty"`
	Description      *string `json:"description" validate:"omitempty"`
	Type             *string `json:"type" validate:"omitempty,oneof=jdbc"`
	IsActive         *bool   `json:"is_active" validate:"omitempty"`
	Sensitivity      *string `json:"sensitivity" validate:"omitempty,oneof=public internal confidential restricted"`
	DBType           *string `json:"db_type" validate:"omitempty,oneof=postgres mysql"`
	Host             *string `json:"host" validate:"omitempty"`
	Port             *int    `json:"port"  validate:"omitempty,min=1"`
	Username         *string `json:"username" validate:"omitempty"`
	DatabaseName     *string `json:"database_name" validate:"omitempty"`
	Params           *string `json:"params" validate:"omitempty"`
	IsUpdatePassword *bool   `json:"is_update_password" validate:"omitempty"`
	Password         *string `json:"password" validate:"omitempty"`
}

type SourcesResponseDTO struct {
	ID           *uuid.UUID             `json:"id"`
	Name         string                 `json:"name"`
	Description  *string                `json:"description"`
	Type         string                 `json:"type"`
	IsActive     bool                   `json:"is_active"`
	Sensitivity  string                 `json:"sensitivity"`
	DBType       string                 `json:"db_type"`
	Host         string                 `json:"host"`
	Port         int                    `json:"port"`
	Username     string                 `json:"username"`
	DatabaseName string                 `json:"database_name"`
	Params       *string                `json:"params"`
	CreatedAt    *helperModel.Timestamp `json:"created_at"`
	UpdatedAt    *helperModel.Timestamp `json:"updated_at"`
}

func SourcesEntityToSourcesResponseDTO(entity []*entity.Sources) []*SourcesResponseDTO {
	response := []*SourcesResponseDTO{}
	for _, e := range entity {
		response = append(response, &SourcesResponseDTO{
			ID:           e.ID,
			Name:         e.Name,
			Description:  e.Description,
			Type:         e.Type,
			IsActive:     e.IsActive,
			Sensitivity:  e.Sensitivity,
			DBType:       e.DBType,
			Host:         e.Host,
			Port:         e.Port,
			Username:     e.Username,
			DatabaseName: e.DatabaseName,
			Params:       e.Params,
			CreatedAt:    e.CreatedAt,
			UpdatedAt:    e.UpdatedAt,
		})
	}
	return response
}

type SourceByIDResponse struct {
	Name         string  `json:"name" validate:"required"`
	Description  *string `json:"description" validate:"omitempty"`
	Type         string  `json:"type" validate:"required,oneof=jdbc"`
	IsActive     bool    `json:"is_active" validate:"required"`
	Sensitivity  string  `json:"sensitivity" validate:"required,oneof=public internal confidential restricted"`
	DBType       string  `json:"db_type" validate:"required,oneof=postgres mysql"`
	Host         string  `json:"host" validate:"required"`
	Port         int     `json:"port"  validate:"required,min=1"`
	Username     string  `json:"username" validate:"required"`
	DatabaseName string  `json:"database_name" validate:"required"`
	Params       *string `json:"params" validate:"omitempty"`
}
