package filter

type ReportingTemplatesFilter struct {
	SearchWord string `query:"search_word" validate:"omitempty"`
	Page       int    `query:"page" validate:"omitempty,min=1"`
	PerPage    int    `query:"per_page" validate:"omitempty,min=1,max=100"`
}
