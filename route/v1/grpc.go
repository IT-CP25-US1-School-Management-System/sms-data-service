package route

import (
	"google.golang.org/grpc"
)

type GrpcRoute struct {
	server *grpc.Server
}

func NewGRPCRoute(server *grpc.Server) *GrpcRoute {
	return &GrpcRoute{server}
}
