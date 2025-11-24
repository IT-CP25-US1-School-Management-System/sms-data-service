package usecase

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"

	helperModel "github.com/GodeFvt/go-backend/helper/models"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/constants"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/errs"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/dto"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/entity"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/filter"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/proto/proto_models"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/data/v1"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/document/v1"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/utils/crypto"
	"github.com/gofrs/uuid"
	"github.com/jung-kurt/gofpdf"
	"github.com/jung-kurt/gofpdf/contrib/gofpdi"
	"github.com/xuri/excelize/v2"
)

type dataUsecase struct {
	dataRepo     data.PsqlDataRepository
	datasetRepo  data.PsqlDatasetRepository
	documentRepo document.GrpcDocumentRepository
	redisRepo    data.RedisRepository
	cryptoSecret string
}

func NewDataUsecase(dataRepo data.PsqlDataRepository, datasetRepo data.PsqlDatasetRepository, documentRepo document.GrpcDocumentRepository, redisRepo data.RedisRepository, cryptoSecret string) data.DataUsecase {
	return &dataUsecase{
		dataRepo:     dataRepo,
		datasetRepo:  datasetRepo,
		documentRepo: documentRepo,
		redisRepo:    redisRepo,
		cryptoSecret: cryptoSecret,
	}
}

// FetchSourceByID implements data.DataUsecase.
func (d *dataUsecase) FetchSourceByID(ctx context.Context, sourceID *uuid.UUID) (*entity.Sources, error) {
	return d.datasetRepo.FetchSourceByID(ctx, sourceID)
}

// DeleteSourceByID implements data.DataUsecase.
func (d *dataUsecase) DeleteSourceByID(ctx context.Context, sourceID *uuid.UUID) error {
	exist, err := d.datasetRepo.ExistSourceByID(ctx, sourceID)
	if err != nil {
		return err
	}
	if !exist {
		return errs.NewNotFoundError(constants.ERR_SOURCE_NOT_FOUND)
	}

	return d.datasetRepo.DeleteSourceByID(ctx, sourceID)
}

// ActivateSourceByID implements data.DataUsecase.
func (d *dataUsecase) ActivateSourceByID(ctx context.Context, sourceID *uuid.UUID) error {
	exist, err := d.datasetRepo.ExistSourceByID(ctx, sourceID)
	if err != nil {
		return err
	}
	if !exist {
		return errs.NewNotFoundError(constants.ERR_SOURCE_NOT_FOUND)
	}

	return d.datasetRepo.ActivateSourceByID(ctx, sourceID)
}

// DeactivateSourceByID implements data.DataUsecase.
func (d *dataUsecase) DeactivateSourceByID(ctx context.Context, sourceID *uuid.UUID) error {
	exist, err := d.datasetRepo.ExistSourceByID(ctx, sourceID)
	if err != nil {
		return err
	}
	if !exist {
		return errs.NewNotFoundError(constants.ERR_SOURCE_NOT_FOUND)
	}

	return d.datasetRepo.DeactivateSourceByID(ctx, sourceID)
}

// InsertSource implements data.DataUsecase.
func (d *dataUsecase) InsertSource(ctx context.Context, source *entity.Sources) error {
	if source == nil {
		return errs.NewBadRequestError(constants.ERR_INVALID_REQUEST_BODY)
	}
	encryptPass, err := crypto.Encrypt(source.Password, d.cryptoSecret)
	if err != nil {
		return err
	}
	source.Password = encryptPass

	now := helperModel.NewTimestampFromNow()
	source.CreatedAt = &now
	source.UpdatedAt = &now

	err = d.datasetRepo.UpsertSource(ctx, source)
	if err != nil {
		return err
	}

	infoSchema, err := d.dataRepo.FetchInformationSchemasBySourceID(ctx, source.DBType, source.ID)
	if err != nil {
		d.datasetRepo.DeleteSourceByID(ctx, source.ID)
		return err
	}
	infoTables, err := d.dataRepo.FetchInformationTablesBySourceID(ctx, source.DBType, source.ID)
	if err != nil {
		d.datasetRepo.DeleteSourceByID(ctx, source.ID)
		return err
	}
	infoColumns, err := d.dataRepo.FetchInformationColumnsBySourceID(ctx, source.DBType, source.ID)
	if err != nil {
		d.datasetRepo.DeleteSourceByID(ctx, source.ID)
		return err
	}
	infoTableRelations, err := d.dataRepo.FetchInformationTableRelationsBySourceID(ctx, source.DBType, source.ID)
	if err != nil {
		d.datasetRepo.DeleteSourceByID(ctx, source.ID)
		return err
	}
	err = d.datasetRepo.BatchInsertInformationDatabase(ctx, infoSchema, infoTables, infoColumns, infoTableRelations)
	if err != nil {
		d.datasetRepo.DeleteSourceByID(ctx, source.ID)
		return err
	}

	return nil
}

// UpdateSource implements data.DataUsecase.
func (d *dataUsecase) UpdateSource(ctx context.Context, sourceID *uuid.UUID, sourceUpdate *dto.UpdateSourcesDTO) error {
	if sourceUpdate == nil {
		return errs.NewBadRequestError(constants.ERR_INVALID_REQUEST_BODY)
	}
	oldSource, err := d.datasetRepo.FetchSourceByID(ctx, sourceID)
	if err != nil {
		return err
	}
	if oldSource == nil {
		return errs.NewNotFoundError(constants.ERR_SOURCE_NOT_FOUND)
	}

	oldSource.ID = sourceID
	if sourceUpdate.Name != nil {
		oldSource.Name = *sourceUpdate.Name
	}
	if sourceUpdate.Description != nil {
		oldSource.Description = sourceUpdate.Description
	}
	if sourceUpdate.Type != nil {
		oldSource.Type = *sourceUpdate.Type
	}
	if sourceUpdate.IsActive != nil {
		oldSource.IsActive = *sourceUpdate.IsActive
	}
	if sourceUpdate.Sensitivity != nil {
		oldSource.Sensitivity = *sourceUpdate.Sensitivity
	}
	if sourceUpdate.DBType != nil {
		oldSource.DBType = *sourceUpdate.DBType
	}
	if sourceUpdate.Host != nil {
		oldSource.Host = *sourceUpdate.Host
	}
	if sourceUpdate.Port != nil {
		oldSource.Port = *sourceUpdate.Port
	}
	if sourceUpdate.Username != nil {
		oldSource.Username = *sourceUpdate.Username
	}
	if sourceUpdate.DatabaseName != nil {
		oldSource.DatabaseName = *sourceUpdate.DatabaseName
	}
	if sourceUpdate.Params != nil {
		oldSource.Params = sourceUpdate.Params
	}
	if sourceUpdate.IsUpdatePassword != nil && *sourceUpdate.IsUpdatePassword && sourceUpdate.Password != nil {
		encryptPass, err := crypto.Encrypt(*sourceUpdate.Password, d.cryptoSecret)
		if err != nil {
			return err
		}
		oldSource.Password = encryptPass
	}

	now := helperModel.NewTimestampFromNow()
	oldSource.UpdatedAt = &now

	return d.datasetRepo.UpsertSource(ctx, oldSource)
}

