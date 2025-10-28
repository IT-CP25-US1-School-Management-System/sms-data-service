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
	var where string
	var conds = make([]string, 0)
	var valArgs = make([]interface{}, 0)
	var orderBy string
	var paginatorSql string

	if filter != nil {
		if filter.Domain != "" {
			conds = append(conds, "domain=?")
			valArgs = append(valArgs, filter.Domain)
		}
		if filter.Owner != "" {
			conds = append(conds, "owner=?")
			valArgs = append(valArgs, filter.Owner)
		}
		if filter.SearchWord != "" {
			filter.SearchWord = fmt.Sprintf("%%%s%%", strings.ToLower(strings.ReplaceAll(filter.SearchWord, " ", "")))
			conds = append(conds, "(LOWER(REPLACE(name, ' ', '')) LIKE ? OR LOWER(REPLACE(description, ' ', '')) LIKE ?)")
			valArgs = append(valArgs, filter.SearchWord, filter.SearchWord)
		}
		if len(filter.Tags) > 0 {
			conds = append(conds, "COALESCE(tags, '{}'::text[]) && ?::text[]")
			valArgs = append(valArgs, pq.Array(filter.Tags))
		}
		if filter.HasPii != nil {
			conds = append(conds, "has_pii=?")
			valArgs = append(valArgs, *filter.HasPii)
		}
		if filter.SortBy != "" && filter.SortOrder != "" {
			if filter.SortOrder == constants.SORT_ORDER_ASC || filter.SortOrder == constants.SORT_ORDER_DESC {
				if filter.SortBy == constants.DATASET_SORT_BY_NAME || filter.SortBy == constants.DATASET_SORT_BY_CREATED_AT || filter.SortBy == constants.DATASET_SORT_BY_UPDATED_AT {
					orderBy = fmt.Sprintf("%s %s", filter.SortBy, filter.SortOrder)
				}
			}
		} else {
			orderBy = "created_at DESC"
		}
	}
	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ")
	}
	if orderBy != "" {
		orderBy = "ORDER BY " + orderBy
	}

	if paginator != nil {
		var limit = paginator.GetLimit()
		var offSet = paginator.GetOffset()
		paginatorSql = fmt.Sprintf(`
			LIMIT %d
			OFFSET %d
			`,
			limit,
			offSet,
		)
	}

	query := fmt.Sprintf(`
    SELECT 
		id, 
		name, 
	    domain, 
	    owner, 
	    sensitivity, 
	    has_pii, 
	    tags, 
	    description, 
	    created_at, 
	    updated_at, 
	    COUNT(*) OVER() as total_row 
	FROM 
	datasets 
	%s 
	%s 
	%s
	`,
		where,
		orderBy,
		paginatorSql,
	)

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

	var totalRow int
	var datasets []*entity.Datasets
	for rows.Next() {
		var dataset entity.Datasets
		if err := rows.Scan(
			&dataset.ID,
			&dataset.Name,
			&dataset.Domain,
			&dataset.Owner,
			&dataset.Sensitivity,
			&dataset.HasPii,
			pq.Array(&dataset.Tags),
			&dataset.Description,
			&dataset.CreatedAt,
			&dataset.UpdatedAt,
			&totalRow,
		); err != nil {
			return nil, err
		}
		datasets = append(datasets, &dataset)
	}

	if paginator != nil && len(datasets) > 0 {
		paginator.SetPaginatorByAllRows(totalRow)
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
