package grpc

import (
	"antibomberman/mego-post/internal/config"
	"antibomberman/mego-post/internal/services"
	postGrpc "github.com/antibomberman/mego-protos/gen/go/post"
	"google.golang.org/grpc"
)

type serverAPI struct {
	postGrpc.UnimplementedPostServiceServer
	service services.PostService
	cfg     *config.Config
}

func Register(gRPC *grpc.Server, cfg *config.Config, service services.PostService) {
	postGrpc.RegisterPostServiceServer(gRPC, &serverAPI{
		service: service,
		cfg:     cfg,
	})
}
