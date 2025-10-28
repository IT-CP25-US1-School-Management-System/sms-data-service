package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/BlackMocca/sqlx"
	"github.com/GodeFvt/go-backend/psql"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/entity"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/filter"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/data/v1"
)

type psqlDatasetRepository struct {
	client *psql.Client
}

// FetchSchemasList implements data.PsqlDatasetRepository.
func (p *psqlDatasetRepository) FetchSchemasList(ctx context.Context, filter *filter.SchemasFilter) ([]*entity.Schemas, error) {
	var where string
	var conds = make([]string, 0)
	var valArgs = make([]interface{}, 0)
	if filter != nil {
		if filter.SourceID != nil {
			conds = append(conds, "source_id=?")
			valArgs = append(valArgs, filter.SourceID)
		}
	}
	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ")
	}

	query := fmt.Sprintf(`SELECT id,source_id,schema,discovered_at FROM physical_schemas %s`, where)
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	stmt, err := p.client.GetClient().PreparexContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryxContext(ctx, valArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var schemas []*entity.Schemas
	for rows.Next() {
		var schema entity.Schemas
		if err := rows.Scan(
			&schema.ID,
			&schema.SourceID,
			&schema.Schema,
			&schema.Discovered_at,
		); err != nil {
			return nil, err
		}
		schemas = append(schemas, &schema)
	}
	return schemas, nil
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
