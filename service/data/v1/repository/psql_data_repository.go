package repository

import (
	"context"
	"fmt"

	helperModel "github.com/GodeFvt/go-backend/helper/models"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/entity"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/data/v1"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/database/v1"
	"github.com/gofrs/uuid"
)

type psqlDataRepository struct {
	dbConnectionManager database.DBConnectionManagerUsecase
}

func NewPsqlDataRepository(dbConnectionManager database.DBConnectionManagerUsecase) data.PsqlDataRepository {
	return &psqlDataRepository{
		dbConnectionManager: dbConnectionManager,
	}
}

func (p *psqlDataRepository) FetchInformationTablesBySourceID(ctx context.Context, dbType string, sourceID *uuid.UUID) ([]*entity.Tables, error) {
	client, err := p.dbConnectionManager.GetConnection(ctx, *sourceID)
	if err != nil {
		return nil, err
	}

	if sourceID == nil {
		return nil, fmt.Errorf("sourceID is nil")
	}

	var query string
	switch dbType {
	case "mysql":
		query = `
			SELECT table_schema, table_name
			FROM information_schema.tables
			WHERE table_type = 'BASE TABLE'
			  AND table_schema NOT IN ('information_schema', 'performance_schema', 'mysql', 'sys')
			ORDER BY table_schema, table_name
		`
	default:
		query = `
			SELECT table_schema, table_name
			FROM information_schema.tables
			WHERE table_type = 'BASE TABLE'
			  AND table_schema NOT IN ('pg_catalog', 'information_schema')
			ORDER BY table_schema, table_name
		`
	}

	rows, err := client.GetClient().QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []*entity.Tables
	for rows.Next() {
		var schema string
		var tableName string
		if err := rows.Scan(&schema, &tableName); err != nil {
			return nil, err
		}
		t := &entity.Tables{
			SourceID:  sourceID,
			Schema:    schema,
			TableName: tableName,
			CreatedAt: nil,
		}
		t.GenUUID()
		now := helperModel.NewTimestampFromNow()
		t.CreatedAt = &now
		tables = append(tables, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tables, nil
}

func (p *psqlDataRepository) FetchInformationColumnsBySourceID(ctx context.Context, dbType string, sourceID *uuid.UUID) ([]*entity.Columns, error) {
	client, err := p.dbConnectionManager.GetConnection(ctx, *sourceID)
	if err != nil {
		return nil, err
	}

	if sourceID == nil {
		return nil, fmt.Errorf("sourceID is nil")
	}

	var query string
	switch dbType {
	case "mysql":
		query = `
			SELECT 
				table_schema, 
				table_name, 
				column_name, 
				data_type, 
				is_nullable, 
				column_default, 
				ordinal_position
			FROM information_schema.columns
			WHERE table_schema NOT IN ('information_schema', 'performance_schema', 'mysql', 'sys')
			ORDER BY table_schema, table_name, ordinal_position
		`
	default:
		query = `
			SELECT 
				table_schema, 
				table_name, 
				column_name, 
				data_type, 
				is_nullable, 
				column_default, 
				ordinal_position
			FROM information_schema.columns
			WHERE table_schema NOT IN ('pg_catalog', 'information_schema')
			ORDER BY table_schema, table_name, ordinal_position
		`
	}

	rows, err := client.GetClient().QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []*entity.Columns
	for rows.Next() {
		var schema, tableName, columnName, dataType, isNullableStr string
		var columnDefault *string
		var ordinalPosition int

		if err := rows.Scan(&schema, &tableName, &columnName, &dataType, &isNullableStr, &columnDefault, &ordinalPosition); err != nil {
			return nil, err
		}

		isNullable := isNullableStr == "YES"

		c := &entity.Columns{
			SourceID:        sourceID,
			Schema:          schema,
			TableName:       tableName,
			ColumnsName:     columnName,
			DataType:        dataType,
			IsNullable:      isNullable,
			ColumnDefault:   columnDefault,
			OrdinalPosition: &ordinalPosition,
			CreatedAt:       nil,
		}
		c.GenUUID()
		now := helperModel.NewTimestampFromNow()
		c.CreatedAt = &now
		columns = append(columns, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return columns, nil
}

func (p *psqlDataRepository) FetchInformationSchemasBySourceID(ctx context.Context, dbType string, sourceID *uuid.UUID) ([]*entity.Schemas, error) {
	client, err := p.dbConnectionManager.GetConnection(ctx, *sourceID)
	if err != nil {
		return nil, err
	}

	if sourceID == nil {
		return nil, fmt.Errorf("sourceID is nil")
	}

	var query string
	switch dbType {
	case "mysql":
		query = `
			SELECT schema_name
			FROM information_schema.schemata
			WHERE schema_name NOT IN ('information_schema', 'performance_schema', 'mysql', 'sys')
			ORDER BY schema_name
		`
	default:
		query = `
			SELECT schema_name
			FROM information_schema.schemata
			WHERE schema_name NOT IN ('pg_catalog', 'information_schema', 'pg_toast', 'pg_temp_1', 'pg_toast_temp_1')
			ORDER BY schema_name
		`
	}

	rows, err := client.GetClient().QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schemas []*entity.Schemas
	for rows.Next() {
		var schemaName string
		if err := rows.Scan(&schemaName); err != nil {
			return nil, err
		}
		s := &entity.Schemas{
			SourceID:  sourceID,
			Schema:    schemaName,
			CreatedAt: nil,
		}
		s.GenUUID()
		now := helperModel.NewTimestampFromNow()
		s.CreatedAt = &now
		schemas = append(schemas, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return schemas, nil
}
