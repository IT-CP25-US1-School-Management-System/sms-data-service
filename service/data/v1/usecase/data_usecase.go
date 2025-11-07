package usecase

import (
	"context"
	"fmt"
	"strings"

	helperModel "github.com/GodeFvt/go-backend/helper/models"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/constants"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/errs"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/dto"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/entity"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/filter"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/data/v1"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/utils/crypto"
	"github.com/gofrs/uuid"
)

type dataUsecase struct {
	dataRepo     data.PsqlDataRepository
	datasetRepo  data.PsqlDatasetRepository
	redisRepo    data.RedisRepository
	cryptoSecret string
}

// DeleteSourceByID implements data.DataUsecase.
func (d *dataUsecase) DeleteSourceByID(ctx context.Context, sourceID *uuid.UUID) error {
	exist, err := d.datasetRepo.ExistSourceByID(ctx, sourceID)
	if err != nil {
		return err
	}
	if !exist {
		return errs.NewNotFoundError(constants.ERR_SOURCE_NOT_FOUND)
	}

	return d.datasetRepo.DeleteSourceByID(ctx, sourceID)
}

// ActivateSourceByID implements data.DataUsecase.
func (d *dataUsecase) ActivateSourceByID(ctx context.Context, sourceID *uuid.UUID) error {
	exist, err := d.datasetRepo.ExistSourceByID(ctx, sourceID)
	if err != nil {
		return err
	}
	if !exist {
		return errs.NewNotFoundError(constants.ERR_SOURCE_NOT_FOUND)
	}

	return d.datasetRepo.ActivateSourceByID(ctx, sourceID)
}

// DeactivateSourceByID implements data.DataUsecase.
func (d *dataUsecase) DeactivateSourceByID(ctx context.Context, sourceID *uuid.UUID) error {
	exist, err := d.datasetRepo.ExistSourceByID(ctx, sourceID)
	if err != nil {
		return err
	}
	if !exist {
		return errs.NewNotFoundError(constants.ERR_SOURCE_NOT_FOUND)
	}

	return d.datasetRepo.DeactivateSourceByID(ctx, sourceID)
}

// InsertSource implements data.DataUsecase.
func (d *dataUsecase) InsertSource(ctx context.Context, source *entity.Sources) error {
	if source == nil {
		return errs.NewBadRequestError(constants.ERR_INVALID_REQUEST_BODY)
	}
	encryptPass, err := crypto.Encrypt(source.Password, d.cryptoSecret)
	if err != nil {
		return err
	}
	source.Password = encryptPass

	now := helperModel.NewTimestampFromNow()
	source.CreatedAt = &now
	source.UpdatedAt = &now

	err = d.datasetRepo.UpsertSource(ctx, source)
	if err != nil {
		return err
	}

	infoSchema, err := d.dataRepo.FetchInformationSchemasBySourceID(ctx, source.DBType, source.ID)
	if err != nil {
		d.datasetRepo.DeleteSourceByID(ctx, source.ID)
		return err
	}
	infoTables, err := d.dataRepo.FetchInformationTablesBySourceID(ctx, source.DBType, source.ID)
	if err != nil {
		d.datasetRepo.DeleteSourceByID(ctx, source.ID)
		return err
	}
	infoColumns, err := d.dataRepo.FetchInformationColumnsBySourceID(ctx, source.DBType, source.ID)
	if err != nil {
		d.datasetRepo.DeleteSourceByID(ctx, source.ID)
		return err
	}
	err = d.datasetRepo.BatchInsertInformationDatabase(ctx, infoSchema, infoTables, infoColumns)
	if err != nil {
		d.datasetRepo.DeleteSourceByID(ctx, source.ID)
		return err
	}

	return nil
}

