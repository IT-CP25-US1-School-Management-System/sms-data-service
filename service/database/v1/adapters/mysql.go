package adapters

import (
	"context"
	"fmt"

	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/entity"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/database/v1/client"
	"github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// MySQLAdapter provides complete MySQL database support
type MySQLAdapter struct{}

// GetName returns the database type identifier
func (m *MySQLAdapter) GetName() string {
	return "mysql"
}

// GetDriverName returns the SQL driver name
func (m *MySQLAdapter) GetDriverName() string {
	return "mysql"
}

// GetDriverType returns the driver type
func (m *MySQLAdapter) GetDriverType() string {
	return "mysql"
}

// Connect creates a new MySQL connection
func (m *MySQLAdapter) Connect(ctx context.Context, config client.ClientConfig) (*client.Client, error) {
	db, err := sqlx.ConnectContext(ctx, m.GetDriverName(), config.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	db.SetMaxOpenConns(4)

	return client.NewClientFromFields(
		db,
		config.ConnectionString,
		m.GetDriverName(),
		m.GetName(),
		config.Tracer,
	), nil
}

// BuildConnectionString creates a MySQL connection string
func (m *MySQLAdapter) BuildConnectionString(source *entity.Sources, decryptedPassword string) (string, error) {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		source.Username, decryptedPassword, source.Host, source.Port, source.DatabaseName)
	return connStr, nil
}

// GetPlaceholderFormat returns MySQL's placeholder format (?, ?, ?, ...)
func (m *MySQLAdapter) GetPlaceholderFormat() squirrel.PlaceholderFormat {
	return squirrel.Question
}

// BuildJSONObject creates a JSON object using MySQL's JSON_OBJECT
func (m *MySQLAdapter) BuildJSONObject(keyValuePairs []string) string {
	return fmt.Sprintf("JSON_OBJECT(%s)", joinStrings(keyValuePairs, ", "))
}

// BuildJSONArrayAgg creates a JSON array aggregation using MySQL's JSON_ARRAYAGG
func (m *MySQLAdapter) BuildJSONArrayAgg(jsonObjectSQL, table, alias, whereClause, resultAlias string) string {
	return fmt.Sprintf(
		"(SELECT COALESCE(JSON_ARRAYAGG(%s), JSON_ARRAY()) FROM %s AS %s WHERE %s) AS %s",
		jsonObjectSQL, table, alias, whereClause, resultAlias,
	)
}

// GetInformationTablesQuery returns MySQL query to fetch table information
func (m *MySQLAdapter) GetInformationTablesQuery() string {
	return `
		SELECT table_schema, table_name
		FROM information_schema.tables
		WHERE table_schema NOT IN ('information_schema', 'performance_schema', 'mysql', 'sys')
		  AND table_type = 'BASE TABLE'
		ORDER BY table_schema, table_name
	`
}

// GetInformationColumnsQuery returns MySQL query to fetch column information
func (m *MySQLAdapter) GetInformationColumnsQuery() string {
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
		WHERE table_schema NOT IN ('information_schema', 'performance_schema', 'mysql', 'sys')
		ORDER BY table_schema, table_name, ordinal_position
	`
}

// GetInformationSchemasQuery returns MySQL query to fetch schema information
func (m *MySQLAdapter) GetInformationSchemasQuery() string {
	return `
		SELECT schema_name
		FROM information_schema.schemata
		WHERE schema_name NOT IN ('information_schema', 'performance_schema', 'mysql', 'sys')
		ORDER BY schema_name
	`
}

// GetInformationTableRelationsQuery returns MySQL query to fetch table relationships
func (m *MySQLAdapter) GetInformationTableRelationsQuery() string {
	return `
		SELECT 
			TABLE_NAME as table_from,
			COLUMN_NAME as column_from,
			REFERENCED_TABLE_NAME as table_to,
			REFERENCED_COLUMN_NAME as column_to
		FROM information_schema.KEY_COLUMN_USAGE
		WHERE REFERENCED_TABLE_NAME IS NOT NULL
		  AND TABLE_SCHEMA NOT IN ('information_schema', 'performance_schema', 'mysql', 'sys')
		ORDER BY TABLE_NAME, COLUMN_NAME
	`
}

// NewMySQLAdapter creates a new MySQL adapter
func NewMySQLAdapter() *MySQLAdapter {
	return &MySQLAdapter{}
}
