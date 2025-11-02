package handler

import (
	"net/http"

	helperModel "github.com/GodeFvt/go-backend/helper/models"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/constants"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/errs"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/dto"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/entity"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/filter"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/data/v1"
	"github.com/labstack/echo/v4"
)

type dataHandler struct {
	dataUs data.DataUsecase
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

	columns, err := d.dataUs.FetchColumnsList(ctx, &filter)
	if err != nil {
		return err
	}
	if columns == nil {
		columns = []*entity.Columns{}
	}
	res := map[string]interface{}{
		"data": columns,
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

	tables, err := d.dataUs.FetchTablesList(ctx, &filter)
	if err != nil {
		return err
	}
	if tables == nil {
		tables = []*entity.Tables{}
	}
	res := map[string]interface{}{
		"data": tables,
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

	schemas, err := d.dataUs.FetchSchemasList(ctx, &filter)
	if err != nil {
		return err
	}
	if schemas == nil {
		schemas = []*entity.Schemas{}
	}
	res := map[string]interface{}{
		"data": schemas,
	}
	return c.JSON(http.StatusOK, res)
}

// FetchSourceList implements data.DataHandler.
func (d *dataHandler) FetchSourceList(c echo.Context) error {
	ctx := c.Request().Context()
	sources, err := d.dataUs.FetchSourceList(ctx)
	if err != nil {
		return err
	}
	if sources == nil {
		sources = []*entity.Sources{}
	}
	res := map[string]interface{}{
		"data": sources,
	}
	return c.JSON(http.StatusOK, res)
}

func NewDataHandler(dataUs data.DataUsecase) data.DataHandler {
	return &dataHandler{
		dataUs: dataUs,
	}
}
