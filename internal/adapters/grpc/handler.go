package grpc

import (
	"antibomberman/mego-post/internal/dto"
	"antibomberman/mego-post/internal/models"
	"context"
	postGrpc "github.com/antibomberman/mego-protos/gen/go/post"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

func (s serverAPI) FindPost(ctx context.Context, req *postGrpc.FindPostRequest) (*postGrpc.FindPostResponse, error) {
	var dateFrom, dateTo *time.Time
	if req.DateFrom != nil {
		dateFromValue := req.DateFrom.AsTime()
		dateFrom = &dateFromValue
	}
	if req.DateTo != nil {
		dateToValue := req.DateTo.AsTime()
		dateTo = &dateToValue
	}
	posts, nextPageToken, err := s.service.Find(int(req.PageSize), req.PageToken, req.SortOrder.String(), req.Search, dateFrom, dateTo)
	if err != nil {
		log.Printf("Error getting posts: %v", err)
		return nil, status.Error(codes.Internal, "Failed to retrieve posts")
	}
	return &postGrpc.FindPostResponse{
		Posts:         dto.ToPbPostDetails(posts),
		NextPageToken: nextPageToken,
	}, nil
}
func (s serverAPI) GetByAuthor(ctx context.Context, req *postGrpc.GetByAuthorRequest) (*postGrpc.GetByAuthorResponse, error) {
	posts, nextPageToken, err := s.service.GetByAuthor(req.AuthorId, int(req.PageSize), req.PageToken, req.SortOrder.String())
	if err != nil {
		log.Printf("Error getting posts: %v", err)
		return nil, status.Error(codes.Internal, "Failed to retrieve posts")
	}

	return &postGrpc.GetByAuthorResponse{
		Posts:         dto.ToPbPostDetails(posts),
		NextPageToken: nextPageToken,
	}, nil
}
func (s serverAPI) GetById(ctx context.Context, req *postGrpc.GetByIdRequest) (*postGrpc.PostDetail, error) {
	post, err := s.service.GetById(req.Id)
	if err != nil {
		log.Printf("Error getting post: %v", err)
		return nil, status.Error(codes.Internal, "Failed to retrieve post")
	}
	return dto.ToPbPostDetail(*post), nil
}
func (s serverAPI) CreatePost(ctx context.Context, req *postGrpc.CreatePostRequest) (*postGrpc.PostDetail, error) {
	postDetail, err := s.service.Create(models.PostCreate{
		AuthorId: req.AuthorId,
		Type:     int(req.Type),
		Image: models.FileCreate{
			FileName:    req.Image.FileName,
			ContentType: req.Image.ContentType,
			Data:        req.Image.Data,
		},
		Contents:   dto.ToPostContentCreateOrUpdate(req.Contents),
		Categories: req.Categories,
	})
	if err != nil {
		log.Printf("Error creating post: %v", err)
		return nil, err
	}
	return dto.ToPbPostDetail(*postDetail), nil
}
func (s serverAPI) UpdatePost(ctx context.Context, req *postGrpc.UpdatePostRequest) (*postGrpc.PostDetail, error) {
	postDetail, err := s.service.Update(models.PostUpdate{
		Id:   req.Id,
		Type: int(req.Type),
		Image: models.FileCreate{
			FileName:    req.Image.FileName,
			ContentType: req.Image.ContentType,
			Data:        req.Image.Data,
		},
		Contents: dto.ToPostContentCreateOrUpdate(req.Contents),
	})
	if err != nil {
		log.Printf("Error updating post: %v", err)
		return nil, err
	}
	return dto.ToPbPostDetail(*postDetail), nil
}
func (s serverAPI) DeletePost(ctx context.Context, req *postGrpc.DeletePostRequest) (*postGrpc.DeletePostResponse, error) {
	err := s.service.Delete(req.Id, req.AuthorId)
	if err != nil {
		log.Printf("Error deleting post: %v", err)
		return nil, status.Error(codes.Internal, "Failed to delete post")
	}
	return &postGrpc.DeletePostResponse{}, nil
}

func (s serverAPI) HidePost(context.Context, *postGrpc.HidePostRequest) (*postGrpc.HidePostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HidePost not implemented")
}
