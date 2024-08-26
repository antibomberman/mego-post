package services

import (
	"antibomberman/mego-post/internal/dto"
	"context"
	"log"
	"sync"
	"time"

	"antibomberman/mego-post/internal/clients"
	"antibomberman/mego-post/internal/models"
	"antibomberman/mego-post/internal/repositories"
	"antibomberman/mego-post/pkg/utils"
	storagePb "github.com/antibomberman/mego-protos/gen/go/storage"
	userPb "github.com/antibomberman/mego-protos/gen/go/user"
)

type postService struct {
	postRepository        repositories.PostRepository
	postContentRepository repositories.PostContentRepository
	userClient            *clients.UserClient
	storageClient         *clients.StorageClient
	favoriteClient        *clients.FavoriteClient
}

func NewPostService(postRepo repositories.PostRepository, postContentRepo repositories.PostContentRepository, userClient *clients.UserClient, storageClient *clients.StorageClient, favoriteClient *clients.FavoriteClient) PostService {
	return &postService{
		postRepository:        postRepo,
		postContentRepository: postContentRepo,
		userClient:            userClient,
		storageClient:         storageClient,
		favoriteClient:        favoriteClient,
	}
}

func (p *postService) Find(pageSize int, pageToken, sort, search string, fromDate, toDate *time.Time) ([]models.PostDetail, string, error) {
	var err error
	if pageSize < 1 {
		pageSize = 10
	}
	startIndex := 0
	if pageToken != "" {
		startIndex, err = utils.DecodePageToken(pageToken)
		if err != nil {
			log.Printf("Error decoding page token: %v", err)
			return nil, "", err
		}
	}

	posts, err := p.postRepository.Find(startIndex, pageSize+1, sort, search, fromDate, toDate)
	if err != nil {
		log.Printf("Error getting posts: %v", err)
		return nil, "", err
	}

	var nextPageToken string
	if len(posts) > pageSize {
		nextPageToken = utils.EncodePageToken(startIndex + pageSize)
		posts = posts[:pageSize]
	}

	postDetails := make([]models.PostDetail, len(posts))
	wg := &sync.WaitGroup{}
	wg.Add(len(posts))
	for i, post := range posts {
		go func(i int, post models.Post) {
			defer wg.Done()
			postDetails[i] = *p.buildPostDetail(post)
		}(i, post)
	}
	wg.Wait()
	return postDetails, nextPageToken, nil
}
func (p *postService) GetByAuthor(authorId string, pageSize int, pageToken, sort string) ([]models.PostDetail, string, error) {
	var err error
	if pageSize < 1 {
		pageSize = 10
	}
	startIndex := 0
	if pageToken != "" {
		startIndex, err = utils.DecodePageToken(pageToken)
	}
	if err != nil {
		log.Printf("Error decoding page token: %v", err)
		return nil, "", err
	}

	posts, err := p.postRepository.GetByAuthor(authorId, startIndex, pageSize+1, sort)
	if err != nil {
		log.Printf("Error getting posts: %v", err)
		return nil, "", err
	}

	var nextPageToken string
	if len(posts) > pageSize {
		nextPageToken = utils.EncodePageToken(startIndex + pageSize)
		posts = posts[:pageSize]
	}

	postDetails := make([]models.PostDetail, len(posts))
	wg := &sync.WaitGroup{}
	wg.Add(len(posts))
	for i, post := range posts {
		go func(i int, post models.Post) {
			defer wg.Done()
			postDetails[i] = *p.buildPostDetail(post)
		}(i, post)
	}
	wg.Wait()

	return postDetails, nextPageToken, nil
}
func (p *postService) GetById(id string) (*models.PostDetail, error) {
	post, err := p.postRepository.GetById(id)
	if err != nil {
		return nil, err
	}
	return p.buildPostDetail(post), nil
}

func (p *postService) Create(data models.PostCreate) (*models.PostDetail, error) {
	postId, err := p.postRepository.Create(data)
	if err != nil {
		return nil, err
	}
	for _, dataContent := range data.Contents {
		_, err := p.postContentRepository.Create(models.PostContentCreate{
			PostId:      postId,
			Title:       dataContent.Title,
			Description: dataContent.Description,
			File:        dataContent.File,
		})
		if err != nil {
			return nil, err
		}
	}

	return p.GetById(postId)
}
func (p *postService) Update(data models.PostUpdate) (*models.PostDetail, error) {
	err := p.postRepository.Update(data)
	if err != nil {
		return nil, err
	}
	err = p.deletePostContents(data.Id)
	if err != nil {
		return nil, err
	}
	for _, dataContent := range data.Contents {
		_, err := p.postContentRepository.Create(models.PostContentCreate{
			PostId:      data.Id,
			Title:       dataContent.Title,
			Description: dataContent.Description,
			File: models.FileCreate{
				FileName:    dataContent.File.FileName,
				ContentType: dataContent.File.ContentType,
				Data:        dataContent.File.Data,
			},
		})
		if err != nil {
			return nil, err
		}

	}

	return p.GetById(data.Id)
}

func (p *postService) Delete(id, authorId string) error {
	err := p.deletePostContents(id)
	if err != nil {
		return err
	}
	err = p.postRepository.Delete(id, authorId)
	if err != nil {
		return err
	}

	return nil
}

func (p *postService) deletePostContents(postId string) error {
	postContents, err := p.postContentRepository.Find(postId)
	if err != nil {
		return err
	}
	for _, postContent := range postContents {
		_, err := p.storageClient.DeleteObject(context.Background(), &storagePb.DeleteObjectRequest{
			FileName: postContent.File.FileName,
		})
		if err != nil {
			log.Printf("Error deleting storage object %s: %v", postContent.File.FileName, err)
		}
		err = p.postContentRepository.Delete(postContent.Id)
		if err != nil {
			log.Printf("Error deleting post content %s: %v", postContent.Id, err)
		}
	}
	return nil
}

func (p *postService) buildPostDetail(post models.Post) *models.PostDetail {
	mediaContents, err := p.getMediaContents(post.Id)
	if err != nil {
		log.Printf("Error getting media contents for post %s: %v", post.Id, err)
	}

	return &models.PostDetail{
		Id:           post.Id,
		Contents:     mediaContents,
		Author:       p.buildPostAuthorDetail(post.AuthorId),
		CommentCount: 0,
		LikeCount:    0,
		ViewCount:    0,
		CreatedAt:    &post.CreatedAt.Time,
		UpdatedAt:    &post.UpdatedAt.Time,
		DeletedAt:    &post.DeletedAt.Time,
	}
}

func (p *postService) buildPostAuthorDetail(authorId string) models.Author {
	user, err := p.userClient.GetById(context.Background(), &userPb.Id{Id: authorId})
	if err != nil {
		log.Printf("Error getting user by id %s: %v", authorId, err)
		return models.Author{}
	}

	return models.Author{
		Id:         user.Id,
		FirstName:  user.FirstName,
		MiddleName: user.MiddleName,
		LastName:   user.LastName,
		Email:      user.Email,
		Phone:      user.Phone,
		Avatar:     dto.ToAvatar(user.Avatar),
	}
}

func (p *postService) getMediaContents(postId string) ([]models.PostContent, error) {
	contents, err := p.postContentRepository.Find(postId)
	if err != nil {
		return []models.PostContent{}, err
	}
	return contents, nil
}
