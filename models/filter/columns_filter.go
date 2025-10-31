package filter

import "github.com/gofrs/uuid"

type ColumnsFilter struct {
	SourceID *uuid.UUID `query:"source_id" validate:"omitempty,uuid"`
	Schema   string     `query:"schema" validate:"omitempty"`
	Table    string     `query:"table" validate:"omitempty"`
}
