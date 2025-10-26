package usecase

import (
	"context"

	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/entity"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/data/v1"
)

type dataUsecase struct {
	dataRepo    data.PsqlDataRepository
	datasetRepo data.PsqlDatasetRepository
	redisRepo   data.RedisRepository
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
