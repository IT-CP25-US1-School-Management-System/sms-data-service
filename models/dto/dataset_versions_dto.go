package dto

import (
	"fmt"

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
	DataType   string   `json:"type" validate:"required"`
	IsNullable bool     `json:"is_nullable"`
	Default    *string  `json:"default,omitempty"`
	Enum       []string `json:"enum,omitempty"`
}

type AccessPolicyDTO struct {
	Role      string   `json:"role" validate:"required"`
	Scope     string   `json:"scope" validate:"required"`
	CanView   bool     `json:"can_view"`
	CanEdit   bool     `json:"can_edit"`
	CanDelete bool     `json:"can_delete"`
	AllowView []string `json:"allow_view" validate:"required,min=1,dive"`
}

type PoliciesDTO struct {
	Runtime *RuntimePolicyDTO   `json:"runtime" validate:"omitempty"`
	Views   map[string][]string `json:"views" validate:"required,min=1"`
	Filters []string            `json:"filters" validate:"omitempty"`
	Write   *WritePolicyDTO     `json:"write,omitempty" validate:"omitempty"`
	Delete  *DeletePolicyDTO    `json:"delete,omitempty" validate:"omitempty"`
}

type RuntimePolicyDTO struct {
	Key         string       `json:"key" validate:"required"`
	DefaultView string       `json:"default_view" validate:"required"`
	KeyField    string       `json:"key_field,omitempty"`
	Query       QueryPlanDTO `json:"query" validate:"required"`
}

type QueryPlanDTO struct {
	From         *FromRefDTO     `json:"from,omitempty" validate:"omitempty"`
	Joins        []JoinRefDTO    `json:"joins,omitempty" validate:"omitempty,dive"`
	Projections  []ProjectionDTO `json:"projections,omitempty" validate:"omitempty,min=1,dive"`
	GroupBy      []ExprDTO       `json:"group_by,omitempty" validate:"omitempty,dive"`
	WhereAllow   []WhereAllowDTO `json:"where_allow,omitempty" validate:"omitempty,dive"`
	OrderAllow   []OrderAllowDTO `json:"order_allow,omitempty" validate:"omitempty,dive"`
	LimitDefault *int            `json:"limit_default,omitempty"`
}

type FromRefDTO struct {
	Table string `json:"table,omitempty"`
	View  string `json:"view,omitempty"`
}

type JoinRefDTO struct {
	Type      string  `json:"type" validate:"required"`
	Table     string  `json:"table" validate:"required"`
	Condition ExprDTO `json:"condition" validate:"required"`
}

type ProjectionDTO struct {
	Column string   `json:"column,omitempty"`
	Alias  string   `json:"alias,omitempty"`
	Expr   *ExprDTO `json:"expr,omitempty"`
}

type ExprDTO struct {
	Field    string `json:"field,omitempty"`
	Operator string `json:"operator,omitempty"`
	Value    string `json:"value,omitempty"`
}

type WhereAllowDTO struct {
	Field     string   `json:"field" validate:"required"`
	Operators []string `json:"operators" validate:"required,min=1,dive,required"`
}

type OrderAllowDTO struct {
	Field      string   `json:"field" validate:"required"`
	Directions []string `json:"directions" validate:"required,min=1,dive,required"`
}

type WritePolicyDTO struct {
	KeyField    string       `json:"key_field,omitempty"`
	DefaultView string       `json:"default_view"`
	AllowEdit   []string     `json:"allow_edit" validate:"required,min=1,dive"`
	Query       QueryPlanDTO `json:"query" validate:"required"`
}

