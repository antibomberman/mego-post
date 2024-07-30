package grpc

import (
	postGrpc "antibomberman/mego-post/gen"
	"antibomberman/mego-post/internal/config"
	"antibomberman/mego-post/internal/services"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type serverAPI struct {
	postGrpc.UnimplementedPostSrvServer
	service services.PostService
	cfg     *config.Config
}

func Register(gRPC *grpc.Server, cfg *config.Config, service services.PostService) {
	postGrpc.RegisterPostSrvServer(gRPC, &serverAPI{
		service: service,
		cfg:     cfg,
	})
}
func (s serverAPI) Index(ctx context.Context, req *postGrpc.IndexRequest) (*postGrpc.IndexResponse, error) {

	posts, nextPageToken, err := s.service.Index(int(req.PageSize), req.PageToken) // +1 для определения наличия следующей страницы
	if err != nil {
		log.Printf("Error getting posts: %v", err)
		return nil, status.Error(codes.Internal, "Failed to retrieve posts")
	}

	postResponses := make([]*postGrpc.Post, len(posts))
	for i, post := range posts {
		postResponses[i] = &postGrpc.Post{
			//Id:        post.Id,
			Title: post.Title,
			//CreatedAt: post.CreatedAt.Unix(),
		}
	}

	return &postGrpc.IndexResponse{
		Posts:         postResponses,
		NextPageToken: nextPageToken,
	}, nil
}
