package grpc

import (
	"antibomberman/mego-post/internal/config"
	"antibomberman/mego-post/internal/services"
	postGrpc "github.com/antibomberman/mego-protos/gen/go/post"
	"google.golang.org/grpc"
)

type serverAPI struct {
	postGrpc.UnimplementedPostServiceServer
	service         services.PostService
	categoryService services.CategoryService
	cfg             *config.Config
}

func Register(gRPC *grpc.Server, cfg *config.Config, service services.PostService, categoryService services.CategoryService) {
	postGrpc.RegisterPostServiceServer(gRPC, &serverAPI{
		service:         service,
		categoryService: categoryService,
		cfg:             cfg,
	})
}
