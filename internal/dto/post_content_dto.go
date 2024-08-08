package dto

import (
	"antibomberman/mego-post/internal/models"
	postGrpc "github.com/antibomberman/mego-protos/gen/go/post"
)

func ToPbPostContent(details []models.PostContentWithFile) []*postGrpc.PostContent {
	pbContents := make([]*postGrpc.PostContent, 0, len(details))
	for _, content := range details {
		pbContents = append(pbContents, &postGrpc.PostContent{
			Title:   content.Title,
			Content: content.Content,
			Files:   ToPbPostContentFiles(content.PostContentFiles),
		})
	}
	return pbContents
}
func ToPbPostContentFiles(details []models.PostContentFile) []*postGrpc.PostContentFile {
	pbFiles := make([]*postGrpc.PostContentFile, 0, len(details))
	for _, file := range details {
		pbFiles = append(pbFiles, &postGrpc.PostContentFile{
			FileName:    file.FileName,
			ContentType: file.ContentType,
			Url:         file.Url,
		})
	}
	return pbFiles
}

func ToPostContentFileCreate(detail []*postGrpc.PostContentFileCreateOrUpdate) []models.PostContentFileCreate {
	files := make([]models.PostContentFileCreate, 0, len(detail))
	for _, file := range detail {
		files = append(files, models.PostContentFileCreate{
			FileName:    file.FileName,
			ContentType: file.ContentType,
			Data:        file.Data,
		})
	}
	return files
}

func ToPostContentCreateOrUpdate(detail []*postGrpc.PostContentCreateOrUpdate) []models.PostContentCreateOrUpdate {
	files := make([]models.PostContentCreateOrUpdate, 0, len(detail))
	for _, reqPostContent := range detail {
		files = append(files, models.PostContentCreateOrUpdate{
			Title:   reqPostContent.Title,
			Content: reqPostContent.Content,
			Files:   ToPostContentFileCreate(reqPostContent.Files),
		})
	}
	return files
}
