package filter

import (
	"encoding/json"
	"fmt"

	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/entity"
)

type ServingDataFilter struct {
	Page           int    `query:"page" validate:"omitempty,min=1"`
	PerPage        int    `query:"per_page" validate:"omitempty,min=1,max=100"`
	View           string `query:"view" validate:"omitempty"`
	WhereLogicalOp string `query:"where_logical_operator" validate:"omitempty,oneof=AND OR and or"`
	WhereRaw       string `query:"where" validate:"omitempty"`
	SortBy         string `query:"sort_by" validate:"omitempty"`
	SortOrder      string `query:"sort_order" validate:"omitempty,oneof=ASC DESC asc desc"`
}

// ParseWhere parses the raw where string into FilterInput slice
func (f *ServingDataFilter) ParseWhere() ([][]entity.FilterInput, error) {
	if f.WhereRaw == "" || f.WhereRaw == "[]" {
		return make([][]entity.FilterInput, 0), nil
	}

	var filterGroups [][]entity.FilterInput
	if err := json.Unmarshal([]byte(f.WhereRaw), &filterGroups); err != nil {
		return nil, fmt.Errorf("invalid 'where' JSON structure: %w", err)
	}

	return filterGroups, nil
}
