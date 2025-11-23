package data

import (
	"context"

	helperModel "github.com/GodeFvt/go-backend/helper/models"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/dto"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/entity"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/filter"
	"github.com/gofrs/uuid"
)

type DataUsecase interface {
	FetchSourceList(ctx context.Context, paginator *helperModel.Paginator) ([]*entity.Sources, error)
	FetchSourceByID(ctx context.Context, sourceID *uuid.UUID) (*entity.Sources, error)
	FetchSchemasList(ctx context.Context, filter *filter.SchemasFilter, paginator *helperModel.Paginator) ([]*entity.Schemas, error)
	FetchTablesList(ctx context.Context, filter *filter.TablesFilter, paginator *helperModel.Paginator) ([]*entity.Tables, error)
	FetchColumnsList(ctx context.Context, filter *filter.ColumnsFilter, paginator *helperModel.Paginator) ([]*entity.Columns, error)
	InsertSource(ctx context.Context, source *entity.Sources) error
	UpdateSource(ctx context.Context, sourceID *uuid.UUID, sourceUpdate *dto.UpdateSourcesDTO) error
	ActivateSourceByID(ctx context.Context, sourceID *uuid.UUID) error
	DeactivateSourceByID(ctx context.Context, sourceID *uuid.UUID) error
	DeleteSourceByID(ctx context.Context, sourceID *uuid.UUID) error

	// Dataset
	FetchDatasetList(ctx context.Context, filter *filter.DatasetsFilter, paginator *helperModel.Paginator) ([]*entity.Datasets, error)
	FetchDatasetByID(ctx context.Context, datasetID string) (*entity.Datasets, error)
	UpsertDataset(ctx context.Context, dataset *entity.Datasets) error
	DeleteDatasetByID(ctx context.Context, datasetID string) error

	// DatasetVersion methods
	FetchDatasetVersionByID(ctx context.Context, datasetID, version string) (*entity.DatasetVersion, error)
	FetchDatasetVersionsList(ctx context.Context, datasetID string, filter *filter.DatasetVersionsFilter, paginator *helperModel.Paginator) ([]*entity.DatasetVersion, error)
	InsertDatasetVersion(ctx context.Context, datasetVersion *entity.DatasetVersion, datasetID string) error
	UpdateDatasetVersion(ctx context.Context, datasetVersion *entity.DatasetVersion, datasetID, version string) error
	UpdateDatasetVersionStatus(ctx context.Context, datasetID, version, status string) error

	// Serving methods
	ServingDatasetVersionData(
		ctx context.Context,
		datasetID string,
		version string,
		paginator *helperModel.Paginator,
		viewName string,
		filterGroups [][]entity.FilterInput,
		logicalOperator string,
		sortBy string,
		sortOrder string,
	) ([]map[string]interface{}, error)
	ServingDatasetVersionDataByKey(
		ctx context.Context,
		datasetID, version, key, viewName string,
	) (map[string]interface{}, error)

	// Data Modification methods (requires write policies)
	CreateDatasetVersionData(ctx context.Context, datasetID, version string, data map[string]interface{}) (map[string]interface{}, error)
	UpdateDatasetVersionDataByKey(ctx context.Context, datasetID, version, key string, data map[string]interface{}) (map[string]interface{}, error)
	DeleteDatasetVersionDataByKey(ctx context.Context, datasetID, version, key string) error

	// Table Data CRUD (direct source access)
	FetchTableData(ctx context.Context, sourceID *uuid.UUID, schemaName, tableName string, filterGroups [][]entity.FilterInput, logicalOperator string, paginator *helperModel.Paginator, sortBy, sortOrder string) ([]map[string]interface{}, error)
	FetchTableDataByKey(ctx context.Context, sourceID *uuid.UUID, schemaName, tableName, keyField string, keyValue interface{}) (map[string]interface{}, error)
	CreateTableData(ctx context.Context, sourceID *uuid.UUID, schemaName, tableName string, data map[string]interface{}) (map[string]interface{}, error)
	UpdateTableData(ctx context.Context, sourceID *uuid.UUID, schemaName, tableName, keyField string, keyValue interface{}, data map[string]interface{}) (map[string]interface{}, error)
	DeleteTableData(ctx context.Context, sourceID *uuid.UUID, schemaName, tableName, keyField string, keyValue interface{}) error

	// Reporting Template methods
	UploadReportingTemplate(ctx context.Context, template *entity.ReportingTemplate, fileData []byte, fileName string) error
	InsertExportJob(ctx context.Context, exportJob *entity.ExportJob) error
	FetchExportJobByID(ctx context.Context, jobID *uuid.UUID) (*entity.ExportJob, error)
}
