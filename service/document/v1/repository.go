package document

import (
	"context"

	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/proto/proto_models"
)

type GrpcDocumentRepository interface {
	UploadFile(ctx context.Context, fileRequest *proto_models.FileRequest) (int, *proto_models.FileResponse, error)
	GetFileByResourceID(ctx context.Context, req *proto_models.GetFileByResourceIDRequest) (int, *proto_models.FileResponse, error)
	Close() error
}
