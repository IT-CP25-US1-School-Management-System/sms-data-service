package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	helperModel "github.com/GodeFvt/go-backend/helper/models"
	"github.com/GodeFvt/go-backend/psql"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/constants"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/entity"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/filter"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/data/v1"
	"github.com/jmoiron/sqlx"

	"github.com/lib/pq"
)

type psqlDatasetRepository struct {
	client *psql.Client
}

// ExistDatasetByID implements data.PsqlDatasetRepository.
func (p *psqlDatasetRepository) ExistDatasetByID(ctx context.Context, datasetID string) (bool, error) {
	query := `
		SELECT
			COUNT(id)
		FROM
			datasets
		WHERE
			id = $1
	`

	var count int
	if err := p.client.GetClient().QueryRowxContext(ctx, query, datasetID).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (p *psqlDatasetRepository) deleteDatasetByID(ctx context.Context, tx *sqlx.Tx, datasetID string) error {
	query := `
		DELETE FROM datasets
		WHERE id = $1
	`
	if _, err := tx.ExecContext(ctx, query, datasetID); err != nil {
		return err
	}
	return nil
}

// DeleteDatasetByID implements data.PsqlDatasetRepository.
func (p *psqlDatasetRepository) DeleteDatasetByID(ctx context.Context, datasetID string) error {
	tx, err := p.client.GetClient().BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := p.deleteDatasetByID(ctx, tx, datasetID); err != nil {
		return err
	}

	return tx.Commit()
}

func (p *psqlDatasetRepository) upsertDataset(ctx context.Context, tx *sqlx.Tx, dataset *entity.Datasets) error {
	query := `
		INSERT INTO datasets (id, name, domain, owner, sensitivity, has_pii, tags, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			domain = EXCLUDED.domain,
			owner = EXCLUDED.owner,
			sensitivity = EXCLUDED.sensitivity,
			has_pii = EXCLUDED.has_pii,
			tags = EXCLUDED.tags,
			description = EXCLUDED.description,
			updated_at = EXCLUDED.updated_at
	`
	if _, err := tx.ExecContext(ctx, query,
		dataset.ID,
		dataset.Name,
		dataset.Domain,
		dataset.Owner,
		dataset.Sensitivity,
		dataset.HasPii,
		pq.Array(dataset.Tags),
		dataset.Description,
		dataset.CreatedAt,
		dataset.UpdatedAt,
	); err != nil {
		return err
	}
	return nil
}

// UpsertDataset implements data.PsqlDatasetRepository.
func (p *psqlDatasetRepository) UpsertDataset(ctx context.Context, dataset *entity.Datasets) error {
	tx, err := p.client.GetClient().BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := p.upsertDataset(ctx, tx, dataset); err != nil {
		return err
	}

	return tx.Commit()
}

// FetchDatasetByID implements data.PsqlDatasetRepository.
func (p *psqlDatasetRepository) FetchDatasetByID(ctx context.Context, datasetID string) (*entity.Datasets, error) {
	query := `
		SELECT
			id, name, domain, owner, sensitivity, has_pii, tags, description, created_at, updated_at
		FROM datasets
		WHERE id = $1
	`
	var data entity.Datasets
	row := p.client.GetClient().QueryRowxContext(ctx, query, datasetID)
	err := row.Scan(
		&data.ID,
		&data.Name,
		&data.Domain,
		&data.Owner,
		&data.Sensitivity,
		&data.HasPii,
		pq.Array(&data.Tags),
		&data.Description,
		&data.CreatedAt,
		&data.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &data, nil
}

// FetchColumnsList implements data.PsqlDatasetRepository.
func (p *psqlDatasetRepository) FetchColumnsList(ctx context.Context, filter *filter.ColumnsFilter, paginator *helperModel.Paginator) ([]*entity.Columns, error) {
	var (
		conds    []string
		valArgs  []interface{}
		where    string
		limitSQL string
	)
	if filter != nil {
		if filter.SourceID != nil {
			conds = append(conds, "source_id=?")
			valArgs = append(valArgs, filter.SourceID)
		}
		if filter.Schema != "" {
			conds = append(conds, "schema=?")
			valArgs = append(valArgs, filter.Schema)
		}
		if filter.Table != "" {
			conds = append(conds, "table_name=?") // table column is table_name
			valArgs = append(valArgs, filter.Table)
		}
	}
	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ")
	}

	query := fmt.Sprintf(`SELECT id, source_id, schema, table_name, column_name, data_type, is_nullable, column_default, ordinal_position, created_at FROM physical_columns %s ORDER BY id ASC %s`, where, limitSQL)
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
	var columns []*entity.Columns
	for rows.Next() {
		var column entity.Columns
		if err := rows.Scan(
			&column.ID,
			&column.SourceID,
			&column.Schema,
			&column.TableName,
			&column.ColumnsName,
			&column.DataType,
			&column.IsNullable,
			&column.ColumnDefault,
			&column.OrdinalPosition,
			&column.CreatedAt,
		); err != nil {
			return nil, err
		}
		columns = append(columns, &column)
	}
	var totalRows int
	countArgs := append([]interface{}(nil), valArgs...)
	if paginator != nil {
		countArgs = countArgs[:len(countArgs)-2]
	}
	countSQL := fmt.Sprintf(`SELECT COUNT(*) FROM physical_columns %s`, where)
	countSQL = sqlx.Rebind(sqlx.DOLLAR, countSQL)

	if err := p.client.GetClient().GetContext(ctx, &totalRows, countSQL, countArgs...); err != nil {
		return nil, err
	}

	if paginator != nil {
		paginator.SetPaginatorByAllRows(totalRows)
	}
	return columns, nil
}

// FetchTablesList implements data.PsqlDatasetRepository.
func (p *psqlDatasetRepository) FetchTablesList(ctx context.Context, filter *filter.TablesFilter, paginator *helperModel.Paginator) ([]*entity.Tables, error) {
	var (
		conds    []string
		valArgs  []interface{}
		where    string
		limitSQL string
	)
	if filter != nil {
		if filter.SourceID != nil {
			conds = append(conds, "source_id=?")
			valArgs = append(valArgs, filter.SourceID)
		}
		if filter.Schema != "" {
			conds = append(conds, "schema=?")
			valArgs = append(valArgs, filter.Schema)
		}
	}
	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ")
	}
	if paginator != nil {
		limitSQL = `
			LIMIT ?
			OFFSET ?
		`
		valArgs = append(valArgs, paginator.GetLimit(), paginator.GetOffset())
	}
	query := fmt.Sprintf(`SELECT id,source_id,schema,table_name,created_at FROM physical_tables %s ORDER BY table_name ASC %s`, where, limitSQL)
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
	var tables []*entity.Tables
	for rows.Next() {
		var table entity.Tables
		if err := rows.Scan(
			&table.ID,
			&table.SourceID,
			&table.Schema,
			&table.TableName,
			&table.CreatedAt,
		); err != nil {
			return nil, err
		}
		tables = append(tables, &table)
	}
	var totalRows int
	countArgs := append([]interface{}(nil), valArgs...)
	if paginator != nil {
		countArgs = countArgs[:len(countArgs)-2]
	}
	countSQL := fmt.Sprintf(`SELECT COUNT(*) FROM physical_tables %s`, where)
	countSQL = sqlx.Rebind(sqlx.DOLLAR, countSQL)

	if err := p.client.GetClient().GetContext(ctx, &totalRows, countSQL, countArgs...); err != nil {
		return nil, err
	}

	if paginator != nil {
		paginator.SetPaginatorByAllRows(totalRows)
	}

	return tables, nil
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
func (p *psqlDatasetRepository) FetchSchemasList(ctx context.Context, filter *filter.SchemasFilter, paginator *helperModel.Paginator) ([]*entity.Schemas, error) {
	var (
		conds    []string
		valArgs  []interface{}
		where    string
		limitSQL string
	)
	if filter != nil {
		if filter.SourceID != nil {
			conds = append(conds, "source_id=?")
			valArgs = append(valArgs, filter.SourceID)
		}
	}
	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ")
	}
	if paginator != nil {
		limitSQL = `
			LIMIT ?
			OFFSET ?
		`
		valArgs = append(valArgs, paginator.GetLimit(), paginator.GetOffset())
	}

	query := fmt.Sprintf(`SELECT id,source_id,schema,created_at FROM physical_schemas %s ORDER BY schema ASC %s`, where, limitSQL)
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
			&schema.CreatedAt,
		); err != nil {
			return nil, err
		}
		schemas = append(schemas, &schema)
	}
	var totalRows int
	countArgs := append([]interface{}(nil), valArgs...)
	if paginator != nil {
		countArgs = countArgs[:len(countArgs)-2]
	}
	countSQL := fmt.Sprintf(`SELECT COUNT(*) FROM physical_schemas %s`, where)
	countSQL = sqlx.Rebind(sqlx.DOLLAR, countSQL)

	if err := p.client.GetClient().GetContext(ctx, &totalRows, countSQL, countArgs...); err != nil {
		return nil, err
	}

	if paginator != nil {
		paginator.SetPaginatorByAllRows(totalRows)
	}

	return schemas, nil
}

