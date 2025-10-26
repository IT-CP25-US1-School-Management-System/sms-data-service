package data

import (
	"context"

	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/entity"
)

type PsqlDataRepository interface {
}

type PsqlDatasetRepository interface {
	FetchSourceList(ctx context.Context) ([]*entity.Sources, error)
}

type RedisRepository interface {
}
