package usecase

import (
	"context"

	"github.com/GodeFvt/go-backend/helper/models"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/entity"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/filter"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/data/v1"
)

type dataUsecase struct {
	dataRepo    data.PsqlDataRepository
	datasetRepo data.PsqlDatasetRepository
	redisRepo   data.RedisRepository
}

// FetchDatasetByID implements data.DataUsecase.
func (d *dataUsecase) FetchDatasetByID(ctx context.Context, datasetID string) (*entity.Datasets, error) {
	return d.datasetRepo.FetchDatasetByID(ctx, datasetID)
}

// FetchDatasetList implements data.DataUsecase.
func (d *dataUsecase) FetchDatasetList(ctx context.Context, filter *filter.DatasetsFilter, paginator *models.Paginator) ([]*entity.Datasets, error) {
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
