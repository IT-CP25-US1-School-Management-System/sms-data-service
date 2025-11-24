package constants

import "regexp"

const (
	// Sort Datasets
	DATASET_SORT_BY_NAME       = "name"
	DATASET_SORT_BY_CREATED_AT = "created_at"
	DATASET_SORT_BY_UPDATED_AT = "updated_at"

	// Sort Order
	SORT_ORDER_ASC  = "ASC"
	SORT_ORDER_DESC = "DESC"
	// Logical Operators
	LOGICAL_OPERATOR_AND = "AND"
	LOGICAL_OPERATOR_OR  = "OR"

	// Export Job Status
	EXPORT_JOB_STATUS_PENDING   = "pending"
	EXPORT_JOB_STATUS_SUCCEEDED = "succeeded"
	EXPORT_JOB_STATUS_FAILED    = "failed"
	EXPORT_JOB_FORMAT_CSV       = "csv"
	EXPORT_JOB_FORMAT_XLSX      = "xlsx"

	DOCUMENT_PATH_REPORTING          = "reporting"
	DOCUMENT_FOLDER_EXPORT_TEMPLATES = "export/templates"
	DOCUMENT_FOLDER_EXPORT           = "exports"
)

var (
	// datasetIDPattern validates dataset ID format: only lowercase english letters, underscore, and hyphen, no spaces
	// Pattern: ^[a-z_-]+$ means start to end with only lowercase letters a-z, underscore, and hyphen
	DATASET_ID_PATTERN = regexp.MustCompile("^[a-z_-]+$")
	// Regex pattern for semantic version: v + digit + . + digit + . + digit
	DATASET_VERSION_PATTERN = regexp.MustCompile("^v\\d+\\.\\d+\\.\\d+$")
)
