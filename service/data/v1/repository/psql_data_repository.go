package repository

import (
	"context"
	"fmt"
	"strings"

	helperModel "github.com/GodeFvt/go-backend/helper/models"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/entity"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/data/v1"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/database/v1"
	"github.com/gofrs/uuid"
)

type psqlDataRepository struct {
	dbConnectionManager database.DBConnectionManagerUsecase
}

// BuildRuntimeSQL builds SQL query based on runtime policy configuration
func (p *psqlDataRepository) BuildRuntimeSQL(ctx context.Context, sourceID *uuid.UUID, runtime *entity.RuntimePolicy) (string, []interface{}, error) {
	_, err := p.dbConnectionManager.GetConnection(ctx, *sourceID)
	if err != nil {
		return "", nil, err
	}

	var query strings.Builder
	var args []interface{}

	// Analyze projections to detect aggregations and collect non-aggregate columns
	var hasAggregation bool
	var nonAggregateColumns []string
	var selectClauses []string

	if len(runtime.Query.Projections) > 0 {
		for _, proj := range runtime.Query.Projections {
			var selectClause string

			if proj.Expr != nil && proj.Expr.Field != "" {
				// Handle expression projections (including aggregations)
				upperExpr := strings.ToUpper(proj.Expr.Field)

				// Check if this is an aggregate function (including complex ones)
				if strings.Contains(upperExpr, "COUNT(") ||
					strings.Contains(upperExpr, "SUM(") ||
					strings.Contains(upperExpr, "AVG(") ||
					strings.Contains(upperExpr, "MAX(") ||
					strings.Contains(upperExpr, "MIN(") ||
					strings.Contains(upperExpr, "COUNT(DISTINCT") ||
					strings.Contains(upperExpr, "SUM(DISTINCT") ||
					strings.Contains(upperExpr, "EXTRACT(") && strings.Contains(upperExpr, "MAX(") {
					hasAggregation = true
					selectClause = proj.Expr.Field
				} else {
					// Non-aggregate expression - should be included in GROUP BY
					selectClause = proj.Expr.Field
					nonAggregateColumns = append(nonAggregateColumns, proj.Expr.Field)
				}

				// Add alias if specified
				if proj.Alias != "" {
					selectClause += " AS " + proj.Alias
				}

			} else if proj.Column != "" {
				// Handle regular column projections - these are NOT aggregates
				nonAggregateColumns = append(nonAggregateColumns, proj.Column)
				selectClause = proj.Column

				// Add alias if specified
				if proj.Alias != "" {
					selectClause += " AS " + proj.Alias
				}
			}

			if selectClause != "" {
				selectClauses = append(selectClauses, selectClause)
			}
		}
	}

	// Build SELECT clause
	if len(selectClauses) > 0 {
		query.WriteString("SELECT " + strings.Join(selectClauses, ", "))
	} else {
		query.WriteString("SELECT *")
	}

	// Build FROM clause
	if runtime.Query.From != nil {
		if runtime.Query.From.Table != "" {
			query.WriteString(" FROM " + runtime.Query.From.Table)
		} else if runtime.Query.From.View != "" {
			query.WriteString(" FROM " + runtime.Query.From.View)
		}
	} else {
		return "", nil, fmt.Errorf("FROM clause is required")
	}

	// Build JOIN clauses
	if len(runtime.Query.Joins) > 0 {
		for _, join := range runtime.Query.Joins {
			if join.Table != "" && join.Condition.Field != "" {
				joinType := strings.ToUpper(join.Type)
				if joinType == "" {
					joinType = "LEFT"
				}
				query.WriteString(fmt.Sprintf(" %s JOIN %s ON %s", joinType, join.Table, join.Condition.Field))
			}
		}
	}

	// Build WHERE clause foundation (conditions will be added dynamically by ExecuteQuery)
	if len(runtime.Query.WhereAllow) > 0 {
		query.WriteString(" WHERE 1=1")
	}

	// Build GROUP BY clause when we have aggregations
	if hasAggregation {
		var groupByFields []string

		if len(runtime.Query.GroupBy) > 0 {
			// Use explicitly specified GROUP BY fields
			for _, group := range runtime.Query.GroupBy {
				if group.Field != "" {
					groupByFields = append(groupByFields, group.Field)
				}
			}
		} else if len(nonAggregateColumns) > 0 {
			// Auto-generate GROUP BY from ALL non-aggregate columns in SELECT
			// This ensures PostgreSQL compliance
			seen := make(map[string]bool)
			for _, col := range nonAggregateColumns {
				if !seen[col] {
					groupByFields = append(groupByFields, col)
					seen[col] = true
				}
			}
		}

		if len(groupByFields) > 0 {
			query.WriteString(" GROUP BY " + strings.Join(groupByFields, ", "))
		}
	}

	// Build ORDER BY clause
	if len(runtime.Query.OrderAllow) > 0 {
		orderClauses := make([]string, 0, len(runtime.Query.OrderAllow))
		for _, order := range runtime.Query.OrderAllow {
			if order.Field != "" {
				direction := "ASC"
				if len(order.Directions) > 0 && strings.ToUpper(order.Directions[0]) == "DESC" {
					direction = "DESC"
				}
				orderClauses = append(orderClauses, order.Field+" "+direction)
			}
		}
		if len(orderClauses) > 0 {
			query.WriteString(" ORDER BY " + strings.Join(orderClauses, ", "))
		}
	}

	// Note: LIMIT and OFFSET will be added by ExecuteQuery function for pagination

	return query.String(), args, nil
}