// UpdateSource implements data.DataUsecase.
func (d *dataUsecase) UpdateSource(ctx context.Context, sourceID *uuid.UUID, sourceUpdate *dto.UpdateSourcesDTO) error {
	if sourceUpdate == nil {
		return errs.NewBadRequestError(constants.ERR_INVALID_REQUEST_BODY)
	}
	oldSource, err := d.datasetRepo.FetchSourceByID(ctx, sourceID)
	if err != nil {
		return err
	}
	if oldSource == nil {
		return errs.NewNotFoundError(constants.ERR_SOURCE_NOT_FOUND)
	}

	oldSource.ID = sourceID
	if sourceUpdate.Name != nil {
		oldSource.Name = *sourceUpdate.Name
	}
	if sourceUpdate.Description != nil {
		oldSource.Description = sourceUpdate.Description
	}
	if sourceUpdate.Type != nil {
		oldSource.Type = *sourceUpdate.Type
	}
	if sourceUpdate.IsActive != nil {
		oldSource.IsActive = *sourceUpdate.IsActive
	}
	if sourceUpdate.Sensitivity != nil {
		oldSource.Sensitivity = *sourceUpdate.Sensitivity
	}
	if sourceUpdate.DBType != nil {
		oldSource.DBType = *sourceUpdate.DBType
	}
	if sourceUpdate.Host != nil {
		oldSource.Host = *sourceUpdate.Host
	}
	if sourceUpdate.Port != nil {
		oldSource.Port = *sourceUpdate.Port
	}
	if sourceUpdate.Username != nil {
		oldSource.Username = *sourceUpdate.Username
	}
	if sourceUpdate.DatabaseName != nil {
		oldSource.DatabaseName = *sourceUpdate.DatabaseName
	}
	if sourceUpdate.Params != nil {
		oldSource.Params = sourceUpdate.Params
	}
	if sourceUpdate.IsUpdatePassword != nil && *sourceUpdate.IsUpdatePassword && sourceUpdate.Password != nil {
		encryptPass, err := crypto.Encrypt(*sourceUpdate.Password, d.cryptoSecret)
		if err != nil {
			return err
		}
		oldSource.Password = encryptPass
	}

	now := helperModel.NewTimestampFromNow()
	oldSource.UpdatedAt = &now

	return d.datasetRepo.UpsertSource(ctx, oldSource)
}

// validateDatasetID validates the dataset ID format
func (d *dataUsecase) validateDatasetID(id string) error {
	if id == "" {
		return errs.NewBadRequestError(constants.ERR_DATASET_ID_IS_REQUIRED)
	}

	if !constants.DATASET_ID_PATTERN.MatchString(id) {
		return errs.NewBadRequestError(constants.ERR_DATASET_ID_INVALID_FORMAT)
	}

	return nil
}

// DeleteDatasetByID implements data.DataUsecase.
func (d *dataUsecase) DeleteDatasetByID(ctx context.Context, datasetID string) error {
	if err := d.validateDatasetID(datasetID); err != nil {
		return err
	}
	exist, err := d.datasetRepo.ExistDatasetByID(ctx, datasetID)
	if err != nil {
		return err
	}
	if !exist {
		return errs.NewNotFoundError(constants.ERR_DATASET_NOT_FOUND)
	}

	return d.datasetRepo.DeleteDatasetByID(ctx, datasetID)
}

// UpsertDataset implements data.DataUsecase.
func (d *dataUsecase) UpsertDataset(ctx context.Context, dataset *entity.Datasets) error {
	if dataset == nil {
		return errs.NewBadRequestError(constants.ERR_INVALID_REQUEST_BODY)
	}

	if err := d.validateDatasetID(dataset.ID); err != nil {
		return err
	}

	now := helperModel.NewTimestampFromNow()
	dataset.CreatedAt = &now
	dataset.UpdatedAt = &now

	return d.datasetRepo.UpsertDataset(ctx, dataset)
}

// FetchDatasetByID implements data.DataUsecase.
func (d *dataUsecase) FetchDatasetByID(ctx context.Context, datasetID string) (*entity.Datasets, error) {
	if err := d.validateDatasetID(datasetID); err != nil {
		return nil, err
	}
	return d.datasetRepo.FetchDatasetByID(ctx, datasetID)
}

// FetchColumnsList implements data.DataUsecase.
func (d *dataUsecase) FetchColumnsList(ctx context.Context, filter *filter.ColumnsFilter, paginator *helperModel.Paginator) ([]*entity.Columns, error) {
	return d.datasetRepo.FetchColumnsList(ctx, filter, paginator)
}

