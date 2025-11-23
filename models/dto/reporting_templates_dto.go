package dto

type ReportingTemplateDTO struct {
	Columns   []*ReportingColumnDTO `json:"columns" db:"columns"`
	Positions []*PositionDTO        `json:"positions" db:"positions"`
}

type PositionDTO struct {
	TableName   string  `json:"table_name" db:"table_name"`
	ColumnsName string  `json:"column_name" db:"column_name"`
	X           float64 `json:"x" db:"x"`
	Y           float64 `json:"y" db:"y"`
}

type ReportingColumnDTO struct {
	ColumnsName string `json:"column_name" validate:"required"`
	TableName   string `json:"table_name" validate:"required"`
}
