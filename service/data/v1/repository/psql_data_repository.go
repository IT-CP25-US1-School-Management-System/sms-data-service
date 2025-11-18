package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	helperModel "github.com/GodeFvt/go-backend/helper/models"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/constants"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/errs"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/entity"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/data/v1"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/database/v1"
	"github.com/Masterminds/squirrel"
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

func (p *psqlDataRepository) FetchInformationTableRelationsBySourceID(ctx context.Context, dbType string, sourceID *uuid.UUID) ([]*entity.TableRelations, error) {
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
				tc.table_name as table_from,
				kcu.column_name as column_from,
				ccu.table_name as table_to,
				ccu.column_name as column_to
			FROM information_schema.table_constraints tc
			JOIN information_schema.key_column_usage kcu 
				ON tc.constraint_name = kcu.constraint_name 
				AND tc.table_schema = kcu.table_schema
			JOIN information_schema.referential_constraints rc 
				ON tc.constraint_name = rc.constraint_name 
				AND tc.table_schema = rc.constraint_schema
			JOIN information_schema.constraint_column_usage ccu 
				ON rc.unique_constraint_name = ccu.constraint_name 
				AND rc.unique_constraint_schema = ccu.constraint_schema
			WHERE tc.constraint_type = 'FOREIGN KEY'
			  AND tc.table_schema NOT IN ('information_schema', 'performance_schema', 'mysql', 'sys')
			ORDER BY tc.table_name, kcu.column_name
		`
	default:
		query = `
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

	rows, err := client.GetClient().QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tableRelations []*entity.TableRelations
	for rows.Next() {
		var tableFrom, columnFrom, tableTo, columnTo string
		if err := rows.Scan(&tableFrom, &columnFrom, &tableTo, &columnTo); err != nil {
			return nil, err
		}

		tr := &entity.TableRelations{
			SourceID:   sourceID,
			TableFrom:  tableFrom,
			ColumnFrom: columnFrom,
			TableTo:    tableTo,
			ColumnTo:   columnTo,
			CreatedAt:  nil,
		}
		tr.GenUUID()
		now := helperModel.NewTimestampFromNow()
		tr.CreatedAt = &now
		tableRelations = append(tableRelations, tr)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tableRelations, nil
}

// psql รันบน Postgres, squirrel.Dollar
var psqlBuilder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

// aliasManager สร้าง alias (t0, t1, t2...)
// map[real_table_name]alias
type aliasManager struct {
	aliasMap  map[string]string
	counter   int
	fromTable string
	fromAlias string
}

func newAliasManager(fromTable string) *aliasManager {
	am := &aliasManager{
		aliasMap:  make(map[string]string),
		counter:   0,
		fromTable: fromTable,
	}
	// สร้าง alias แรก (t0) สำหรับตาราง FROM
	am.fromAlias = am.generate(fromTable)
	return am
}

// generate สร้าง alias (t0, t1, t2...) และเก็บไว้
func (am *aliasManager) generate(tableName string) string {
	alias := fmt.Sprintf("t%d", am.counter)
	am.counter++
	if _, ok := am.aliasMap[tableName]; !ok {
		am.aliasMap[tableName] = alias
	}
	return alias
}

// Get ดึง alias ของตาราง
func (am *aliasManager) Get(tableName string) (string, bool) {
	alias, ok := am.aliasMap[tableName]
	return alias, ok
}

// buildOnClause สร้าง ON clause
//
//	t1.column_to = t0.column_from ต้องมีการทำ alias เเล้ว
func buildOnClause(cond entity.Condition, fromAlias string, toAlias string) (string, error) {
	if cond.ColumnFrom == "" || cond.ColumnTo == "" || cond.Operator == "" {
		return "", fmt.Errorf("condition object must have 'column_from', 'op', and 'column_to'")
	}

	// สร้าง SQL: toAlias.ColumnTo OP fromAlias.ColumnFrom
	// e.g., "t1.id = t0.role_id"
	// e.g., "t1.person_id = t0.id" (สำหรับ subquery)
	return fmt.Sprintf("%s.%s %s %s.%s",
		toAlias,         // e.g., "t1"
		cond.ColumnTo,   // e.g., "id"
		cond.Operator,   // e.g., "="
		fromAlias,       // e.g., "t0"
		cond.ColumnFrom, // e.g., "role_id"
	), nil
}

// buildAllowedWhereMap สร้าง map["t0.id"]["="] = true
func buildAllowedWhereMap(whereAllows []entity.WhereAllow, am *aliasManager) (map[string]map[string]bool, error) {
	allowedMap := make(map[string]map[string]bool)
	for _, w := range whereAllows {
		tableAlias, ok := am.Get(w.TableName)
		if !ok {
			// ข้ามไป ถ้าตารางนี้ไม่ได้ Join
			continue
		}
		fieldWithAlias := fmt.Sprintf("%s.%s", tableAlias, w.Field) // "t0.id"

		allowedMap[fieldWithAlias] = make(map[string]bool)
		for _, op := range w.Operators {
			allowedMap[fieldWithAlias][op] = true
		}
	}
	return allowedMap, nil
}

// buildSquirrelExpr แปลง FilterInput เป็น Squirrel expression
func buildSquirrelExpr(field string, op string, val interface{}) (squirrel.Sqlizer, error) {
	switch strings.ToUpper(op) {
	case "=":
		return squirrel.Eq{field: val}, nil
	case "!=":
		return squirrel.NotEq{field: val}, nil
	case ">":
		return squirrel.Gt{field: val}, nil
	case ">=":
		return squirrel.GtOrEq{field: val}, nil
	case "<":
		return squirrel.Lt{field: val}, nil
	case "<=":
		return squirrel.LtOrEq{field: val}, nil
	case "IN":
		return squirrel.Eq{field: val}, nil
	case "NOT IN":
		return squirrel.NotEq{field: val}, nil
	case "LIKE":
		return squirrel.Like{field: val}, nil
	case "NOT LIKE":
		return squirrel.NotLike{field: val}, nil
	case "IS NULL":
		return squirrel.Eq{field: nil}, nil
	case "IS NOT NULL":
		return squirrel.NotEq{field: nil}, nil
	default:
		return nil, fmt.Errorf("unsupported operator: %s", op)
	}
}

// validateAndPrepareData ทำหน้าที่:
// 1. กรอง 'data' ให้เหลือเฉพาะ field ที่อยู่ใน 'allowEdit'
// 2. ตรวจสอบ 'data' ที่เหลือ เทียบกับ 'schema' (IsNullable, DataType, Enum)
// 3. ใส่ค่า 'default' ถ้า field นั้นไม่ได้ถูกส่งมา
// 4. คืนค่า map[string]interface{} ที่ "สะอาดแล้ว" (Cleaned)
func (p *psqlDataRepository) validateAndPrepareData(
	schema entity.Schema,
	writePolicy *entity.WritePolicy,
	data map[string]interface{},
) (map[string]interface{}, error) {

	// 1. สร้าง Map ของ Schema Columns เพื่อให้ค้นหาได้เร็ว
	// (กรองเฉพาะตารางที่จะ INSERT/UPDATE)
	targetTable := writePolicy.Query.From.Table
	schemaMap := make(map[string]entity.Column)
	for _, col := range schema.Columns {
		if col.TableName == targetTable {
			schemaMap[col.Name] = col
		}
	}

	// 2. สร้าง Map ของ field ที่อนุญาตให้แก้ไข AllowEdit
	allowedEditMap := make(map[string]bool)
	for _, fieldName := range writePolicy.AllowEdit {
		allowedEditMap[fieldName] = true
	}

	// 3. วน Loop ตรวจสอบทุก field ที่อนุญาต AllowEdit
	validatedData := make(map[string]interface{})

	for fieldName := range allowedEditMap {
		schemaCol, ok := schemaMap[fieldName]
		if !ok {
			return nil, fmt.Errorf("config error: field '%s' in AllowEdit not found in schema for table '%s'", fieldName, targetTable)
		}

		val, dataExists := data[fieldName]

		// --- 3a. จัดการกรณีข้อมูลเป็น nil หรือไม่ได้ส่งมา ---
		if !dataExists || val == nil {
			if schemaCol.Default != nil {
				// มีค่า Default: ใส่ค่า Default ให้
				//TODO (ต้องแปลง Default ที่เป็น string กลับไปเป็น type ที่ถูกต้อง)
				validatedData[fieldName] = *schemaCol.Default // (ตอนนี้ยังเป็น string)
			} else if !schemaCol.IsNullable {
				// ไม่มีค่า Default และ Not Null: Error
				return nil, fmt.Errorf("validation failed: field '%s' is required (not nullable) but was not provided", fieldName)
			} else {
				// ไม่มีค่า Default แต่ Nullable: nil คือค่าที่ถูกต้อง
				validatedData[fieldName] = nil
			}
			continue // ไปตรวจสอบ field ถัดไป
		}

		// --- 3b. ข้อมูลมีค่า (ไม่ nil), ตรวจสอบความถูกต้อง ---

		// Rule 1: ตรวจสอบ Enum
		if len(schemaCol.Enum) > 0 {
			if err := validateEnum(schemaCol.Enum, val); err != nil {
				return nil, fmt.Errorf("validation failed for field '%s': %w", fieldName, err)
			}
		}

		// Rule 2: ตรวจสอบ DataType
		if err := validateDataType(schemaCol.DataType, val); err != nil {
			return nil, fmt.Errorf("validation failed for field '%s': %w", fieldName, err)
		}

		// ถ้าผ่านทุกอย่าง
		validatedData[fieldName] = val
	}

	// จงใจ "ทิ้ง" field ที่อยูใน 'data' แต่ไม่อยู่ใน 'allowEditMap'
	return validatedData, nil
}

// validateEnum (Helper)
func validateEnum(enum []string, value interface{}) error {
	valStr, ok := value.(string)
	if !ok {
		return fmt.Errorf("expects a string for enum check, but got %T", value)
	}

	for _, enumVal := range enum {
		if valStr == enumVal {
			return nil
		}
	}

	return fmt.Errorf("value '%v' is not in the allowed enum list %v", value, enum)
}

// validateDataType (Helper) - ตรวจสอบ Go Type เทียบกับ SQL Type (แบบพื้นฐาน)
func validateDataType(dataType string, value interface{}) error {
	kind := reflect.TypeOf(value).Kind()
	// (อนุญาต float64 สำหรับ number ทุกชนิด เพราะ JSON Unmarshal จะได้ float64 เสมอ)
	switch dataType {
	case "int", "serial", "int4", "int8", "integer", "bigint", "smallint":
		if kind != reflect.Float64 && kind != reflect.Int {
			return fmt.Errorf("expects a number (int/float64), but got %T", value)
		}
	case "varchar", "text", "genders", "blood_types", "honor_types", "character varying", "USER-DEFINED", "longtext":
		if kind != reflect.String {
			return fmt.Errorf("expects a string, but got %T", value)
		}
	case "uuid":
		if kind != reflect.String {
			return fmt.Errorf("expects a string (for uuid), but got %T", value)
		}
		if _, err := uuid.FromString(value.(string)); err != nil {
			return fmt.Errorf("invalid uuid format: %w", err)
		}
	case "date", "timestamp", "timestamp without time zone", "timestamp with time zone":
		if kind != reflect.String {
			return fmt.Errorf("expects a string (for date/timestamp), but got %T", value)
		}
		if dataType == "date" {
			if _, err := time.Parse("2006-01-02", value.(string)); err != nil {
				return fmt.Errorf("expects a date in 'YYYY-MM-DD' format: %w , but got %T", err, value)
			}
		}
		if dataType == "timestamp" {
			_, err := helperModel.ParseTimestampFromString(value.(string))
			if err != nil {
				return fmt.Errorf("expects a valid timestamp format: %w , but got %T", err, value)
			}
		}
	case "bool":
		if kind != reflect.Bool {
			return fmt.Errorf("expects a boolean, but got %T", value)
		}
	case "decimal(12,2)", "numeric", "float", "float4", "float8", "double precision":
		if kind != reflect.Float64 && kind != reflect.Int {
			return fmt.Errorf("expects a number (decimal/float64), but got %T", value)
		}
	default:
		// ไม่รู้จัก DataType, อนุญาตให้ผ่านไปก่อน
		return nil
	}
	return nil
}

// สร้าง Map["table_name"]["field_name"] -> entity.Column
func createSchemaMap(schema *entity.Schema) map[string]map[string]entity.Column {
	schemaMap := make(map[string]map[string]entity.Column)
	for _, col := range schema.Columns {
		if _, ok := schemaMap[col.TableName]; !ok {
			schemaMap[col.TableName] = make(map[string]entity.Column)
		}
		schemaMap[col.TableName][col.Name] = col
	}
	return schemaMap
}

// แปลง []entity.View (จาก config) ให้เป็น Map
// viewMap["person_data"]["id"] = true
func createViewMap(viewConfigs []entity.View) map[string]map[string]bool {
	viewMap := make(map[string]map[string]bool)
	for _, view := range viewConfigs {
		colMap := make(map[string]bool)
		for _, colName := range view.Columns {
			colMap[colName] = true
		}
		viewMap[view.TableName] = colMap
	}
	return viewMap
}

func buildRuntimeSQLBuilder(
	ctx context.Context,
	schemaMap map[string]map[string]entity.Column,
	queryPlan *entity.QueryPlan,
	filterGroups [][]entity.FilterInput,
	logicalOperator string,
	pagination *helperModel.Paginator,
	sortBy string,
	sortOrder string,
	viewMap map[string]map[string]bool,
) (squirrel.SelectBuilder, error) {

	// --- 0. Alias Management & Mode Detection ---
	if queryPlan.From == nil || queryPlan.From.Table == "" {
		return squirrel.SelectBuilder{}, fmt.Errorf("QueryPlan.From.Table is required")
	}
	am := newAliasManager(queryPlan.From.Table)
	builder := psqlBuilder.Select().From(fmt.Sprintf("%s AS %s", am.fromTable, am.fromAlias))

	// [NEW] ตรวจสอบโหมด
	isAggregateQuery := len(queryPlan.GroupBy) > 0

	// (ดึง View ของตาราง From)
	fromViewCols, ok := viewMap[am.fromTable]
	if !ok || len(fromViewCols) == 0 {
		return builder, fmt.Errorf("view is missing configuration for base table '%s'", am.fromTable)
	}

	// --- 1. Joins (ต้องทำก่อน Projections) ---
	// (ในโหมด Aggregate, เราต้อง Join 1:1 ก่อน เพื่อให้ Projections และ GroupBy หา t1.name เจอ)
	// (ในโหมด Nesting, เราก็ต้อง Join 1:1 ก่อน เพื่อให้ GroupBy (สำหรับ 1:1) ทำงาน)

	// `groupByColumns` จะถูกใช้ในโหมด Nesting เท่านั้น
	var groupByColumns []string
	// `joinedTables` จะเก็บ alias ของตาราง 1:1 ที่ Join แล้ว
	joinedTables := make(map[string]string) // map[real_table_name] -> alias

	for _, j := range queryPlan.Joins {
		// (ข้าม Join ที่ View ไม่ได้เลือก)
		joinViewCols, shouldJoin := viewMap[j.TableTo]
		if !shouldJoin || len(joinViewCols) == 0 {
			continue
		}

		fromAlias, _ := am.Get(j.TableFrom)
		toAlias := am.generate(j.TableTo)
		onClause, err := buildOnClause(j.Condition, fromAlias, toAlias)
		if err != nil {
			return builder, fmt.Errorf("failed to build ON clause for join '%s': %w", j.Alias, err)
		}

		if j.Relation == "one_to_one" || j.Relation == "many_to_one" {
			// (Join 1:1 ทำงานได้ทั้ง 2 โหมด)
			joinTableSQL := fmt.Sprintf("%s AS %s", j.TableTo, toAlias)
			builder = builder.LeftJoin(fmt.Sprintf("%s ON %s", joinTableSQL, onClause))
			joinedTables[j.TableTo] = toAlias

			if !isAggregateQuery {
				// (ถ้าเป็น Nesting Mode, สร้าง JSON_BUILD_OBJECT)
				var joinProjections []string
				for _, p := range j.Projections {
					if !joinViewCols[p.Column] {
						continue
					}
					jsonKey := fmt.Sprintf("'%s'", p.Alias)
					jsonValue := fmt.Sprintf("%s.%s", toAlias, p.Column)
					joinProjections = append(joinProjections, jsonKey, jsonValue)
				}
				if len(joinProjections) == 0 {
					continue
				}

				jsonBuild := fmt.Sprintf("COALESCE(JSON_BUILD_OBJECT(%s), NULL) AS %s", strings.Join(joinProjections, ", "), j.Alias)
				builder = builder.Column(jsonBuild)

				// (เพิ่มเข้า GroupBy สำหรับ Nesting Mode)
				for _, p := range j.Projections {
					if joinViewCols[p.Column] && p.Column != "" {
						groupByColumns = append(groupByColumns, fmt.Sprintf("%s.%s", toAlias, p.Column))
					}
				}
			}
		} else if j.Relation == "one_to_many" {
			// (Join 1:N ทำงาน *เฉพาะ* Nesting Mode)
			if !isAggregateQuery {
				var joinProjections []string
				for _, p := range j.Projections {
					if !joinViewCols[p.Column] {
						continue
					}
					jsonKey := fmt.Sprintf("'%s'", p.Alias)
					jsonValue := fmt.Sprintf("%s.%s", toAlias, p.Column)
					joinProjections = append(joinProjections, jsonKey, jsonValue)
				}
				if len(joinProjections) == 0 {
					continue
				}

				subQuery := fmt.Sprintf(
					`(SELECT COALESCE(JSON_AGG(JSON_BUILD_OBJECT(%s)), '[]') FROM %s AS %s WHERE %s) AS %s`,
					strings.Join(joinProjections, ", "), j.TableTo, toAlias, onClause, j.Alias,
				)
				builder = builder.Column(subQuery)
			}
			// (ถ้าเป็น Aggregate Mode, เราจะ *ข้าม* 1:N Join ทั้งหมด)
		}
	}

	// --- 2. Projections (Base Table & Aggregates) ---
	allowedSortAliases := make(map[string]string)
	if len(queryPlan.Projections) == 0 {
		return builder, fmt.Errorf("QueryPlan.Projections (for base table) must not be empty")
	}

	for _, p := range queryPlan.Projections {
		var colSQL string
		var sortCol string

		if p.Expr != nil {
			// --- Aggregate Projection (e.g., COUNT(*)) ---
			if !isAggregateQuery {
				// (ข้าม Expr ถ้าไม่ได้อยู่ในโหมด Aggregate)
				continue
			}

			// หา Alias ของตาราง
			tableAlias, ok := am.Get(p.Expr.TableName)
			if !ok {
				return builder, fmt.Errorf("expr projection table '%s' not found in alias map", p.Expr.TableName)
			}

			field := p.Expr.Field
			if field != "*" {
				field = fmt.Sprintf("%s.%s", tableAlias, field) // "t0.id"
			}

			colSQL = fmt.Sprintf("%s(%s) AS %s", p.Expr.Operator, field, p.Alias) // "COUNT(t0.id) AS total"
			sortCol = p.Alias                                                     // (Aggregates สามารถ sort ด้วย Alias ได้)

		} else {
			// --- Simple Column Projection (e.g., person_data.gender_id) ---

			// (หา Alias ของตาราง)
			tableName := p.TableName
			if tableName == "" {
				tableName = am.fromTable // (ถ้าไม่ระบุ, ให้เป็นตาราง From)
			}
			tableAlias, ok := am.Get(tableName)
			if !ok {
				continue // ข้าม (อาจจะ Join ไม่ได้เลือก)
			}

			// (เช็ค ViewMap)
			viewCols, ok := viewMap[tableName]
			if !ok || !viewCols[p.Column] {
				continue
			}

			colSQL = fmt.Sprintf("%s.%s", tableAlias, p.Column)
			sortCol = colSQL // (Sort ด้วยชื่อคอลัมน์จริง)
			if p.Alias != "" {
				colSQL = fmt.Sprintf("%s AS %s", colSQL, p.Alias)
				sortCol = p.Alias // (ถ้ามี Alias, Sort ด้วย Alias)
			}

			if isAggregateQuery {
				// (ถ้าเป็น Agg Mode, ไม่ต้องเพิ่มเข้า groupByColumns
				// เพราะ GroupBy จะถูกกำหนดโดย queryPlan.GroupBy)
			} else {
				// (ถ้าเป็น Nesting Mode, เพิ่มเข้า groupByColumns)
				groupByColumns = append(groupByColumns, fmt.Sprintf("%s.%s", tableAlias, p.Column))
			}
		}

		builder = builder.Column(colSQL)
		allowedSortAliases[p.Alias] = sortCol
	}

	// --- 3. Where ---
	// (ต้องสร้าง allowedWhereMap *หลังจาก* Join เพื่อให้ am.Get() ทำงาน)
	allowedWhereMap, err := buildAllowedWhereMap(queryPlan.WhereAllow, am)
	if err != nil {
		return builder, err
	}
	var masterWhereClause squirrel.Sqlizer
	if strings.ToUpper(logicalOperator) == "OR" {
		masterWhereClause = squirrel.Or{}
	} else {
		masterWhereClause = squirrel.And{}
	}
	if filterGroups != nil {
		for _, group := range filterGroups {
			if len(group) == 0 {
				continue
			}
			orClause := squirrel.Or{}
			for _, f := range group {
				tableAlias, ok := am.Get(f.TableName)
				if !ok {
					continue // ข้าม filter นี้, ตารางนี้ไม่ได้ถูก Join ใน View
				}
				fieldWithAlias := fmt.Sprintf("%s.%s", tableAlias, f.Field)

				// Validate Data Type
				col, ok := schemaMap[f.TableName][f.Field]
				if !ok {
					return builder, errs.NewBadRequestError(fmt.Sprintf("filter field '%s.%s' not found in schema", f.TableName, f.Field))
				}
				if col.Enum != nil {
					if err := validateEnum(col.Enum, f.Value); err != nil {
						return builder, errs.NewBadRequestError(fmt.Sprintf("invalid filter value for '%s.%s': %v", f.TableName, f.Field, err))
					}
				}
				if err := validateDataType(col.DataType, f.Value); err != nil {
					return builder, errs.NewBadRequestError(fmt.Sprintf("invalid filter value for '%s.%s': %v", f.TableName, f.Field, err))
				}

				// ตรวจสอบ Allow
				isAllowed := false
				if ops, ok := allowedWhereMap[fieldWithAlias]; ok {
					if _, ok := ops[f.Operator]; ok {
						isAllowed = true
					}
				}
				if !isAllowed {
					return builder, errs.NewBadRequestError(fmt.Sprintf("filter is not allowed: %s.%s %s", f.TableName, f.Field, f.Operator))
				}

				// สร้าง Expression
				expr, err := buildSquirrelExpr(fieldWithAlias, f.Operator, f.Value)
				if err != nil {
					return builder, errs.NewBadRequestError(err.Error())
				}
				orClause = append(orClause, expr)
			}
			if strings.ToUpper(logicalOperator) == "OR" {
				masterWhereClause = append(masterWhereClause.(squirrel.Or), orClause)
			} else {
				masterWhereClause = append(masterWhereClause.(squirrel.And), orClause)
			}
		}
	}
	addWhere := false
	if op, ok := masterWhereClause.(squirrel.Or); ok && len(op) > 0 {
		addWhere = true
	} else if op, ok := masterWhereClause.(squirrel.And); ok && len(op) > 0 {
		addWhere = true
	}
	if addWhere {
		builder = builder.Where(masterWhereClause)
	}

	// --- 4. Group By ---
	if isAggregateQuery {
		// (AGGREGATE MODE)
		var finalGroupBy []string
		for _, gb := range queryPlan.GroupBy {
			tableAlias, ok := am.Get(gb.TableName)
			if !ok {
				continue // ข้าม GroupBy ที่อ้างอิงตารางที่ไม่ได้ Join
			}
			finalGroupBy = append(finalGroupBy, fmt.Sprintf("%s.%s", tableAlias, gb.Field))
		}

		if len(finalGroupBy) > 0 {
			builder = builder.GroupBy(finalGroupBy...)
		} else {
			return builder, fmt.Errorf("aggregate query requires at least one valid GroupBy field")
		}
	} else {
		if len(groupByColumns) > 0 {
			builder = builder.GroupBy(groupByColumns...)
		}
	}

	// --- 5. Sorting & Pagination---
	if sortBy != "" {
		sortColumn, ok := allowedSortAliases[sortBy]
		if !ok {
			return builder, errs.NewBadRequestError(fmt.Sprintf("sort_by field '%s' is not an allowed projection alias for sorting", sortBy))
		}
		order := constants.SORT_ORDER_ASC
		if strings.ToUpper(sortOrder) == constants.SORT_ORDER_DESC {
			order = constants.SORT_ORDER_DESC
		}
		builder = builder.OrderBy(fmt.Sprintf("%s %s", sortColumn, order))
	}
	if pagination != nil {
		builder = builder.Limit(uint64(pagination.GetLimit()))
		builder = builder.Offset(uint64(pagination.GetOffset()))
	}
	return builder, nil
}

func buildCountSQLBuilder(
	ctx context.Context,
	schemaMap map[string]map[string]entity.Column,
	queryPlan *entity.QueryPlan,
	filterGroups [][]entity.FilterInput,
	logicalOperator string,
) (squirrel.SelectBuilder, error) {

	// --- 0. Alias Management ---
	if queryPlan.From == nil || queryPlan.From.Table == "" {
		return squirrel.SelectBuilder{}, fmt.Errorf("QueryPlan.From.Table is required")
	}
	am := newAliasManager(queryPlan.From.Table)
	countCol := fmt.Sprintf("COUNT(DISTINCT %s.id)", am.fromAlias)
	builder := psqlBuilder.Select(countCol).From(fmt.Sprintf("%s AS %s", am.fromTable, am.fromAlias))

	// --- 3. Where ---
	allowedWhereMap, err := buildAllowedWhereMap(queryPlan.WhereAllow, am)
	if err != nil {
		return builder, err
	}
	tablesToJoin := make(map[string]bool)
	if filterGroups != nil {
		for _, group := range filterGroups {
			for _, f := range group {
				if f.TableName != am.fromTable {
					tablesToJoin[f.TableName] = true
				}
			}
		}
	}
	for _, j := range queryPlan.Joins {
		if (j.Relation == "one_to_one" || j.Relation == "many_to_one") && tablesToJoin[j.TableTo] {
			fromAlias, _ := am.Get(j.TableFrom)
			toAlias := am.generate(j.TableTo)
			onClause, err := buildOnClause(j.Condition, fromAlias, toAlias)
			if err != nil {
				return builder, fmt.Errorf("count: failed ON clause for join '%s': %w", j.Alias, err)
			}
			joinTableSQL := fmt.Sprintf("%s AS %s", j.TableTo, toAlias)
			builder = builder.LeftJoin(fmt.Sprintf("%s ON %s", joinTableSQL, onClause))
		}
	}

	var masterWhereClause squirrel.Sqlizer
	if strings.ToUpper(logicalOperator) == "OR" {
		masterWhereClause = squirrel.Or{}
	} else {
		masterWhereClause = squirrel.And{}
	}
	if filterGroups != nil {
		for _, group := range filterGroups {
			if len(group) == 0 {
				continue
			}
			orClause := squirrel.Or{}
			for _, f := range group {
				tableAlias, ok := am.Get(f.TableName)
				if !ok {
					continue
				}
				fieldWithAlias := fmt.Sprintf("%s.%s", tableAlias, f.Field)

				//  Validate Data Type
				col, ok := schemaMap[f.TableName][f.Field]
				if !ok {
					continue // (Count ข้าม)
				}
				if col.Enum != nil {
					if err := validateEnum(col.Enum, f.Value); err != nil {
						continue
					}
				}
				if err := validateDataType(col.DataType, f.Value); err != nil {
					continue // (Count ข้าม)
				}

				// ตรวจสอบ Allow
				isAllowed := false
				if ops, ok := allowedWhereMap[fieldWithAlias]; ok {
					if _, ok := ops[f.Operator]; ok {
						isAllowed = true
					}
				}
				if !isAllowed {
					continue
				}
				// Expression
				expr, err := buildSquirrelExpr(fieldWithAlias, f.Operator, f.Value)
				if err != nil {
					continue
				}
				orClause = append(orClause, expr)
			}
			if strings.ToUpper(logicalOperator) == "OR" {
				masterWhereClause = append(masterWhereClause.(squirrel.Or), orClause)
			} else {
				masterWhereClause = append(masterWhereClause.(squirrel.And), orClause)
			}
		}
	}
	addWhere := false
	if op, ok := masterWhereClause.(squirrel.Or); ok && len(op) > 0 {
		addWhere = true
	} else if op, ok := masterWhereClause.(squirrel.And); ok && len(op) > 0 {
		addWhere = true
	}
	if addWhere {
		builder = builder.Where(masterWhereClause)
	}

	return builder, nil
}

func (p *psqlDataRepository) ExecuteQuery(
	ctx context.Context,
	sourceID *uuid.UUID,
	schema *entity.Schema,
	policies *entity.Policies,
	filterGroups [][]entity.FilterInput,
	logicalOperator string,
	paginator *helperModel.Paginator,
	viewName string,
	sortBy string,
	sortOrder string,
) ([]map[string]interface{}, error) {
	client, err := p.dbConnectionManager.GetConnection(ctx, *sourceID)
	if err != nil {
		return nil, err
	}
	runtime := policies.Runtime
	if runtime == nil {
		return nil, fmt.Errorf("policies.runtime is missing")
	}

	// --- View Logic ---
	activeViewName := runtime.DefaultView
	if viewName != "" {
		activeViewName = viewName
	}
	viewConfigs, ok := policies.Views[activeViewName]
	if !ok || len(viewConfigs) == 0 {
		return nil, errs.NewNotFoundError(fmt.Sprintf("view '%s' not found or is empty in policies", activeViewName))
	}
	viewMap := createViewMap(viewConfigs)

	// สร้าง Schema Map
	schemaMap := createSchemaMap(schema)

	// 1. สร้าง Base Query Builder
	baseBuilder, err := buildRuntimeSQLBuilder(ctx, schemaMap, &runtime.Query, filterGroups, logicalOperator, paginator, sortBy, sortOrder, viewMap)
	if err != nil {
		return nil, fmt.Errorf("failed to build base query: %w", err)
	}

	// 2. สร้าง Count Query Builder
	countBuilder, err := buildCountSQLBuilder(ctx, schemaMap, &runtime.Query, filterGroups, logicalOperator)
	if err != nil {
		return nil, fmt.Errorf("failed to build count query: %w", err)
	}

	countSQL, countArgs, err := countBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build count sql: %w", err)
	}
	var total int64
	if err := client.GetClient().GetContext(ctx, &total, countSQL, countArgs...); err != nil {
		return nil, fmt.Errorf("failed to execute count query: %w", err)
	}

	// 3. สร้าง Query สำหรับดึงข้อมูลจริง
	querySQL, queryArgs, err := baseBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build final query: %w", err)
	}

	// 4. Execute และ Scan
	rows, err := client.GetClient().QueryxContext(ctx, querySQL, queryArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute data query: %w", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}
	finalItems := make([]map[string]interface{}, 0)
	for rows.Next() {
		scanArgs := make([]interface{}, len(cols))
		for i := range scanArgs {
			scanArgs[i] = new(sql.RawBytes)
		}
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		rowMap := make(map[string]interface{})
		for i, colName := range cols {
			rawBytes := *(scanArgs[i].(*sql.RawBytes))
			if rawBytes == nil {
				rowMap[colName] = nil
				continue
			}
			if len(rawBytes) > 0 && (rawBytes[0] == '{' || rawBytes[0] == '[') {
				var v interface{}
				if err := json.Unmarshal(rawBytes, &v); err == nil {
					rowMap[colName] = v
				} else {
					rowMap[colName] = string(rawBytes)
				}
			} else {
				rowMap[colName] = string(rawBytes)
			}
		}
		finalItems = append(finalItems, rowMap)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	if paginator != nil {
		paginator.SetPaginatorByAllRows(int(total))
	}
	return finalItems, nil
}

func (p *psqlDataRepository) ExecuteQueryByKey(
	ctx context.Context,
	sourceID *uuid.UUID,
	schema *entity.Schema,
	policies *entity.Policies,
	key interface{},
	viewName string,
) (map[string]interface{}, error) {

	runtime := policies.Runtime
	if runtime.KeyField == "" {
		return nil, fmt.Errorf("RuntimePolicy.KeyField is not defined")
	}
	if runtime.Query.From.Table == "" {
		return nil, fmt.Errorf("RuntimePolicy.Query.From.Table is not defined")
	}

	// --- View Logic ---
	activeViewName := runtime.DefaultView
	if viewName != "" {
		activeViewName = viewName
	}
	viewConfigs, ok := policies.Views[activeViewName]
	if !ok || len(viewConfigs) == 0 {
		return nil, errs.NewNotFoundError(fmt.Sprintf("view '%s' not found or is empty in policies", activeViewName))
	}
	viewMap := createViewMap(viewConfigs)

	filterGroups := [][]entity.FilterInput{
		{{
			TableName: runtime.Query.From.Table,
			Field:     runtime.KeyField,
			Operator:  "=",
			Value:     key,
		}},
	}
	client, err := p.dbConnectionManager.GetConnection(ctx, *sourceID)
	if err != nil {
		return nil, err
	}

	// สร้าง Schema Map
	schemaMap := createSchemaMap(schema)

	// 1. สร้าง Base Query Builder
	builder, err := buildRuntimeSQLBuilder(ctx, schemaMap, &runtime.Query, filterGroups, "AND", nil, "", "", viewMap)
	if err != nil {
		return nil, fmt.Errorf("failed to build key query: %w", err)
	}
	builder = builder.Limit(1)

	querySQL, queryArgs, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build final key query: %w", err)
	}
	// 3. Execute
	row := client.GetClient().QueryRowxContext(ctx, querySQL, queryArgs...)
	if row.Err() != nil {
		return nil, fmt.Errorf("failed to execute key query: %w", row.Err())
	}
	cols, err := row.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns for key query: %w", err)
	}
	scanArgs := make([]interface{}, len(cols))
	for i := range scanArgs {
		scanArgs[i] = new(any)
	}
	if err := row.Scan(scanArgs...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan key query row: %w", err)
	}
	resultMap := make(map[string]interface{})
	for i, colName := range cols {
		val := *(scanArgs[i].(*any))
		if val == nil {
			resultMap[colName] = nil
			continue
		}
		if rawBytes, ok := val.([]byte); ok {
			if len(rawBytes) > 0 && (rawBytes[0] == '{' || rawBytes[0] == '[') {
				var v interface{}
				if err := json.Unmarshal(rawBytes, &v); err == nil {
					resultMap[colName] = v
				} else {
					resultMap[colName] = string(rawBytes)
				}
			} else {
				resultMap[colName] = string(rawBytes)
			}
		} else {
			resultMap[colName] = val
		}
	}
	return resultMap, nil
}