// FetchSourceList implements data.PsqlDatasetRepository.
func (p *psqlDatasetRepository) FetchSourceList(ctx context.Context, paginator *helperModel.Paginator) ([]*entity.Sources, error) {
	var (
		valArgs  []interface{}
		limitSQL string
	)
	if paginator != nil {
		limitSQL = `
			LIMIT ?
			OFFSET ?
		`
		valArgs = append(valArgs, paginator.GetLimit(), paginator.GetOffset())
	}
	query := fmt.Sprintf(`SELECT id, name, type, description, created_at FROM sources %s`, limitSQL)
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
	var sources []*entity.Sources
	for rows.Next() {
		var source entity.Sources
		if err := rows.Scan(
			&source.ID,
			&source.Name,
			&source.Type,
			&source.Description,
			&source.CreatedAt,
		); err != nil {
			return nil, err
		}
		sources = append(sources, &source)
	}
	var totalRows int
	countSQL := `SELECT COUNT(*) FROM sources`
	if err := p.client.GetClient().GetContext(ctx, &totalRows, countSQL); err != nil {
		return nil, err
	}
	if paginator != nil {
		paginator.SetPaginatorByAllRows(totalRows)
	}

	return sources, nil
}

func NewPsqlDatasetRepository(client *psql.Client) data.PsqlDatasetRepository {
	return &psqlDatasetRepository{
		client: client,
	}
}
