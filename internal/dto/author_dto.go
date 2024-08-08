package dto

import (
	"antibomberman/mego-post/internal/models"
	postGrpc "github.com/antibomberman/mego-protos/gen/go/post"
	userGrpc "github.com/antibomberman/mego-protos/gen/go/user"
)

func ToPbAuthorDetail(author models.Author, avatar *models.Avatar) *postGrpc.Author {
	return &postGrpc.Author{
		Id:        author.Id,
		FirstName: author.FirstName,
		LastName:  author.LastName,
		Email:     author.Email,
		Phone:     author.Phone,
		Avatar:    ToPbAvatar(avatar),
	}
}
func ToPbAvatar(avatar *models.Avatar) *postGrpc.Avatar {
	if avatar == nil {
		return nil
	}
	return &postGrpc.Avatar{
		FileName: avatar.FileName,
		Url:      avatar.Url,
	}
}
func ToAvatar(avatar *userGrpc.Avatar) *models.Avatar {
	if avatar == nil {
		return nil
	}
	return &models.Avatar{
		FileName: avatar.FileName,
		Url:      avatar.Url,
	}
}
