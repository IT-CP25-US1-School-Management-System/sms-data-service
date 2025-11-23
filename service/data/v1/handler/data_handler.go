package handler

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/GodeFvt/go-backend/helper"
	helperModel "github.com/GodeFvt/go-backend/helper/models"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/constants"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/errs"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/dto"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/entity"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/filter"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/data/v1"
	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
)

type dataHandler struct {
	dataUs data.DataUsecase
}

func NewDataHandler(dataUs data.DataUsecase) data.DataHandler {
	return &dataHandler{
		dataUs: dataUs,
	}
}

// FetchSourceByID implements data.DataHandler.
func (d *dataHandler) FetchSourceByID(c echo.Context) error {
	ctx := c.Request().Context()
	sourceIDParam := c.Param("id")
	sourceIDUUID, err := uuid.FromString(sourceIDParam)
	if err != nil {
		return errs.NewBadRequestError(constants.ERR_INVALID_SOURCE_ID)
	}
	source, err := d.dataUs.FetchSourceByID(ctx, &sourceIDUUID)
	if err != nil {
		return err
	}
	if source == nil {
		return errs.NewNotFoundError(constants.ERR_SOURCE_NOT_FOUND)
	}
	sourceDTO, err := helperModel.ConvertStruct[entity.Sources, dto.SourcesResponseDTO](*source)
	if err != nil {
		return err
	}

	res := map[string]interface{}{
		"data": sourceDTO,
	}
	return c.JSON(http.StatusOK, res)
}

// InsertDatasetVersion implements data.DataHandler.
func (d *dataHandler) InsertDatasetVersion(c echo.Context) error {
	ctx := c.Request().Context()
	datasetID := c.Param("id")

	if datasetID == "" {
		return errs.NewBadRequestError(constants.ERR_DATASET_ID_IS_REQUIRED)
	}

	var datasetVersion dto.InsertDatasetVersionDTO
	if err := c.Bind(&datasetVersion); err != nil {
		return errs.ErrBadRequest(err)
	}
	if err := c.Validate(&datasetVersion); err != nil {
		return errs.ErrBadRequest(err)
	}
	datasetVersionEntity, err := datasetVersion.InsertDatasetVersionDTOToEntity()
	if err != nil {
		return errs.ErrBadRequest(err)
	}

	if err := d.dataUs.InsertDatasetVersion(ctx, datasetVersionEntity, datasetID); err != nil {
		return err
	}

	res := map[string]interface{}{
		"message": "success",
	}
	return c.JSON(http.StatusOK, res)
}

// UpdateDatasetVersion implements data.DataHandler.
func (d *dataHandler) UpdateDatasetVersion(c echo.Context) error {
	ctx := c.Request().Context()
	datasetID := c.Param("id")
	datasetVersion := c.Param("version")

	if datasetID == "" {
		return errs.NewBadRequestError(constants.ERR_DATASET_ID_IS_REQUIRED)
	}

	var datasetVersionDTO dto.UpdateDatasetVersionDTO
	if err := c.Bind(&datasetVersionDTO); err != nil {
		return errs.ErrBadRequest(err)
	}
	if err := c.Validate(&datasetVersionDTO); err != nil {
		return errs.ErrBadRequest(err)
	}
	datasetVersionEntity, err := datasetVersionDTO.UpdateDatasetVersionDTOToEntity()
	if err != nil {
		return errs.ErrBadRequest(err)
	}

	if err := d.dataUs.UpdateDatasetVersion(ctx, datasetVersionEntity, datasetID, datasetVersion); err != nil {
		return err
	}

	res := map[string]interface{}{
		"message": "success",
	}
	return c.JSON(http.StatusOK, res)
}

// DeleteSourceByID implements data.DataHandler.
func (d *dataHandler) DeleteSourceByID(c echo.Context) error {
	ctx := c.Request().Context()
	sourceIDParam := c.Param("id")
	sourceIDUUID, err := uuid.FromString(sourceIDParam)
	if err != nil {
		return errs.NewBadRequestError(constants.ERR_INVALID_SOURCE_ID)
	}
	err = d.dataUs.DeleteSourceByID(ctx, &sourceIDUUID)
	if err != nil {
		return err
	}

	res := map[string]interface{}{
		"message": "success",
	}
	return c.JSON(http.StatusOK, res)
}

