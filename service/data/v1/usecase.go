package data

import (
	"context"

	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/entity"
)

type DataUsecase interface {
	FetchSourceList(ctx context.Context) ([]*entity.Sources, error)
}
