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
	err = d.datasetRepo.BatchInsertInformationDatabase(ctx, infoSchema, infoTables, infoColumns)
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

func NewDataUsecase(dataRepo data.PsqlDataRepository, datasetRepo data.PsqlDatasetRepository, redisRepo data.RedisRepository, cryptoSecret string) data.DataUsecase {
	return &dataUsecase{
		dataRepo:     dataRepo,
		datasetRepo:  datasetRepo,
		redisRepo:    redisRepo,
		cryptoSecret: cryptoSecret,
	}
}