// FetchTablesList implements data.DataUsecase.
func (d *dataUsecase) FetchTablesList(ctx context.Context, filter *filter.TablesFilter, paginator *helperModel.Paginator) ([]*entity.Tables, error) {
	return d.datasetRepo.FetchTablesList(ctx, filter, paginator)
}

// FetchDatasetList implements data.DataUsecase.
func (d *dataUsecase) FetchDatasetList(ctx context.Context, filter *filter.DatasetsFilter, paginator *helperModel.Paginator) ([]*entity.Datasets, error) {
	return d.datasetRepo.FetchDatasetList(ctx, filter, paginator)
}

// FetchSchemasList implements data.DataUsecase.
func (d *dataUsecase) FetchSchemasList(ctx context.Context, filter *filter.SchemasFilter, paginator *helperModel.Paginator) ([]*entity.Schemas, error) {
	return d.datasetRepo.FetchSchemasList(ctx, filter, paginator)
}

// FetchSourceList implements data.DataUsecase.
func (d *dataUsecase) FetchSourceList(ctx context.Context, paginator *helperModel.Paginator) ([]*entity.Sources, error) {
	return d.datasetRepo.FetchSourceList(ctx, paginator)
}

// FetchDatasetVersionByID implements data.DataUsecase.
func (d *dataUsecase) FetchDatasetVersionByID(ctx context.Context, datasetID string, version string) (*entity.DatasetVersion, error) {
	if err := d.validateDatasetID(datasetID); err != nil {
		return nil, err
	}
	return d.datasetRepo.FetchDatasetVersionByID(ctx, datasetID, version)
}

// FetchDatasetVersionsList implements data.DataUsecase.
func (d *dataUsecase) FetchDatasetVersionsList(ctx context.Context, datasetID string, paginator *helperModel.Paginator) ([]*entity.DatasetVersion, error) {
	if err := d.validateDatasetID(datasetID); err != nil {
		return nil, err
	}
	return d.datasetRepo.FetchDatasetVersionsList(ctx, datasetID, paginator)
}

// UpsertDatasetVersion implements data.DataUsecase.
func (d *dataUsecase) UpsertDatasetVersion(ctx context.Context, datasetVersion *entity.DatasetVersion) error {
	if datasetVersion == nil {
		return errs.NewBadRequestError(constants.ERR_INVALID_REQUEST_BODY)
	}

	if err := d.validateDatasetID(datasetVersion.DatasetID); err != nil {
		return err
	}

	// Validate that the parent dataset exists
	exists, err := d.datasetRepo.ExistDatasetByID(ctx, datasetVersion.DatasetID)
	if err != nil {
		return err
	}
	if !exists {
		return errs.NewNotFoundError(constants.ERR_DATASET_NOT_FOUND)
	}
	now := helperModel.NewTimestampFromNow()
	datasetVersion.CreatedAt = &now
	datasetVersion.UpdatedAt = &now

	return d.datasetRepo.UpsertDatasetVersion(ctx, datasetVersion)
}

// DeleteDatasetVersionByID implements data.DataUsecase.
func (d *dataUsecase) DeleteDatasetVersionByID(ctx context.Context, datasetID string, version string) error {
	if err := d.validateDatasetID(datasetID); err != nil {
		return err
	}

	exists, err := d.datasetRepo.ExistDatasetVersionByID(ctx, datasetID, version)
	if err != nil {
		return err
	}
	if !exists {
		return errs.NewNotFoundError("dataset version not found")
	}

	return d.datasetRepo.DeleteDatasetVersionByID(ctx, datasetID, version)
}

