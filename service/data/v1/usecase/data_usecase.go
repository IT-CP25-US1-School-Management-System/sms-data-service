package usecase

import (
	"context"

	helperModel "github.com/GodeFvt/go-backend/helper/models"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/constants"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/errs"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/dto"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/entity"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/filter"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/data/v1"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/utils/crypto"
	"github.com/gofrs/uuid"
)

type dataUsecase struct {
	dataRepo     data.PsqlDataRepository
	datasetRepo  data.PsqlDatasetRepository
	redisRepo    data.RedisRepository
	cryptoSecret string
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

func NewDataUsecase(dataRepo data.PsqlDataRepository, datasetRepo data.PsqlDatasetRepository, redisRepo data.RedisRepository, cryptoSecret string) data.DataUsecase {
	return &dataUsecase{
		dataRepo:     dataRepo,
		datasetRepo:  datasetRepo,
		redisRepo:    redisRepo,
		cryptoSecret: cryptoSecret,
	}
}
