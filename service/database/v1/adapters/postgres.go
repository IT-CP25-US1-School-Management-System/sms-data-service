package adapters

import (
	"context"
	"fmt"
	"net/url"

	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/entity"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/database/v1/client"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// PostgresAdapter provides complete PostgreSQL database support
type PostgresAdapter struct{}

// GetName returns the database type identifier
func (p *PostgresAdapter) GetName() string {
	return "postgres"
}

// GetDriverName returns the SQL driver name
func (p *PostgresAdapter) GetDriverName() string {
	return "pgx"
}

// GetDriverType returns the driver type
func (p *PostgresAdapter) GetDriverType() string {
	return "postgres"
}

// Connect creates a new PostgreSQL connection
func (p *PostgresAdapter) Connect(ctx context.Context, config client.ClientConfig) (*client.Client, error) {
	pool, err := pgxpool.New(ctx, config.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db := sqlx.NewDb(stdlib.OpenDBFromPool(pool), p.GetDriverName())
	db.SetMaxOpenConns(4)

	return client.NewClientFromFields(
		db,
		config.ConnectionString,
		p.GetDriverName(),
		p.GetName(),
		config.Tracer,
	), nil
}

// BuildConnectionString creates a PostgreSQL connection string
func (p *PostgresAdapter) BuildConnectionString(source *entity.Sources, decryptedPassword string) (string, error) {
	user := url.QueryEscape(source.Username)
	pass := url.QueryEscape(decryptedPassword)
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		user, pass, source.Host, source.Port, source.DatabaseName)
	return connStr, nil
}

// GetPlaceholderFormat returns PostgreSQL's placeholder format ($1, $2, $3, ...)
func (p *PostgresAdapter) GetPlaceholderFormat() squirrel.PlaceholderFormat {
	return squirrel.Dollar
}

// BuildJSONObject creates a JSON object using PostgreSQL's JSON_BUILD_OBJECT
func (p *PostgresAdapter) BuildJSONObject(keyValuePairs []string) string {
	return fmt.Sprintf("JSON_BUILD_OBJECT(%s)", joinStrings(keyValuePairs, ", "))
}

// BuildJSONArrayAgg creates a JSON array aggregation using PostgreSQL's JSON_AGG
func (p *PostgresAdapter) BuildJSONArrayAgg(jsonObjectSQL, table, alias, whereClause, resultAlias string) string {
	return fmt.Sprintf(
		"(SELECT COALESCE(JSON_AGG(%s), '[]') FROM %s AS %s WHERE %s) AS %s",
		jsonObjectSQL, table, alias, whereClause, resultAlias,
	)
}

// GetInformationTablesQuery returns PostgreSQL query to fetch table information
func (p *PostgresAdapter) GetInformationTablesQuery() string {
	return `
		SELECT table_schema, table_name
		FROM information_schema.tables
		WHERE table_schema NOT IN ('pg_catalog', 'information_schema', 'pg_toast', 'pg_temp_1', 'pg_toast_temp_1')
		  AND table_type = 'BASE TABLE'
		ORDER BY table_schema, table_name
	`
}

// GetInformationColumnsQuery returns PostgreSQL query to fetch column information
func (p *PostgresAdapter) GetInformationColumnsQuery() string {
	return `
		SELECT 
			table_schema,
			table_name,
			column_name,
			data_type,
			is_nullable,
			column_default,
			ordinal_position
		FROM information_schema.columns
		WHERE table_schema NOT IN ('pg_catalog', 'information_schema', 'pg_toast', 'pg_temp_1', 'pg_toast_temp_1')
		ORDER BY table_schema, table_name, ordinal_position
	`
}

// GetInformationSchemasQuery returns PostgreSQL query to fetch schema information
func (p *PostgresAdapter) GetInformationSchemasQuery() string {
	return `
		SELECT schema_name
		FROM information_schema.schemata
		WHERE schema_name NOT IN ('pg_catalog', 'information_schema', 'pg_toast', 'pg_temp_1', 'pg_toast_temp_1')
		ORDER BY schema_name
	`
}

// GetInformationTableRelationsQuery returns PostgreSQL query to fetch table relationships
func (p *PostgresAdapter) GetInformationTableRelationsQuery() string {
	return `
		SELECT 
			tc.table_name as table_from,
			kcu.column_name as column_from,
			ccu.table_name as table_to,
			ccu.column_name as column_to
		FROM information_schema.table_constraints tc
		JOIN information_schema.key_column_usage kcu 
			ON tc.constraint_name = kcu.constraint_name 
			AND tc.table_schema = kcu.table_schema
		JOIN information_schema.constraint_column_usage ccu 
			ON tc.constraint_name = ccu.constraint_name 
			AND tc.table_schema = ccu.table_schema
		WHERE tc.constraint_type = 'FOREIGN KEY'
		  AND tc.table_schema NOT IN ('pg_catalog', 'information_schema', 'pg_toast', 'pg_temp_1', 'pg_toast_temp_1')
		ORDER BY tc.table_name, kcu.column_name
	`
}

// NewPostgresAdapter creates a new PostgreSQL adapter
func NewPostgresAdapter() *PostgresAdapter {
	return &PostgresAdapter{}
}