// validateDatasetID validates the dataset ID format
func (d *dataUsecase) validateDatasetID(id string) error {
	if id == "" {
		return errs.NewBadRequestError(constants.ERR_DATASET_ID_IS_REQUIRED)
	}

	if !constants.DATASET_ID_PATTERN.MatchString(id) {
		return errs.NewBadRequestError(constants.ERR_DATASET_ID_INVALID_FORMAT)
	}

	return nil
}

func (d *dataUsecase) validateVersionFormat(version string) error {
	if version == "" {
		return errs.NewBadRequestError(constants.ERR_DATASET_VERSION_IS_REQUIRED)
	}

	if !constants.DATASET_VERSION_PATTERN.MatchString(version) {
		return errs.NewBadRequestError(constants.ERR_DATASET_VERSION_INVALID_FORMAT)
	}

	return nil
}

// DeleteDatasetByID implements data.DataUsecase.
func (d *dataUsecase) DeleteDatasetByID(ctx context.Context, datasetID string) error {
	if err := d.validateDatasetID(datasetID); err != nil {
		return err
	}
	exist, err := d.datasetRepo.ExistDatasetByID(ctx, datasetID)
	if err != nil {
		return err
	}
	if !exist {
		return errs.NewNotFoundError(constants.ERR_DATASET_NOT_FOUND)
	}

	return d.datasetRepo.DeleteDatasetByID(ctx, datasetID)
}

// UpsertDataset implements data.DataUsecase.
func (d *dataUsecase) UpsertDataset(ctx context.Context, dataset *entity.Datasets) error {
	if dataset == nil {
		return errs.NewBadRequestError(constants.ERR_INVALID_REQUEST_BODY)
	}

	if err := d.validateDatasetID(dataset.ID); err != nil {
		return err
	}

	now := helperModel.NewTimestampFromNow()
	dataset.CreatedAt = &now
	dataset.UpdatedAt = &now

	return d.datasetRepo.UpsertDataset(ctx, dataset)
}

// FetchDatasetByID implements data.DataUsecase.
func (d *dataUsecase) FetchDatasetByID(ctx context.Context, datasetID string) (*entity.Datasets, error) {
	if err := d.validateDatasetID(datasetID); err != nil {
		return nil, err
	}
	return d.datasetRepo.FetchDatasetByID(ctx, datasetID)
}

// FetchColumnsList implements data.DataUsecase.
func (d *dataUsecase) FetchColumnsList(ctx context.Context, filter *filter.ColumnsFilter, paginator *helperModel.Paginator) ([]*entity.Columns, error) {
	return d.datasetRepo.FetchColumnsList(ctx, filter, paginator)
}

// FetchTablesList implements data.DataUsecase.
func (d *dataUsecase) FetchTablesList(ctx context.Context, filter *filter.TablesFilter, paginator *helperModel.Paginator) ([]*entity.Tables, error) {
	return d.datasetRepo.FetchTablesList(ctx, filter, paginator)
}

// FetchDatasetList implements data.DataUsecase.
func (d *dataUsecase) FetchDatasetList(ctx context.Context, filter *filter.DatasetsFilter, paginator *helperModel.Paginator) ([]*entity.Datasets, error) {
	return d.datasetRepo.FetchDatasetList(ctx, filter, paginator)
}

// FetchSchemasList implements data.DataUsecase.
func (d *dataUsecase) FetchSchemasList(ctx context.Context, filter *filter.SchemasFilter, paginator *helperModel.Paginator) ([]*entity.Schemas, error) {
	return d.datasetRepo.FetchSchemasList(ctx, filter, paginator)
}

// FetchSourceList implements data.DataUsecase.
func (d *dataUsecase) FetchSourceList(ctx context.Context, paginator *helperModel.Paginator) ([]*entity.Sources, error) {
	return d.datasetRepo.FetchSourceList(ctx, paginator)
}

// FetchDatasetVersionByID implements data.DataUsecase.
func (d *dataUsecase) FetchDatasetVersionByID(ctx context.Context, datasetID string, version string) (*entity.DatasetVersion, error) {
	if err := d.validateDatasetID(datasetID); err != nil {
		return nil, err
	}
	if err := d.validateVersionFormat(version); err != nil {
		return nil, err
	}
	return d.datasetRepo.FetchDatasetVersionByID(ctx, datasetID, version)
}

// FetchDatasetVersionsList implements data.DataUsecase.
func (d *dataUsecase) FetchDatasetVersionsList(ctx context.Context, datasetID string, filter *filter.DatasetVersionsFilter, paginator *helperModel.Paginator) ([]*entity.DatasetVersion, error) {
	if err := d.validateDatasetID(datasetID); err != nil {
		return nil, err
	}
	exist, err := d.datasetRepo.ExistDatasetByID(ctx, datasetID)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errs.NewNotFoundError(constants.ERR_DATASET_NOT_FOUND)
	}
	return d.datasetRepo.FetchDatasetVersionsList(ctx, datasetID, filter, paginator)
}

// UpsertDatasetVersion implements data.DataUsecase.
func (d *dataUsecase) UpsertDatasetVersion(ctx context.Context, datasetVersion *entity.DatasetVersion) error {
	if datasetVersion == nil {
		return errs.NewBadRequestError(constants.ERR_INVALID_REQUEST_BODY)
	}

	if err := d.validateDatasetID(datasetVersion.DatasetID); err != nil {
		return err
	}

	// Validate that the parent dataset exists
	exists, err := d.datasetRepo.ExistDatasetByID(ctx, datasetVersion.DatasetID)
	if err != nil {
		return err
	}
	if !exists {
		return errs.NewNotFoundError(constants.ERR_DATASET_NOT_FOUND)
	}
	now := helperModel.NewTimestampFromNow()
	datasetVersion.CreatedAt = &now
	datasetVersion.UpdatedAt = &now

	return d.datasetRepo.UpsertDatasetVersion(ctx, datasetVersion)
}

