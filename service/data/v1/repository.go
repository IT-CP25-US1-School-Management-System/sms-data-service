package data

import (
	"context"
	"database/sql"

	helperModel "github.com/GodeFvt/go-backend/helper/models"

	"github.com/GodeFvt/go-backend/helper/models"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/entity"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/filter"
	"github.com/gofrs/uuid"
)

type PsqlDataRepository interface {
	FetchInformationTablesBySourceID(ctx context.Context, dbType string, sourceID *uuid.UUID) ([]*entity.Tables, error)
	FetchInformationColumnsBySourceID(ctx context.Context, dbType string, sourceID *uuid.UUID) ([]*entity.Columns, error)
	FetchInformationSchemasBySourceID(ctx context.Context, dbType string, sourceID *uuid.UUID) ([]*entity.Schemas, error)
	FetchInformationTableRelationsBySourceID(ctx context.Context, dbType string, sourceID *uuid.UUID) ([]*entity.TableRelations, error)

	ExecuteQueryByKey(
		ctx context.Context,
		sourceID *uuid.UUID,
		schema *entity.Schema,
		policies *entity.Policies,
		key interface{},
		viewName string,
		ownerFilter *entity.OwnerFilter,
	) (map[string]interface{}, error)

	ExecuteQuery(
		ctx context.Context,
		sourceID *uuid.UUID,
		schema *entity.Schema,
		policies *entity.Policies,
		filterGroups [][]entity.FilterInput,
		logicalOperator string,
		paginator *models.Paginator,
		viewName string,
		sortBy string,
		sortOrder string,
		ownerFilter *entity.OwnerFilter,
	) ([]map[string]interface{}, error)

	// Data modification functions
	ExecuteCreate(ctx context.Context, sourceID *uuid.UUID, schema entity.Schema, writePolicy *entity.WritePolicy, data map[string]interface{}, ownerFilter *entity.OwnerFilter) (map[string]interface{}, error)
	ExecuteUpdate(ctx context.Context, sourceID *uuid.UUID, schema entity.Schema, writePolicy *entity.WritePolicy, key interface{}, data map[string]interface{}, ownerFilter *entity.OwnerFilter) (map[string]interface{}, error)
	ExecuteDelete(ctx context.Context, sourceID *uuid.UUID, deletePolicy *entity.DeletePolicy, key interface{}, ownerFilter *entity.OwnerFilter) (sql.Result, error)
	ExecuteBatchCreate(ctx context.Context, sourceID *uuid.UUID, schema entity.Schema, writePolicy *entity.WritePolicy, batchData []map[string]interface{}) (int64, error)

	// Table Data CRUD (direct table access)
	FetchTableData(ctx context.Context, sourceID *uuid.UUID, schemaName, tableName string, filterGroups [][]entity.FilterInput, logicalOperator string, paginator *models.Paginator, sortBy, sortOrder string) ([]map[string]interface{}, error)
	FetchTableDataByKey(ctx context.Context, sourceID *uuid.UUID, schemaName, tableName, keyField string, keyValue interface{}) (map[string]interface{}, error)
	CreateTableData(ctx context.Context, sourceID *uuid.UUID, schemaName, tableName string, columns []*entity.Columns, data map[string]interface{}) (map[string]interface{}, error)
	UpdateTableData(ctx context.Context, sourceID *uuid.UUID, schemaName, tableName, keyField string, keyValue interface{}, columns []*entity.Columns, data map[string]interface{}) (map[string]interface{}, error)
	DeleteTableData(ctx context.Context, sourceID *uuid.UUID, schemaName, tableName, keyField string, keyValue interface{}) (sql.Result, error)
}

