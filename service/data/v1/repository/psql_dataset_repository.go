package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/BlackMocca/sqlx"
	helperModel "github.com/GodeFvt/go-backend/helper/models"
	"github.com/GodeFvt/go-backend/psql"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/constants"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/entity"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/filter"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/data/v1"
	"github.com/lib/pq"
)

type psqlDatasetRepository struct {
	client *psql.Client
}

// FetchDatasetList implements data.PsqlDatasetRepository.
func (p *psqlDatasetRepository) FetchDatasetList(ctx context.Context, filter *filter.DatasetsFilter, paginator *helperModel.Paginator) ([]*entity.Datasets, error) {
	var (
		conds    []string
		valArgs  []interface{}
		orderBy  string
		where    string
		limitSQL string
	)

	if filter != nil {
		if filter.Domain != "" {
			conds = append(conds, "domain = ?")
			valArgs = append(valArgs, filter.Domain)
		}
		if filter.Owner != "" {
			conds = append(conds, "owner = ?")
			valArgs = append(valArgs, filter.Owner)
		}
		if filter.SearchWord != "" {
			sw := fmt.Sprintf("%%%s%%", strings.ToLower(strings.ReplaceAll(filter.SearchWord, " ", "")))
			conds = append(conds, "(LOWER(REPLACE(name,' ','')) LIKE ? OR LOWER(REPLACE(description,' ','')) LIKE ?)")
			valArgs = append(valArgs, sw, sw)
		}
		if len(filter.Tags) > 0 {
			conds = append(conds, "COALESCE(tags, '{}'::text[]) && ?::text[]")
			valArgs = append(valArgs, pq.Array(filter.Tags))
		}
		if filter.HasPii != nil {
			conds = append(conds, "has_pii = ?")
			valArgs = append(valArgs, *filter.HasPii)
		}

		validSortColumns := map[string]bool{
			constants.DATASET_SORT_BY_NAME:       true,
			constants.DATASET_SORT_BY_CREATED_AT: true,
			constants.DATASET_SORT_BY_UPDATED_AT: true,
		}
		validSortOrders := map[string]bool{
			constants.SORT_ORDER_ASC:  true,
			constants.SORT_ORDER_DESC: true,
		}
		if filter.SortBy != "" && validSortColumns[filter.SortBy] {
			order := constants.SORT_ORDER_DESC
			if validSortOrders[filter.SortOrder] {
				order = filter.SortOrder
			}
			orderBy = fmt.Sprintf("ORDER BY %s %s", filter.SortBy, order)
		}
	}

	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ")
	}
	if orderBy == "" {
		orderBy = "ORDER BY created_at DESC"
	}
	if paginator != nil {
		limitSQL = `
			LIMIT ?
			OFFSET ?
		`
		valArgs = append(valArgs, paginator.GetLimit(), paginator.GetOffset())
	}

	dataSQL := fmt.Sprintf(`
		SELECT
			id, name, domain, owner, sensitivity, has_pii, tags, description, created_at, updated_at
		FROM datasets
		%s
		%s
		%s
	`, where, orderBy, limitSQL)
	dataSQL = sqlx.Rebind(sqlx.DOLLAR, dataSQL)

	rows, err := p.client.GetClient().QueryxContext(ctx, dataSQL, valArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var datasets []*entity.Datasets
	for rows.Next() {
		var d entity.Datasets
		if err := rows.Scan(
			&d.ID,
			&d.Name,
			&d.Domain,
			&d.Owner,
			&d.Sensitivity,
			&d.HasPii,
			pq.Array(&d.Tags),
			&d.Description,
			&d.CreatedAt,
			&d.UpdatedAt,
		); err != nil {
			return nil, err
		}
		datasets = append(datasets, &d)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	var totalRows int
	countArgs := append([]interface{}(nil), valArgs...)
	if paginator != nil {
		countArgs = countArgs[:len(countArgs)-2]
	}
	countSQL := fmt.Sprintf(`SELECT COUNT(*) FROM datasets %s`, where)
	countSQL = sqlx.Rebind(sqlx.DOLLAR, countSQL)

	if err := p.client.GetClient().GetContext(ctx, &totalRows, countSQL, countArgs...); err != nil {
		return nil, err
	}

	if paginator != nil {
		paginator.SetPaginatorByAllRows(totalRows)
	}

	return datasets, nil
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
