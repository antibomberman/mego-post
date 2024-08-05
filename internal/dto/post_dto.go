package dto

import (
	"antibomberman/mego-post/internal/models"
	postGrpc "github.com/antibomberman/mego-protos/gen/go/post"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToPbPostDetail(details models.PostDetail) *postGrpc.PostDetail {
	pbUserDetails := &postGrpc.PostDetail{
		Id:        details.Id,
		Title:     details.Title,
		CreatedAt: timestamppb.New(*details.CreatedAt),
		UpdatedAt: timestamppb.New(*details.UpdatedAt),
		DeletedAt: timestamppb.New(*details.DeletedAt),
		Author:    ToPbAuthorDetail(details.Author),
		Contents:  ToPbPostContent(details.Contents),
	}
	return pbUserDetails
}
func ToPbPostDetails(details []models.PostDetail) []*postGrpc.PostDetail {
	pbPostDetails := make([]*postGrpc.PostDetail, len(details))
	for _, detail := range details {
		pbPostDetails = append(pbPostDetails, ToPbPostDetail(detail))
	}
	return pbPostDetails
}
