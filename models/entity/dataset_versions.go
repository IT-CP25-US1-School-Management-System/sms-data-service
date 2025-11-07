package entity

import (
	helperModel "github.com/GodeFvt/go-backend/helper/models"
	"github.com/gofrs/uuid"
)

type DatasetVersion struct {
	DatasetID      string                 `json:"dataset_id"`
	Version        string                 `json:"version"`
	SourceID       *uuid.UUID             `json:"source_id" db:"source_id"`
	Status         string                 `json:"status"`
	Schema         Schema                 `json:"schema"`
	AccessPolicies []AccessPolicies       `json:"access_policies"`
	Policies       Policies               `json:"policies"`
	CreatedAt      *helperModel.Timestamp `json:"created_at" db:"created_at"`
	UpdatedAt      *helperModel.Timestamp `json:"updated_at" db:"updated_at"`
}

type Schema struct {
	Columns []Column `json:"columns"`
}

type Column struct {
	Name       string   `json:"name"`
	DataType   string   `json:"type"`
	IsNullable bool     `json:"is_nullable"`
	Default    *string  `json:"default,omitempty"`
	Enum       []string `json:"enum,omitempty"`
}

type AccessPolicies struct {
	Role      string   `json:"role"`
	Scope     string   `json:"scope"`
	CanView   bool     `json:"can_view"`
	CanEdit   bool     `json:"can_edit"`
	CanDelete bool     `json:"can_delete"`
	AllowView []string `json:"allow_view"`
}

type Policies struct {
	Runtime RuntimePolicy       `json:"runtime"`
	Views   map[string][]string `json:"views,omitempty"`
	Filters []string            `json:"filters,omitempty"`
	Write   *WritePolicy        `json:"write,omitempty"`
	Delete  *DeletePolicy       `json:"delete,omitempty"`
}

// RuntimePolicy defines the runtime access policy for a dataset version
type RuntimePolicy struct {
	Key         string    `json:"key" jsonb:"key"`
	DefaultView string    `json:"default_view" jsonb:"default_view"`
	KeyField    string    `json:"key_field,omitempty" jsonb:"key_field"`
	Query       QueryPlan `json:"query" jsonb:"query"`
}

type QueryPlan struct {
	From         *FromRef     `json:"from,omitempty"`
	Joins        []JoinRef    `json:"joins,omitempty"`
	Projections  []Projection `json:"projections,omitempty"`
	GroupBy      []Expr       `json:"group_by,omitempty"`
	WhereAllow   []WhereAllow `json:"where_allow,omitempty"`
	OrderAllow   []OrderAllow `json:"order_allow,omitempty"`
	LimitDefault *int         `json:"limit_default,omitempty"`
}

type FromRef struct {
	Table string `json:"table,omitempty"`
	View  string `json:"view,omitempty"`
}

type JoinRef struct {
	Type      string `json:"type"`
	Table     string `json:"table"`
	Condition Expr   `json:"condition"`
}

type Projection struct {
	Column string `json:"column,omitempty"`
	Alias  string `json:"alias,omitempty"`
	Expr   *Expr  `json:"expr,omitempty"`
}

type Expr struct {
	Field    string `json:"field,omitempty"`
	Operator string `json:"operator,omitempty"`
	Value    string `json:"value,omitempty"`
}

type WhereAllow struct {
	Field     string   `json:"field"`
	Operators []string `json:"operators"`
}

type OrderAllow struct {
	Field      string   `json:"field"`
	Directions []string `json:"directions"`
}

type WritePolicy struct {
	KeyField    string    `json:"key_field,omitempty" jsonb:"key_field"`
	DefaultView string    `json:"default_view"`
	AllowEdit   []string  `json:"allow_edit"`
	Query       QueryPlan `json:"query"`
}

type DeletePolicy struct {
	KeyField    string    `json:"key_field,omitempty" jsonb:"key_field"`
	DefaultView string    `json:"default_view"`
	Query       QueryPlan `json:"query"`
}