// filterRuntimePolicyByViewAndColumns filters runtime policy based on view and requested columns
func (d *dataUsecase) filterRuntimePolicyByViewAndColumns(runtime entity.RuntimePolicy, viewName string, requestedColumns []string, views map[string][]string) (*entity.RuntimePolicy, error) {
	// Make a copy of the runtime policy
	filteredRuntime := runtime

	// Debug: Log input parameters
	fmt.Printf("DEBUG: viewName='%s', requestedColumns=%v\n", viewName, requestedColumns)
	fmt.Printf("DEBUG: Available views: %v\n", views)

	// If view is specified, validate and get allowed columns
	var allowedColumns []string
	if viewName != "" {
		if viewColumns, exists := views[viewName]; exists {
			allowedColumns = viewColumns
			fmt.Printf("DEBUG: Found view '%s' with columns: %v\n", viewName, allowedColumns)
		} else {
			return nil, errs.NewBadRequestError(fmt.Sprintf("view '%s' not found", viewName))
		}
	}

	// If specific columns are requested, validate against view (if specified) or all available columns
	if len(requestedColumns) > 0 {
		var validColumns []string

		if len(allowedColumns) > 0 {
			// Validate against view columns
			allowedSet := make(map[string]bool)
			for _, col := range allowedColumns {
				allowedSet[col] = true
			}

			for _, col := range requestedColumns {
				if allowedSet[col] {
					validColumns = append(validColumns, col)
				} else {
					return nil, errs.NewBadRequestError(fmt.Sprintf("column '%s' is not allowed in view '%s'", col, viewName))
				}
			}
		} else {
			// If no view specified, validate against all available projections
			availableColumns := make(map[string]bool)
			for _, proj := range runtime.Query.Projections {
				if proj.Alias != "" {
					availableColumns[proj.Alias] = true
				} else if proj.Column != "" {
					availableColumns[proj.Column] = true
				}
			}

			for _, col := range requestedColumns {
				if availableColumns[col] {
					validColumns = append(validColumns, col)
				} else {
					return nil, errs.NewBadRequestError(fmt.Sprintf("column '%s' is not available", col))
				}
			}
		}

		// Filter projections based on valid columns
		filteredProjections := []entity.Projection{}
		for _, proj := range runtime.Query.Projections {
			include := false

			// Get the projection's output name (alias or column name)
			var projOutputName string
			if proj.Alias != "" {
				projOutputName = proj.Alias
			} else if proj.Column != "" {
				// Extract column name without table prefix
				parts := strings.Split(proj.Column, ".")
				if len(parts) > 1 {
					projOutputName = parts[len(parts)-1] // Take the last part (column name)
				} else {
					projOutputName = proj.Column
				}
			} else if proj.Expr != nil && proj.Expr.Field != "" {
				// For expressions, use the expression as is or extract meaningful name
				projOutputName = proj.Expr.Field
			}

			// Check if this projection should be included
			for _, col := range validColumns {
				if projOutputName == col {
					include = true
					break
				}
			}

			if include {
				filteredProjections = append(filteredProjections, proj)
			}
		}

		filteredRuntime.Query.Projections = filteredProjections
	} else if len(allowedColumns) > 0 {
		// If view is specified but no specific columns, filter by view columns
		fmt.Printf("DEBUG: Filtering projections by view columns: %v\n", allowedColumns)

		filteredProjections := []entity.Projection{}
		for _, proj := range runtime.Query.Projections {
			include := false

			// Get the projection's output name (alias or column name)
			var projOutputName string
			if proj.Alias != "" {
				projOutputName = proj.Alias
			} else if proj.Column != "" {
				// Extract column name without table prefix
				parts := strings.Split(proj.Column, ".")
				if len(parts) > 1 {
					projOutputName = parts[len(parts)-1] // Take the last part (column name)
				} else {
					projOutputName = proj.Column
				}
			} else if proj.Expr != nil && proj.Expr.Field != "" {
				// For expressions, use the expression as is or extract meaningful name
				projOutputName = proj.Expr.Field
			}

			fmt.Printf("DEBUG: Checking projection - Column: '%s', Alias: '%s', OutputName: '%s'\n", proj.Column, proj.Alias, projOutputName)

			// Check if this projection should be included in the view
			for _, col := range allowedColumns {
				if projOutputName == col {
					include = true
					fmt.Printf("DEBUG: Including projection '%s' (matches view column '%s')\n", projOutputName, col)
					break
				}
			}

			if !include {
				fmt.Printf("DEBUG: Excluding projection '%s'\n", projOutputName)
			}

			if include {
				filteredProjections = append(filteredProjections, proj)
			}
		}

		fmt.Printf("DEBUG: Original projections count: %d, Filtered count: %d\n", len(runtime.Query.Projections), len(filteredProjections))
		filteredRuntime.Query.Projections = filteredProjections
	}

	return &filteredRuntime, nil
}