func (d *dataUsecase) InsertDatasetVersion(ctx context.Context, datasetVersion *entity.DatasetVersion, datasetID string) error {
	if err := d.validateDatasetID(datasetID); err != nil {
		return err
	}
	if err := d.validateVersionFormat(datasetVersion.Version); err != nil {
		return err
	}
	exists, err := d.datasetRepo.ExistDatasetByID(ctx, datasetID)
	if err != nil {
		return err
	}
	if !exists {
		return errs.NewNotFoundError(constants.ERR_DATASET_NOT_FOUND)
	}
	exists, err = d.datasetRepo.ExistDatasetVersionByID(ctx, datasetID, datasetVersion.Version)
	if err != nil {
		return err
	}
	if exists {
		return errs.NewConflictError(constants.ERR_DATASET_VERSION_ALREADY_EXISTS)
	}

	datasetVersion.DatasetID = datasetID
	now := helperModel.NewTimestampFromNow()
	datasetVersion.CreatedAt = &now
	datasetVersion.UpdatedAt = &now
	return d.datasetRepo.UpsertDatasetVersion(ctx, datasetVersion)
}

func (d *dataUsecase) UpdateDatasetVersion(ctx context.Context, datasetVersion *entity.DatasetVersion, datasetID, version string) error {
	if err := d.validateDatasetID(datasetID); err != nil {
		return err
	}
	if err := d.validateVersionFormat(version); err != nil {
		return err
	}

	exists, err := d.datasetRepo.ExistDatasetVersionByID(ctx, datasetID, version)
	if err != nil {
		return err
	}
	if !exists {
		return errs.NewNotFoundError(constants.ERR_DATASET_VERSION_NOT_FOUND)
	}
	datasetVersion.DatasetID = datasetID
	datasetVersion.Version = version
	now := helperModel.NewTimestampFromNow()
	datasetVersion.CreatedAt = &now
	datasetVersion.UpdatedAt = &now
	return d.datasetRepo.UpsertDatasetVersion(ctx, datasetVersion)
}

// UpdateDatasetVersionStatus implements data.DataUsecase.
func (d *dataUsecase) UpdateDatasetVersionStatus(ctx context.Context, datasetID string, version string, status string) error {
	if err := d.validateDatasetID(datasetID); err != nil {
		return err
	}
	if err := d.validateVersionFormat(version); err != nil {
		return err
	}
	exists, err := d.datasetRepo.ExistDatasetVersionByID(ctx, datasetID, version)
	if err != nil {
		return err
	}
	if !exists {
		return errs.NewNotFoundError("dataset version not found")
	}

	return d.datasetRepo.UpdateDatasetVersionStatus(ctx, datasetID, version, status)
}

