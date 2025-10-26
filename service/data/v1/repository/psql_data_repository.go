package repository

import (
	"github.com/GodeFvt/go-backend/psql"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/data/v1"
)

type psqlDataRepository struct {
	client *psql.Client
}

func NewPsqlDataRepository(client *psql.Client) data.PsqlDataRepository {
	return &psqlDataRepository{
		client: client,
	}
}
