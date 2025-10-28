package data

import (
	"context"

	helperModel "github.com/GodeFvt/go-backend/helper/models"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/entity"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/filter"
)

type PsqlDataRepository interface {
}

type PsqlDatasetRepository interface {
	FetchSourceList(ctx context.Context) ([]*entity.Sources, error)
	FetchSchemasList(ctx context.Context, filter *filter.SchemasFilter) ([]*entity.Schemas, error)

	// Dataset
	FetchDatasetList(ctx context.Context, filter *filter.DatasetsFilter, paginator *helperModel.Paginator) ([]*entity.Datasets, error)
}

type RedisRepository interface {
}