// BuildCreateSQL builds INSERT SQL based on write policy configuration
func (p *psqlDataRepository) BuildCreateSQL(ctx context.Context, sourceID *uuid.UUID, writePolicy *entity.WritePolicy, data map[string]interface{}) (string, []interface{}, error) {
	_, err := p.dbConnectionManager.GetConnection(ctx, *sourceID)
	if err != nil {
		return "", nil, err
	}

	var query strings.Builder
	var args []interface{}
	argIndex := 1

	// Get target table from write policy
	tableName := ""
	if writePolicy.Query.From != nil {
		if writePolicy.Query.From.Table != "" {
			tableName = writePolicy.Query.From.Table
		} else if writePolicy.Query.From.View != "" {
			tableName = writePolicy.Query.From.View
		}
	}

	if tableName == "" {
		return "", nil, fmt.Errorf("table name not specified in write policy")
	}

	// Build INSERT statement
	var columns []string
	var placeholders []string

	for _, allowedCol := range writePolicy.AllowEdit {
		if value, exists := data[allowedCol]; exists {
			columns = append(columns, allowedCol)
			placeholders = append(placeholders, fmt.Sprintf("$%d", argIndex))
			args = append(args, value)
			argIndex++
		}
	}

	if len(columns) == 0 {
		return "", nil, fmt.Errorf("no valid columns to insert")
	}

	query.WriteString(fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", ")))

	return query.String(), args, nil
}

// BuildUpdateSQL builds UPDATE SQL based on write policy configuration
func (p *psqlDataRepository) BuildUpdateSQL(ctx context.Context, sourceID *uuid.UUID, writePolicy *entity.WritePolicy, data map[string]interface{}, whereConditions map[string]interface{}) (string, []interface{}, error) {
	_, err := p.dbConnectionManager.GetConnection(ctx, *sourceID)
	if err != nil {
		return "", nil, err
	}

	var query strings.Builder
	var args []interface{}
	argIndex := 1

	// Get target table from write policy
	tableName := ""
	if writePolicy.Query.From != nil {
		if writePolicy.Query.From.Table != "" {
			tableName = writePolicy.Query.From.Table
		} else if writePolicy.Query.From.View != "" {
			tableName = writePolicy.Query.From.View
		}
	}

	if tableName == "" {
		return "", nil, fmt.Errorf("table name not specified in write policy")
	}

	// Build UPDATE statement
	query.WriteString(fmt.Sprintf("UPDATE %s SET ", tableName))

	var setClauses []string
	for _, allowedCol := range writePolicy.AllowEdit {
		if value, exists := data[allowedCol]; exists {
			setClauses = append(setClauses, fmt.Sprintf("%s = $%d", allowedCol, argIndex))
			args = append(args, value)
			argIndex++
		}
	}

	if len(setClauses) == 0 {
		return "", nil, fmt.Errorf("no valid columns to update")
	}

	query.WriteString(strings.Join(setClauses, ", "))

	// WHERE clause
	if len(whereConditions) > 0 {
		query.WriteString(" WHERE ")
		var whereClauses []string
		for field, value := range whereConditions {
			whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", field, argIndex))
			args = append(args, value)
			argIndex++
		}
		query.WriteString(strings.Join(whereClauses, " AND "))
	}

	return query.String(), args, nil
}