// ActivateSourceByID implements data.DataHandler.
func (d *dataHandler) ActivateSourceByID(c echo.Context) error {
	ctx := c.Request().Context()
	sourceIDParam := c.Param("id")
	sourceIDUUID, err := uuid.FromString(sourceIDParam)
	if err != nil {
		return errs.NewBadRequestError(constants.ERR_INVALID_SOURCE_ID)
	}
	err = d.dataUs.ActivateSourceByID(ctx, &sourceIDUUID)
	if err != nil {
		return err
	}

	res := map[string]interface{}{
		"message": "success",
	}
	return c.JSON(http.StatusOK, res)
}

// DeactivateSourceByID implements data.DataHandler.
func (d *dataHandler) DeactivateSourceByID(c echo.Context) error {
	ctx := c.Request().Context()
	sourceIDParam := c.Param("id")
	sourceIDUUID, err := uuid.FromString(sourceIDParam)
	if err != nil {
		return errs.NewBadRequestError(constants.ERR_INVALID_SOURCE_ID)
	}
	err = d.dataUs.DeactivateSourceByID(ctx, &sourceIDUUID)
	if err != nil {
		return err
	}

	res := map[string]interface{}{
		"message": "success",
	}
	return c.JSON(http.StatusOK, res)
}

// InsertSource implements data.DataHandler.
func (d *dataHandler) InsertSource(c echo.Context) error {
	ctx := c.Request().Context()
	var sourceDTO dto.SourcesDTO
	if err := c.Bind(&sourceDTO); err != nil {
		return errs.ErrBadRequest(err)
	}
	if err := c.Validate(&sourceDTO); err != nil {
		return errs.ErrBadRequest(err)
	}

	sourceEntity := sourceDTO.SourcesDTOToEntity()
	sourceEntity.GenUUID()
	if err := d.dataUs.InsertSource(ctx, sourceEntity); err != nil {
		return err
	}

	res := map[string]interface{}{
		"message": "success",
		"id":      sourceEntity.ID,
	}
	return c.JSON(http.StatusOK, res)
}

// UpdateSource implements data.DataHandler.
func (d *dataHandler) UpdateSource(c echo.Context) error {
	ctx := c.Request().Context()
	sourceIDParam := c.Param("id")
	sourceIDUUID, err := uuid.FromString(sourceIDParam)
	if err != nil {
		return errs.NewBadRequestError(constants.ERR_INVALID_SOURCE_ID)
	}

	var updateSourceDTO dto.UpdateSourcesDTO
	if err := c.Bind(&updateSourceDTO); err != nil {
		return errs.ErrBadRequest(err)
	}
	if err := c.Validate(&updateSourceDTO); err != nil {
		return errs.ErrBadRequest(err)
	}

	if err := d.dataUs.UpdateSource(ctx, &sourceIDUUID, &updateSourceDTO); err != nil {
		return err
	}

	res := map[string]interface{}{
		"message": "success",
		"id":      sourceIDUUID,
	}
	return c.JSON(http.StatusOK, res)
}

// DeleteDatasetByID implements data.DataHandler.
func (d *dataHandler) DeleteDatasetByID(c echo.Context) error {
	ctx := c.Request().Context()
	datasetID := c.Param("id")
	if datasetID == "" {
		return errs.NewBadRequestError(constants.ERR_DATASET_ID_IS_REQUIRED)
	}
	err := d.dataUs.DeleteDatasetByID(ctx, datasetID)
	if err != nil {
		return err
	}

	res := map[string]interface{}{
		"message": "success",
	}
	return c.JSON(http.StatusOK, res)
}

// UpsertDataset implements data.DataHandler.
func (d *dataHandler) UpsertDataset(c echo.Context) error {
	ctx := c.Request().Context()
	var datasetDTO dto.UpsertDatasetsDTO
	if err := c.Bind(&datasetDTO); err != nil {
		return errs.ErrBadRequest(err)
	}
	if err := c.Validate(&datasetDTO); err != nil {
		return errs.ErrBadRequest(err)
	}

	datasetEntity := datasetDTO.UpsertDatasetsDTOToEntity()
	if err := d.dataUs.UpsertDataset(ctx, datasetEntity); err != nil {
		return err
	}

	res := map[string]interface{}{
		"message": "success",
	}
	return c.JSON(http.StatusOK, res)

}

