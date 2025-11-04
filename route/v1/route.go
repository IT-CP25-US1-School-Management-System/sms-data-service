package route

import (
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/middleware"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/data/v1"
	"github.com/labstack/echo/v4"
)

type Route struct {
	e     *echo.Echo
	middl middleware.GoMiddlewareInf
}

func NewRoute(e *echo.Echo, middl middleware.GoMiddlewareInf) *Route {
	return &Route{e: e, middl: middl}
}

func (r *Route) RegisterDataRoute(handler data.DataHandler) {
	introspectGroup := r.e.Group("/v1/introspect")
	introspectGroup.GET("/sources", handler.FetchSourceList)
	introspectGroup.GET("/schemas", handler.FetchSchemasList)
	introspectGroup.GET("/tables", handler.FetchTablesList)
	introspectGroup.GET("/columns", handler.FetchColumnsList)
	introspectGroup.POST("/sources", handler.InsertSource)
	introspectGroup.PUT("/sources/:id", handler.UpdateSource)
	introspectGroup.PATCH("/sources/:id/activate", handler.ActivateSourceByID, r.middl.ValidateParamId("id"))
	introspectGroup.PATCH("/sources/:id/deactivate", handler.DeactivateSourceByID, r.middl.ValidateParamId("id"))

	// Datasets Route
	datasetsGroup := r.e.Group("/v1/datasets")
	datasetsGroup.GET("", handler.FetchDatasetList)
	datasetsGroup.GET("/:id", handler.FetchDatasetByID)
	datasetsGroup.POST("", handler.UpsertDataset)
	datasetsGroup.DELETE("/:id", handler.DeleteDatasetByID)
}
