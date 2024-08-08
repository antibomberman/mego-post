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
}
