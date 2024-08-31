package dto

import (
	"antibomberman/mego-post/internal/models"
	postGrpc "github.com/antibomberman/mego-protos/gen/go/post"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToPbPostDetail(details models.PostDetail) *postGrpc.PostDetail {

	image := &postGrpc.File{}
	if details.Image != nil {
		image.FileName = details.Image.FileName
		image.Url = details.Image.Url
		image.ContentType = details.Image.ContentType
	}
	pbUserDetails := &postGrpc.PostDetail{
		Id:         details.Id,
		Type:       postGrpc.PostType(details.Type),
		Image:      image,
		CreatedAt:  timestamppb.New(*details.CreatedAt),
		UpdatedAt:  timestamppb.New(*details.UpdatedAt),
		DeletedAt:  timestamppb.New(*details.DeletedAt),
		Author:     ToPbAuthorDetail(details.Author),
		Contents:   ToPbPostContent(details.Contents),
		Categories: ToPbCategories(details.Categories),
	}
	return pbUserDetails
}
func ToPbPostDetails(details []models.PostDetail) []*postGrpc.PostDetail {
	pbPostDetails := make([]*postGrpc.PostDetail, len(details))
	for i, detail := range details {
		pbPostDetails[i] = ToPbPostDetail(detail)
	}
	return pbPostDetails
}
