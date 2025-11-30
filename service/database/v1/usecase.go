package database

import (
	"context"

	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/database/v1/client"
	"github.com/gofrs/uuid"
)

type DBConnectionManagerUsecase interface {
	GetConnection(ctx context.Context, sourceID uuid.UUID) (*client.Client, error)
	GetDBType(sourceID uuid.UUID) (string, error)
	Close(sourceID uuid.UUID)
	CloseAll()
}
