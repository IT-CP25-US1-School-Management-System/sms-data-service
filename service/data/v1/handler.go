package data

import (
	"github.com/labstack/echo/v4"
)

type DataHandler interface {
	FetchSourceList(c echo.Context) error
	FetchSchemasList(c echo.Context) error
	FetchTablesList(c echo.Context) error
	FetchColumnsList(c echo.Context) error
	InsertSource(c echo.Context) error
	UpdateSource(c echo.Context) error
	ActivateSourceByID(c echo.Context) error
	DeactivateSourceByID(c echo.Context) error
	DeleteSourceByID(c echo.Context) error
	// Dataset
	FetchDatasetList(c echo.Context) error
	FetchDatasetByID(c echo.Context) error
	UpsertDataset(c echo.Context) error
	DeleteDatasetByID(c echo.Context) error

	// Dataset Version
	FetchDatasetVersionByID(c echo.Context) error
	FetchDatasetVersionsList(c echo.Context) error
	UpsertDatasetVersion(c echo.Context) error
	DeleteDatasetVersionByID(c echo.Context) error

	// Serving
	ServingDatasetVersionData(c echo.Context) error
	ServingDatasetVersionDataByKey(c echo.Context) error

	// Data Modification (requires write policies)
	CreateDatasetVersionData(c echo.Context) error
	UpdateDatasetVersionDataByKey(c echo.Context) error
	DeleteDatasetVersionDataByKey(c echo.Context) error
}