// FetchDatasetByID implements data.DataHandler.
func (d *dataHandler) FetchDatasetByID(c echo.Context) error {
	ctx := c.Request().Context()
	datasetID := c.Param("id")
	if datasetID == "" {
		return errs.NewBadRequestError(constants.ERR_DATASET_ID_IS_REQUIRED)
	}
	dataset, err := d.dataUs.FetchDatasetByID(ctx, datasetID)
	if err != nil {
		return err
	}
	if dataset == nil {
		return errs.NewNotFoundError(constants.ERR_DATASET_NOT_FOUND)
	}
	res := map[string]interface{}{
		"data": dataset,
	}
	return c.JSON(http.StatusOK, res)
}

// FetchColumnsList implements data.DataHandler.
func (d *dataHandler) FetchColumnsList(c echo.Context) error {
	ctx := c.Request().Context()
	var filter filter.ColumnsFilter
	if err := c.Bind(&filter); err != nil {
		return errs.ErrBadRequest(err)
	}
	if err := c.Validate(&filter); err != nil {
		return errs.ErrBadRequest(err)
	}
	pageStr := c.QueryParam("page")
	perPageStr := c.QueryParam("per_page")
	if pageStr == "0" {
		return errs.NewBadRequestError("page parameter must be greater than 0")
	}
	if perPageStr == "0" {
		return errs.NewBadRequestError("per_page parameter must be greater than 0")
	}
	paginator := helperModel.NewPaginator()
	if filter.Page > 0 && filter.PerPage > 0 {
		paginator.Page = filter.Page
		paginator.PerPage = filter.PerPage
	}

	columns, err := d.dataUs.FetchColumnsList(ctx, &filter, &paginator)
	if err != nil {
		return err
	}
	if len(columns) == 0 {
		columns = []*entity.Columns{}
	}
	res := map[string]interface{}{
		"data":        columns,
		"page":        paginator.Page,
		"per_page":    paginator.PerPage,
		"total_pages": paginator.TotalPages,
		"total_rows":  paginator.TotalEntrySizes,
	}
	return c.JSON(http.StatusOK, res)
}

// FetchTablesList implements data.DataHandler.
func (d *dataHandler) FetchTablesList(c echo.Context) error {
	ctx := c.Request().Context()
	var filter filter.TablesFilter
	if err := c.Bind(&filter); err != nil {
		return errs.ErrBadRequest(err)
	}
	if err := c.Validate(&filter); err != nil {
		return errs.ErrBadRequest(err)
	}
	pageStr := c.QueryParam("page")
	perPageStr := c.QueryParam("per_page")
	if pageStr == "0" {
		return errs.NewBadRequestError("page parameter must be greater than 0")
	}
	if perPageStr == "0" {
		return errs.NewBadRequestError("per_page parameter must be greater than 0")
	}
	paginator := helperModel.NewPaginator()
	if filter.Page > 0 && filter.PerPage > 0 {
		paginator.Page = filter.Page
		paginator.PerPage = filter.PerPage
	}

	tables, err := d.dataUs.FetchTablesList(ctx, &filter, &paginator)
	if err != nil {
		return err
	}
	if len(tables) == 0 {
		tables = []*entity.Tables{}
	}
	res := map[string]interface{}{
		"data":        tables,
		"page":        paginator.Page,
		"per_page":    paginator.PerPage,
		"total_pages": paginator.TotalPages,
		"total_rows":  paginator.TotalEntrySizes,
	}
	return c.JSON(http.StatusOK, res)
}

// FetchDatasetList implements data.DataHandler.
func (d *dataHandler) FetchDatasetList(c echo.Context) error {
	ctx := c.Request().Context()
	var filter filter.DatasetsFilter
	paginator := helperModel.NewPaginator()

	if err := c.Bind(&filter); err != nil {
		return errs.ErrBadRequest(err)
	}
	if err := c.Validate(&filter); err != nil {
		return errs.ErrBadRequest(err)
	}
	pageStr := c.QueryParam("page")
	perPageStr := c.QueryParam("per_page")
	if pageStr == "0" {
		return errs.NewBadRequestError("page parameter must be greater than 0")
	}
	if perPageStr == "0" {
		return errs.NewBadRequestError("per_page parameter must be greater than 0")
	}
	if filter.Page > 0 && filter.PerPage > 0 {
		paginator.Page = filter.Page
		paginator.PerPage = filter.PerPage
	}

	datasets, err := d.dataUs.FetchDatasetList(ctx, &filter, &paginator)
	if err != nil {
		return err
	}
	if datasets == nil {
		datasets = []*entity.Datasets{}
	}
	res := map[string]interface{}{
		"data":        datasets,
		"page":        paginator.Page,
		"per_page":    paginator.PerPage,
		"total_pages": paginator.TotalPages,
		"total_rows":  paginator.TotalEntrySizes,
	}
	return c.JSON(http.StatusOK, res)
}

