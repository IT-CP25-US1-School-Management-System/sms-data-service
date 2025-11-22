package repository

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/proto/proto_models"
	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/service/document/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type grpcDocumentRepositoryInf struct {
	grpcAddress string
	timeout     int
	conn        *grpc.ClientConn
	mu          sync.RWMutex
}

func NewGRPCDocumentRepository(grpcAddress string, timeout int) document.GrpcDocumentRepository {
	return &grpcDocumentRepositoryInf{
		grpcAddress: grpcAddress,
		timeout:     timeout,
	}
}

func (g *grpcDocumentRepositoryInf) getConnection() (*grpc.ClientConn, error) {
	g.mu.RLock()
	if g.conn != nil {
		g.mu.RUnlock()
		return g.conn, nil
	}
	g.mu.RUnlock()

	g.mu.Lock()
	defer g.mu.Unlock()

	// Double-check after acquiring write lock
	if g.conn != nil {
		return g.conn, nil
	}

	conn, err := grpc.NewClient(
		g.grpcAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	g.conn = conn
	return g.conn, nil
}

// GetFileByResourceID implements document.GrpcDocumentRepository.
func (g *grpcDocumentRepositoryInf) GetFileByResourceID(ctx context.Context, req *proto_models.GetFileByResourceIDRequest) (int, *proto_models.FileResponse, error) {
	conn, err := g.getConnection()
	if err != nil {
		return http.StatusBadGateway, nil, err
	}

	client := proto_models.NewDocumentClient(conn)

	ctx, cancel := context.WithTimeout(ctx, time.Duration(g.timeout)*time.Second)
	defer cancel()

	response, err := client.GetFileByResourceID(ctx, req)
	if err != nil {
		grpcCode := status.Code(err)
		httpStatus := http.StatusInternalServerError

		switch grpcCode {
		case codes.InvalidArgument:
			httpStatus = http.StatusBadRequest
		case codes.NotFound:
			httpStatus = http.StatusNotFound
		case codes.Unavailable:
			httpStatus = http.StatusServiceUnavailable
		case codes.Unauthenticated:
			httpStatus = http.StatusNetworkAuthenticationRequired
		}
		return httpStatus, nil, err
	}

	if response == nil {
		return http.StatusBadGateway, nil, errors.New(http.StatusText(http.StatusBadGateway))
	}

	return http.StatusOK, response, nil
}

// UploadFile implements document.GrpcDocumentRepository.
func (g *grpcDocumentRepositoryInf) UploadFile(ctx context.Context, fileRequest *proto_models.FileRequest) (int, *proto_models.FileResponse, error) {
	conn, err := g.getConnection()
	if err != nil {
		return http.StatusBadGateway, nil, err
	}

	client := proto_models.NewDocumentClient(conn)

	ctx, cancel := context.WithTimeout(ctx, time.Duration(g.timeout)*time.Second)
	defer cancel()

	stream, err := client.UploadFile(ctx)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if err := stream.Send(fileRequest); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	response, err := stream.CloseAndRecv()

	if err != nil {
		grpcCode := status.Code(err)
		httpStatus := http.StatusInternalServerError

		switch grpcCode {
		case codes.InvalidArgument:
			httpStatus = http.StatusBadRequest
		case codes.NotFound:
			httpStatus = http.StatusNotFound
		case codes.Unavailable:
			httpStatus = http.StatusServiceUnavailable
		case codes.Unauthenticated:
			httpStatus = http.StatusNetworkAuthenticationRequired
		}
		return httpStatus, nil, err
	}

	if response == nil {
		return http.StatusBadGateway, nil, errors.New(http.StatusText(http.StatusBadGateway))
	}

	return http.StatusOK, response, nil
}