func (d *dataUsecase) ServingDatasetVersionData(
	ctx context.Context,
	datasetID string,
	version string,
	paginator *helperModel.Paginator,
	viewName string,
	filterGroups [][]entity.FilterInput,
	logicalOperator string,
	sortBy string,
	sortOrder string,
) ([]map[string]interface{}, error) {
	if err := d.validateDatasetID(datasetID); err != nil {
		return nil, err
	}
	if err := d.validateVersionFormat(version); err != nil {
		return nil, err
	}
	datasetVersion, err := d.datasetRepo.FetchDatasetVersionByID(ctx, datasetID, version)
	if err != nil {
		return nil, err
	}
	if datasetVersion == nil {
		return nil, errs.NewNotFoundError("dataset version not found")
	}
	if datasetVersion.Policies.Runtime == nil {
		return nil, errs.NewConflictError("runtime policy is not configured for this dataset version")
	}

	results, err := d.dataRepo.ExecuteQuery(
		ctx,
		datasetVersion.SourceID,
		&datasetVersion.Schema,
		&datasetVersion.Policies,
		filterGroups,
		logicalOperator,
		paginator,
		viewName,
		sortBy,
		sortOrder,
	)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// ServingDatasetVersionDataByKey implements data.DataUsecase.
func (d *dataUsecase) ServingDatasetVersionDataByKey(ctx context.Context, datasetID, version, key, viewName string) (map[string]interface{}, error) {
	if err := d.validateDatasetID(datasetID); err != nil {
		return nil, err
	}
	if err := d.validateVersionFormat(version); err != nil {
		return nil, err
	}

	// Get dataset version to get policies
	datasetVersion, err := d.datasetRepo.FetchDatasetVersionByID(ctx, datasetID, version)
	if err != nil {
		return nil, err
	}
	if datasetVersion == nil {
		return nil, errs.NewNotFoundError("dataset version not found")
	}

	// Use runtime policy and prepare for key-based filtering
	if datasetVersion.Policies.Runtime == nil {
		return nil, errs.NewConflictError("runtime policy is not configured for this dataset version")
	}

	results, err := d.dataRepo.ExecuteQueryByKey(ctx, datasetVersion.SourceID, &datasetVersion.Schema, &datasetVersion.Policies, key, viewName)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// CreateDatasetVersionData implements data.DataUsecase.
func (d *dataUsecase) CreateDatasetVersionData(ctx context.Context, datasetID string, version string, data map[string]interface{}) (map[string]interface{}, error) {
	if err := d.validateDatasetID(datasetID); err != nil {
		return nil, err
	}
	if err := d.validateVersionFormat(version); err != nil {
		return nil, err
	}

	// Get dataset version to get policies
	datasetVersion, err := d.datasetRepo.FetchDatasetVersionByID(ctx, datasetID, version)
	if err != nil {
		return nil, err
	}
	if datasetVersion == nil {
		return nil, errs.NewNotFoundError("dataset version not found")
	}

	// Check if WritePolicy exists
	if datasetVersion.Policies.Write == nil {
		return nil, errs.NewConflictError("WritePolicy is not configured for data creation")
	}

	// Validate data is not null or empty
	if len(data) == 0 {
		return nil, errs.NewBadRequestError("data cannot be null or empty")
	}

	result, err := d.dataRepo.ExecuteCreate(ctx, datasetVersion.SourceID, datasetVersion.Schema, datasetVersion.Policies.Write, data)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateDatasetVersionDataByKey implements data.DataUsecase.
func (d *dataUsecase) UpdateDatasetVersionDataByKey(ctx context.Context, datasetID string, version string, key string, data map[string]interface{}) (map[string]interface{}, error) {
	if err := d.validateDatasetID(datasetID); err != nil {
		return nil, err
	}
	if err := d.validateVersionFormat(version); err != nil {
		return nil, err
	}
	datasetVersion, err := d.datasetRepo.FetchDatasetVersionByID(ctx, datasetID, version)
	if err != nil {
		return nil, err
	}
	if datasetVersion == nil {
		return nil, errs.NewNotFoundError("dataset version not found")
	}

	if datasetVersion.Policies.Write == nil {
		return nil, errs.NewConflictError("WritePolicy is not configured for data update")
	}

	if datasetVersion.Policies.Write.KeyField == "" {
		return nil, errs.NewConflictError("KeyField is not configured in write policy for key-based update")
	}

	if len(data) == 0 {
		return nil, errs.NewBadRequestError("data cannot be null or empty")
	}

	result, err := d.dataRepo.ExecuteUpdate(ctx, datasetVersion.SourceID, datasetVersion.Schema, datasetVersion.Policies.Write, key, data)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errs.NewNotFoundError("data with the specified key not found")
	}

	return result, nil
}

// DeleteDatasetVersionDataByKey implements data.DataUsecase.
func (d *dataUsecase) DeleteDatasetVersionDataByKey(ctx context.Context, datasetID string, version string, key string) error {
	if err := d.validateDatasetID(datasetID); err != nil {
		return err
	}
	if err := d.validateVersionFormat(version); err != nil {
		return err
	}
	datasetVersion, err := d.datasetRepo.FetchDatasetVersionByID(ctx, datasetID, version)
	if err != nil {
		return err
	}
	if datasetVersion == nil {
		return errs.NewNotFoundError("dataset version not found")
	}
	if datasetVersion.Policies.Delete == nil {
		return errs.NewConflictError("DeletePolicy is not configured for data deletion")
	}

	sqlResult, err := d.dataRepo.ExecuteDelete(ctx, datasetVersion.SourceID, datasetVersion.Policies.Delete, key)
	if err != nil {
		return err
	}
	if sqlResult == nil {
		return errs.NewNotFoundError("data with the specified key not found")
	}
	rowsAffected, err := sqlResult.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errs.NewNotFoundError("data with the specified key not found")
	}
	return nil
}

func (d *dataUsecase) validateExistSource(ctx context.Context, sourceID *uuid.UUID, schemaName, tableName string) error {
	exists, err := d.datasetRepo.ExistSourceByID(ctx, sourceID)
	if err != nil {
		return err
	}
	if !exists {
		return errs.NewNotFoundError(constants.ERR_SOURCE_NOT_FOUND)
	}
	exists, err = d.datasetRepo.ExistSchemaByName(ctx, sourceID, schemaName)
	if err != nil {
		return err
	}
	if !exists {
		return errs.NewNotFoundError(constants.ERR_SCHEMA_NOT_FOUND)
	}
	exists, err = d.datasetRepo.ExistTableByName(ctx, sourceID, schemaName, tableName)
	if err != nil {
		return err
	}
	if !exists {
		return errs.NewNotFoundError(constants.ERR_TABLE_NOT_FOUND)
	}

	return nil
}

// FetchTableData implements data.DataUsecase.
func (d *dataUsecase) FetchTableData(
	ctx context.Context,
	sourceID *uuid.UUID,
	schemaName, tableName string,
	filterGroups [][]entity.FilterInput,
	logicalOperator string,
	paginator *helperModel.Paginator,
	sortBy, sortOrder string,
) ([]map[string]interface{}, error) {
	// Validate exists
	err := d.validateExistSource(ctx, sourceID, schemaName, tableName)
	if err != nil {
		return nil, err
	}

	return d.dataRepo.FetchTableData(ctx, sourceID, schemaName, tableName, filterGroups, logicalOperator, paginator, sortBy, sortOrder)
}

// FetchTableDataByKey implements data.DataUsecase.
func (d *dataUsecase) FetchTableDataByKey(
	ctx context.Context,
	sourceID *uuid.UUID,
	schemaName, tableName, keyField string,
	keyValue interface{},
) (map[string]interface{}, error) {
	// Validate source exists
	err := d.validateExistSource(ctx, sourceID, schemaName, tableName)
	if err != nil {
		return nil, err
	}
	exists, err := d.datasetRepo.ExistColumnByName(ctx, sourceID, schemaName, tableName, keyField)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errs.NewNotFoundError(constants.ERR_COLUMN_NOT_FOUND)
	}

	return d.dataRepo.FetchTableDataByKey(ctx, sourceID, schemaName, tableName, keyField, keyValue)
}

// CreateTableData implements data.DataUsecase.
func (d *dataUsecase) CreateTableData(
	ctx context.Context,
	sourceID *uuid.UUID,
	schemaName, tableName string,
	data map[string]interface{},
) (map[string]interface{}, error) {
	// Validate source exists
	// Validate source exists
	err := d.validateExistSource(ctx, sourceID, schemaName, tableName)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, errs.NewBadRequestError("data cannot be null or empty")
	}

	// Fetch columns information for validation
	columnsFilter := &filter.ColumnsFilter{
		SourceID: sourceID,
		Schema:   schemaName,
		Table:    tableName,
	}
	paginator := helperModel.NewPaginator()
	paginator.PerPage = 1000 // Get all columns

	columns, err := d.datasetRepo.FetchColumnsList(ctx, columnsFilter, &paginator)
	if err != nil {
		return nil, err
	}
	if len(columns) == 0 {
		return nil, errs.NewBadRequestError("table not found or has no columns")
	}

	return d.dataRepo.CreateTableData(ctx, sourceID, schemaName, tableName, columns, data)
}

// UpdateTableData implements data.DataUsecase.
func (d *dataUsecase) UpdateTableData(
	ctx context.Context,
	sourceID *uuid.UUID,
	schemaName, tableName, keyField string,
	keyValue interface{},
	data map[string]interface{},
) (map[string]interface{}, error) {
	// Validate source exists
	// Validate source exists
	err := d.validateExistSource(ctx, sourceID, schemaName, tableName)
	if err != nil {
		return nil, err
	}
	exists, err := d.datasetRepo.ExistColumnByName(ctx, sourceID, schemaName, tableName, keyField)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errs.NewNotFoundError(constants.ERR_COLUMN_NOT_FOUND)
	}

	if len(data) == 0 {
		return nil, errs.NewBadRequestError("data cannot be null or empty")
	}

	// Fetch columns information for validation
	columnsFilter := &filter.ColumnsFilter{
		SourceID: sourceID,
		Schema:   schemaName,
		Table:    tableName,
	}
	paginator := helperModel.NewPaginator()
	paginator.PerPage = 1000 // Get all columns

	columns, err := d.datasetRepo.FetchColumnsList(ctx, columnsFilter, &paginator)
	if err != nil {
		return nil, err
	}
	if len(columns) == 0 {
		return nil, errs.NewBadRequestError("table not found or has no columns")
	}

	result, err := d.dataRepo.UpdateTableData(ctx, sourceID, schemaName, tableName, keyField, keyValue, columns, data)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errs.NewNotFoundError("data with the specified key not found")
	}

	return result, nil
}

