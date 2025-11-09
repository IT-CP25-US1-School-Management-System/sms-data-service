package constants

import "regexp"

const (
	// Sort Datasets
	DATASET_SORT_BY_NAME       = "name"
	DATASET_SORT_BY_CREATED_AT = "created_at"
	DATASET_SORT_BY_UPDATED_AT = "updated_at"

	// Sort Order
	SORT_ORDER_ASC  = "asc"
	SORT_ORDER_DESC = "desc"
)

var (
	// datasetIDPattern validates dataset ID format: only lowercase english letters, underscore, and hyphen, no spaces
	// Pattern: ^[a-z_-]+$ means start to end with only lowercase letters a-z, underscore, and hyphen
	DATASET_ID_PATTERN = regexp.MustCompile("^[a-z_-]+$")
	// Regex pattern for semantic version: v + digit + . + digit + . + digit
	DATASET_VERSION_PATTERN = regexp.MustCompile("^v\\d+\\.\\d+\\.\\d+$")
)