// FetchSchemasList implements data.DataHandler.
func (d *dataHandler) FetchSchemasList(c echo.Context) error {
	ctx := c.Request().Context()
	var filter filter.SchemasFilter
	if err := c.Bind(&filter); err != nil {
		return errs.ErrBadRequest(err)
	}
	if err := c.Validate(&filter); err != nil {
		return errs.ErrBadRequest(err)
	}
	pageStr := c.QueryParam("page")
	perPageStr := c.QueryParam("per_page")
	if pageStr == "0" {
		return errs.NewBadRequestError("page parameter must be greater than 0")
	}
	if perPageStr == "0" {
		return errs.NewBadRequestError("per_page parameter must be greater than 0")
	}
	paginator := helperModel.NewPaginator()
	if filter.Page > 0 && filter.PerPage > 0 {
		paginator.Page = filter.Page
		paginator.PerPage = filter.PerPage
	}

	schemas, err := d.dataUs.FetchSchemasList(ctx, &filter, &paginator)
	if err != nil {
		return err
	}
	if len(schemas) == 0 {
		schemas = []*entity.Schemas{}
	}
	res := map[string]interface{}{
		"data":        schemas,
		"page":        paginator.Page,
		"per_page":    paginator.PerPage,
		"total_pages": paginator.TotalPages,
		"total_rows":  paginator.TotalEntrySizes,
	}
	return c.JSON(http.StatusOK, res)
}

// FetchSourceList implements data.DataHandler.
func (d *dataHandler) FetchSourceList(c echo.Context) error {
	ctx := c.Request().Context()
	paginator := helperModel.NewPaginator()
	page := c.QueryParam("page")
	perPage := c.QueryParam("per_page")
	if page != "" && perPage != "" {
		p, err := strconv.Atoi(page)
		if err != nil || p < 1 {
			return errs.ErrBadRequest(fmt.Errorf("invalid page"))
		}
		pp, err := strconv.Atoi(perPage)
		if err != nil || pp < 1 || pp > 100 {
			return errs.ErrBadRequest(fmt.Errorf("invalid per_page"))
		}
		paginator.Page = p
		paginator.PerPage = pp
	}

	sources, err := d.dataUs.FetchSourceList(ctx, &paginator)
	if err != nil {
		return err
	}

	res := map[string]interface{}{
		"data":        dto.SourcesEntityToSourcesResponseDTO(sources),
		"page":        paginator.Page,
		"per_page":    paginator.PerPage,
		"total_pages": paginator.TotalPages,
		"total_rows":  paginator.TotalEntrySizes,
	}
	return c.JSON(http.StatusOK, res)
}

// FetchDatasetVersionByID implements data.DataHandler.
func (d *dataHandler) FetchDatasetVersionByID(c echo.Context) error {
	ctx := c.Request().Context()
	datasetID := c.Param("id")
	version := c.Param("version")

	if datasetID == "" {
		return errs.NewBadRequestError(constants.ERR_DATASET_ID_IS_REQUIRED)
	}
	if version == "" {
		return errs.NewBadRequestError("version is required")
	}

	datasetVersion, err := d.dataUs.FetchDatasetVersionByID(ctx, datasetID, version)
	if err != nil {
		return err
	}
	if datasetVersion == nil {
		return errs.NewNotFoundError("dataset version not found")
	}

	res := map[string]interface{}{
		"data": datasetVersion,
	}
	return c.JSON(http.StatusOK, res)
}

// FetchDatasetVersionsList implements data.DataHandler.
func (d *dataHandler) FetchDatasetVersionsList(c echo.Context) error {
	ctx := c.Request().Context()
	datasetID := c.Param("id")
	if datasetID == "" {
		return errs.NewBadRequestError(constants.ERR_DATASET_ID_IS_REQUIRED)
	}
	var filter filter.DatasetVersionsFilter
	if err := c.Bind(&filter); err != nil {
		return errs.ErrBadRequest(err)
	}
	if err := c.Validate(&filter); err != nil {
		return errs.ErrBadRequest(err)
	}
	pageStr := c.QueryParam("page")
	perPageStr := c.QueryParam("per_page")
	if pageStr == "0" {
		return errs.NewBadRequestError("page parameter must be greater than 0")
	}
	if perPageStr == "0" {
		return errs.NewBadRequestError("per_page parameter must be greater than 0")
	}
	paginator := helperModel.NewPaginator()
	if filter.Page > 0 && filter.PerPage > 0 {
		paginator.Page = filter.Page
		paginator.PerPage = filter.PerPage
	}

	versions, err := d.dataUs.FetchDatasetVersionsList(ctx, datasetID, &filter, &paginator)
	if err != nil {
		return err
	}
	if versions == nil {
		versions = []*entity.DatasetVersion{}
	}

	res := map[string]interface{}{
		"data":        versions,
		"page":        paginator.Page,
		"per_page":    paginator.PerPage,
		"total_pages": paginator.TotalPages,
		"total_rows":  paginator.TotalEntrySizes,
	}
	return c.JSON(http.StatusOK, res)
}

