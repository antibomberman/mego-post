package grpc

import (
	"antibomberman/mego-post/internal/dto"
	"antibomberman/mego-post/internal/models"
	"context"
	postGrpc "github.com/antibomberman/mego-protos/gen/go/post"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func (s serverAPI) FindCategory(ctx context.Context, req *postGrpc.FindCategoryRequest) (*postGrpc.FindCategoryResponse, error) {
	categories, err := s.categoryService.Find()
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	pbCategories := make([]*postGrpc.Category, 0, len(categories))
	for _, category := range categories {
		pbCategories = append(pbCategories, &postGrpc.Category{
			Id:   category.Id,
			Name: category.Name,
			Icon: &postGrpc.File{
				FileName:    category.Icon.FileName,
				ContentType: category.Icon.ContentType,
				Url:         category.Icon.Url,
			},
		})
	}
	return &postGrpc.FindCategoryResponse{
		Categories: pbCategories,
	}, nil

}
func (s serverAPI) CreateCategory(ctx context.Context, req *postGrpc.CreateCategoryRequest) (*postGrpc.Category, error) {
	image := models.FileCreate{}
	if req.Icon != nil {
		image.FileName = req.Icon.FileName
		image.ContentType = req.Icon.ContentType
		image.Data = req.Icon.Data
	}
	categoryDetail, err := s.categoryService.Create(models.CategoryCreate{
		Name: req.Name,
		Icon: &image,
	})
	if err != nil {
		log.Printf("Error creating post: %v", err)
		return nil, err
	}
	return dto.ToPbCategory(*categoryDetail), nil
}
func (s serverAPI) UpdateCategory(ctx context.Context, req *postGrpc.UpdateCategoryRequest) (*postGrpc.Category, error) {
	image := models.FileCreate{}
	if req.Icon != nil {
		image.FileName = req.Icon.FileName
		image.ContentType = req.Icon.ContentType
		image.Data = req.Icon.Data
	}
	categoryDetail, err := s.categoryService.Update(models.CategoryUpdate{
		Id:   req.Id,
		Name: req.Name,
		Icon: &image,
	})
	if err != nil {
		log.Printf("Error updating post: %v", err)
		return nil, err
	}
	return dto.ToPbCategory(*categoryDetail), nil
}
func (s serverAPI) DeleteCategory(ctx context.Context, req *postGrpc.DeleteCategoryRequest) (*postGrpc.Category, error) {
	err := s.categoryService.Delete(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	return nil, nil
}