func buildCreateSQLBuilder(ctx context.Context, queryPlan *entity.QueryPlan, validatedData map[string]interface{}) (squirrel.InsertBuilder, error) {
	if queryPlan.From == nil || queryPlan.From.Table == "" {
		return squirrel.InsertBuilder{}, fmt.Errorf("WritePolicy.Query.From.Table is required for CREATE")
	}

	builder := psqlBuilder.Insert(queryPlan.From.Table)

	var columns []string
	var values []interface{}
	for key, val := range validatedData {
		columns = append(columns, key)
		values = append(values, val)
	}

	builder = builder.Columns(columns...).Values(values...)
	return builder, nil
}

func (p *psqlDataRepository) ExecuteCreate(
	ctx context.Context,
	sourceID *uuid.UUID,
	schema entity.Schema,
	writePolicy *entity.WritePolicy,
	data map[string]interface{},
) (map[string]interface{}, error) {

	// 1. Validate และ Prepare ข้อมูลก่อน
	validatedData, err := p.validateAndPrepareData(schema, writePolicy, data)
	if err != nil {
		return nil, errs.NewBadRequestError(fmt.Sprintf("create validation failed: %v", err))
	}

	// (ถ้า validatedData ไม่มี field เลย ก็ไม่ควร Insert)
	if len(validatedData) == 0 {
		return nil, errs.NewBadRequestError("no valid fields provided for creation")
	}

	// 2. ดึง Connection
	client, err := p.dbConnectionManager.GetConnection(ctx, *sourceID)
	if err != nil {
		return nil, err
	}

	// 3. สร้าง Builder (โดยใช้ข้อมูลที่ Validate แล้ว)
	builder, err := buildCreateSQLBuilder(ctx, &writePolicy.Query, validatedData)
	if err != nil {
		return nil, err
	}

	// 4. สร้าง SQL
	querySQL, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	// 5. สร้าง RETURNING clause
	colsToReturn := make(map[string]bool)
	if writePolicy.KeyField != "" {
		colsToReturn[writePolicy.KeyField] = true
	}
	for fieldName := range validatedData {
		colsToReturn[fieldName] = true
	}

	quotedCols := make([]string, 0, len(colsToReturn))
	for col := range colsToReturn {
		quotedCols = append(quotedCols, fmt.Sprintf("\"%s\"", col))
	}
	returningSQL := querySQL + " RETURNING " + strings.Join(quotedCols, ", ")

	// 6. Execute
	row := make(map[string]interface{})
	err = client.GetClient().QueryRowxContext(ctx, returningSQL, args...).MapScan(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to execute create with returning: %w", err)
	}

	return row, nil
}

