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
}
