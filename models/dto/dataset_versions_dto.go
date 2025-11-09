package dto

import (
	"fmt"

	helperModel "github.com/GodeFvt/go-backend/helper/models"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/entity"

	"github.com/gofrs/uuid"
)

// DTO used for inserting a new dataset version.
type InsertDatasetVersionDTO struct {
	Version        string            `json:"version" validate:"required"`
	SourceID       string            `json:"source_id" validate:"required,uuid"`
	Status         string            `json:"status" validate:"required,oneof=active preview deprecated"`
	Schema         SchemaDTO         `json:"schema" validate:"required"`
	AccessPolicies []AccessPolicyDTO `json:"access_policies" validate:"required,min=1,dive"`
	Policies       PoliciesDTO       `json:"policies" validate:"required"`
}

type UpdateDatasetVersionDTO struct {
	SourceID       string            `json:"source_id" validate:"required,uuid"`
	Status         string            `json:"status" validate:"required,oneof=active preview deprecated"`
	Schema         SchemaDTO         `json:"schema" validate:"required"`
	AccessPolicies []AccessPolicyDTO `json:"access_policies" validate:"required,min=1,dive"`
	Policies       PoliciesDTO       `json:"policies" validate:"required"`
}

type UpdateDatasetVersionStatusDTO struct {
	Status string `json:"status" validate:"required,oneof=active preview deprecated"`
}

type SchemaDTO struct {
	Columns []ColumnDTO `json:"columns" validate:"required,min=1,dive"`
}

type ColumnDTO struct {
	Name       string   `json:"name" validate:"required"`
	TableName  string   `json:"table_name" validate:"required"`
	Alias      string   `json:"alias,omitempty"`
	DataType   string   `json:"type" validate:"required"`
	IsNullable bool     `json:"is_nullable"`
	Default    *string  `json:"default" validate:"omitempty"`
	Enum       []string `json:"enum,omitempty"`
}

type AccessPolicyDTO struct {
	Role      string   `json:"role" validate:"required"`
	Scope     string   `json:"scope" validate:"required"`
	CanView   bool     `json:"can_view"`
	CanEdit   bool     `json:"can_edit"`
	CanDelete bool     `json:"can_delete"`
	AllowView []string `json:"allow_view" validate:"required,min=1"`
}

type PoliciesDTO struct {
	Runtime *RuntimePolicyDTO    `json:"runtime" validate:"omitempty"`
	Views   map[string][]ViewDTO `json:"views" validate:"required,min=1,dive"`
	Write   *WritePolicyDTO      `json:"write" validate:"omitempty"`
	Delete  *DeletePolicyDTO     `json:"delete" validate:"omitempty"`
}

type ViewDTO struct {
	TableName string   `json:"table_name" validate:"required"`
	Columns   []string `json:"columns" validate:"required,min=1"`
}

type RuntimePolicyDTO struct {
	DefaultView string       `json:"default_view" validate:"required"`
	KeyField    string       `json:"key_field" validate:"required"`
	Query       QueryPlanDTO `json:"query" validate:"required"`
}

type QueryPlanDTO struct {
	From        FromDTO         `json:"from" validate:"required"`
	Joins       []JoinRefDTO    `json:"joins" validate:"omitempty,dive"`
	Projections []ProjectionDTO `json:"projections" validate:"required,min=1,dive"`
	GroupBy     []GroupByDTO    `json:"group_by" validate:"omitempty,dive"`
	WhereAllow  []WhereAllowDTO `json:"where_allow" validate:"omitempty,dive"`
}

type FromDTO struct {
	Table string `json:"table" validate:"required"`
}

type JoinRefDTO struct {
	Type      string       `json:"type" validate:"required"`
	TableFrom string       `json:"table_from" validate:"required"`
	TableTo   string       `json:"table_to" validate:"required"`
	Condition ConditionDTO `json:"condition" validate:"required"`
	Relation  string       `json:"relation" validate:"required,oneof=one_to_one many_to_one one_to_many"`
	Alias     string       `json:"alias" validate:"required"`

	// คอลัมน์ที่จะ SELECT เข้าไปใน JSON
	Projections []ProjectionDTO `json:"projections" validate:"required,min=1,dive"`
}

type ProjectionDTO struct {
	Column string `json:"column" validate:"required"`
	Alias  string `json:"alias" validate:"required"`
}

type ConditionDTO struct {
	ColumnFrom string `json:"column_from" validate:"required"`
	ColumnTo   string `json:"column_to" validate:"required"`
	Operator   string `json:"operator" validate:"required"`
}

type WhereAllowDTO struct {
	TableName string   `json:"table_name" validate:"required"`
	Field     string   `json:"field" validate:"required"`
	Operators []string `json:"operators" validate:"required,min=1"`
}

type GroupByDTO struct {
	Field     string `json:"field" validate:"required"`
	TableName string `json:"table_name" validate:"omitempty"`
}

type WritePolicyDTO struct {
	KeyField  string       `json:"key_field" validate:"required"`
	AllowEdit []string     `json:"allow_edit" validate:"required,min=1"`
	Query     QueryPlanDTO `json:"query" validate:"required"`
}

type DeletePolicyDTO struct {
	KeyField string       `json:"key_field" validate:"required"`
	Query    QueryPlanDTO `json:"query" validate:"required"`
}

func (dto *InsertDatasetVersionDTO) InsertDatasetVersionDTOToEntity() (*entity.DatasetVersion, error) {
	// Convert source_id to UUID
	var sourceID *uuid.UUID
	if dto.SourceID != "" {
		id, err := uuid.FromString(dto.SourceID)
		if err != nil {
			return nil, fmt.Errorf("invalid source_id: %w", err)
		}
		sourceID = &id
	}

	entityResult, err := helperModel.ConvertStruct[InsertDatasetVersionDTO, entity.DatasetVersion](*dto)
	if err != nil {
		return nil, fmt.Errorf("failed to convert DTO to entity: %w", err)
	}

	entityResult.SourceID = sourceID

	return &entityResult, nil
}

func (dto *UpdateDatasetVersionDTO) UpdateDatasetVersionDTOToEntity() (*entity.DatasetVersion, error) {
	// Convert source_id to UUID
	var sourceID *uuid.UUID
	if dto.SourceID != "" {
		id, err := uuid.FromString(dto.SourceID)
		if err != nil {
			return nil, fmt.Errorf("invalid source_id: %w", err)
		}
		sourceID = &id
	}

	entityResult, err := helperModel.ConvertStruct[UpdateDatasetVersionDTO, entity.DatasetVersion](*dto)
	if err != nil {
		return nil, fmt.Errorf("failed to convert DTO to entity: %w", err)
	}

	entityResult.SourceID = sourceID

	return &entityResult, nil
}