func buildUpdateSQLBuilder(ctx context.Context, queryPlan *entity.QueryPlan, validatedData map[string]interface{}, whereConditions map[string]interface{}) (squirrel.UpdateBuilder, error) {
	if queryPlan.From == nil || queryPlan.From.Table == "" {
		return squirrel.UpdateBuilder{}, fmt.Errorf("WritePolicy.Query.From.Table is required for UPDATE")
	}

	builder := psqlBuilder.Update(queryPlan.From.Table)

	// [CHANGED] 1. SET clause (ไม่ต้องกรอง/Error)
	// (เราเช็ค len(validatedData) ใน ExecuteUpdate ไปแล้ว)
	builder = builder.SetMap(validatedData)

	// 2. WHERE clause
	if len(whereConditions) == 0 {
		return builder, fmt.Errorf("update requires at least one where condition")
	}
	builder = builder.Where(squirrel.Eq(whereConditions))

	return builder, nil
}

func (p *psqlDataRepository) ExecuteUpdate(
	ctx context.Context,
	sourceID *uuid.UUID,
	schema entity.Schema,
	writePolicy *entity.WritePolicy,
	key interface{},
	data map[string]interface{},
) (map[string]interface{}, error) {

	// 1. Validate และ Prepare ข้อมูลก่อน
	validatedData, err := p.validateAndPrepareData(schema, writePolicy, data)
	if err != nil {
		return nil, errs.NewBadRequestError(fmt.Sprintf("update validation failed: %v", err))
	}

	// (ถ้า validatedData ไม่มี field เลย ก็ไม่อัปเดต)
	if len(validatedData) == 0 {
		return nil, errs.NewBadRequestError("no valid fields provided for update")
	}

	// 2. ตรวจสอบ KeyField
	if writePolicy.KeyField == "" {
		return nil, errs.NewBadRequestError("WritePolicy.KeyField is not defined")
	}

	// 3. ดึง Connection
	client, err := p.dbConnectionManager.GetConnection(ctx, *sourceID)
	if err != nil {
		return nil, err
	}

	// 4. สร้าง whereConditions map จาก key
	whereConditions := map[string]interface{}{
		writePolicy.KeyField: key,
	}

	// 5. สร้าง SQL query
	builder, err := buildUpdateSQLBuilder(ctx, &writePolicy.Query, validatedData, whereConditions)
	if err != nil {
		return nil, fmt.Errorf("failed to build update query: %w", err)
	}

	querySQL, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to get update sql: %w", err)
	}

	// 6. สร้าง RETURNING clause
	colsToReturn := make(map[string]bool)
	colsToReturn[writePolicy.KeyField] = true
	for fieldName := range validatedData {
		colsToReturn[fieldName] = true
	}

	quotedCols := make([]string, 0, len(colsToReturn))
	for col := range colsToReturn {
		quotedCols = append(quotedCols, fmt.Sprintf("\"%s\"", col))
	}
	returningSQL := querySQL + " RETURNING " + strings.Join(quotedCols, ", ")

	// 7. Execute
	row := make(map[string]interface{})
	err = client.GetClient().QueryRowxContext(ctx, returningSQL, args...).MapScan(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errs.NewNotFoundError(fmt.Sprintf("failed to execute update with returning: %v", err))
	}

	return row, nil
}