type PsqlDatasetRepository interface {
	FetchSourceList(ctx context.Context, paginator *models.Paginator) ([]*entity.Sources, error)
	FetchSchemasList(ctx context.Context, filter *filter.SchemasFilter, paginator *models.Paginator) ([]*entity.Schemas, error)
	FetchTablesList(ctx context.Context, filter *filter.TablesFilter, paginator *models.Paginator) ([]*entity.Tables, error)
	FetchColumnsList(ctx context.Context, filter *filter.ColumnsFilter, paginator *models.Paginator) ([]*entity.Columns, error)
	FetchSourceByID(ctx context.Context, sourceID *uuid.UUID) (*entity.Sources, error)
	UpsertSource(ctx context.Context, source *entity.Sources) error
	ActivateSourceByID(ctx context.Context, sourceID *uuid.UUID) error
	DeactivateSourceByID(ctx context.Context, sourceID *uuid.UUID) error
	ExistSourceByID(ctx context.Context, sourceID *uuid.UUID) (bool, error)
	DeleteSourceByID(ctx context.Context, sourceID *uuid.UUID) error
	BatchInsertInformationDatabase(ctx context.Context, schemas []*entity.Schemas, tables []*entity.Tables, columns []*entity.Columns, tableRelations []*entity.TableRelations) error
	ExistSchemaByName(ctx context.Context, sourceID *uuid.UUID, schemaName string) (bool, error)
	ExistTableByName(ctx context.Context, sourceID *uuid.UUID, schemaName, tableName string) (bool, error)
	ExistColumnByName(ctx context.Context, sourceID *uuid.UUID, schemaName, tableName, columnName string) (bool, error)
	ExistSourceByName(ctx context.Context, sourceName string) (bool, error)
	ExistSourceByNameAndNotID(ctx context.Context, sourceID *uuid.UUID, sourceName string) (bool, error)

	// Dataset
	FetchDatasetList(ctx context.Context, filter *filter.DatasetsFilter, paginator *models.Paginator) ([]*entity.Datasets, error)
	FetchDatasetByID(ctx context.Context, datasetID string) (*entity.Datasets, error)
	UpsertDataset(ctx context.Context, dataset *entity.Datasets) error
	DeleteDatasetByID(ctx context.Context, datasetID string) error
	ExistDatasetByID(ctx context.Context, datasetID string) (bool, error)

	// Dataset Version
	FetchDatasetVersionByID(ctx context.Context, datasetID string, version string) (*entity.DatasetVersion, error)
	FetchDatasetVersionsList(ctx context.Context, datasetID string, filter *filter.DatasetVersionsFilter, paginator *models.Paginator) ([]*entity.DatasetVersion, error)
	UpsertDatasetVersion(ctx context.Context, datasetVersion *entity.DatasetVersion) error
	DeleteDatasetVersionByID(ctx context.Context, datasetID string, version string) error
	ExistDatasetVersionByID(ctx context.Context, datasetID string, version string) (bool, error)
	UpdateDatasetVersionStatus(ctx context.Context, datasetID string, version string, status string) error

	// Reporting Template
	FetchReportingTemplateByID(ctx context.Context, templateID *uuid.UUID) (*entity.ReportingTemplate, error)
	FetchReportingTemplatesList(ctx context.Context, filter *filter.ReportingTemplatesFilter, paginator *models.Paginator) ([]*entity.ReportingTemplate, error)
	UpsertReportingTemplate(ctx context.Context, template *entity.ReportingTemplate) error
	DeleteReportingTemplateByID(ctx context.Context, templateID *uuid.UUID) error
	ExistReportingTemplateByID(ctx context.Context, templateID *uuid.UUID) (bool, error)
	UpdateReportingExportStatusSuccess(ctx context.Context, jobID *uuid.UUID, completedAt *helperModel.Timestamp, resourceID string) error
	UpdateReportingExportStatusFail(ctx context.Context, jobID *uuid.UUID, errorMessage string) error
	UpsertReportingTemplateExportJob(ctx context.Context, job *entity.ReportingTemplateExportJob) error
	FetchReportingExportJobByID(ctx context.Context, jobID *uuid.UUID) (*entity.ReportingTemplateExportJob, error)
	InsertExportJob(ctx context.Context, exportJob *entity.ExportJob) error
	FetchExportJobByID(ctx context.Context, jobID *uuid.UUID) (*entity.ExportJob, error)
	UpdateStatusSuccess(ctx context.Context, jobId *uuid.UUID, destinationUri string, completedAt *helperModel.Timestamp) error
	UpdateStatusFail(ctx context.Context, jobId *uuid.UUID, errorMessage string) error

	// Import Job
	InsertImportJob(ctx context.Context, importJob *entity.ImportJob) error
	FetchImportJobByID(ctx context.Context, jobID *uuid.UUID) (*entity.ImportJob, error)
	UpdateImportJobStatusSuccess(ctx context.Context, jobID *uuid.UUID, completedAt *helperModel.Timestamp) error
	UpdateImportJobStatusFail(ctx context.Context, jobID *uuid.UUID, errorMessage string) error
}

type RedisRepository interface {
}
