package usecase

import (
	"context"

	"github.com/GodeFvt/go-backend/psql"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/database/v1"
	"github.com/gofrs/uuid"
)

type dbConnectionManagerUsecase struct {
	dbConRepo database.DBConnectionManagerRepository
}

// Close implements database.DBConnectionManagerUsecase.
func (d *dbConnectionManagerUsecase) Close(sourceID uuid.UUID) {
	d.dbConRepo.Close(sourceID)
}

// CloseAll implements database.DBConnectionManagerUsecase.
func (d *dbConnectionManagerUsecase) CloseAll() {
	d.dbConRepo.CloseAll()
}

// GetConnection implements database.DBConnectionManagerUsecase.
func (d *dbConnectionManagerUsecase) GetConnection(ctx context.Context, sourceID uuid.UUID) (*psql.Client, error) {
	return d.dbConRepo.GetConnection(ctx, sourceID)
}

func NewDBConnectionManagerUsecase(dbConRepo database.DBConnectionManagerRepository) database.DBConnectionManagerUsecase {
	return &dbConnectionManagerUsecase{
		dbConRepo: dbConRepo,
	}
}