func buildDeleteSQLBuilder(ctx context.Context, queryPlan *entity.QueryPlan, whereConditions map[string]interface{}) (squirrel.DeleteBuilder, error) {
	if queryPlan.From == nil || queryPlan.From.Table == "" {
		return squirrel.DeleteBuilder{}, fmt.Errorf("DeletePolicy.Query.From.Table is required for DELETE")
	}

	builder := psqlBuilder.Delete(queryPlan.From.Table)

	// WHERE clause
	if len(whereConditions) == 0 {
		return builder, fmt.Errorf("delete requires at least one where condition")
	}
	builder = builder.Where(squirrel.Eq(whereConditions))

	return builder, nil
}

func (p *psqlDataRepository) ExecuteDelete(
	ctx context.Context,
	sourceID *uuid.UUID,
	deletePolicy *entity.DeletePolicy,
	key interface{},
) (sql.Result, error) {
	// ตรวจสอบว่ามี KeyField ใน Policy
	if deletePolicy.KeyField == "" {
		return nil, fmt.Errorf("DeletePolicy.KeyField is not defined")
	}

	// ดึง Connection
	client, err := p.dbConnectionManager.GetConnection(ctx, *sourceID)
	if err != nil {
		return nil, err
	}

	// สร้าง whereConditions map จาก key ที่รับเข้ามา
	whereConditions := map[string]interface{}{
		deletePolicy.KeyField: key,
	}

	// สร้าง SQL query
	builder, err := buildDeleteSQLBuilder(ctx, &deletePolicy.Query, whereConditions)
	if err != nil {
		return nil, fmt.Errorf("failed to build delete query: %w", err)
	}

	// แปลงเป็น SQL string และ arguments
	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to get delete sql: %w", err)
	}

	// Execute
	sqlResult, err := client.GetClient().ExecContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute delete: %w", err)
	}

	return sqlResult, nil
}

