package dto

import (
	"antibomberman/mego-post/internal/models"
	postGrpc "github.com/antibomberman/mego-protos/gen/go/post"
)

func ToPbCategories(details []models.CategoryDetails) []*postGrpc.Category {
	pbContents := make([]*postGrpc.Category, 0, len(details))
	for _, category := range details {
		pbContents = append(pbContents, ToPbCategory(category))
	}
	return pbContents
}

func ToPbCategory(details models.CategoryDetails) *postGrpc.Category {
	return &postGrpc.Category{
		Id:   details.Id,
		Name: details.Name,
		Icon: &postGrpc.File{
			FileName:    details.Icon.FileName,
			ContentType: details.Icon.ContentType,
			Url:         details.Icon.Url,
		},
	}
}
