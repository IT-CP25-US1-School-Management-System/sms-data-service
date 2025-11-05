package handler

import (
	"fmt"
	"net/http"
	"strconv"

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

func NewDataHandler(dataUs data.DataUsecase) data.DataHandler {
	return &dataHandler{
		dataUs: dataUs,
	}
}
