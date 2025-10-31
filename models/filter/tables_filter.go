package filter

import "github.com/gofrs/uuid"

type TablesFilter struct {
	SourceID *uuid.UUID `query:"source_id" validate:"omitempty,uuid"`
	Schema   string     `query:"schema" validate:"omitempty"`
}
