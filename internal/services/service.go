package services

import (
	"antibomberman/mego-post/internal/models"
	"time"
)

type PostService interface {
	Find(pageSize int, pageToken string, sort string, search string, dateFrom, dateTo *time.Time) ([]models.PostDetail, string, error)
	Create(post models.PostCreate) (*models.PostDetail, error)
	Update(post models.PostUpdate) (*models.PostDetail, error)
	GetById(id string) (*models.PostDetail, error)
	Delete(id, authorId string) error
	GetByAuthor(authorId string, pageSize int, pageToken string, sort string) ([]models.PostDetail, string, error)
}
type PostContentService interface {
}

type CategoryService interface {
	Find() ([]models.CategoryDetails, error)
	Create(models.CategoryCreate) (*models.CategoryDetails, error)
	Update(models.CategoryUpdate) (*models.CategoryDetails, error)
	GetById(string) (*models.CategoryDetails, error)
	Delete(string) error
}
