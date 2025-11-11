package entity

// FilterInput สำหรับส่งเงื่อนไข (Filter) เข้าไปใน ExecuteQuery ไม่มีใน database
type FilterInput struct {
	TableName string      `json:"table_name"` // e.g., "person_data"
	Field     string      `json:"field"`      // e.g., "gender_id"
	Operator  string      `json:"operator"`   // e.g., "=", "!=", "IN", "LIKE"
	Value     interface{} `json:"value"`      // e.g., "ชาย" or []string{"ชาย", "หญิง"}
}
