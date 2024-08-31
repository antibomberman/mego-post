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
	categoryRepository    repositories.CategoryRepository
	userClient            *clients.UserClient
	storageClient         *clients.StorageClient
	favoriteClient        *clients.FavoriteClient
}

func NewPostService(postRepo repositories.PostRepository, postContentRepo repositories.PostContentRepository, catRepo repositories.CategoryRepository, userClient *clients.UserClient, storageClient *clients.StorageClient, favoriteClient *clients.FavoriteClient) PostService {
	return &postService{
		postRepository:        postRepo,
		postContentRepository: postContentRepo,
		categoryRepository:    catRepo,
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
	//upload main image
	if data.Image != nil && data.Image.FileName != "" {
		log.Printf("Uploading main image: %s", data.Image.FileName)

		rsp, err := p.storageClient.PutObject(context.Background(), &storagePb.PutObjectRequest{
			FileName:    data.Image.FileName,
			ContentType: data.Image.ContentType,
			Data:        data.Image.Data,
		})
		if err != nil {
			log.Printf("Error uploading main image: %v", err)
			return nil, err
		}
		log.Println(rsp)
		data.Image.FileName = rsp.FileName
	}

	postId, err := p.postRepository.Create(data.AuthorId, data.Type, data.Image.FileName)
	if err != nil {
		return nil, err
	}
	for _, dataContent := range data.Contents {
		imageName := ""
		if dataContent.Image.FileName != "" {
			rsp, err := p.storageClient.PutObject(context.Background(), &storagePb.PutObjectRequest{
				FileName:    dataContent.Image.FileName,
				ContentType: dataContent.Image.ContentType,
				Data:        dataContent.Image.Data,
			})
			if err != nil {
				return nil, err
			}
			imageName = rsp.FileName
		}

		_, err = p.postContentRepository.Create(postId, dataContent.Title, dataContent.Description, imageName)
		if err != nil {
			return nil, err
		}
	}
	if len(data.Categories) > 0 {
		err = p.categoryRepository.AddToPost(postId, data.Categories)
		if err != nil {
			log.Printf("Error adding categories to post: %v", err)
		}
	}

	return p.GetById(postId)
}
func (p *postService) Update(data models.PostUpdate) (*models.PostDetail, error) {
	if data.Image.Data != nil {
		err := p.DeletePostImage(data.Id)
		if err != nil {
			return nil, err
		}
		rsp, err := p.storageClient.PutObject(context.Background(), &storagePb.PutObjectRequest{
			FileName:    data.Image.FileName,
			ContentType: data.Image.ContentType,
			Data:        data.Image.Data,
		})
		if err != nil {
			return nil, err
		}
		data.Image.FileName = rsp.FileName
	}
	err := p.postRepository.Update(data.Id, data.Type, data.Image.FileName)
	if err != nil {
		return nil, err
	}
	err = p.deletePostContents(data.Id)
	if err != nil {
		return nil, err
	}
	for _, dataContent := range data.Contents {
		rsp, err := p.storageClient.PutObject(context.Background(), &storagePb.PutObjectRequest{
			FileName:    dataContent.Image.FileName,
			ContentType: dataContent.Image.ContentType,
			Data:        dataContent.Image.Data,
		})
		if err != nil {
			return nil, err
		}
		_, err = p.postContentRepository.Create(data.Id, dataContent.Title, dataContent.Description, rsp.FileName)
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
func (p *postService) DeletePostImage(id string) error {
	post, err := p.postRepository.GetById(id)
	if err != nil {
		return err
	}

	_, err = p.storageClient.DeleteObject(context.Background(), &storagePb.DeleteObjectRequest{
		FileName: post.Image,
	})
	if err != nil {
		log.Printf("Error deleting storage object %s: %v", post.Image, err)
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
			FileName: postContent.Image,
		})
		if err != nil {
			log.Printf("Error deleting storage object %s: %v", postContent.Image, err)
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
	image := models.File{}
	if post.Image != "" {
		rsp, err := p.storageClient.GetObjectUrl(context.Background(), &storagePb.GetObjectUrlRequest{
			FileName: post.Image,
		})
		if err != nil {
			log.Printf("Error getting storage object %s: %v", post.Image, err)
		}
		image.Url = rsp.Url
		image.FileName = rsp.FileName
		image.ContentType = rsp.ContentType

	}

	return &models.PostDetail{
		Id:           post.Id,
		Contents:     mediaContents,
		Categories:   p.buildCategoryDetail(post.Id),
		Author:       p.buildPostAuthorDetail(post.AuthorId),
		Image:        &image,
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
func (p *postService) buildCategoryDetail(postId string) []models.CategoryDetails {
	categories, err := p.categoryRepository.ByPostId(postId)
	if err != nil {
		return []models.CategoryDetails{}
	}
	var categoryDetails []models.CategoryDetails
	for _, category := range categories {
		file := &models.File{}
		if category.Icon != "" {
			rsp, err := p.storageClient.GetObjectUrl(context.Background(), &storagePb.GetObjectUrlRequest{
				FileName: category.Icon,
			})
			if err == nil {
				file.FileName = rsp.Url
				file.ContentType = rsp.ContentType
				file.Url = rsp.Url
			}
		}
		categoryDetail := models.CategoryDetails{
			Id:   category.Id,
			Name: category.Name,
			Icon: file,
		}
		categoryDetails = append(categoryDetails, categoryDetail)
	}
	return categoryDetails

}

func (p *postService) getMediaContents(postId string) ([]models.PostContentDetails, error) {
	contents, err := p.postContentRepository.Find(postId)
	if err != nil {
		return []models.PostContentDetails{}, err
	}
	var contentDetails []models.PostContentDetails
	for _, content := range contents {
		log.Printf("storage object %s", content.Image)
		image := models.File{}
		if content.Image != "" {
			rsp, err := p.storageClient.GetObjectUrl(context.Background(), &storagePb.GetObjectUrlRequest{
				FileName: content.Image,
			})
			if err == nil {
				image.FileName = rsp.FileName
				image.ContentType = rsp.ContentType
				image.Url = rsp.Url
			}
		}
		contentDetail := models.PostContentDetails{
			Id:          content.Id,
			Title:       content.Title,
			Description: content.Description,
			Image:       &image,
		}
		contentDetails = append(contentDetails, contentDetail)
	}
	return contentDetails, nil
}
