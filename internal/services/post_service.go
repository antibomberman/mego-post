package services

import (
	"antibomberman/mego-post/internal/models"
	"antibomberman/mego-post/pkg/utils"
	"log"
)
import "antibomberman/mego-post/internal/repositories"

type postService struct {
	postRepository repositories.PostRepository
}

func NewPostService(repo repositories.PostRepository) PostService {
	return &postService{postRepository: repo}
}

func (p *postService) Find(pageSize int, pageToken string, search string) ([]models.Post, string, error) {
	if pageSize == 0 {
		pageSize = 10
	}
	startIndex := 0
	if pageToken != "" {
		var err error
		startIndex, err = utils.DecodePageToken(pageToken)
		if err != nil {
			log.Printf("Error decoding page token: %v", err)
		}
	}
	posts, err := p.postRepository.Find(startIndex, pageSize+1, "", 0)
	if err != nil {
		log.Printf("Error getting posts: %v", err)
	}

	var nextPageToken string
	if len(posts) > pageSize {
		nextPageToken = utils.EncodePageToken(startIndex + pageSize)
		posts = posts[:pageSize]
	}
	return posts, nextPageToken, nil
}
