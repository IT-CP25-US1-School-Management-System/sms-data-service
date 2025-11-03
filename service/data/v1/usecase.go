package data

import (
	"context"

	helperModel "github.com/GodeFvt/go-backend/helper/models"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/entity"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/filter"
)

type DataUsecase interface {
	FetchSourceList(ctx context.Context, paginator *helperModel.Paginator) ([]*entity.Sources, error)
	FetchSchemasList(ctx context.Context, filter *filter.SchemasFilter, paginator *helperModel.Paginator) ([]*entity.Schemas, error)
	FetchTablesList(ctx context.Context, filter *filter.TablesFilter, paginator *helperModel.Paginator) ([]*entity.Tables, error)
	FetchColumnsList(ctx context.Context, filter *filter.ColumnsFilter, paginator *helperModel.Paginator) ([]*entity.Columns, error)
	// Dataset
	FetchDatasetList(ctx context.Context, filter *filter.DatasetsFilter, paginator *helperModel.Paginator) ([]*entity.Datasets, error)
	FetchDatasetByID(ctx context.Context, datasetID string) (*entity.Datasets, error)
	UpsertDataset(ctx context.Context, dataset *entity.Datasets) error
	DeleteDatasetByID(ctx context.Context, datasetID string) error
}
