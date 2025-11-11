package filter

type DatasetVersionsFilter struct {
	SourceID   string `query:"source_id" validate:"omitempty"`
	SearchWord string `query:"search_word" validate:"omitempty"`
	Status     string `query:"status" validate:"omitempty,oneof=active preview deprecated"`
	Page       int    `query:"page" validate:"omitempty,min=1"`
	PerPage    int    `query:"per_page" validate:"omitempty,min=1,max=100"`
}
