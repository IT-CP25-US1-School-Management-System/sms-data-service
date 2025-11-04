package data

import (
	"context"

	helperModel "github.com/GodeFvt/go-backend/helper/models"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/entity"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/models/filter"
	"github.com/gofrs/uuid"
)

type PsqlDataRepository interface {
	FetchInformationTablesBySourceID(ctx context.Context, dbType string, sourceID *uuid.UUID) ([]*entity.Tables, error)
	FetchInformationColumnsBySourceID(ctx context.Context, dbType string, sourceID *uuid.UUID) ([]*entity.Columns, error)
	FetchInformationSchemasBySourceID(ctx context.Context, dbType string, sourceID *uuid.UUID) ([]*entity.Schemas, error)
}

type PsqlDatasetRepository interface {
	FetchSourceList(ctx context.Context, paginator *helperModel.Paginator) ([]*entity.Sources, error)
	FetchSchemasList(ctx context.Context, filter *filter.SchemasFilter, paginator *helperModel.Paginator) ([]*entity.Schemas, error)
	FetchTablesList(ctx context.Context, filter *filter.TablesFilter, paginator *helperModel.Paginator) ([]*entity.Tables, error)
	FetchColumnsList(ctx context.Context, filter *filter.ColumnsFilter, paginator *helperModel.Paginator) ([]*entity.Columns, error)
	FetchSourceByID(ctx context.Context, sourceID *uuid.UUID) (*entity.Sources, error)
	UpsertSource(ctx context.Context, source *entity.Sources) error
	ActivateSourceByID(ctx context.Context, sourceID *uuid.UUID) error
	DeactivateSourceByID(ctx context.Context, sourceID *uuid.UUID) error
	ExistSourceByID(ctx context.Context, sourceID *uuid.UUID) (bool, error)
	DeleteSourceByID(ctx context.Context, sourceID *uuid.UUID) error

	BatchInsertInformationDatabase(ctx context.Context, schemas []*entity.Schemas, tables []*entity.Tables, columns []*entity.Columns) error

	// Dataset
	FetchDatasetList(ctx context.Context, filter *filter.DatasetsFilter, paginator *helperModel.Paginator) ([]*entity.Datasets, error)
	FetchDatasetByID(ctx context.Context, datasetID string) (*entity.Datasets, error)
	UpsertDataset(ctx context.Context, dataset *entity.Datasets) error
	DeleteDatasetByID(ctx context.Context, datasetID string) error
	ExistDatasetByID(ctx context.Context, datasetID string) (bool, error)
}

type RedisRepository interface {
}
