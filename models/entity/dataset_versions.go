package entity

import (
	"encoding/json"
	"fmt"

	helperModel "github.com/GodeFvt/go-backend/helper/models"
	"github.com/gofrs/uuid"
)

type DatasetVersion struct {
	DatasetID      string                 `json:"dataset_id" db:"dataset_id"`
	Version        string                 `json:"version" db:"version"`
	SourceID       *uuid.UUID             `json:"source_id" db:"source_id"`
	Status         string                 `json:"status" db:"status"`
	Schema         Schema                 `json:"schema" db:"schema"`
	AccessPolicies []AccessPolicies       `json:"access_policies" db:"access_policies"`
	Policies       Policies               `json:"policies" db:"policies"`
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
	TableName  string   `json:"table_name"`
	Alias      string   `json:"alias,omitempty"`
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
	Runtime *RuntimePolicy    `json:"runtime"`
	Views   map[string][]View `json:"views,omitempty"`
	Filters []string          `json:"filters,omitempty"`
	Write   *WritePolicy      `json:"write,omitempty"`
	Delete  *DeletePolicy     `json:"delete,omitempty"`
}
type View struct {
	TableName string   `json:"table_name"`
	Columns   []string `json:"columns"`
}

type RuntimePolicy struct {
	DefaultView string    `json:"default_view"`
	KeyField    string    `json:"key_field,omitempty"`
	Query       QueryPlan `json:"query" jsonb:"query"`
}

type QueryPlan struct {
	From        *FromRef     `json:"from"`
	Joins       []JoinRef    `json:"joins,omitempty"`
	Projections []Projection `json:"projections"`
	GroupBy     []GroupBy    `json:"group_by,omitempty"`
	WhereAllow  []WhereAllow `json:"where_allow,omitempty"`
}

type FromRef struct {
	Table string `json:"table,omitempty"`
	View  string `json:"view,omitempty"`
}

type JoinRef struct {
	Type      string    `json:"type"`
	TableFrom string    `json:"table_from"`
	TableTo   string    `json:"table_to"`
	Condition Condition `json:"condition"`
	Relation  string    `json:"relation"`
	Alias     string    `json:"alias"`
	// คอลัมน์ที่จะ SELECT เข้าไปใน JSON
	Projections []Projection `json:"projections" validate:"required,min=1,dive"`
}

type Projection struct {
	Column string `json:"column"`
	Alias  string `json:"alias"`
}

type Condition struct {
	ColumnFrom string `json:"column_from"`
	ColumnTo   string `json:"column_to"`
	Operator   string `json:"operator"`
}

type WhereAllow struct {
	TableName string   `json:"table_name"`
	Field     string   `json:"field"`
	Operators []string `json:"operators"`
}

type GroupBy struct {
	Field     string `json:"field"`
	TableName string `json:"table_name,omitempty"`
}

type WritePolicy struct {
	KeyField  string    `json:"key_field"`
	AllowEdit []string  `json:"allow_edit"`
	Query     QueryPlan `json:"query"`
}

type DeletePolicy struct {
	KeyField string    `json:"key_field"`
	Query    QueryPlan `json:"query"`
}

// Custom UnmarshalJSON for Policies to handle both formats of Views
func (p *Policies) UnmarshalJSON(data []byte) error {
	// Define a temporary struct for unmarshaling
	type Alias Policies
	aux := &struct {
		Views interface{} `json:"views,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	// Handle Views field specially
	if aux.Views != nil {
		switch v := aux.Views.(type) {
		case map[string]interface{}:
			p.Views = make(map[string][]View)
			for key, value := range v {
				switch val := value.(type) {
				case []interface{}:
					// Check if it's array of strings (simplified format)
					if len(val) > 0 {
						if _, ok := val[0].(string); ok {
							// It's array of strings - convert to View format
							columns := make([]string, len(val))
							for i, col := range val {
								if colStr, ok := col.(string); ok {
									columns[i] = colStr
								}
							}
							p.Views[key] = []View{{
								TableName: "", // Will be empty for simplified format
								Columns:   columns,
							}}
						} else {
							// It's array of View objects - unmarshal normally
							viewsJSON, err := json.Marshal(val)
							if err != nil {
								return fmt.Errorf("failed to marshal views for key %s: %w", key, err)
							}
							var views []View
							if err := json.Unmarshal(viewsJSON, &views); err != nil {
								return fmt.Errorf("failed to unmarshal views for key %s: %w", key, err)
							}
							p.Views[key] = views
						}
					}
				}
			}
		}
	}

	return nil
}
