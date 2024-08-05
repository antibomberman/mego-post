package dto

import (
	"antibomberman/mego-post/internal/models"
	postGrpc "github.com/antibomberman/mego-protos/gen/go/post"
)

func ToPbPostContent(details []models.PostContentWithFile) []*postGrpc.MediaContents {
	pbContents := make([]*postGrpc.MediaContents, 0, len(details))
	for _, content := range details {
		pbContents = append(pbContents, &postGrpc.MediaContents{
			Title:   content.Title,
			Content: content.Content,
			Files:   ToPbPostContentFiles(content.PostContentFiles),
		})
	}
	return pbContents
}
func ToPbPostContentFiles(details []models.PostContentFile) []*postGrpc.MediaContentFiles {
	pbFiles := make([]*postGrpc.MediaContentFiles, 0, len(details))
	for _, file := range details {
		pbFiles = append(pbFiles, &postGrpc.MediaContentFiles{
			FileName: file.FileName,
			Size:     file.Size,
			Url:      file.Url,
		})
	}
	return pbFiles
}
