package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	helperModel "github.com/GodeFvt/go-backend/helper/models"
	"github.com/GodeFvt/go-backend/psql"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/constants"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/entity"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/filter"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/data/v1"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/lib/pq"
)

type psqlDatasetRepository struct {
	client *psql.Client
}

func (p *psqlDatasetRepository) deleteSourceByID(ctx context.Context, tx *sqlx.Tx, sourceID string) error {
	query := `
		DELETE FROM sources WHERE id = $1
	`
	if _, err := tx.ExecContext(ctx, query, sourceID); err != nil {
		return err
	}
	return nil
}

// DeleteSourceByID implements data.PsqlDatasetRepository.
func (p *psqlDatasetRepository) DeleteSourceByID(ctx context.Context, sourceID *uuid.UUID) error {
	tx, err := p.client.GetClient().BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := p.deleteSourceByID(ctx, tx, sourceID.String()); err != nil {
		return err
	}

	return tx.Commit()
}

func (p *psqlDatasetRepository) batchInsertPhysicalSchemas(ctx context.Context, tx *sqlx.Tx, schemas []*entity.Schemas) error {
	valueStrings := make([]string, 0, len(schemas))
	valueArgs := make([]interface{}, 0, len(schemas)*4)

	for _, schema := range schemas {
		valueStrings = append(valueStrings, "(?, ?, ?, ?)")
		valueArgs = append(valueArgs, schema.ID, schema.SourceID, schema.Schema, schema.CreatedAt)
	}

	query := `INSERT INTO physical_schemas (id, source_id, schema, created_at) VALUES ` + strings.Join(valueStrings, ",") + `;`

	query = sqlx.Rebind(sqlx.DOLLAR, query)
	if _, err := tx.ExecContext(ctx, query, valueArgs...); err != nil {
		return err
	}
	return nil
}

func (p *psqlDatasetRepository) batchInsertPhysicalTables(ctx context.Context, tx *sqlx.Tx, tables []*entity.Tables) error {
	valueStrings := make([]string, 0, len(tables))
	valueArgs := make([]interface{}, 0, len(tables)*5)

	for _, table := range tables {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, table.ID, table.SourceID, table.Schema, table.TableName, table.CreatedAt)
	}

	query := `INSERT INTO physical_tables (id, source_id, schema, table_name, created_at) VALUES ` + strings.Join(valueStrings, ",") + `;`

	query = sqlx.Rebind(sqlx.DOLLAR, query)
	if _, err := tx.ExecContext(ctx, query, valueArgs...); err != nil {
		return err
	}
	return nil
}

func (p *psqlDatasetRepository) batchInsertPhysicalColumns(ctx context.Context, tx *sqlx.Tx, columns []*entity.Columns) error {
	valueStrings := make([]string, 0, len(columns))
	valueArgs := make([]interface{}, 0, len(columns)*10)

	for _, column := range columns {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, column.ID, column.SourceID, column.Schema, column.TableName, column.ColumnsName, column.DataType, column.IsNullable, column.ColumnDefault, column.OrdinalPosition, column.CreatedAt)
	}

	query := `INSERT INTO physical_columns (id, source_id, schema, table_name, column_name, data_type, is_nullable, column_default, ordinal_position, created_at) VALUES ` + strings.Join(valueStrings, ",") + `;`

	query = sqlx.Rebind(sqlx.DOLLAR, query)
	if _, err := tx.ExecContext(ctx, query, valueArgs...); err != nil {
		return err
	}
	return nil
}

func (p *psqlDatasetRepository) batchInsertPhysicalTableRelations(ctx context.Context, tx *sqlx.Tx, tableRelations []*entity.TableRelations) error {
	valueStrings := make([]string, 0, len(tableRelations))
	valueArgs := make([]interface{}, 0, len(tableRelations)*7)

	for _, tr := range tableRelations {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, tr.ID, tr.SourceID, tr.TableFrom, tr.ColumnFrom, tr.TableTo, tr.ColumnTo, tr.CreatedAt)
	}

	query := `INSERT INTO table_relations (id, source_id, table_from, column_from, table_to, column_to, created_at) VALUES ` + strings.Join(valueStrings, ",") + `;`

	query = sqlx.Rebind(sqlx.DOLLAR, query)
	if _, err := tx.ExecContext(ctx, query, valueArgs...); err != nil {
		return err
	}
	return nil
}

