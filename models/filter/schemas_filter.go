package filter

import "github.com/gofrs/uuid"

type SchemasFilter struct {
	SourceID *uuid.UUID `query:"source_id" validate:"omitempty,uuid"`
}
