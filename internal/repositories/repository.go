package repositories

import (
	"antibomberman/mego-post/internal/models"
	"time"
)

type PostRepository interface {
	Find(startIndex int, size int, sort string, search string, dateFrom *time.Time, dateTo *time.Time) ([]models.Post, error)
	GetByAuthor(authorId string, startIndex int, size int, sort string) ([]models.Post, error)
	GetById(string) (models.Post, error)
	Create(AuthorId string, Type int, FileName string) (string, error)
	Delete(id, authorId string) error
	Update(id string, Type int, FileName string) error
	CountByAuthor(string) (int, error)
}

type PostContentRepository interface {
	Find(postId string) ([]models.PostContent, error)
	Create(postId, Title, Description, FileName string) (id string, err error)
	Delete(postId string) error
}

type CategoryRepository interface {
	Find() ([]models.Category, error)
	Create(name, FileName string) (id string, err error)
	Delete(id string) error
	ByPostId(postId string) ([]models.Category, error)
	AddToPost(postId string, ids []string) error
	RemoveFromPost(postId string) error
}