// DeleteTableData implements data.DataUsecase.
func (d *dataUsecase) DeleteTableData(
	ctx context.Context,
	sourceID *uuid.UUID,
	schemaName, tableName, keyField string,
	keyValue interface{},
) error {
	// Validate source exists
	err := d.validateExistSource(ctx, sourceID, schemaName, tableName)
	if err != nil {
		return err
	}
	exists, err := d.datasetRepo.ExistColumnByName(ctx, sourceID, schemaName, tableName, keyField)
	if err != nil {
		return err
	}
	if !exists {
		return errs.NewNotFoundError(constants.ERR_COLUMN_NOT_FOUND)
	}

	sqlResult, err := d.dataRepo.DeleteTableData(ctx, sourceID, schemaName, tableName, keyField, keyValue)
	if err != nil {
		return err
	}
	if sqlResult == nil {
		return errs.NewNotFoundError("data with the specified key not found")
	}
	rowsAffected, err := sqlResult.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errs.NewNotFoundError("data with the specified key not found")
	}
	return nil
}

func (d *dataUsecase) UploadReportingTemplate(ctx context.Context, template *entity.ReportingTemplate, fileData []byte, fileName string) error {
	if template == nil {
		return errs.NewBadRequestError(constants.ERR_INVALID_REQUEST_BODY)
	}
	if err := d.validateDatasetID(template.DatasetID); err != nil {
		return err
	}
	if err := d.validateVersionFormat(template.Version); err != nil {
		return err
	}
	if template.Name == "" {
		return errs.NewBadRequestError("template name is required")
	} else if len(template.Positions) == 0 || len(template.Positions) < 1 {
		return errs.NewBadRequestError("template positions are required")
	}
	datasetVersion, err := d.datasetRepo.FetchDatasetVersionByID(ctx, template.DatasetID, template.Version)
	if err != nil {
		return err
	}
	if datasetVersion == nil {
		return errs.NewNotFoundError("dataset version not found")
	}
	if datasetVersion.Policies.Runtime == nil {
		return errs.NewConflictError("runtime policy is not configured for this dataset version")
	}
	if datasetVersion.Policies.Views == nil {
		return errs.NewConflictError("views policy is not configured for this dataset version")
	}
	if template.View == "" {
		template.View = datasetVersion.Policies.Runtime.DefaultView
	} else {
		viewConfigs, ok := datasetVersion.Policies.Views[template.View]
		if !ok || len(viewConfigs) == 0 {
			return errs.NewNotFoundError(fmt.Sprintf("view '%s' not found or is empty in policies", template.View))
		}
	}
	datasetViews := datasetVersion.Policies.Views[template.View]

	// Build a map of available columns per table from the view definition
	viewColumnsMap := make(map[string]map[string]bool)
	for _, v := range datasetViews {
		table := v.TableName
		if _, ok := viewColumnsMap[table]; !ok {
			viewColumnsMap[table] = make(map[string]bool)
		}
		for _, c := range v.Columns {
			viewColumnsMap[table][c] = true
		}
	}

	// Validate template positions exist in the view
	missingPos := []string{}
	for _, pos := range template.Positions {
		if pos.TableName == "" || pos.ColumnsName == "" {
			return errs.NewBadRequestError("template positions have invalid format")
		}
		if colsMap, ok := viewColumnsMap[pos.TableName]; ok {
			if !colsMap[pos.ColumnsName] {
				missingPos = append(missingPos, fmt.Sprintf("%s.%s", pos.TableName, pos.ColumnsName))
			}
		} else {
			missingPos = append(missingPos, fmt.Sprintf("%s.%s", pos.TableName, pos.ColumnsName))
		}
	}

	if len(missingPos) > 0 {
		return errs.NewConflictError(fmt.Sprintf("template references not found in view: missing_positions=%v", missingPos))
	}
	now := helperModel.NewTimestampFromNow()
	template.CreatedAt = &now
	template.UpdatedAt = &now
	template.GenUUID()

	fileRequest := proto_models.FileRequest{
		Path:               "reporting",
		Folder:             "templates",
		OriginalFilename:   fileName,
		IsGenerateFilename: true,
		Body:               fileData,
	}
	status, data, err := d.documentRepo.UploadFile(ctx, &fileRequest)
	if err != nil || status != http.StatusOK {
		if status == http.StatusServiceUnavailable {
			return errs.NewInternalServerError("document service is unavailable")
		}
		return errs.NewInternalServerError("failed to upload reporting template file")
	}
	if data == nil {
		return errs.NewInternalServerError("failed to upload reporting template file")
	}
	template.ResourceID = data.ResourceId

	return d.datasetRepo.UpsertReportingTemplate(ctx, template)
}

func (d *dataUsecase) InsertReportingJob(ctx context.Context, job *entity.ReportingTemplateExportJob, key string) error {
	if job == nil {
		return errs.NewBadRequestError(constants.ERR_INVALID_REQUEST_BODY)
	}
	exists, err := d.datasetRepo.ExistReportingTemplateByID(ctx, job.ReportingTemplateID)
	if err != nil {
		return err
	}
	if !exists {
		return errs.NewNotFoundError("reporting template not found")
	}
	now := helperModel.NewTimestampFromNow()
	job.CreatedAt = &now
	job.Status = constants.EXPORT_JOB_STATUS_PENDING
	job.ResourceID = ""

	template, err := d.datasetRepo.FetchReportingTemplateByID(ctx, job.ReportingTemplateID)
	if err != nil {
		return err
	}
	if template == nil {
		return errs.NewNotFoundError("reporting template not found")
	}

	go d.processReportingTemplateExportJob(job.JobID, template, key)

	return d.datasetRepo.UpsertReportingTemplateExportJob(ctx, job)
}

