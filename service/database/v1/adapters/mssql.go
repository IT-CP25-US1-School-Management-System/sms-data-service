package adapters

import (
	"context"
	"fmt"

	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/entity"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/database/v1/client"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/microsoft/go-mssqldb"
)

// MSSQLAdapter provides complete SQL Server database support
type MSSQLAdapter struct{}

// GetName returns the database type identifier
func (ms *MSSQLAdapter) GetName() string {
	return "mssql"
}

// GetDriverName returns the SQL driver name
func (ms *MSSQLAdapter) GetDriverName() string {
	return "sqlserver"
}

// GetDriverType returns the driver type
func (ms *MSSQLAdapter) GetDriverType() string {
	return "mssql"
}

// Connect creates a new SQL Server connection
func (ms *MSSQLAdapter) Connect(ctx context.Context, config client.ClientConfig) (*client.Client, error) {
	db, err := sqlx.ConnectContext(ctx, ms.GetDriverName(), config.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQL Server: %w", err)
	}

	db.SetMaxOpenConns(4)

	return client.NewClientFromFields(
		db,
		config.ConnectionString,
		ms.GetDriverName(),
		ms.GetName(),
		config.Tracer,
	), nil
}

// BuildConnectionString creates a SQL Server connection string
func (ms *MSSQLAdapter) BuildConnectionString(source *entity.Sources, decryptedPassword string) (string, error) {
	connStr := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
		source.Username, decryptedPassword, source.Host, source.Port, source.DatabaseName)
	return connStr, nil
}

// GetPlaceholderFormat returns SQL Server's placeholder format (@p1, @p2, @p3, ...)
func (ms *MSSQLAdapter) GetPlaceholderFormat() squirrel.PlaceholderFormat {
	return squirrel.AtP
}

// BuildJSONObject creates a JSON object using SQL Server's FOR JSON PATH
func (ms *MSSQLAdapter) BuildJSONObject(keyValuePairs []string) string {
	var pairs []string
	for i := 0; i < len(keyValuePairs); i += 2 {
		if i+1 < len(keyValuePairs) {
			key := trimQuotes(keyValuePairs[i])
			value := keyValuePairs[i+1]
			pairs = append(pairs, fmt.Sprintf("%s AS %s", value, key))
		}
	}
	return fmt.Sprintf("(SELECT %s FOR JSON PATH, WITHOUT_ARRAY_WRAPPER)", joinStrings(pairs, ", "))
}

// BuildJSONArrayAgg creates a JSON array aggregation using SQL Server's FOR JSON PATH
func (ms *MSSQLAdapter) BuildJSONArrayAgg(jsonObjectSQL, table, alias, whereClause, resultAlias string) string {
	return fmt.Sprintf(
		"(SELECT COALESCE((SELECT * FROM %s AS %s WHERE %s FOR JSON PATH), '[]')) AS %s",
		table, alias, whereClause, resultAlias,
	)
}

// GetInformationTablesQuery returns SQL Server query to fetch table information
func (ms *MSSQLAdapter) GetInformationTablesQuery() string {
	return `
		SELECT table_schema, table_name
		FROM information_schema.tables
		WHERE table_schema NOT IN ('sys', 'INFORMATION_SCHEMA')
		  AND table_type = 'BASE TABLE'
		ORDER BY table_schema, table_name
	`
}

// GetInformationColumnsQuery returns SQL Server query to fetch column information
func (ms *MSSQLAdapter) GetInformationColumnsQuery() string {
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
		WHERE table_schema NOT IN ('sys', 'INFORMATION_SCHEMA')
		ORDER BY table_schema, table_name, ordinal_position
	`
}

// GetInformationSchemasQuery returns SQL Server query to fetch schema information
func (ms *MSSQLAdapter) GetInformationSchemasQuery() string {
	return `
		SELECT schema_name
		FROM information_schema.schemata
		WHERE schema_name NOT IN ('sys', 'INFORMATION_SCHEMA', 'guest', 'db_owner', 'db_accessadmin', 'db_securityadmin', 
			'db_ddladmin', 'db_backupoperator', 'db_datareader', 'db_datawriter', 'db_denydatareader', 'db_denydatawriter')
		ORDER BY schema_name
	`
}

// GetInformationTableRelationsQuery returns SQL Server query to fetch table relationships
func (ms *MSSQLAdapter) GetInformationTableRelationsQuery() string {
	return `
		SELECT 
			fk.name AS constraint_name,
			tp.name AS table_from,
			cp.name AS column_from,
			tr.name AS table_to,
			cr.name AS column_to
		FROM sys.foreign_keys AS fk
		INNER JOIN sys.tables AS tp ON fk.parent_object_id = tp.object_id
		INNER JOIN sys.tables AS tr ON fk.referenced_object_id = tr.object_id
		INNER JOIN sys.foreign_key_columns AS fkc ON fk.object_id = fkc.constraint_object_id
		INNER JOIN sys.columns AS cp ON fkc.parent_column_id = cp.column_id AND fkc.parent_object_id = cp.object_id
		INNER JOIN sys.columns AS cr ON fkc.referenced_column_id = cr.column_id AND fkc.referenced_object_id = cr.object_id
		WHERE SCHEMA_NAME(tp.schema_id) NOT IN ('sys', 'INFORMATION_SCHEMA')
		ORDER BY tp.name, cp.name
	`
}

// NewMSSQLAdapter creates a new SQL Server adapter
func NewMSSQLAdapter() *MSSQLAdapter {
	return &MSSQLAdapter{}
}
