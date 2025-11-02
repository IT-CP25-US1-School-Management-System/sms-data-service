package usecase

import (
	"context"
	"regexp"
	"strings"

	helperModel "github.com/GodeFvt/go-backend/helper/models"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/constants"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/errs"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/entity"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/filter"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/data/v1"
)

type dataUsecase struct {
	dataRepo    data.PsqlDataRepository
	datasetRepo data.PsqlDatasetRepository
	redisRepo   data.RedisRepository
}

// validateDatasetID validates the dataset ID format
// Returns error if ID is invalid (empty, has spaces, or contains invalid characters)
func (d *dataUsecase) validateDatasetID(id string) error {
	if id == "" {
		return errs.NewBadRequestError(constants.ERR_DATASET_ID_IS_REQUIRED)
	}

	// Validate dataset ID format: only lowercase english letters, underscore, and hyphen, no spaces
	// Pattern: ^[a-z_-]+$ means start to end with only lowercase letters a-z, underscore, and hyphen
	idPattern := regexp.MustCompile("^[a-z_-]+$")

	// Check for spaces or invalid characters
	if strings.Contains(id, " ") || !idPattern.MatchString(id) {
		return errs.NewBadRequestError(constants.ERR_DATASET_ID_INVALID_FORMAT)
	}

	return nil
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
func (d *dataUsecase) FetchColumnsList(ctx context.Context, filter *filter.ColumnsFilter) ([]*entity.Columns, error) {
	return d.datasetRepo.FetchColumnsList(ctx, filter)
}

// FetchTablesList implements data.DataUsecase.
func (d *dataUsecase) FetchTablesList(ctx context.Context, filter *filter.TablesFilter) ([]*entity.Tables, error) {
	return d.datasetRepo.FetchTablesList(ctx, filter)
}

// FetchDatasetList implements data.DataUsecase.
func (d *dataUsecase) FetchDatasetList(ctx context.Context, filter *filter.DatasetsFilter, paginator *helperModel.Paginator) ([]*entity.Datasets, error) {
	return d.datasetRepo.FetchDatasetList(ctx, filter, paginator)
}

// FetchSchemasList implements data.DataUsecase.
func (d *dataUsecase) FetchSchemasList(ctx context.Context, filter *filter.SchemasFilter) ([]*entity.Schemas, error) {
	return d.datasetRepo.FetchSchemasList(ctx, filter)
}

// FetchSourceList implements data.DataUsecase.
func (d *dataUsecase) FetchSourceList(ctx context.Context) ([]*entity.Sources, error) {
	return d.datasetRepo.FetchSourceList(ctx)
}

func NewDataUsecase(dataRepo data.PsqlDataRepository, datasetRepo data.PsqlDatasetRepository, redisRepo data.RedisRepository) data.DataUsecase {
	return &dataUsecase{
		dataRepo:    dataRepo,
		datasetRepo: datasetRepo,
		redisRepo:   redisRepo,
	}
}
