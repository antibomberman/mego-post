package services

import (
	"antibomberman/mego-post/internal/models"
)

type PostService interface {
	Find(int, string, string) ([]models.Post, string, error)
}
