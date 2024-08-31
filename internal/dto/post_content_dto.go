package dto

import (
	"antibomberman/mego-post/internal/models"
	postGrpc "github.com/antibomberman/mego-protos/gen/go/post"
)

func ToPbPostContent(details []models.PostContentDetails) []*postGrpc.PostContent {
	pbContents := make([]*postGrpc.PostContent, 0, len(details))
	for _, content := range details {

		pbContents = append(pbContents, &postGrpc.PostContent{
			Title:       content.Title,
			Description: content.Description,
			Image: &postGrpc.File{
				FileName:    content.Image.FileName,
				ContentType: content.Image.ContentType,
				Url:         content.Image.Url,
			},
		})
	}
	return pbContents
}

func ToPostContentCreateOrUpdate(detail []*postGrpc.PostContentCreateOrUpdate) []models.PostContentCreateOrUpdate {
	files := make([]models.PostContentCreateOrUpdate, 0, len(detail))
	for _, reqPostContent := range detail {
		image := models.FileCreate{}
		if reqPostContent.File != nil {
			image.FileName = reqPostContent.File.FileName
			image.ContentType = reqPostContent.File.ContentType
			image.Data = reqPostContent.File.Data
		}
		files = append(files, models.PostContentCreateOrUpdate{
			Title:       reqPostContent.Title,
			Description: reqPostContent.Description,
			Image:       &image,
		})
	}
	return files
}
func ToCategoriesCreateOrUpdate(detail []*postGrpc.Category) []models.Category {
	categories := make([]models.Category, 0, len(detail))
	for _, reqPostContent := range detail {
		categories = append(categories, models.Category{
			Id:   reqPostContent.Id,
			Name: reqPostContent.Name,
			Icon: reqPostContent.Icon.FileName,
		})
	}
	return categories
}
