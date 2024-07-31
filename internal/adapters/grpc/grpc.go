package grpc

import (
	"antibomberman/mego-post/internal/config"
	"antibomberman/mego-post/internal/services"
	"context"
	postGrpc "github.com/antibomberman/mego-protos/gen/go/post"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
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
func (s serverAPI) Find(ctx context.Context, req *postGrpc.FindPostRequest) (*postGrpc.FindPostResponse, error) {
	posts, nextPageToken, err := s.service.Index(int(req.PageSize), req.PageToken, req.Search)
	if err != nil {
		log.Printf("Error getting posts: %v", err)
		return nil, status.Error(codes.Internal, "Failed to retrieve posts")
	}

	postResponses := make([]*postGrpc.PostDetail, len(posts))
	for i, post := range posts {
		postResponses[i] = &postGrpc.PostDetail{
			//Id:        post.Id,
			Title: post.Title,
			//CreatedAt: post.CreatedAt.Unix(),
		}
	}

	return &postGrpc.FindPostResponse{
		Posts:         postResponses,
		NextPageToken: nextPageToken,
	}, nil
}