// UpdateDatasetVersionStatus implements data.DataHandler.
func (d *dataHandler) UpdateDatasetVersionStatus(c echo.Context) error {
	ctx := c.Request().Context()
	datasetID := c.Param("id")
	version := c.Param("version")
	var status dto.UpdateDatasetVersionStatusDTO
	if err := c.Bind(&status); err != nil {
		return errs.ErrBadRequest(err)
	}
	if err := c.Validate(&status); err != nil {
		return errs.ErrBadRequest(err)
	}

	if err := d.dataUs.UpdateDatasetVersionStatus(ctx, datasetID, version, status.Status); err != nil {
		return err
	}

	res := map[string]interface{}{
		"message": "success",
	}
	return c.JSON(http.StatusOK, res)
}

// ServingDatasetVersionData implements data.DataHandler.
func (d *dataHandler) ServingDatasetVersionData(c echo.Context) error {
	ctx := c.Request().Context()
	datasetID := c.Param("id")
	version := c.Param("version")
	if datasetID == "" {
		return errs.NewBadRequestError("dataset parameter is required")
	}
	if version == "" {
		return errs.NewBadRequestError("version parameter is required")
	}

	var servingFilter filter.ServingDataFilter
	if err := c.Bind(&servingFilter); err != nil {
		return errs.ErrBadRequest(err)
	}
	if err := c.Validate(&servingFilter); err != nil {
		return errs.ErrBadRequest(err)
	}
	pageStr := c.QueryParam("page")
	perPageStr := c.QueryParam("per_page")
	if pageStr == "0" {
		return errs.NewBadRequestError("page parameter must be greater than 0")
	}
	if perPageStr == "0" {
		return errs.NewBadRequestError("per_page parameter must be greater than 0")
	}
	paginator := helperModel.NewPaginator()
	if servingFilter.Page > 0 && servingFilter.PerPage > 0 {
		paginator.Page = servingFilter.Page
		paginator.PerPage = servingFilter.PerPage
	}

	// Convert logical operator to uppercase for consistency
	logicalOperator := servingFilter.WhereLogicalOp
	if logicalOperator != "" {
		logicalOperator = strings.ToUpper(logicalOperator)
	}

	// Parse the where filters from the raw string
	filterGroups, err := servingFilter.ParseWhere()
	if err != nil {
		return errs.ErrBadRequest(err)
	}

	// Convert sort order to uppercase for consistency
	sortOrder := servingFilter.SortOrder
	if sortOrder != "" {
		sortOrder = strings.ToUpper(sortOrder)
	}

	data, err := d.dataUs.ServingDatasetVersionData(ctx, datasetID, version, &paginator, servingFilter.View, filterGroups, logicalOperator, servingFilter.SortBy, sortOrder)
	if err != nil {
		return err
	}
	if data == nil {
		data = []map[string]interface{}{}
	}

	res := map[string]interface{}{
		"data":        data,
		"page":        paginator.Page,
		"per_page":    paginator.PerPage,
		"total_pages": paginator.TotalPages,
		"total_rows":  paginator.TotalEntrySizes,
	}
	return c.JSON(http.StatusOK, res)
}

// ServingDatasetVersionDataByKey implements data.DataHandler.
func (d *dataHandler) ServingDatasetVersionDataByKey(c echo.Context) error {
	ctx := c.Request().Context()
	datasetID := c.Param("id")
	version := c.Param("version")
	key := c.Param("key")

	if datasetID == "" {
		return errs.NewBadRequestError("dataset parameter is required")
	}
	if version == "" {
		return errs.NewBadRequestError("version parameter is required")
	}
	if key == "" {
		return errs.NewBadRequestError("key parameter is required")
	}

	viewName := c.QueryParam("view")

	data, err := d.dataUs.ServingDatasetVersionDataByKey(ctx, datasetID, version, key, viewName)
	if err != nil {
		return err
	}
	if data == nil {
		return errs.NewNotFoundError("data not found for the given key")
	}

	res := map[string]interface{}{
		"data": data,
	}
	return c.JSON(http.StatusOK, res)
}

