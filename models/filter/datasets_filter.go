package filter

type DatasetsFilter struct {
	SearchWord string   `query:"search_word" validate:"omitempty"`
	Tags       []string `query:"tags" validate:"omitempty,dive,required"`
	Domain     string   `query:"domain" validate:"omitempty"`
	HasPii     *bool    `query:"has_pii" validate:"omitempty"`
	Owner      string   `query:"owner" validate:"omitempty"`
	Page       int      `query:"page" validate:"omitempty,min=1"`
	PerPage    int      `query:"per_page" validate:"omitempty,min=1,max=100"`
	SortBy     string   `query:"sort_by" validate:"omitempty,oneof=name created_at updated_at"`
	SortOrder  string   `query:"sort_order" validate:"omitempty,oneof=asc desc"`
}