type DeletePolicyDTO struct {
	KeyField    string       `json:"key_field" validate:"required"`
	DefaultView string       `json:"default_view"`
	Query       QueryPlanDTO `json:"query" validate:"required"`
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

	// Convert Schema
	schema := entity.Schema{
		Columns: make([]entity.Column, 0, len(dto.Schema.Columns)),
	}
	for _, col := range dto.Schema.Columns {
		schema.Columns = append(schema.Columns, entity.Column{
			Name:       col.Name,
			DataType:   col.DataType,
			IsNullable: col.IsNullable,
			Default:    col.Default,
			Enum:       col.Enum,
		})
	}

	// Convert AccessPolicies
	accessPolicies := make([]entity.AccessPolicies, 0, len(dto.AccessPolicies))
	for _, ap := range dto.AccessPolicies {
		accessPolicies = append(accessPolicies, entity.AccessPolicies{
			Role:      ap.Role,
			Scope:     ap.Scope,
			CanView:   ap.CanView,
			CanEdit:   ap.CanEdit,
			CanDelete: ap.CanDelete,
			AllowView: ap.AllowView,
		})
	}

	// Convert Policies
	policies := entity.Policies{
		Views:   dto.Policies.Views,
		Filters: dto.Policies.Filters,
	}

	// Convert Runtime Policy if exists
	if dto.Policies.Runtime != nil {
		runtime := entity.RuntimePolicy{
			Key:         dto.Policies.Runtime.Key,
			DefaultView: dto.Policies.Runtime.DefaultView,
			KeyField:    dto.Policies.Runtime.KeyField,
		}

		// Convert Query
		query := entity.QueryPlan{
			LimitDefault: dto.Policies.Runtime.Query.LimitDefault,
		}

		// Convert From
		if dto.Policies.Runtime.Query.From != nil {
			query.From = &entity.FromRef{
				Table: dto.Policies.Runtime.Query.From.Table,
				View:  dto.Policies.Runtime.Query.From.View,
			}
		}

		// Convert Joins
		if len(dto.Policies.Runtime.Query.Joins) > 0 {
			query.Joins = make([]entity.JoinRef, 0, len(dto.Policies.Runtime.Query.Joins))
			for _, join := range dto.Policies.Runtime.Query.Joins {
				query.Joins = append(query.Joins, entity.JoinRef{
					Type:  join.Type,
					Table: join.Table,
					Condition: entity.Expr{
						Field:    join.Condition.Field,
						Operator: join.Condition.Operator,
						Value:    join.Condition.Value,
					},
				})
			}
		}

		// Convert Projections
		if len(dto.Policies.Runtime.Query.Projections) > 0 {
			query.Projections = make([]entity.Projection, 0, len(dto.Policies.Runtime.Query.Projections))
			for _, proj := range dto.Policies.Runtime.Query.Projections {
				p := entity.Projection{
					Column: proj.Column,
					Alias:  proj.Alias,
				}
				if proj.Expr != nil {
					p.Expr = &entity.Expr{
						Field:    proj.Expr.Field,
						Operator: proj.Expr.Operator,
						Value:    proj.Expr.Value,
					}
				}
				query.Projections = append(query.Projections, p)
			}
		}

		// Convert GroupBy
		if len(dto.Policies.Runtime.Query.GroupBy) > 0 {
			query.GroupBy = make([]entity.Expr, 0, len(dto.Policies.Runtime.Query.GroupBy))
			for _, gb := range dto.Policies.Runtime.Query.GroupBy {
				query.GroupBy = append(query.GroupBy, entity.Expr{
					Field:    gb.Field,
					Operator: gb.Operator,
					Value:    gb.Value,
				})
			}
		}

		// Convert WhereAllow
		if len(dto.Policies.Runtime.Query.WhereAllow) > 0 {
			query.WhereAllow = make([]entity.WhereAllow, 0, len(dto.Policies.Runtime.Query.WhereAllow))
			for _, wa := range dto.Policies.Runtime.Query.WhereAllow {
				query.WhereAllow = append(query.WhereAllow, entity.WhereAllow{
					Field:     wa.Field,
					Operators: wa.Operators,
				})
			}
		}

		// Convert OrderAllow
		if len(dto.Policies.Runtime.Query.OrderAllow) > 0 {
			query.OrderAllow = make([]entity.OrderAllow, 0, len(dto.Policies.Runtime.Query.OrderAllow))
			for _, oa := range dto.Policies.Runtime.Query.OrderAllow {
				query.OrderAllow = append(query.OrderAllow, entity.OrderAllow{
					Field:      oa.Field,
					Directions: oa.Directions,
				})
			}
		}

		runtime.Query = query
		policies.Runtime = &runtime
	}

	// Convert Write Policy if exists
	if dto.Policies.Write != nil {
		writeQuery := entity.QueryPlan{
			LimitDefault: dto.Policies.Write.Query.LimitDefault,
		}
		// Convert Write Query (copy from Runtime Query conversion above)
		if dto.Policies.Write.Query.From != nil {
			writeQuery.From = &entity.FromRef{
				Table: dto.Policies.Write.Query.From.Table,
				View:  dto.Policies.Write.Query.From.View,
			}
		}
		// ... (repeat similar conversion for Joins, Projections, etc.)

		writePolicy := entity.WritePolicy{
			KeyField:    dto.Policies.Write.KeyField,
			DefaultView: dto.Policies.Write.DefaultView,
			AllowEdit:   dto.Policies.Write.AllowEdit,
			Query:       writeQuery, // แก้: เพิ่ม Query conversion
		}
		policies.Write = &writePolicy
	}

	// Convert Delete Policy if exists
	if dto.Policies.Delete != nil {
		deleteQuery := entity.QueryPlan{
			LimitDefault: dto.Policies.Delete.Query.LimitDefault,
		}
		// Convert Delete Query (copy from Runtime Query conversion above)
		if dto.Policies.Delete.Query.From != nil {
			deleteQuery.From = &entity.FromRef{
				Table: dto.Policies.Delete.Query.From.Table,
				View:  dto.Policies.Delete.Query.From.View,
			}
		}
		// ... (repeat similar conversion)

		deletePolicy := entity.DeletePolicy{
			KeyField:    dto.Policies.Delete.KeyField,
			DefaultView: dto.Policies.Delete.DefaultView,
			Query:       deleteQuery, // แก้: เพิ่ม Query conversion
		}
		policies.Delete = &deletePolicy
	}

	return &entity.DatasetVersion{
		Version:        dto.Version,
		Status:         dto.Status,
		Schema:         schema,
		AccessPolicies: accessPolicies,
		Policies:       policies,
		SourceID:       sourceID,
	}, nil
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

	// Convert Schema
	schema := entity.Schema{
		Columns: make([]entity.Column, 0, len(dto.Schema.Columns)),
	}
	for _, col := range dto.Schema.Columns {
		schema.Columns = append(schema.Columns, entity.Column{
			Name:       col.Name,
			DataType:   col.DataType,
			IsNullable: col.IsNullable,
			Default:    col.Default,
			Enum:       col.Enum,
		})
	}

	// Convert AccessPolicies
	accessPolicies := make([]entity.AccessPolicies, 0, len(dto.AccessPolicies))
	for _, ap := range dto.AccessPolicies {
		accessPolicies = append(accessPolicies, entity.AccessPolicies{
			Role:      ap.Role,
			Scope:     ap.Scope,
			CanView:   ap.CanView,
			CanEdit:   ap.CanEdit,
			CanDelete: ap.CanDelete,
			AllowView: ap.AllowView,
		})
	}

	// Convert Policies
	policies := entity.Policies{
		Views:   dto.Policies.Views,
		Filters: dto.Policies.Filters,
	}

	// Convert Runtime Policy if exists
	if dto.Policies.Runtime != nil {
		runtime := entity.RuntimePolicy{
			Key:         dto.Policies.Runtime.Key,
			DefaultView: dto.Policies.Runtime.DefaultView,
			KeyField:    dto.Policies.Runtime.KeyField,
		}

		// Convert Query
		query := entity.QueryPlan{
			LimitDefault: dto.Policies.Runtime.Query.LimitDefault,
		}

		// Convert From
		if dto.Policies.Runtime.Query.From != nil {
			query.From = &entity.FromRef{
				Table: dto.Policies.Runtime.Query.From.Table,
				View:  dto.Policies.Runtime.Query.From.View,
			}
		}

		// Convert Joins
		if len(dto.Policies.Runtime.Query.Joins) > 0 {
			query.Joins = make([]entity.JoinRef, 0, len(dto.Policies.Runtime.Query.Joins))
			for _, join := range dto.Policies.Runtime.Query.Joins {
				query.Joins = append(query.Joins, entity.JoinRef{
					Type:  join.Type,
					Table: join.Table,
					Condition: entity.Expr{
						Field:    join.Condition.Field,
						Operator: join.Condition.Operator,
						Value:    join.Condition.Value,
					},
				})
			}
		}

		// Convert Projections
		if len(dto.Policies.Runtime.Query.Projections) > 0 {
			query.Projections = make([]entity.Projection, 0, len(dto.Policies.Runtime.Query.Projections))
			for _, proj := range dto.Policies.Runtime.Query.Projections {
				p := entity.Projection{
					Column: proj.Column,
					Alias:  proj.Alias,
				}
				if proj.Expr != nil {
					p.Expr = &entity.Expr{
						Field:    proj.Expr.Field,
						Operator: proj.Expr.Operator,
						Value:    proj.Expr.Value,
					}
				}
				query.Projections = append(query.Projections, p)
			}
		}

		// Convert GroupBy
		if len(dto.Policies.Runtime.Query.GroupBy) > 0 {
			query.GroupBy = make([]entity.Expr, 0, len(dto.Policies.Runtime.Query.GroupBy))
			for _, gb := range dto.Policies.Runtime.Query.GroupBy {
				query.GroupBy = append(query.GroupBy, entity.Expr{
					Field:    gb.Field,
					Operator: gb.Operator,
					Value:    gb.Value,
				})
			}
		}

		// Convert WhereAllow
		if len(dto.Policies.Runtime.Query.WhereAllow) > 0 {
			query.WhereAllow = make([]entity.WhereAllow, 0, len(dto.Policies.Runtime.Query.WhereAllow))
			for _, wa := range dto.Policies.Runtime.Query.WhereAllow {
				query.WhereAllow = append(query.WhereAllow, entity.WhereAllow{
					Field:     wa.Field,
					Operators: wa.Operators,
				})
			}
		}

		// Convert OrderAllow
		if len(dto.Policies.Runtime.Query.OrderAllow) > 0 {
			query.OrderAllow = make([]entity.OrderAllow, 0, len(dto.Policies.Runtime.Query.OrderAllow))
			for _, oa := range dto.Policies.Runtime.Query.OrderAllow {
				query.OrderAllow = append(query.OrderAllow, entity.OrderAllow{
					Field:      oa.Field,
					Directions: oa.Directions,
				})
			}
		}

		runtime.Query = query
		policies.Runtime = &runtime
	}

	// Convert Write Policy if exists
	if dto.Policies.Write != nil {
		writeQuery := entity.QueryPlan{
			LimitDefault: dto.Policies.Write.Query.LimitDefault,
		}
		// Convert Write Query (copy from Runtime Query conversion above)
		if dto.Policies.Write.Query.From != nil {
			writeQuery.From = &entity.FromRef{
				Table: dto.Policies.Write.Query.From.Table,
				View:  dto.Policies.Write.Query.From.View,
			}
		}
		// ... (repeat similar conversion for Joins, Projections, etc.)

		writePolicy := entity.WritePolicy{
			KeyField:    dto.Policies.Write.KeyField,
			DefaultView: dto.Policies.Write.DefaultView,
			AllowEdit:   dto.Policies.Write.AllowEdit,
			Query:       writeQuery, // แก้: เพิ่ม Query conversion
		}
		policies.Write = &writePolicy
	}

	// Convert Delete Policy if exists
	if dto.Policies.Delete != nil {
		deleteQuery := entity.QueryPlan{
			LimitDefault: dto.Policies.Delete.Query.LimitDefault,
		}
		// Convert Delete Query (copy from Runtime Query conversion above)
		if dto.Policies.Delete.Query.From != nil {
			deleteQuery.From = &entity.FromRef{
				Table: dto.Policies.Delete.Query.From.Table,
				View:  dto.Policies.Delete.Query.From.View,
			}
		}
		// ... (repeat similar conversion)

		deletePolicy := entity.DeletePolicy{
			KeyField:    dto.Policies.Delete.KeyField,
			DefaultView: dto.Policies.Delete.DefaultView,
			Query:       deleteQuery, // แก้: เพิ่ม Query conversion
		}
		policies.Delete = &deletePolicy
	}

	return &entity.DatasetVersion{
		Status:         dto.Status,
		Schema:         schema,
		AccessPolicies: accessPolicies,
		Policies:       policies,
		SourceID:       sourceID,
	}, nil
}