// ServingDatasetVersionData implements data.DataUsecase.
func (d *dataUsecase) ServingDatasetVersionData(ctx context.Context, datasetID string, version string, paginator *helperModel.Paginator, viewName string, requestedColumns []string) ([]map[string]interface{}, error) {
	if err := d.validateDatasetID(datasetID); err != nil {
		return nil, err
	}

	// 1. Get dataset version to get policies
	datasetVersion, err := d.datasetRepo.FetchDatasetVersionByID(ctx, datasetID, version)
	if err != nil {
		return nil, err
	}
	if datasetVersion == nil {
		return nil, errs.NewNotFoundError("dataset version not found")
	}

	// 2. Validate view and columns
	filteredRuntime, err := d.filterRuntimePolicyByViewAndColumns(datasetVersion.Policies.Runtime, viewName, requestedColumns, datasetVersion.Policies.Views)
	if err != nil {
		return nil, err
	}

	// 3. Build SQL from filtered runtime policy
	query, args, err := d.dataRepo.BuildRuntimeSQL(ctx, datasetVersion.SourceID, filteredRuntime)
	if err != nil {
		return nil, err
	}

	// 4. Execute query
	results, err := d.dataRepo.ExecuteQuery(ctx, datasetVersion.SourceID, query, args, paginator)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// ServingDatasetVersionDataByKey implements data.DataUsecase.
func (d *dataUsecase) ServingDatasetVersionDataByKey(ctx context.Context, datasetID string, version string, key string, paginator *helperModel.Paginator, viewName string, requestedColumns []string) ([]map[string]interface{}, error) {
	if err := d.validateDatasetID(datasetID); err != nil {
		return nil, err
	}

	// 1. Get dataset version to get policies
	datasetVersion, err := d.datasetRepo.FetchDatasetVersionByID(ctx, datasetID, version)
	if err != nil {
		return nil, err
	}
	if datasetVersion == nil {
		return nil, errs.NewNotFoundError("dataset version not found")
	}

	// 2. Use runtime policy and prepare for key-based filtering
	selectedPolicy := &datasetVersion.Policies.Runtime

	// 3. Filter policy by view and columns
	filteredPolicy, err := d.filterRuntimePolicyByViewAndColumns(*selectedPolicy, viewName, requestedColumns, datasetVersion.Policies.Views)
	if err != nil {
		return nil, err
	}

	// 4. Build SQL from filtered policy
	query, args, err := d.dataRepo.BuildRuntimeSQL(ctx, datasetVersion.SourceID, filteredPolicy)
	if err != nil {
		return nil, err
	}

	// 5. Check if KeyField is configured in the runtime policy
	if filteredPolicy.KeyField == "" {
		return nil, errs.NewConflictError("KeyField is not configured in runtime policy for key-based access")
	}

	keyField := filteredPolicy.KeyField

	// 6. Execute query with key-based WHERE condition
	whereConditions := map[string]interface{}{
		keyField: key,
	}

	results, err := d.dataRepo.ExecuteQueryByKey(ctx, datasetVersion.SourceID, query, args, whereConditions)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// CreateDatasetVersionData implements data.DataUsecase.
func (d *dataUsecase) CreateDatasetVersionData(ctx context.Context, datasetID string, version string, data map[string]interface{}) (map[string]interface{}, error) {
	if err := d.validateDatasetID(datasetID); err != nil {
		return nil, err
	}

	// 1. Get dataset version to get policies
	datasetVersion, err := d.datasetRepo.FetchDatasetVersionByID(ctx, datasetID, version)
	if err != nil {
		return nil, err
	}
	if datasetVersion == nil {
		return nil, errs.NewNotFoundError("dataset version not found")
	}

	// 2. Check if WritePolicy exists
	if datasetVersion.Policies.Write == nil {
		return nil, errs.NewConflictError("WritePolicy is not configured for data creation")
	}

	// 3. Validate data is not null or empty
	if len(data) == 0 {
		return nil, errs.NewConflictError("data cannot be null or empty")
	}

	// 4. Build SQL from write policy
	query, args, err := d.dataRepo.BuildCreateSQL(ctx, datasetVersion.SourceID, datasetVersion.Policies.Write, data)
	if err != nil {
		return nil, err
	}

	// 5. Execute insert
	result, err := d.dataRepo.ExecuteInsert(ctx, datasetVersion.SourceID, query, args)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateDatasetVersionDataByKey implements data.DataUsecase.
func (d *dataUsecase) UpdateDatasetVersionDataByKey(ctx context.Context, datasetID string, version string, key string, data map[string]interface{}) (map[string]interface{}, error) {
	if err := d.validateDatasetID(datasetID); err != nil {
		return nil, err
	}

	// 1. Get dataset version to get policies
	datasetVersion, err := d.datasetRepo.FetchDatasetVersionByID(ctx, datasetID, version)
	if err != nil {
		return nil, err
	}
	if datasetVersion == nil {
		return nil, errs.NewNotFoundError("dataset version not found")
	}

	// 2. Check if WritePolicy exists
	if datasetVersion.Policies.Write == nil {
		return nil, errs.NewConflictError("WritePolicy is not configured for data update")
	}

	// 3. Check if KeyField is configured in the runtime policy
	if datasetVersion.Policies.Runtime.KeyField == "" {
		return nil, errs.NewConflictError("KeyField is not configured in runtime policy for key-based update")
	}

	// 4. Validate data is not null or empty
	if len(data) == 0 {
		return nil, errs.NewConflictError("data cannot be null or empty")
	}

	// 5. Build WHERE condition using KeyField
	whereConditions := map[string]interface{}{
		datasetVersion.Policies.Runtime.KeyField: key,
	}

	// 6. Build SQL from write policy
	query, args, err := d.dataRepo.BuildUpdateSQL(ctx, datasetVersion.SourceID, datasetVersion.Policies.Write, data, whereConditions)
	if err != nil {
		return nil, err
	}

	// 7. Execute update
	rowsAffected, err := d.dataRepo.ExecuteUpdate(ctx, datasetVersion.SourceID, query, args)
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, errs.NewNotFoundError("no record found with the provided key")
	}

	// 8. Return the updated record by key
	updatedData, err := d.ServingDatasetVersionDataByKey(ctx, datasetID, version, key, &helperModel.Paginator{Page: 1, PerPage: 1}, "", nil)
	if err != nil {
		return nil, err
	}

	if len(updatedData) == 0 {
		return nil, errs.NewNotFoundError("updated record not found")
	}

	return updatedData[0], nil
}

// DeleteDatasetVersionDataByKey implements data.DataUsecase.
func (d *dataUsecase) DeleteDatasetVersionDataByKey(ctx context.Context, datasetID string, version string, key string) error {
	if err := d.validateDatasetID(datasetID); err != nil {
		return err
	}

	// 1. Get dataset version to get policies
	datasetVersion, err := d.datasetRepo.FetchDatasetVersionByID(ctx, datasetID, version)
	if err != nil {
		return err
	}
	if datasetVersion == nil {
		return errs.NewNotFoundError("dataset version not found")
	}

	// 2. Check if DeletePolicy exists
	if datasetVersion.Policies.Delete == nil {
		return errs.NewConflictError("DeletePolicy is not configured for data deletion")
	}

	// 3. Check if KeyField is configured in the runtime policy
	if datasetVersion.Policies.Runtime.KeyField == "" {
		return errs.NewConflictError("KeyField is not configured in runtime policy for key-based deletion")
	}

	// 5. Build WHERE condition using KeyField
	whereConditions := map[string]interface{}{
		datasetVersion.Policies.Runtime.KeyField: key,
	}

	// 6. Build SQL from delete policy
	query, args, err := d.dataRepo.BuildDeleteSQL(ctx, datasetVersion.SourceID, datasetVersion.Policies.Delete, whereConditions)
	if err != nil {
		return err
	}

	// 7. Execute delete
	rowsAffected, err := d.dataRepo.ExecuteDelete(ctx, datasetVersion.SourceID, query, args)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errs.NewNotFoundError("no record found with the provided key")
	}

	return nil
}

func NewDataUsecase(dataRepo data.PsqlDataRepository, datasetRepo data.PsqlDatasetRepository, redisRepo data.RedisRepository, cryptoSecret string) data.DataUsecase {
	return &dataUsecase{
		dataRepo:     dataRepo,
		datasetRepo:  datasetRepo,
		redisRepo:    redisRepo,
		cryptoSecret: cryptoSecret,
	}
}
