package database

import (
	"context"

	"github.com/GodeFvt/go-backend/psql"
	"github.com/gofrs/uuid"
)

type DBConnectionManagerRepository interface {
	GetConnection(ctx context.Context, sourceID uuid.UUID) (*psql.Client, error)
	Close(sourceID uuid.UUID)
	CloseAll()
}