// CreateDatasetVersionData implements data.DataHandler.
func (d *dataHandler) CreateDatasetVersionData(c echo.Context) error {
	ctx := c.Request().Context()
	datasetID := c.Param("id")
	version := c.Param("version")

	if datasetID == "" {
		return errs.NewBadRequestError("dataset parameter is required")
	}
	if version == "" {
		return errs.NewBadRequestError("version parameter is required")
	}

	var data map[string]interface{}
	if err := c.Bind(&data); err != nil {
		return errs.ErrBadRequest(err)
	}

	result, err := d.dataUs.CreateDatasetVersionData(ctx, datasetID, version, data)
	if err != nil {
		return err
	}

	res := map[string]interface{}{
		"message": "success",
		"data":    result,
	}
	return c.JSON(http.StatusCreated, res)
}

// UpdateDatasetVersionDataByKey implements data.DataHandler.
func (d *dataHandler) UpdateDatasetVersionDataByKey(c echo.Context) error {
	ctx := c.Request().Context()
	datasetID := c.Param("id")
	version := c.Param("version")
	key := c.Param("key")

	if datasetID == "" {
		return errs.NewBadRequestError("dataset parameter is required")
	}
	if version == "" {
		return errs.NewBadRequestError("version parameter is required")
	}
	if key == "" {
		return errs.NewBadRequestError("key parameter is required")
	}

	var data map[string]interface{}
	if err := c.Bind(&data); err != nil {
		return errs.ErrBadRequest(err)
	}

	result, err := d.dataUs.UpdateDatasetVersionDataByKey(ctx, datasetID, version, key, data)
	if err != nil {
		return err
	}

	res := map[string]interface{}{
		"message": "success",
		"data":    result,
	}
	return c.JSON(http.StatusOK, res)
}

// DeleteDatasetVersionDataByKey implements data.DataHandler.
func (d *dataHandler) DeleteDatasetVersionDataByKey(c echo.Context) error {
	ctx := c.Request().Context()
	datasetID := c.Param("id")
	version := c.Param("version")
	key := c.Param("key")

	if datasetID == "" {
		return errs.NewBadRequestError("dataset parameter is required")
	}
	if version == "" {
		return errs.NewBadRequestError("version parameter is required")
	}
	if key == "" {
		return errs.NewBadRequestError("key parameter is required")
	}

	err := d.dataUs.DeleteDatasetVersionDataByKey(ctx, datasetID, version, key)
	if err != nil {
		return err
	}

	res := map[string]interface{}{
		"message": "success",
	}
	return c.JSON(http.StatusOK, res)
}

// FetchTableData implements data.DataHandler.
func (d *dataHandler) FetchTableData(c echo.Context) error {
	ctx := c.Request().Context()
	sourceIDParam := c.Param("id")
	schemaName := c.Param("schema")
	tableName := c.Param("table")

	sourceIDUUID, err := uuid.FromString(sourceIDParam)
	if err != nil {
		return errs.NewBadRequestError(constants.ERR_INVALID_SOURCE_ID)
	}

	var servingFilter filter.ServingDataFilter
	if err := c.Bind(&servingFilter); err != nil {
		return errs.ErrBadRequest(err)
	}
	if err := c.Validate(&servingFilter); err != nil {
		return errs.ErrBadRequest(err)
	}
	pageStr := c.QueryParam("page")
	perPageStr := c.QueryParam("per_page")
	if pageStr == "0" {
		return errs.NewBadRequestError("page parameter must be greater than 0")
	}
	if perPageStr == "0" {
		return errs.NewBadRequestError("per_page parameter must be greater than 0")
	}
	paginator := helperModel.NewPaginator()
	if servingFilter.Page > 0 && servingFilter.PerPage > 0 {
		paginator.Page = servingFilter.Page
		paginator.PerPage = servingFilter.PerPage
	}

	// Convert logical operator to uppercase
	logicalOperator := servingFilter.WhereLogicalOp
	if logicalOperator != "" {
		logicalOperator = strings.ToUpper(logicalOperator)
	}

	// Parse filters
	filterGroups, err := servingFilter.ParseWhere()
	if err != nil {
		return errs.ErrBadRequest(err)
	}

	// Convert sort order to uppercase
	sortOrder := servingFilter.SortOrder
	if sortOrder != "" {
		sortOrder = strings.ToUpper(sortOrder)
	}

	data, err := d.dataUs.FetchTableData(ctx, &sourceIDUUID, schemaName, tableName, filterGroups, logicalOperator, &paginator, servingFilter.SortBy, sortOrder)
	if err != nil {
		return err
	}
	if data == nil {
		data = []map[string]interface{}{}
	}

	res := map[string]interface{}{
		"data":        data,
		"page":        paginator.Page,
		"per_page":    paginator.PerPage,
		"total_pages": paginator.TotalPages,
		"total_rows":  paginator.TotalEntrySizes,
	}
	return c.JSON(http.StatusOK, res)
}