// validateTableDataFields validates data against table columns
func (p *psqlDataRepository) validateTableDataFields(
	columns []*entity.Columns,
	data map[string]interface{},
	isCreate bool,
) error {
	// Build column map for validation
	columnMap := make(map[string]*entity.Columns)
	for _, col := range columns {
		columnMap[col.ColumnsName] = col
	}

	// Validate data against columns
	for fieldName, fieldValue := range data {
		col, exists := columnMap[fieldName]
		if !exists {
			return fmt.Errorf("column '%s' does not exist in table", fieldName)
		}

		// Check if value is nil
		if fieldValue == nil {
			if !col.IsNullable {
				if isCreate && col.ColumnDefault == nil {
					return fmt.Errorf("column '%s' cannot be null and has no default value", fieldName)
				} else if !isCreate {
					return fmt.Errorf("column '%s' cannot be null", fieldName)
				}
			}
			continue
		}

		// Validate data type
		if err := validateDataType(col.DataType, fieldValue); err != nil {
			return fmt.Errorf("column '%s': %v", fieldName, err)
		}
	}

	// For CREATE: Check required fields (not nullable and no default)
	if isCreate {
		for _, col := range columns {
			if !col.IsNullable && col.ColumnDefault == nil {
				if _, exists := data[col.ColumnsName]; !exists {
					return fmt.Errorf("required column '%s' is missing", col.ColumnsName)
				}
			}
		}
	}

	return nil
}

