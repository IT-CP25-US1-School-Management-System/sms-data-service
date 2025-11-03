package filter

import "github.com/gofrs/uuid"

type TablesFilter struct {
	SourceID *uuid.UUID `query:"source_id" validate:"omitempty,uuid"`
	Schema   string     `query:"schema" validate:"omitempty"`
	Page     int        `query:"page" validate:"omitempty,min=1"`
	PerPage  int        `query:"per_page" validate:"omitempty,min=1,max=100"`
}
