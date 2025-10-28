package data

import (
	"context"

	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/entity"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/filter"
)

type PsqlDataRepository interface {
}

type PsqlDatasetRepository interface {
	FetchSourceList(ctx context.Context) ([]*entity.Sources, error)
	FetchSchemasList(ctx context.Context, filter *filter.SchemasFilter) ([]*entity.Schemas, error)
}

type RedisRepository interface {
}
