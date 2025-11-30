package usecase

import (
	"context"
	"strings"

	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/errs"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/database/v1"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/database/v1/client"
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
func (d *dbConnectionManagerUsecase) GetConnection(ctx context.Context, sourceID uuid.UUID) (*client.Client, error) {
	conn, err := d.dbConRepo.GetConnection(ctx, sourceID)
	if err != nil {
		if strings.Contains(err.Error(), "failed to create connection") {
			return nil, errs.NewConflictError("failed connection: " + err.Error())
		} else if strings.Contains(err.Error(), "unsupported database type") {
			return nil, errs.NewBadRequestError("unsupported database type")
		} else if strings.Contains(err.Error(), "source not found") {
			return nil, errs.NewNotFoundError("source not found")
		}
		return nil, err
	}
	return conn, nil

}

// GetDBType implements database.DBConnectionManagerUsecase.
func (d *dbConnectionManagerUsecase) GetDBType(sourceID uuid.UUID) (string, error) {
	return d.dbConRepo.GetDBType(sourceID)
}

func NewDBConnectionManagerUsecase(dbConRepo database.DBConnectionManagerRepository) database.DBConnectionManagerUsecase {
	return &dbConnectionManagerUsecase{
		dbConRepo: dbConRepo,
	}
}
