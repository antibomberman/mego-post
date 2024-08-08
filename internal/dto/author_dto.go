package dto

import (
	"antibomberman/mego-post/internal/models"
	postGrpc "github.com/antibomberman/mego-protos/gen/go/post"
)

func ToPbAuthorDetail(author models.Author) *postGrpc.Author {
	return &postGrpc.Author{
		Id:        author.Id,
		FirstName: author.FirstName,
		LastName:  author.LastName,
		Email:     author.Email,
		Phone:     author.Phone,
		//Avatar:    author.Avatar,
	}
}
