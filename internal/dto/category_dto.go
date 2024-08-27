package dto

import (
	"antibomberman/mego-post/internal/models"
	postGrpc "github.com/antibomberman/mego-protos/gen/go/post"
)

func ToPbCategory(details []models.CategoryDetails) []*postGrpc.Category {
	pbContents := make([]*postGrpc.Category, 0, len(details))
	for _, category := range details {
		pbContents = append(pbContents, &postGrpc.Category{
			Id:   category.Id,
			Name: category.Name,
			Icon: &postGrpc.File{
				FileName:    category.Icon.FileName,
				ContentType: category.Icon.ContentType,
				Url:         category.Icon.Url,
			},
		})
	}
	return pbContents
}
