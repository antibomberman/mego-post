package services

import (
	postGrpc "github.com/antibomberman/mego-protos/gen/go/post"
)

type PostService interface {
	Find(int, string, string) ([]*postGrpc.PostDetail, string, error)
}
