package data

import (
	"github.com/labstack/echo/v4"
)

type DataHandler interface {
	FetchSourceList(c echo.Context) error
	FetchSchemasList(c echo.Context) error

	// Dataset
	FetchDatasetList(c echo.Context) error
}