// BatchInsertInformationDatabase implements data.PsqlDatasetRepository.
func (p *psqlDatasetRepository) BatchInsertInformationDatabase(ctx context.Context, schemas []*entity.Schemas, tables []*entity.Tables, columns []*entity.Columns, tableRelations []*entity.TableRelations) error {
	tx, err := p.client.GetClient().BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if len(schemas) > 0 {
		if err := p.batchInsertPhysicalSchemas(ctx, tx, schemas); err != nil {
			return err
		}
	}
	if len(tables) > 0 {
		if err := p.batchInsertPhysicalTables(ctx, tx, tables); err != nil {
			return err
		}
	}
	if len(columns) > 0 {
		if err := p.batchInsertPhysicalColumns(ctx, tx, columns); err != nil {
			return err
		}
	}
	if len(tableRelations) > 0 {
		if err := p.batchInsertPhysicalTableRelations(ctx, tx, tableRelations); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// ExistSourceByID implements data.PsqlDatasetRepository.
func (p *psqlDatasetRepository) ExistSourceByID(ctx context.Context, sourceID *uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM sources
			WHERE id = $1
		)
	`
	var exists bool
	if err := p.client.GetClient().QueryRowxContext(ctx, query, sourceID).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (p *psqlDatasetRepository) ExistSchemaByName(ctx context.Context, sourceID *uuid.UUID, schemaName string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM physical_schemas
			WHERE source_id = $1 AND schema = $2
		)
	`
	var exists bool
	if err := p.client.GetClient().QueryRowxContext(ctx, query, sourceID, schemaName).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (p *psqlDatasetRepository) ExistTableByName(ctx context.Context, sourceID *uuid.UUID, schemaName, tableName string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM physical_tables
			WHERE source_id = $1 AND schema = $2 AND table_name = $3
		)
	`
	var exists bool
	if err := p.client.GetClient().QueryRowxContext(ctx, query, sourceID, schemaName, tableName).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (p *psqlDatasetRepository) ExistColumnByName(ctx context.Context, sourceID *uuid.UUID, schemaName, tableName, columnName string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM physical_columns
			WHERE source_id = $1 AND schema = $2 AND table_name = $3 AND column_name = $4
		)
	`
	var exists bool
	if err := p.client.GetClient().QueryRowxContext(ctx, query, sourceID, schemaName, tableName, columnName).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

// ActivateSourceByID implements data.PsqlDatasetRepository.
func (p *psqlDatasetRepository) ActivateSourceByID(ctx context.Context, sourceID *uuid.UUID) error {
	query := `
		UPDATE sources
		SET is_active = TRUE
		WHERE id = $1
	`
	if _, err := p.client.GetClient().ExecContext(ctx, query, sourceID); err != nil {
		return err
	}
	return nil
}

// DeactivateSourceByID implements data.PsqlDatasetRepository.
func (p *psqlDatasetRepository) DeactivateSourceByID(ctx context.Context, sourceID *uuid.UUID) error {
	query := `
		UPDATE sources
		SET is_active = FALSE
		WHERE id = $1
	`
	if _, err := p.client.GetClient().ExecContext(ctx, query, sourceID); err != nil {
		return err
	}
	return nil
}

func (p *psqlDatasetRepository) upsertSource(ctx context.Context, tx *sqlx.Tx, source *entity.Sources) error {
	query := `
		INSERT INTO sources (id, name, description, is_active, type, db_type, created_at, updated_at,
			host, port, username, password, database_name, params, sensitivity)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			is_active = EXCLUDED.is_active,
			type = EXCLUDED.type,
			db_type = EXCLUDED.db_type,
			updated_at = EXCLUDED.updated_at,
			host = EXCLUDED.host,
			port = EXCLUDED.port,
			username = EXCLUDED.username,
			password = EXCLUDED.password,
			database_name = EXCLUDED.database_name,
			params = EXCLUDED.params,
			sensitivity = EXCLUDED.sensitivity
	`
	if _, err := tx.ExecContext(ctx, query,
		source.ID,
		source.Name,
		source.Description,
		source.IsActive,
		source.Type,
		source.DBType,
		source.CreatedAt,
		source.UpdatedAt,
		source.Host,
		source.Port,
		source.Username,
		source.Password,
		source.DatabaseName,
		source.Params,
		source.Sensitivity,
	); err != nil {
		return err
	}
	return nil
}

// UpsertSource implements data.PsqlDatasetRepository.
func (p *psqlDatasetRepository) UpsertSource(ctx context.Context, source *entity.Sources) error {
	tx, err := p.client.GetClient().BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := p.upsertSource(ctx, tx, source); err != nil {
		return err
	}

	return tx.Commit()
}

// FetchSourceByID implements data.PsqlDatasetRepository.
func (p *psqlDatasetRepository) FetchSourceByID(ctx context.Context, sourceID *uuid.UUID) (*entity.Sources, error) {
	query := `
		SELECT
			id, name, description, is_active, sensitivity, type, db_type, created_at, updated_at,
			host, port, username, password, database_name, params
		FROM sources
		WHERE id = $1
	`
	var source entity.Sources
	row := p.client.GetClient().QueryRowxContext(ctx, query, sourceID)
	err := row.Scan(
		&source.ID,
		&source.Name,
		&source.Description,
		&source.IsActive,
		&source.Sensitivity,
		&source.Type,
		&source.DBType,
		&source.CreatedAt,
		&source.UpdatedAt,
		&source.Host,
		&source.Port,
		&source.Username,
		&source.Password,
		&source.DatabaseName,
		&source.Params,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &source, nil
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
		if paginator != nil {
			limitSQL = `
			LIMIT ?
			OFFSET ?
		`
			valArgs = append(valArgs, paginator.GetLimit(), paginator.GetOffset())
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
			conds = append(conds, "COALESCE(tags, '{}') @> ?")
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
			sortOrder := strings.ToUpper(filter.SortOrder)
			if validSortOrders[sortOrder] {
				order = sortOrder
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
	query := fmt.Sprintf(`SELECT id, name, description, type, is_active, sensitivity, db_type, host, port, username, database_name, params, created_at, updated_at FROM sources %s`, limitSQL)
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
			&source.Description,
			&source.Type,
			&source.IsActive,
			&source.Sensitivity,
			&source.DBType,
			&source.Host,
			&source.Port,
			&source.Username,
			&source.DatabaseName,
			&source.Params,
			&source.CreatedAt,
			&source.UpdatedAt,
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

// FetchDatasetVersionByID implements data.PsqlDatasetRepository.
func (p *psqlDatasetRepository) FetchDatasetVersionByID(ctx context.Context, datasetID string, version string) (*entity.DatasetVersion, error) {
	query := `
		SELECT
			dataset_id, version, status, schema_json, access_policies, policies, source_id
		FROM dataset_versions
		WHERE dataset_id = $1 AND version = $2
	`
	var data entity.DatasetVersion
	var schemaJSON, accessPoliciesJSON, policiesJSON []byte

	row := p.client.GetClient().QueryRowxContext(ctx, query, datasetID, version)
	err := row.Scan(
		&data.DatasetID,
		&data.Version,
		&data.Status,
		&schemaJSON,
		&accessPoliciesJSON,
		&policiesJSON,
		&data.SourceID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Unmarshal JSON fields
	if err := json.Unmarshal(schemaJSON, &data.Schema); err != nil {
		return nil, fmt.Errorf("failed to unmarshal schema: %w", err)
	}

	if err := json.Unmarshal(accessPoliciesJSON, &data.AccessPolicies); err != nil {
		return nil, fmt.Errorf("failed to unmarshal access_policies: %w", err)
	}

	if err := json.Unmarshal(policiesJSON, &data.Policies); err != nil {
		return nil, fmt.Errorf("failed to unmarshal policies: %w", err)
	}

	return &data, nil
}

// FetchDatasetVersionsList implements data.PsqlDatasetRepository.
func (p *psqlDatasetRepository) FetchDatasetVersionsList(ctx context.Context, datasetID string, filter *filter.DatasetVersionsFilter, paginator *helperModel.Paginator) ([]*entity.DatasetVersion, error) {
	var (
		conds    []string
		valArgs  []interface{}
		where    string
		limitSQL string
	)
	if datasetID != "" {
		conds = append(conds, "dataset_id = ?")
		valArgs = append(valArgs, datasetID)
	}
	if filter != nil {
		if filter.SourceID != "" {
			conds = append(conds, "source_id = ?")
			valArgs = append(valArgs, filter.SourceID)
		}
		if filter.SearchWord != "" {
			sw := fmt.Sprintf("%%%s%%", strings.ToLower(strings.ReplaceAll(filter.SearchWord, " ", "")))
			conds = append(conds, "LOWER(REPLACE(version,' ','')) LIKE ?")
			valArgs = append(valArgs, sw)
		}
		if filter.Status != "" {
			conds = append(conds, "status = ?")
			valArgs = append(valArgs, filter.Status)
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
	query := fmt.Sprintf(`
		SELECT dataset_id, version, status, schema_json, access_policies, policies, source_id
		FROM dataset_versions
		%s
		ORDER BY version DESC %s
	`, where, limitSQL)
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	stmt, err := p.client.GetClient().PreparexContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := p.client.GetClient().QueryxContext(ctx, query, valArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []*entity.DatasetVersion
	for rows.Next() {
		var data entity.DatasetVersion
		var schemaJSON, accessPoliciesJSON, policiesJSON []byte

		if err := rows.Scan(
			&data.DatasetID,
			&data.Version,
			&data.Status,
			&schemaJSON,
			&accessPoliciesJSON,
			&policiesJSON,
			&data.SourceID,
		); err != nil {
			return nil, err
		}

		// Unmarshal JSON fields
		if len(schemaJSON) > 0 {
			if err := json.Unmarshal(schemaJSON, &data.Schema); err != nil {
				return nil, fmt.Errorf("failed to unmarshal schema for version %s: %w", data.Version, err)
			}
		}
		if len(accessPoliciesJSON) > 0 {
			if err := json.Unmarshal(accessPoliciesJSON, &data.AccessPolicies); err != nil {
				return nil, fmt.Errorf("failed to unmarshal access_policies for version %s: %w", data.Version, err)
			}
		}
		if len(policiesJSON) > 0 {
			if err := json.Unmarshal(policiesJSON, &data.Policies); err != nil {
				return nil, fmt.Errorf("failed to unmarshal policies for version %s: %w", data.Version, err)
			}
		}

		versions = append(versions, &data)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Get total count for pagination
	if paginator != nil {
		var totalRows int
		countArgs := append([]interface{}(nil), valArgs...)
		// remove LIMIT/OFFSET from args if present
		if len(countArgs) >= 2 {
			countArgs = countArgs[:len(countArgs)-2]
		}
		countSQL := fmt.Sprintf(`SELECT COUNT(*) FROM dataset_versions %s`, where)
		countSQL = sqlx.Rebind(sqlx.DOLLAR, countSQL)
		if err := p.client.GetClient().GetContext(ctx, &totalRows, countSQL, countArgs...); err != nil {
			return nil, err
		}
		paginator.SetPaginatorByAllRows(totalRows)
	}

	return versions, nil
}

func (p *psqlDatasetRepository) upsertDatasetVersion(ctx context.Context, tx *sqlx.Tx, datasetVersion *entity.DatasetVersion) error {
	// Marshal JSON fields
	schemaJSON, err := json.Marshal(datasetVersion.Schema)
	if err != nil {
		return fmt.Errorf("failed to marshal schema: %w", err)
	}

	accessPoliciesJSON, err := json.Marshal(datasetVersion.AccessPolicies)
	if err != nil {
		return fmt.Errorf("failed to marshal access_policies: %w", err)
	}

	policiesJSON, err := json.Marshal(datasetVersion.Policies)
	if err != nil {
		return fmt.Errorf("failed to marshal policies: %w", err)
	}

	query := `
		INSERT INTO dataset_versions (dataset_id, version, status, schema_json, access_policies, policies, source_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (dataset_id, version) DO UPDATE SET
			status = EXCLUDED.status,
			schema_json = EXCLUDED.schema_json,
			access_policies = EXCLUDED.access_policies,
			policies = EXCLUDED.policies,
			source_id = EXCLUDED.source_id,
			updated_at = EXCLUDED.updated_at
	`

	if _, err := tx.ExecContext(ctx, query,
		datasetVersion.DatasetID,
		datasetVersion.Version,
		datasetVersion.Status,
		schemaJSON,
		accessPoliciesJSON,
		policiesJSON,
		datasetVersion.SourceID,
		datasetVersion.CreatedAt,
		datasetVersion.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}

// UpsertDatasetVersion implements data.PsqlDatasetRepository.
func (p *psqlDatasetRepository) UpsertDatasetVersion(ctx context.Context, datasetVersion *entity.DatasetVersion) error {
	tx, err := p.client.GetClient().BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := p.upsertDatasetVersion(ctx, tx, datasetVersion); err != nil {
		return err
	}

	return tx.Commit()
}

func (p *psqlDatasetRepository) deleteDatasetVersionByID(ctx context.Context, tx *sqlx.Tx, datasetID string, version string) error {
	query := `
		DELETE FROM dataset_versions
		WHERE dataset_id = $1 AND version = $2
	`
	if _, err := tx.ExecContext(ctx, query, datasetID, version); err != nil {
		return err
	}
	return nil
}
func (p *psqlDatasetRepository) UpdateDatasetVersionStatus(ctx context.Context, datasetID string, version string, status string) error {
	query := `
		UPDATE dataset_versions
		SET status = $1
		WHERE dataset_id = $2 AND version = $3
	`
	if _, err := p.client.GetClient().ExecContext(ctx, query, status, datasetID, version); err != nil {
		return err
	}
	return nil
}

// DeleteDatasetVersionByID implements data.PsqlDatasetRepository.
func (p *psqlDatasetRepository) DeleteDatasetVersionByID(ctx context.Context, datasetID string, version string) error {
	tx, err := p.client.GetClient().BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := p.deleteDatasetVersionByID(ctx, tx, datasetID, version); err != nil {
		return err
	}

	return tx.Commit()
}

// ExistDatasetVersionByID implements data.PsqlDatasetRepository.
func (p *psqlDatasetRepository) ExistDatasetVersionByID(ctx context.Context, datasetID string, version string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM dataset_versions
			WHERE dataset_id = $1 AND version = $2
		)
	`
	var exists bool
	if err := p.client.GetClient().QueryRowxContext(ctx, query, datasetID, version).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

// FetchExportJobByID implements data.PsqlDatasetRepository.
func (p *psqlDatasetRepository) FetchExportJobByID(ctx context.Context, jobID *uuid.UUID) (*entity.ExportJob, error) {
	query := `
		SELECT
			job_id, dataset_id, view, format, version, destination_uri, status, created_at, completed_at
		FROM export_jobs
		WHERE job_id = $1
	`
	var job entity.ExportJob
	row := p.client.GetClient().QueryRowxContext(ctx, query, jobID)
	err := row.Scan(
		&job.JobId,
		&job.DatasetId,
		&job.View,
		&job.Format,
		&job.Version,
		&job.DestinationUri,
		&job.Status,
		&job.CreatedAt,
		&job.CompletedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &job, nil
}


func (p *psqlDatasetRepository) insertExportJob(ctx context.Context, tx *sqlx.Tx, exportJob *entity.ExportJob) error {
	query := `
		INSERT INTO export_jobs (job_id, dataset_id, view, format, version, destination_uri, status, created_at, completed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	if _, err := tx.ExecContext(ctx, query,
		exportJob.JobId,
		exportJob.DatasetId,
		exportJob.View,
		exportJob.Format,
		exportJob.Version,
		exportJob.DestinationUri,
		exportJob.Status,
		exportJob.CreatedAt,
		exportJob.CompletedAt,
	); err != nil {
		return err
	}
	return nil
}

// InsertExportJob implements data.PsqlDatasetRepository.
func (p *psqlDatasetRepository) InsertExportJob(ctx context.Context, exportJob *entity.ExportJob) error {
	tx, err := p.client.GetClient().BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := p.insertExportJob(ctx, tx, exportJob); err != nil {
		return err
	}

	return tx.Commit()
}


func NewPsqlDatasetRepository(client *psql.Client) data.PsqlDatasetRepository {
	return &psqlDatasetRepository{
		client: client,
	}
}
