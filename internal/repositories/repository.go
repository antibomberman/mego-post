package repositories

import (
	"antibomberman/mego-post/internal/models"
)

type PostRepository interface {
	Find(startIndex int, size int, search string, sort int) ([]models.Post, error)
	GetByAuthor(string, int, int, int) ([]models.Post, error)
	GetById(string) (models.Post, error)
	Create(models.PostCreate) (int, error)
	Delete(string) error
	Update(string, models.PostUpdate) error
}
