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
	introspectGroup.GET("/sources/:id", handler.FetchSourceByID)
	introspectGroup.GET("/schemas", handler.FetchSchemasList)
	introspectGroup.GET("/tables", handler.FetchTablesList)
	introspectGroup.GET("/columns", handler.FetchColumnsList)
	introspectGroup.POST("/sources", handler.InsertSource)
	introspectGroup.PUT("/sources/:id", handler.UpdateSource)
	introspectGroup.DELETE("/sources/:id", handler.DeleteSourceByID, r.middl.ValidateParamId("id"))
	introspectGroup.PATCH("/sources/:id/activate", handler.ActivateSourceByID, r.middl.ValidateParamId("id"))
	introspectGroup.PATCH("/sources/:id/deactivate", handler.DeactivateSourceByID, r.middl.ValidateParamId("id"))
	// Table Data CRUD (direct source access)
	tableDataGroup := introspectGroup.Group("/sources/:id/schemas/:schema/tables/:table", r.middl.ValidateParamId("id"))
	tableDataGroup.GET("/data", handler.FetchTableData)
	tableDataGroup.GET("/data/key/:key", handler.FetchTableDataByKey)
	tableDataGroup.POST("/data", handler.CreateTableData)
	tableDataGroup.PUT("/data/key/:key", handler.UpdateTableData)
	tableDataGroup.DELETE("/data/key/:key", handler.DeleteTableData)

	// Datasets Route
	datasetsGroup := r.e.Group("/v1/datasets")
	datasetsGroup.GET("", handler.FetchDatasetList)
	datasetsGroup.GET("/:id", handler.FetchDatasetByID)
	datasetsGroup.POST("", handler.UpsertDataset)
	datasetsGroup.DELETE("/:id", handler.DeleteDatasetByID)

	// Dataset Versions Route
	datasetVersionsGroup := r.e.Group("/v1/datasets/:id/versions")
	datasetVersionsGroup.GET("", handler.FetchDatasetVersionsList) // filter validate // finish
	datasetVersionsGroup.GET("/:version", handler.FetchDatasetVersionByID)
	datasetVersionsGroup.POST("", handler.InsertDatasetVersion)                 //insert + DTO validate //REPO function ไม่ต้องแก้
	datasetVersionsGroup.PUT("/:version", handler.UpdateDatasetVersion)         //update + DTO validate
	datasetVersionsGroup.PATCH("/:version", handler.UpdateDatasetVersionStatus) // patch update status active,preview,deprecated + DTO validate รับ status

	// Serving Routes
	servingGroup := r.e.Group("/v1/datasets/:id/versions/:version")
	servingGroup.GET("/data", handler.ServingDatasetVersionData)
	servingGroup.GET("/data/key/:key", handler.ServingDatasetVersionDataByKey)
	servingGroup.POST("/data", handler.CreateDatasetVersionData)
	servingGroup.PUT("/data/key/:key", handler.UpdateDatasetVersionDataByKey)
	servingGroup.DELETE("/data/key/:key", handler.DeleteDatasetVersionDataByKey)

	// Reporting Template Route
	reportingGroup := r.e.Group("/v1/reporting")
	reportingTemplateGroup := reportingGroup.Group("/templates")
	reportingTemplateGroup.POST("/upload", handler.UploadReportingTemplate, r.middl.InputForm)
	reportingTemplateGroup.POST("/:reporting_template_id/export/key/:key", handler.ExportReportingJob)
	reportingTemplateGroup.GET("/export/job/:job_id", handler.FetchReportingExportJobByID, r.middl.ValidateParamId("job_id"))
	reportingGroup.GET("/export/job/:job_id", handler.FetchExportJobByJobId, r.middl.ValidateParamId("job_id"))
	reportingGroup.POST("/export/job", handler.ExportJob)

	// Import Routes
	reportingGroup.POST("/import/template", handler.CreateImportTemplate)
	reportingGroup.POST("/import/job", handler.CreateImportJob)
	reportingGroup.GET("/import/job/:job_id", handler.FetchImportJobByID, r.middl.ValidateParamId("job_id"))
}