// FetchTableDataByKey implements data.DataHandler.
func (d *dataHandler) FetchTableDataByKey(c echo.Context) error {
	ctx := c.Request().Context()
	sourceIDParam := c.Param("id")
	schemaName := c.Param("schema")
	tableName := c.Param("table")
	key := c.Param("key")

	sourceIDUUID, err := uuid.FromString(sourceIDParam)
	if err != nil {
		return errs.NewBadRequestError(constants.ERR_INVALID_SOURCE_ID)
	}

	if schemaName == "" {
		return errs.NewBadRequestError("schema parameter is required")
	}
	if tableName == "" {
		return errs.NewBadRequestError("table parameter is required")
	}
	if key == "" {
		return errs.NewBadRequestError("key parameter is required")
	}

	keyField := c.QueryParam("key_field")
	if keyField == "" {
		keyField = "id" // default key field
	}

	data, err := d.dataUs.FetchTableDataByKey(ctx, &sourceIDUUID, schemaName, tableName, keyField, key)
	if err != nil {
		return err
	}
	if data == nil {
		return errs.NewNotFoundError("data not found for the given key")
	}

	res := map[string]interface{}{
		"data": data,
	}
	return c.JSON(http.StatusOK, res)
}

// CreateTableData implements data.DataHandler.
func (d *dataHandler) CreateTableData(c echo.Context) error {
	ctx := c.Request().Context()
	sourceIDParam := c.Param("id")
	schemaName := c.Param("schema")
	tableName := c.Param("table")

	sourceIDUUID, err := uuid.FromString(sourceIDParam)
	if err != nil {
		return errs.NewBadRequestError(constants.ERR_INVALID_SOURCE_ID)
	}

	if schemaName == "" {
		return errs.NewBadRequestError("schema parameter is required")
	}
	if tableName == "" {
		return errs.NewBadRequestError("table parameter is required")
	}

	var data map[string]interface{}
	if err := c.Bind(&data); err != nil {
		return errs.ErrBadRequest(err)
	}

	// Remove reserved fields that come from path parameters
	delete(data, "schema")
	delete(data, "table")
	delete(data, "source_id")

	result, err := d.dataUs.CreateTableData(ctx, &sourceIDUUID, schemaName, tableName, data)
	if err != nil {
		return err
	}

	res := map[string]interface{}{
		"message": "success",
		"data":    result,
	}
	return c.JSON(http.StatusCreated, res)
}

// UpdateTableData implements data.DataHandler.
func (d *dataHandler) UpdateTableData(c echo.Context) error {
	ctx := c.Request().Context()
	sourceIDParam := c.Param("id")
	schemaName := c.Param("schema")
	tableName := c.Param("table")
	key := c.Param("key")

	sourceIDUUID, err := uuid.FromString(sourceIDParam)
	if err != nil {
		return errs.NewBadRequestError(constants.ERR_INVALID_SOURCE_ID)
	}

	if schemaName == "" {
		return errs.NewBadRequestError("schema parameter is required")
	}
	if tableName == "" {
		return errs.NewBadRequestError("table parameter is required")
	}
	if key == "" {
		return errs.NewBadRequestError("key parameter is required")
	}

	keyField := c.QueryParam("key_field")
	if keyField == "" {
		keyField = "id" // default key field
	}

	var data map[string]interface{}
	if err := c.Bind(&data); err != nil {
		return errs.ErrBadRequest(err)
	}

	delete(data, "schema")
	delete(data, "table")
	delete(data, "source_id")
	delete(data, "key")

	result, err := d.dataUs.UpdateTableData(ctx, &sourceIDUUID, schemaName, tableName, keyField, key, data)
	if err != nil {
		return err
	}

	res := map[string]interface{}{
		"message": "success",
		"data":    result,
	}
	return c.JSON(http.StatusOK, res)
}

