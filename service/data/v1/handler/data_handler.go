package handler

import (
	"net/http"

	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/errs"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/data/v1"
	"github.com/labstack/echo/v4"
)

type dataHandler struct {
	dataUs data.DataUsecase
}

// FetchSourceList implements data.DataHandler.
func (d *dataHandler) FetchSourceList(c echo.Context) error {
	ctx := c.Request().Context()
	sources, err := d.dataUs.FetchSourceList(ctx)
	if err != nil {
		return errs.ErrInternalServer(err)
	}
	if sources == nil {
		return errs.ErrNoContent()
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