func (d *dataUsecase) processReportingTemplateExportJob(JobID *uuid.UUID, template *entity.ReportingTemplate, key string) error {
	ctx := context.Background()
	exportData, err := d.generateReportingTemplateExportFile(ctx, JobID, template, key)
	if err != nil {
		d.datasetRepo.UpdateReportingExportStatusFail(ctx, JobID, err.Error())
		return err
	}

	now := helperModel.NewTimestampFromNow()
	fileName := fmt.Sprintf("%s_%s_%s.pdf", template.Name, key, now.Format("20060102150405"))
	fileRequest := proto_models.FileRequest{
		Path:               constants.DOCUMENT_PATH_REPORTING,
		Folder:             constants.DOCUMENT_FOLDER_EXPORT_TEMPLATES,
		OriginalFilename:   fileName,
		IsGenerateFilename: true,
		Body:               exportData,
	}
	status, data, err := d.documentRepo.UploadFile(ctx, &fileRequest)
	if err != nil || status != http.StatusOK {
		if status == http.StatusServiceUnavailable {
			d.datasetRepo.UpdateReportingExportStatusFail(ctx, JobID, "document service is unavailable")
			return err
		}
		d.datasetRepo.UpdateReportingExportStatusFail(ctx, JobID, "failed to upload export file")
		return err
	}
	if data == nil {
		d.datasetRepo.UpdateReportingExportStatusFail(ctx, JobID, "failed to upload export file")
		return errs.NewInternalServerError("failed to upload export file")
	}

	err = d.datasetRepo.UpdateReportingExportStatusSuccess(ctx, JobID, &now, data.ResourceId)
	if err != nil {
		return err
	}

	return nil
}

func (d *dataUsecase) generateReportingTemplateExportFile(ctx context.Context, JobID *uuid.UUID, template *entity.ReportingTemplate, key string) ([]byte, error) {
	datasetVersion, err := d.datasetRepo.FetchDatasetVersionByID(ctx, template.DatasetID, template.Version)
	if err != nil {
		return nil, err
	}

	var fileData []byte
	if datasetVersion != nil && datasetVersion.Policies.Runtime != nil && datasetVersion.Policies.Views != nil {
		// Fetch data using the key
		data, err := d.ServingDatasetVersionDataByKey(ctx, template.DatasetID, template.Version, key, template.View)
		if err != nil {
			return nil, err
		}
		if data == nil {
			return nil, errs.NewNotFoundError("data not found for the provided key")
		}

		// Get the PDF template file
		resourceReq := &proto_models.GetFileByResourceIDRequest{
			ResourceId: template.ResourceID,
		}
		status, resp, err := d.documentRepo.GetFileByResourceID(ctx, resourceReq)
		if err != nil || status != http.StatusOK {
			if status == http.StatusServiceUnavailable {
				return nil, errs.NewInternalServerError("document service is unavailable")
			}
			return nil, errs.NewInternalServerError("failed to get reporting template file")
		}
		if resp == nil {
			return nil, errs.NewNotFoundError("reporting template file not found")
		}

		// Download the PDF template from URL
		httpResp, err := http.Get(resp.Url)
		if err != nil {
			return nil, errs.NewInternalServerError("failed to download PDF template")
		}
		defer httpResp.Body.Close()

		if httpResp.StatusCode != http.StatusOK {
			return nil, errs.NewInternalServerError("failed to download PDF template")
		}

		templateData, err := io.ReadAll(httpResp.Body)
		if err != nil {
			return nil, errs.NewInternalServerError("failed to read PDF template")
		}

		// Generate PDF with data
		fileData, err = d.generatePDFFromTemplate(templateData, template, data)
		if err != nil {
			return nil, err
		}
	}
	return fileData, nil
}

func (d *dataUsecase) generatePDFFromTemplate(templateData []byte, template *entity.ReportingTemplate, data map[string]interface{}) ([]byte, error) {
	// Initialize PDF
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Import the template PDF
	templateReader := bytes.NewReader(templateData)
	imp := gofpdi.NewImporter()

	// Read the template
	var rs io.ReadSeeker = templateReader
	tpl := imp.ImportPageFromStream(pdf, &rs, 1, "/MediaBox")

	// Add a page and use the template
	pdf.AddPage()
	imp.UseImportedTemplate(pdf, tpl, 0, 0, 210, 297) // A4 size in mm

	// Set font for text
	pdf.AddUTF8Font("THSarabunNew", "", "./assets/fonts/THSarabunNew/THSarabunNew.ttf")
	pdf.AddUTF8Font("THSarabunNew Bold", "B", "./assets/fonts/THSarabunNew/THSarabunNew Bold.ttf")

	pdf.SetFont("THSarabunNew", "", 16)
	pdf.SetTextColor(0, 0, 0)
	// Process positions and add text to PDF
	for _, pos := range template.Positions {
		var value string
		// Extract value based on position type
		if pos.TableName != "" && pos.ColumnsName != "" {
			// Case 1: Key matches TableName
			if tableData, ok := data[pos.TableName]; ok {
				fmt.Printf("Extracting value from table '%s' for alias '%s'\n", pos.TableName, pos.Alias)
				value = d.extractValueFromTableData(tableData, pos.Alias)
			}
		}
		if pos.ColumnsName != "" {
			// Case 2: Key matches ColumnsName or Alias - use value directly
			if val, ok := data[pos.ColumnsName]; ok {
				value = fmt.Sprintf("%v", val)
			} else if pos.Alias != "" {
				if val, ok := data[pos.Alias]; ok {
					value = fmt.Sprintf("%v", val)
				}
			}
		}

		// Set text position and write
		if value != "" {
			pdf.SetXY(pos.X, pos.Y)
			pdf.Cell(0, 0, value)
		}
	}

	// Output PDF to buffer
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, errs.NewInternalServerError("failed to generate PDF")
	}

	return buf.Bytes(), nil
}

func (d *dataUsecase) FetchReportingExportJobByID(ctx context.Context, jobID *uuid.UUID) (*dto.ReportingExportJobResponseDTO, error) {
	job, err := d.datasetRepo.FetchReportingExportJobByID(ctx, jobID)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, errs.NewNotFoundError("reporting export job not found")
	}
	req := proto_models.GetFileByResourceIDRequest{
		ResourceId: job.ResourceID,
	}
	status, response, err := d.documentRepo.GetFileByResourceID(ctx, &req)
	if err != nil || status != http.StatusOK {
		if status == http.StatusServiceUnavailable {
			return nil, errs.NewInternalServerError("document service is unavailable")
		}
		return nil, errs.NewInternalServerError("failed to get export file information")
	}
	if response == nil {
		return nil, errs.NewNotFoundError("export file not found")
	}
	resp, err := helperModel.ConvertStruct[*entity.ReportingTemplateExportJob, *dto.ReportingExportJobResponseDTO](job)
	if err != nil {
		return nil, err
	}
	resp.Url = response.Url
	resp.OriginalFilename = response.OriginalFilename
	resp.FileSize = response.Size
	resp.ContentType = response.ContentType
	return resp, nil
}

func (d *dataUsecase) extractValueFromTableData(tableData interface{}, alias string) string {
	switch v := tableData.(type) {
	case map[string]interface{}:
		// Single object - extract value by alias
		if val, ok := v[alias]; ok {
			return fmt.Sprintf("%v", val)
		}
	case []interface{}:
		// Array of objects - take the first one
		if len(v) > 0 {
			if firstObj, ok := v[0].(map[string]interface{}); ok {
				if val, ok := firstObj[alias]; ok {
					return fmt.Sprintf("%v", val)
				}
			}
		}
	}
	return ""
}

