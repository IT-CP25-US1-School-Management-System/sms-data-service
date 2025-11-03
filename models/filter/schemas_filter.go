package filter

import "github.com/gofrs/uuid"

type SchemasFilter struct {
	SourceID *uuid.UUID `query:"source_id" validate:"omitempty,uuid"`
	Page     int        `query:"page" validate:"omitempty,min=1"`
	PerPage  int        `query:"per_page" validate:"omitempty,min=1,max=100"`
}
