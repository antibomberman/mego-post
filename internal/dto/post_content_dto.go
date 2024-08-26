package dto

import (
	"antibomberman/mego-post/internal/models"
	postGrpc "github.com/antibomberman/mego-protos/gen/go/post"
)

func ToPbPostContent(details []models.PostContent) []*postGrpc.PostContent {
	pbContents := make([]*postGrpc.PostContent, 0, len(details))
	for _, content := range details {
		pbContents = append(pbContents, &postGrpc.PostContent{
			Title:       content.Title,
			Description: content.Description,
			File: &postGrpc.File{
				FileName:    content.File.FileName,
				ContentType: content.File.ContentType,
				Url:         content.File.Url,
			},
		})
	}
	return pbContents
}

func ToPostContentCreateOrUpdate(detail []*postGrpc.PostContentCreateOrUpdate) []models.PostContentCreateOrUpdate {
	files := make([]models.PostContentCreateOrUpdate, 0, len(detail))
	for _, reqPostContent := range detail {
		files = append(files, models.PostContentCreateOrUpdate{
			Title:       reqPostContent.Title,
			Description: reqPostContent.Description,
			File: models.FileCreate{
				FileName:    reqPostContent.File.FileName,
				ContentType: reqPostContent.File.ContentType,
				Data:        reqPostContent.File.Data,
			},
		})
	}
	return files
}