func (d *dataUsecase) FetchExportJobByID(ctx context.Context, jobID *uuid.UUID) (*entity.ExportJob, error) {
	job, err := d.datasetRepo.FetchExportJobByID(ctx, jobID)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, errs.NewNotFoundError("export job not found")
	}
	req := proto_models.GetFileByResourceIDRequest{
		ResourceId: job.DestinationUri,
	}
	status, response, err := d.documentRepo.GetFileByResourceID(ctx, &req)
	if err != nil || status != http.StatusOK {
		if status == http.StatusServiceUnavailable {
			return nil, errs.NewInternalServerError("document service is unavailable")
		}
		return nil, errs.NewInternalServerError("failed to get export file information")
	}
	if response == nil {
		return nil, errs.NewNotFoundError("export file not found")
	}
	job.DestinationUri = response.Url
	job.OriginalFilename = response.OriginalFilename
	job.FileSize = response.Size
	return job, nil
}

func (d *dataUsecase) exportDatasetExcel(ctx context.Context, exportJob *entity.ExportJob) ([]byte, error) {
	// Create Excel file first
	f := excelize.NewFile()
	sheetName := "Sheet1"

	// Get or create sheet
	index, err := f.GetSheetIndex(sheetName)
	if err != nil || index == -1 {
		index, err = f.NewSheet(sheetName)
		if err != nil {
			return nil, err
		}
	}
	f.SetActiveSheet(index)

	page := 1
	currentRow := 1 // Start from row 1
	columnNames := []string{}
	headerWritten := false
	hasData := false

	for {
		// Setup paginator for each batch
		paginator := helperModel.NewPaginator()
		paginator.Page = page
		paginator.PerPage = 100

		// Fetch data for this page
		datas, err := d.ServingDatasetVersionData(ctx, exportJob.DatasetId, exportJob.Version, &paginator, exportJob.View, nil, "", "", "")
		if err != nil {
			return nil, err
		}

		// If no data returned, break the loop
		if len(datas) == 0 {
			break
		}

		hasData = true

		// Process and write data immediately
		for _, data := range datas {
			rows := d.flattenDataForExcel(data)

			// Write header on first batch only
			if !headerWritten && len(rows) > 0 {
				columnNames = d.getUniqueColumnNames(rows)
				for colIdx, colName := range columnNames {
					cellName, _ := excelize.CoordinatesToCellName(colIdx+1, currentRow)
					f.SetCellValue(sheetName, cellName, colName)
				}
				currentRow++
				headerWritten = true
			}

			// Update column names if new columns appear
			if headerWritten {
				newColumns := d.getUniqueColumnNames(rows)
				for _, newCol := range newColumns {
					found := false
					for _, existingCol := range columnNames {
						if existingCol == newCol {
							found = true
							break
						}
					}
					if !found {
						columnNames = append(columnNames, newCol)
						// Add new column to header
						cellName, _ := excelize.CoordinatesToCellName(len(columnNames), 1)
						f.SetCellValue(sheetName, cellName, newCol)
					}
				}
			}

			// Write data rows
			for _, row := range rows {
				for colIdx, colName := range columnNames {
					cellName, _ := excelize.CoordinatesToCellName(colIdx+1, currentRow)
					value := row[colName]
					f.SetCellValue(sheetName, cellName, value)
				}
				currentRow++
			}
		}

		// If we got less than 100 records, this is the last page
		if len(datas) < 100 {
			break
		}

		// Move to next page
		page++
	}

	// Check if we have any data
	if !hasData {
		return nil, errs.ErrNoContent()
	}

	// Save to buffer and build the Excel file
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (d *dataUsecase) exportDatasetCSV(ctx context.Context, exportJob *entity.ExportJob) ([]byte, error) {
	// Create buffer for CSV
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	page := 1
	columnNames := []string{}
	headerWritten := false
	hasData := false

	for {
		// Setup paginator for each batch
		paginator := helperModel.NewPaginator()
		paginator.Page = page
		paginator.PerPage = 100

		// Fetch data for this page
		datas, err := d.ServingDatasetVersionData(ctx, exportJob.DatasetId, exportJob.Version, &paginator, exportJob.View, nil, "", "", "")
		if err != nil {
			return nil, err
		}

		// If no data returned, break the loop
		if len(datas) == 0 {
			break
		}

		hasData = true

		// Process and write data immediately
		for _, data := range datas {
			rows := d.flattenDataForExcel(data)

			// Write header on first batch only
			if !headerWritten && len(rows) > 0 {
				columnNames = d.getUniqueColumnNames(rows)
				if err := writer.Write(columnNames); err != nil {
					return nil, err
				}
				headerWritten = true
			}

			// Update column names if new columns appear
			if headerWritten {
				newColumns := d.getUniqueColumnNames(rows)
				for _, newCol := range newColumns {
					found := false
					for _, existingCol := range columnNames {
						if existingCol == newCol {
							found = true
							break
						}
					}
					if !found {
						columnNames = append(columnNames, newCol)
						// Note: CSV doesn't support adding columns dynamically like Excel
						// New columns will only appear in subsequent rows
					}
				}
			}

			// Write data rows
			for _, row := range rows {
				record := make([]string, len(columnNames))
				for colIdx, colName := range columnNames {
					value := row[colName]
					record[colIdx] = d.formatValueForCSV(value)
				}
				if err := writer.Write(record); err != nil {
					return nil, err
				}
			}
		}

		// If we got less than 100 records, this is the last page
		if len(datas) < 100 {
			break
		}

		// Move to next page
		page++
	}

	// Check if we have any data
	if !hasData {
		return nil, errs.ErrNoContent()
	}

	// Flush any buffered data
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// formatValueForCSV converts values to CSV-compatible string format
func (d *dataUsecase) formatValueForCSV(value interface{}) string {
	if value == nil {
		return ""
	}

	switch v := value.(type) {
	case string:
		return v
	case map[string]interface{}:
		return fmt.Sprintf("%v", v)
	case []interface{}:
		return fmt.Sprintf("%v", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// flattenDataForExcel flattens nested objects and expands arrays
func (d *dataUsecase) flattenDataForExcel(data map[string]interface{}) []map[string]interface{} {
	// First pass: separate array fields from non-array fields
	arrayFields := make(map[string][]interface{})
	maxArrayLen := 1

	for key, value := range data {
		if arr, ok := value.([]interface{}); ok {
			arrayFields[key] = arr
			if len(arr) > maxArrayLen {
				maxArrayLen = len(arr)
			}
		}
	}

	// Create rows based on the maximum array length
	rows := make([]map[string]interface{}, maxArrayLen)
	for i := 0; i < maxArrayLen; i++ {
		rows[i] = make(map[string]interface{})
	}

	// Process each field
	for key, value := range data {
		switch v := value.(type) {
		case []interface{}:
			// Case 2: Array of objects
			// Skip empty arrays
			if len(v) == 0 {
				continue
			}

			for i, item := range v {
				if i < maxArrayLen {
					if obj, ok := item.(map[string]interface{}); ok {
						// Flatten object in array with key_subkey format
						for subKey, subValue := range obj {
							columnName := fmt.Sprintf("%s_%s", key, subKey)
							rows[i][columnName] = d.formatValue(subValue)
						}
					} else {
						// Simple array value
						rows[i][key] = d.formatValue(item)
					}
				}
			}
			// Fill empty rows with nil for this field
			for i := len(v); i < maxArrayLen; i++ {
				if obj, ok := v[0].(map[string]interface{}); ok {
					for subKey := range obj {
						columnName := fmt.Sprintf("%s_%s", key, subKey)
						rows[i][columnName] = nil
					}
				} else {
					rows[i][key] = nil
				}
			}
		case map[string]interface{}:
			// Case 1: Nested object - flatten with key_subkey format
			for subKey, subValue := range v {
				columnName := fmt.Sprintf("%s_%s", key, subKey)
				for i := 0; i < maxArrayLen; i++ {
					rows[i][columnName] = d.formatValue(subValue)
				}
			}
		default:
			// Case 3: Simple string or primitive value
			for i := 0; i < maxArrayLen; i++ {
				rows[i][key] = d.formatValue(v)
			}
		}
	}

	return rows
}

// formatValue converts values to appropriate Excel format
func (d *dataUsecase) formatValue(value interface{}) interface{} {
	if value == nil {
		return ""
	}

	switch v := value.(type) {
	case map[string]interface{}:
		// If still nested, convert to string representation
		return fmt.Sprintf("%v", v)
	case []interface{}:
		// If still array, convert to string representation
		return fmt.Sprintf("%v", v)
	default:
		return v
	}
}

// getUniqueColumnNames extracts all unique column names from rows
func (d *dataUsecase) getUniqueColumnNames(rows []map[string]interface{}) []string {
	columnSet := make(map[string]bool)
	var columns []string

	for _, row := range rows {
		for key := range row {
			if !columnSet[key] {
				columnSet[key] = true
				columns = append(columns, key)
			}
		}
	}

	return columns
}

func (d *dataUsecase) processJob(exportJob *entity.ExportJob) error {
	ctx := context.Background()
	var fileByte []byte
	if exportJob.Format == constants.EXPORT_JOB_FORMAT_XLSX {
		excelByte, err := d.exportDatasetExcel(ctx, exportJob)
		if err != nil {
			err := d.datasetRepo.UpdateStatusFail(ctx, exportJob.JobId, err.Error())
			return err
		}
		fileByte = excelByte
	} else if exportJob.Format == constants.EXPORT_JOB_FORMAT_CSV {
		csvByte, err := d.exportDatasetCSV(ctx, exportJob)
		if err != nil {
			err := d.datasetRepo.UpdateStatusFail(ctx, exportJob.JobId, err.Error())
			return err
		}
		fileByte = csvByte
	} else {
		err := d.datasetRepo.UpdateStatusFail(ctx, exportJob.JobId, "unsupported export format")
		return err
	}

	now := helperModel.NewTimestampFromNow()
	fileName := "exported_data_" + now.Format("20060102150405") + "." + exportJob.Format
	fileReq := proto_models.FileRequest{
		Path:               constants.DOCUMENT_PATH_REPORTING,
		Folder:             constants.DOCUMENT_FOLDER_EXPORT,
		OriginalFilename:   fileName,
		IsGenerateFilename: true,
		Body:               fileByte,
	}
	_, response, err := d.documentRepo.UploadFile(ctx, &fileReq)
	if err != nil {
		err := d.datasetRepo.UpdateStatusFail(ctx, exportJob.JobId, err.Error())
		return err
	}
	if response == nil || response.ResourceId == "" {
		err := d.datasetRepo.UpdateStatusFail(ctx, exportJob.JobId, "No resource file")
		return err
	}
	err = d.datasetRepo.UpdateStatusSuccess(ctx, exportJob.JobId, response.ResourceId, &now)
	if err != nil {
		err := d.datasetRepo.UpdateStatusFail(ctx, exportJob.JobId, err.Error())
		return err
	}
	return nil
}

// ExportJob implements data.DataUsecase.
func (d *dataUsecase) InsertExportJob(ctx context.Context, exportJob *entity.ExportJob) error {
	if exportJob == nil {
		return errs.NewBadRequestError(constants.ERR_INVALID_REQUEST_BODY)
	}
	if err := d.validateDatasetID(exportJob.DatasetId); err != nil {
		return err
	}
	if err := d.validateVersionFormat(exportJob.Version); err != nil {
		return err
	}
	datasetVersion, err := d.datasetRepo.FetchDatasetVersionByID(ctx, exportJob.DatasetId, exportJob.Version)
	if err != nil {
		return err
	}
	if datasetVersion == nil {
		return errs.NewNotFoundError("dataset version not found")
	}
	if datasetVersion.Policies.Runtime == nil {
		return errs.NewConflictError("runtime policy is not configured for this dataset version")
	}
	if datasetVersion.Policies.Views == nil {
		return errs.NewConflictError("views policy is not configured for this dataset version")
	}
	if exportJob.View == "" {
		exportJob.View = datasetVersion.Policies.Runtime.DefaultView
	} else {
		viewConfigs, ok := datasetVersion.Policies.Views[exportJob.View]
		if !ok || len(viewConfigs) == 0 {
			return errs.NewNotFoundError(fmt.Sprintf("view '%s' not found or is empty in policies", exportJob.View))
		}
	}
	now := helperModel.NewTimestampFromNow()
	exportJob.CreatedAt = &now
	exportJob.DestinationUri = ""
	exportJob.Status = constants.EXPORT_JOB_STATUS_PENDING
	//exportJob.CompletedAt = &now

	err = d.datasetRepo.InsertExportJob(ctx, exportJob)
	if err != nil {
		err := d.datasetRepo.UpdateStatusFail(ctx, exportJob.JobId, err.Error())
		return err
	}
	go d.processJob(exportJob)
	return nil
}
