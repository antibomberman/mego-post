package grpc

import (
	"context"
	postGrpc "github.com/antibomberman/mego-protos/gen/go/post"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func (s serverAPI) Find(ctx context.Context, req *postGrpc.FindPostRequest) (*postGrpc.FindPostResponse, error) {
	posts, nextPageToken, err := s.service.Find(int(req.PageSize), req.PageToken, req.Search)
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
func (s serverAPI) CreatePost(context.Context, *postGrpc.CreatePostRequest) (*postGrpc.PostDetail, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePost not implemented")
}
func (s serverAPI) UpdatePost(context.Context, *postGrpc.UpdatePostRequest) (*postGrpc.PostDetail, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePost not implemented")
}
func (s serverAPI) DeletePost(context.Context, *postGrpc.DeletePostRequest) (*postGrpc.PostDetail, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePost not implemented")
}
func (s serverAPI) HidePost(context.Context, *postGrpc.HidePostRequest) (*postGrpc.PostDetail, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HidePost not implemented")
}
