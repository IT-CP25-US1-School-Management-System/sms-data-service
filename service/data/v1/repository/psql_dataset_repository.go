package repository

import (
	"context"

	"github.com/GodeFvt/go-backend/psql"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/entity"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/data/v1"
)

type psqlDatasetRepository struct {
	client *psql.Client
}

// FetchSourceList implements data.PsqlDatasetRepository.
func (p *psqlDatasetRepository) FetchSourceList(ctx context.Context) ([]*entity.Sources, error) {
	query := `SELECT id, name, type, connection_ref, sensitivity,config,created_at FROM sources`
	stmt, err := p.client.GetClient().PreparexContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryxContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sources []*entity.Sources
	for rows.Next() {
		var source entity.Sources
		if err := rows.Scan(
			&source.ID,
			&source.Name,
			&source.Type,
			&source.ConnectionRef,
			&source.Sensitivity,
			&source.Config,
			&source.CreatedAt,
		); err != nil {
			return nil, err
		}
		sources = append(sources, &source)
	}
	return sources, nil
}

func NewPsqlDatasetRepository(client *psql.Client) data.PsqlDatasetRepository {
	return &psqlDatasetRepository{
		client: client,
	}
}