// FetchTableData implements data.PsqlDataRepository.
func (p *psqlDataRepository) FetchTableData(
	ctx context.Context,
	sourceID *uuid.UUID,
	schemaName, tableName string,
	filterGroups [][]entity.FilterInput,
	logicalOperator string,
	paginator *helperModel.Paginator,
	sortBy, sortOrder string,
) ([]map[string]interface{}, error) {
	client, err := p.dbConnectionManager.GetConnection(ctx, *sourceID)
	if err != nil {
		return nil, err
	}

	// Build SELECT query
	fullTableName := fmt.Sprintf("\"%s\".\"%s\"", schemaName, tableName)
	builder := psqlBuilder.Select("*").From(fullTableName)

	// Apply WHERE filters
	var masterWhereClause squirrel.Sqlizer
	if strings.ToUpper(logicalOperator) == "OR" {
		masterWhereClause = squirrel.Or{}
	} else {
		masterWhereClause = squirrel.And{}
	}

	for _, group := range filterGroups {
		if len(group) == 0 {
			continue
		}
		orClause := squirrel.Or{}
		for _, f := range group {
			expr, err := buildSquirrelExpr(f.Field, f.Operator, f.Value)
			if err != nil {
				return nil, errs.NewBadRequestError(err.Error())
			}
			orClause = append(orClause, expr)
		}
		if strings.ToUpper(logicalOperator) == "OR" {
			masterWhereClause = append(masterWhereClause.(squirrel.Or), orClause)
		} else {
			masterWhereClause = append(masterWhereClause.(squirrel.And), orClause)
		}
	}

	addWhere := false
	if op, ok := masterWhereClause.(squirrel.Or); ok && len(op) > 0 {
		addWhere = true
	} else if op, ok := masterWhereClause.(squirrel.And); ok && len(op) > 0 {
		addWhere = true
	}
	if addWhere {
		builder = builder.Where(masterWhereClause)
	}

	// Apply sorting
	if sortBy != "" {
		order := constants.SORT_ORDER_ASC
		if strings.ToUpper(sortOrder) == constants.SORT_ORDER_DESC {
			order = constants.SORT_ORDER_DESC
		}
		builder = builder.OrderBy(fmt.Sprintf("%s %s", sortBy, order))
	}

	// Count total rows
	if len(filterGroups) > 0 {
		// Build count query with same WHERE filters
		countBuilder := psqlBuilder.Select("COUNT(*)").From(fullTableName)

		// Apply same WHERE clause
		if addWhere {
			countBuilder = countBuilder.Where(masterWhereClause)
		}

		countSQL, countArgs, err := countBuilder.ToSql()
		if err == nil {
			var total int64
			if err := client.GetClient().GetContext(ctx, &total, countSQL, countArgs...); err != nil {
				return nil, err
			}
			if paginator != nil {
				paginator.SetPaginatorByAllRows(int(total))
			}
		}
	} else {
		countBuilder := psqlBuilder.Select("COUNT(*)").From(fullTableName)
		countSQL, countArgs, _ := countBuilder.ToSql()
		var total int64
		if err := client.GetClient().GetContext(ctx, &total, countSQL, countArgs...); err != nil {
			return nil, err
		}
		if paginator != nil {
			paginator.SetPaginatorByAllRows(int(total))
		}
	}

	// Apply pagination
	if paginator != nil {
		builder = builder.Limit(uint64(paginator.GetLimit()))
		builder = builder.Offset(uint64(paginator.GetOffset()))
	}

	// Execute query
	querySQL, queryArgs, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := client.GetClient().QueryxContext(ctx, querySQL, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	results := make([]map[string]interface{}, 0)
	for rows.Next() {
		scanArgs := make([]interface{}, len(cols))
		for i := range scanArgs {
			scanArgs[i] = new(sql.RawBytes)
		}
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
		}
		rowMap := make(map[string]interface{})
		for i, colName := range cols {
			rawBytes := *(scanArgs[i].(*sql.RawBytes))
			if rawBytes == nil {
				rowMap[colName] = nil
				continue
			}
			if len(rawBytes) > 0 && (rawBytes[0] == '{' || rawBytes[0] == '[') {
				var v interface{}
				if err := json.Unmarshal(rawBytes, &v); err == nil {
					rowMap[colName] = v
				} else {
					rowMap[colName] = string(rawBytes)
				}
			} else {
				rowMap[colName] = string(rawBytes)
			}
		}
		results = append(results, rowMap)
	}

	return results, nil
}