// DeleteTableData implements data.DataHandler.
func (d *dataHandler) DeleteTableData(c echo.Context) error {
	ctx := c.Request().Context()
	sourceIDParam := c.Param("id")
	schemaName := c.Param("schema")
	tableName := c.Param("table")
	key := c.Param("key")

	sourceIDUUID, err := uuid.FromString(sourceIDParam)
	if err != nil {
		return errs.NewBadRequestError(constants.ERR_INVALID_SOURCE_ID)
	}

	if schemaName == "" {
		return errs.NewBadRequestError("schema parameter is required")
	}
	if tableName == "" {
		return errs.NewBadRequestError("table parameter is required")
	}
	if key == "" {
		return errs.NewBadRequestError("key parameter is required")
	}

	keyField := c.QueryParam("key_field")
	if keyField == "" {
		keyField = "id" // default key field
	}

	err = d.dataUs.DeleteTableData(ctx, &sourceIDUUID, schemaName, tableName, keyField, key)
	if err != nil {
		return err
	}

	res := map[string]interface{}{
		"message": "success",
	}
	return c.JSON(http.StatusOK, res)
}

// FetchExportJobByJobId implements data.DataHandler.
func (d *dataHandler) FetchExportJobByJobId(c echo.Context) error {
	ctx := c.Request().Context()
	jobIDParam := c.Param("job_id")

	jobIDUUID, err := uuid.FromString(jobIDParam)
	if err != nil {
		return errs.NewBadRequestError(constants.ERR_INVALID_EXPORT_JOB_ID)
	}
	export_job, err := d.dataUs.FetchExportJobByID(ctx, &jobIDUUID)
	if err != nil {
		return err
	}
	if export_job == nil {
		return errs.NewNotFoundError(constants.ERR_EXPORT_JOB_NOT_FOUND)
	}
	exportJobDTO, err := helperModel.ConvertStruct[entity.ExportJob, dto.ExportJobResponseDTO](*export_job)
	if err != nil {
		return err
	}

	res := map[string]interface{}{
		"data": exportJobDTO,
	}
	return c.JSON(http.StatusOK, res)

}


// InsertExportJob implements data.DataHandler.
func (d *dataHandler) ExportJob(c echo.Context) error {
	ctx := c.Request().Context()
	var exportJobDTO dto.ExportJobDTO
	if err := c.Bind(&exportJobDTO); err != nil {
		return errs.ErrBadRequest(err)
	}
	if err := c.Validate(&exportJobDTO); err != nil {
		return errs.ErrBadRequest(err)
	}

	exportJobEntity := exportJobDTO.ExportJobDTOToEntity()
	exportJobEntity.GenUUID()
	if err := d.dataUs.InsertExportJob(ctx, exportJobEntity); err != nil {
		return err
	}

	res := map[string]interface{}{
		"job_id":      exportJobEntity.JobId,
	}
	return c.JSON(http.StatusOK, res)
}

func (d *dataHandler) UploadReportingTemplate(c echo.Context) error {
	ctx := c.Request().Context()
	if c.Get("file") == nil {
		return errs.NewBadRequestError("file not found")
	}
	if c.Get("params") == nil {
		return errs.NewBadRequestError("params not found")
	}
	file := c.Get("file").([]*multipart.FileHeader)[0]
	paramsRaw := c.Get("params").(map[string]interface{})
	paramsData, ok := paramsRaw["params"].(map[string]interface{})
	if !ok {
		return errs.NewBadRequestError("invalid params structure")
	}
	entityResult, err := helperModel.ConvertStruct[map[string]interface{}, entity.ReportingTemplate](paramsData)
	if err != nil {
		return errs.ErrBadRequest(err)
	}

	reader, err := file.Open()
	if err != nil {
		return err
	}
	buf, contentType, _, err := helper.GetMimeType(reader)
	if err != nil {
		return err
	}
	if contentType != "application/pdf" {
		return errs.NewBadRequestError("invalid file type, only .pdf files are allowed")
	}

	err = d.dataUs.UploadReportingTemplate(ctx, &entityResult, buf.Bytes(), file.Filename)
	if err != nil {
		return err
	}

	res := map[string]interface{}{
		"message": "success",
	}
	return c.JSON(http.StatusOK, res)

}
