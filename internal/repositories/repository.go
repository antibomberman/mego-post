package repositories

import (
	"antibomberman/mego-post/internal/models"
	"time"
)

type PostRepository interface {
	Find(startIndex int, size int, sort string, search string, dateFrom *time.Time, dateTo *time.Time) ([]models.Post, error)
	GetByAuthor(authorId string, startIndex int, size int, sort string) ([]models.Post, error)
	GetById(string) (models.Post, error)
	Create(models.PostCreate) (string, error)
	Delete(id, authorId string) error
	Update(models.PostUpdate) error
	CountByAuthor(string) (int, error)
}

type PostContentRepository interface {
	Find(postId string) ([]models.PostContent, error)
	Create(create models.PostContentCreate) (id string, err error)
	Delete(postId string) error
}
