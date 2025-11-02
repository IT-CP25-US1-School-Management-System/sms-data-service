package data

import (
	"github.com/labstack/echo/v4"
)

type DataHandler interface {
	FetchSourceList(c echo.Context) error
	FetchSchemasList(c echo.Context) error
	FetchTablesList(c echo.Context) error
	FetchColumnsList(c echo.Context) error

	// Dataset
	FetchDatasetList(c echo.Context) error
	FetchDatasetByID(c echo.Context) error
	UpsertDataset(c echo.Context) error
	DeleteDatasetByID(c echo.Context) error
}