// BuildDeleteSQL builds DELETE SQL based on delete policy configuration
func (p *psqlDataRepository) BuildDeleteSQL(ctx context.Context, sourceID *uuid.UUID, deletePolicy *entity.DeletePolicy, whereConditions map[string]interface{}) (string, []interface{}, error) {
	_, err := p.dbConnectionManager.GetConnection(ctx, *sourceID)
	if err != nil {
		return "", nil, err
	}

	var query strings.Builder
	var args []interface{}
	argIndex := 1

	// Get target table from delete policy
	tableName := ""
	if deletePolicy.Query.From != nil {
		if deletePolicy.Query.From.Table != "" {
			tableName = deletePolicy.Query.From.Table
		} else if deletePolicy.Query.From.View != "" {
			tableName = deletePolicy.Query.From.View
		}
	}

	if tableName == "" {
		return "", nil, fmt.Errorf("table name not specified in delete policy")
	}

	// Build DELETE statement
	query.WriteString(fmt.Sprintf("DELETE FROM %s", tableName))

	// WHERE clause
	if len(whereConditions) > 0 {
		query.WriteString(" WHERE ")
		var whereClauses []string
		for field, value := range whereConditions {
			whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", field, argIndex))
			args = append(args, value)
			argIndex++
		}
		query.WriteString(strings.Join(whereClauses, " AND "))
	} else {
		return "", nil, fmt.Errorf("delete operation requires WHERE conditions for safety")
	}

	return query.String(), args, nil
}