// FetchTableDataByKey implements data.PsqlDataRepository.
func (p *psqlDataRepository) FetchTableDataByKey(
	ctx context.Context,
	sourceID *uuid.UUID,
	schemaName, tableName, keyField string,
	keyValue interface{},
) (map[string]interface{}, error) {
	client, err := p.dbConnectionManager.GetConnection(ctx, *sourceID)
	if err != nil {
		return nil, err
	}

	fullTableName := fmt.Sprintf("\"%s\".\"%s\"", schemaName, tableName)
	builder := psqlBuilder.Select("*").
		From(fullTableName).
		Where(squirrel.Eq{keyField: keyValue}).
		Limit(1)

	querySQL, queryArgs, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	row := client.GetClient().QueryRowxContext(ctx, querySQL, queryArgs...)
	if row.Err() != nil {
		return nil, row.Err()
	}

	cols, err := row.Columns()
	if err != nil {
		return nil, err
	}

	scanArgs := make([]interface{}, len(cols))
	for i := range scanArgs {
		scanArgs[i] = new(any)
	}

	if err := row.Scan(scanArgs...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	resultMap := make(map[string]interface{})
	for i, colName := range cols {
		val := *(scanArgs[i].(*any))
		if val == nil {
			resultMap[colName] = nil
			continue
		}
		if rawBytes, ok := val.([]byte); ok {
			if len(rawBytes) > 0 && (rawBytes[0] == '{' || rawBytes[0] == '[') {
				var v interface{}
				if err := json.Unmarshal(rawBytes, &v); err == nil {
					resultMap[colName] = v
				} else {
					resultMap[colName] = string(rawBytes)
				}
			} else {
				resultMap[colName] = string(rawBytes)
			}
		} else {
			resultMap[colName] = val
		}
	}

	return resultMap, nil
}

// CreateTableData implements data.PsqlDataRepository.
func (p *psqlDataRepository) CreateTableData(
	ctx context.Context,
	sourceID *uuid.UUID,
	schemaName, tableName string,
	tableColumns []*entity.Columns,
	data map[string]interface{},
) (map[string]interface{}, error) {
	if len(data) == 0 {
		return nil, errs.NewBadRequestError("data cannot be empty")
	}

	// Validate data against columns
	if err := p.validateTableDataFields(tableColumns, data, true); err != nil {
		return nil, errs.NewBadRequestError(fmt.Sprintf("validation failed: %v", err))
	}

	client, err := p.dbConnectionManager.GetConnection(ctx, *sourceID)
	if err != nil {
		return nil, err
	}

	fullTableName := fmt.Sprintf("\"%s\".\"%s\"", schemaName, tableName)
	builder := psqlBuilder.Insert(fullTableName)

	var columns []string
	var values []interface{}
	for key, val := range data {
		fmt.Println("Inserting:", key, "=", val)
		columns = append(columns, key)
		values = append(values, val)
	}

	builder = builder.Columns(columns...).Values(values...)

	querySQL, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	// Add RETURNING clause to get all columns
	quotedCols := make([]string, len(columns))
	for i, col := range columns {
		quotedCols[i] = fmt.Sprintf("\"%s\"", col)
	}
	returningSQL := querySQL + " RETURNING " + strings.Join(quotedCols, ", ")
	fmt.Println("Create SQL:", returningSQL)
	row := make(map[string]interface{})
	err = client.GetClient().QueryRowxContext(ctx, returningSQL, args...).MapScan(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to execute create: %w", err)
	}

	return row, nil
}

// UpdateTableData implements data.PsqlDataRepository.
func (p *psqlDataRepository) UpdateTableData(
	ctx context.Context,
	sourceID *uuid.UUID,
	schemaName, tableName, keyField string,
	keyValue interface{},
	tableColumns []*entity.Columns,
	data map[string]interface{},
) (map[string]interface{}, error) {
	if len(data) == 0 {
		return nil, errs.NewBadRequestError("data cannot be empty")
	}

	// Validate data against columns
	if err := p.validateTableDataFields(tableColumns, data, false); err != nil {
		return nil, errs.NewBadRequestError(fmt.Sprintf("validation failed: %v", err))
	}

	client, err := p.dbConnectionManager.GetConnection(ctx, *sourceID)
	if err != nil {
		return nil, err
	}

	fullTableName := fmt.Sprintf("\"%s\".\"%s\"", schemaName, tableName)
	builder := psqlBuilder.Update(fullTableName).
		SetMap(data).
		Where(squirrel.Eq{keyField: keyValue})

	querySQL, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	// Create RETURNING clause with updated fields + key field
	colsToReturn := make(map[string]bool)
	colsToReturn[keyField] = true
	for fieldName := range data {
		colsToReturn[fieldName] = true
	}

	quotedCols := make([]string, 0, len(colsToReturn))
	for col := range colsToReturn {
		quotedCols = append(quotedCols, fmt.Sprintf("\"%s\"", col))
	}
	returningSQL := querySQL + " RETURNING " + strings.Join(quotedCols, ", ")

	row := make(map[string]interface{})
	err = client.GetClient().QueryRowxContext(ctx, returningSQL, args...).MapScan(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to execute update: %w", err)
	}

	return row, nil
}

// DeleteTableData implements data.PsqlDataRepository.
func (p *psqlDataRepository) DeleteTableData(
	ctx context.Context,
	sourceID *uuid.UUID,
	schemaName, tableName, keyField string,
	keyValue interface{},
) (sql.Result, error) {
	client, err := p.dbConnectionManager.GetConnection(ctx, *sourceID)
	if err != nil {
		return nil, err
	}

	fullTableName := fmt.Sprintf("\"%s\".\"%s\"", schemaName, tableName)
	builder := psqlBuilder.Delete(fullTableName).
		Where(squirrel.Eq{keyField: keyValue})

	querySQL, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	sqlResult, err := client.GetClient().ExecContext(ctx, querySQL, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute delete: %w", err)
	}

	return sqlResult, nil
}