// ExecuteQuery executes a query and returns the results as map[string]interface{}
func (p *psqlDataRepository) ExecuteQuery(ctx context.Context, sourceID *uuid.UUID, query string, args []interface{}, paginator *helperModel.Paginator) ([]map[string]interface{}, error) {
	client, err := p.dbConnectionManager.GetConnection(ctx, *sourceID)
	if err != nil {
		return nil, err
	}

	// Prepare the final query with pagination
	finalQuery := query
	finalArgs := args

	if paginator != nil {
		finalQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
		finalArgs = append(finalArgs, paginator.GetLimit(), paginator.GetOffset())
	}

	rows, err := client.GetClient().QueryxContext(ctx, finalQuery, finalArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		row := make(map[string]interface{})
		err = rows.MapScan(row)
		if err != nil {
			return nil, err
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Set pagination info if paginator is provided
	if paginator != nil {
		// Create count query by wrapping the original query
		countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s) as count_query", query)

		var totalRows int
		if err := client.GetClient().GetContext(ctx, &totalRows, countQuery, args...); err != nil {
			return nil, err
		}
		paginator.SetPaginatorByAllRows(totalRows)
	}

	return results, nil
}

// ExecuteQueryByKey executes a query with dynamic WHERE conditions based on key-value pairs
func (p *psqlDataRepository) ExecuteQueryByKey(ctx context.Context, sourceID *uuid.UUID, baseQuery string, baseArgs []interface{}, whereConditions map[string]interface{}) ([]map[string]interface{}, error) {
	client, err := p.dbConnectionManager.GetConnection(ctx, *sourceID)
	if err != nil {
		return nil, err
	}

	// Build WHERE clause from conditions
	query := baseQuery
	args := baseArgs
	argIndex := len(baseArgs) + 1

	// Debug logging
	fmt.Printf("DEBUG ExecuteQueryByKey - Base Query: %s\n", baseQuery)
	fmt.Printf("DEBUG ExecuteQueryByKey - Base Args: %v\n", baseArgs)
	fmt.Printf("DEBUG ExecuteQueryByKey - Where Conditions: %v\n", whereConditions)

	if len(whereConditions) > 0 {
		var whereClauses []string
		for field, value := range whereConditions {
			whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", field, argIndex))
			args = append(args, value)
			argIndex++
		}

		// Find the position to insert WHERE clause
		// We need to insert before ORDER BY, GROUP BY, HAVING, LIMIT, etc.
		upperQuery := strings.ToUpper(query)
		var insertPos int = len(query) // default to end of query

		// Find ORDER BY position
		if orderPos := strings.Index(upperQuery, " ORDER BY"); orderPos != -1 {
			insertPos = orderPos
		}
		// Find GROUP BY position (if it comes before ORDER BY)
		if groupPos := strings.Index(upperQuery, " GROUP BY"); groupPos != -1 && groupPos < insertPos {
			insertPos = groupPos
		}
		// Find HAVING position (if it comes before ORDER BY/GROUP BY)
		if havingPos := strings.Index(upperQuery, " HAVING"); havingPos != -1 && havingPos < insertPos {
			insertPos = havingPos
		}
		// Find LIMIT position (if it comes before others)
		if limitPos := strings.Index(upperQuery, " LIMIT"); limitPos != -1 && limitPos < insertPos {
			insertPos = limitPos
		}

		// Split query at insert position
		beforeClause := query[:insertPos]
		afterClause := query[insertPos:]

		// Check if WHERE clause already exists in beforeClause
		if strings.Contains(strings.ToUpper(beforeClause), "WHERE") {
			query = beforeClause + " AND " + strings.Join(whereClauses, " AND ") + afterClause
		} else {
			query = beforeClause + " WHERE " + strings.Join(whereClauses, " AND ") + afterClause
		}
	}

	// Debug final query
	fmt.Printf("DEBUG ExecuteQueryByKey - Final Query: %s\n", query)
	fmt.Printf("DEBUG ExecuteQueryByKey - Final Args: %v\n", args)

	rows, err := client.GetClient().QueryxContext(ctx, query, args...)
	if err != nil {
		fmt.Printf("DEBUG ExecuteQueryByKey - SQL Error: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		row := make(map[string]interface{})
		err = rows.MapScan(row)
		if err != nil {
			return nil, err
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	fmt.Printf("DEBUG ExecuteQueryByKey - Results count: %d\n", len(results))
	return results, nil
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

// ExecuteInsert implements data.PsqlDataRepository.
func (p *psqlDataRepository) ExecuteInsert(ctx context.Context, sourceID *uuid.UUID, query string, args []interface{}) (map[string]interface{}, error) {
	client, err := p.dbConnectionManager.GetConnection(ctx, *sourceID)
	if err != nil {
		return nil, err
	}

	fmt.Printf("DEBUG ExecuteInsert - Query: %s\n", query)
	fmt.Printf("DEBUG ExecuteInsert - Args: %v\n", args)

	var result map[string]interface{}
	rows, err := client.GetClient().QueryxContext(ctx, query, args...)
	if err != nil {
		fmt.Printf("DEBUG ExecuteInsert - SQL Error: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		result = make(map[string]interface{})
		err = rows.MapScan(result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// ExecuteUpdate implements data.PsqlDataRepository.
func (p *psqlDataRepository) ExecuteUpdate(ctx context.Context, sourceID *uuid.UUID, query string, args []interface{}) (int64, error) {
	client, err := p.dbConnectionManager.GetConnection(ctx, *sourceID)
	if err != nil {
		return 0, err
	}

	fmt.Printf("DEBUG ExecuteUpdate - Query: %s\n", query)
	fmt.Printf("DEBUG ExecuteUpdate - Args: %v\n", args)

	result, err := client.GetClient().ExecContext(ctx, query, args...)
	if err != nil {
		fmt.Printf("DEBUG ExecuteUpdate - SQL Error: %v\n", err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

// ExecuteDelete implements data.PsqlDataRepository.
func (p *psqlDataRepository) ExecuteDelete(ctx context.Context, sourceID *uuid.UUID, query string, args []interface{}) (int64, error) {
	client, err := p.dbConnectionManager.GetConnection(ctx, *sourceID)
	if err != nil {
		return 0, err
	}

	fmt.Printf("DEBUG ExecuteDelete - Query: %s\n", query)
	fmt.Printf("DEBUG ExecuteDelete - Args: %v\n", args)

	result, err := client.GetClient().ExecContext(ctx, query, args...)
	if err != nil {
		fmt.Printf("DEBUG ExecuteDelete - SQL Error: %v\n", err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
