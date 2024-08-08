package services

import (
	"antibomberman/mego-post/internal/dto"
	"context"
	pb "github.com/antibomberman/mego-protos/gen/go/storage"
	"log"
	"time"

	"antibomberman/mego-post/internal/clients"
	"antibomberman/mego-post/internal/models"
	"antibomberman/mego-post/internal/repositories"
	"antibomberman/mego-post/pkg/utils"
	userPb "github.com/antibomberman/mego-protos/gen/go/user"
)

type postService struct {
	postRepository            repositories.PostRepository
	postContentRepository     repositories.PostContentRepository
	postContentFileRepository repositories.PostContentFileRepository
	userClient                *clients.UserClient
	storageClient             *clients.StorageClient
}

func NewPostService(postRepo repositories.PostRepository, postContentRepo repositories.PostContentRepository, postContentFileRepo repositories.PostContentFileRepository, userClient *clients.UserClient, storageClient *clients.StorageClient) PostService {
	return &postService{
		postRepository:            postRepo,
		postContentRepository:     postContentRepo,
		postContentFileRepository: postContentFileRepo,
		userClient:                userClient,
		storageClient:             storageClient,
	}
}

func (p *postService) Find(pageSize int, pageToken, sort, search string, fromDate, toDate *time.Time) ([]models.PostDetail, string, error) {
	if pageSize < 1 {
		pageSize = 10
	}
	startIndex, err := utils.DecodePageToken(pageToken)
	if err != nil {
		log.Printf("Error decoding page token: %v", err)
		return nil, "", err
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
	for i, post := range posts {
		postDetails[i] = *p.buildPostDetail(post)
	}

	return postDetails, nextPageToken, nil
}
func (p *postService) GetByAuthor(authorId string, pageSize int, pageToken, sort string) ([]models.PostDetail, string, error) {
	if pageSize < 1 {
		pageSize = 10
	}
	startIndex, err := utils.DecodePageToken(pageToken)
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
	for i, post := range posts {
		postDetails[i] = *p.buildPostDetail(post)
	}

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
		postContentId, err := p.postContentRepository.Create(models.PostContentCreate{
			PostId:  postId,
			Title:   dataContent.Title,
			Content: dataContent.Content,
		})
		if err != nil {
			return nil, err
		}
		for _, file := range dataContent.Files {
			object, err := p.storageClient.PutObject(context.Background(), &pb.PutObjectRequest{FileName: file.FileName, Data: file.Data, ContentType: file.ContentType})
			if err != nil {
				log.Printf("Error uploading file %s: %v", file.FileName, err)
				return nil, err
			}

			_, err = p.postContentFileRepository.Create(models.PostContentFileCreate{
				PostContentId: postContentId,
				FileName:      object.FileName,
				ContentType:   file.ContentType,
			})
			if err != nil {
				log.Printf("Error creating post content file: %v", err)
				return nil, err
			}
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
		postContentId, err := p.postContentRepository.Create(models.PostContentCreate{
			PostId:  data.Id,
			Title:   dataContent.Title,
			Content: dataContent.Content,
		})
		if err != nil {
			return nil, err
		}
		for _, file := range dataContent.Files {
			object, err := p.storageClient.PutObject(context.Background(), &pb.PutObjectRequest{FileName: file.FileName, Data: file.Data, ContentType: file.ContentType})
			if err != nil {
				log.Printf("Error uploading file %s: %v", file.FileName, err)
				return nil, err
			}

			_, err = p.postContentFileRepository.Create(models.PostContentFileCreate{
				PostContentId: postContentId,
				FileName:      object.FileName,
				ContentType:   file.ContentType,
			})
			if err != nil {
				log.Printf("Error creating post content file: %v", err)
				return nil, err
			}
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
		err := p.deletePostContentFiles(postContent.Id)
		if err != nil {
			log.Printf("Error deleting post content files for post %s: %v", postContent.Id, err)
			continue
		}
		err = p.postContentRepository.Delete(postContent.Id)
		if err != nil {
			log.Printf("Error deleting post content %s: %v", postContent.Id, err)
		}
	}
	return nil
}
func (p *postService) deletePostContentFiles(postContentId string) error {
	postContentFiles, err := p.postContentFileRepository.Find(postContentId)
	if err != nil {
		return err
	}
	for _, postContentFile := range postContentFiles {
		_, err := p.storageClient.DeleteObject(context.Background(), &pb.DeleteObjectRequest{FileName: postContentFile.FileName})
		if err != nil {
			log.Printf("Error deleting file %s: %v", postContentFile.FileName, err)
			continue
		}
		err = p.postContentFileRepository.Delete(postContentFile.Id)
		if err != nil {
			log.Printf("Error deleting post content file %s: %v", postContentFile.Id, err)
			continue
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
		Title:        post.Title,
		Contents:     mediaContents,
		Author:       p.buildPostAuthorDetail(post.AuthorId),
		CommentCount: 0,
		LikeCount:    0,
		RepostCount:  0,
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

func (p *postService) getMediaContents(postId string) ([]models.PostContentWithFile, error) {
	contents, err := p.postContentRepository.Find(postId)
	if err != nil {
		return []models.PostContentWithFile{}, err
	}

	mediaContents := make([]models.PostContentWithFile, 0, len(contents))
	for _, content := range contents {
		mediaContentFiles, err := p.getMediaContentFiles(content.Id)
		if err != nil {
			log.Printf("Error getting media content files for content %s: %v", content.Id, err)
			continue
		}
		mediaContents = append(mediaContents, models.PostContentWithFile{
			PostContentFiles: mediaContentFiles,
		})
	}
	return mediaContents, nil
}

func (p *postService) getMediaContentFiles(contentId string) ([]models.PostContentFile, error) {
	contentFiles, err := p.postContentFileRepository.Find(contentId)
	if err != nil {
		return []models.PostContentFile{}, err
	}

	mediaContentFiles := make([]models.PostContentFile, len(contentFiles))
	for i, contentFile := range contentFiles {
		mediaContentFiles[i] = models.PostContentFile{
			FileName:    contentFile.FileName,
			ContentType: contentFile.ContentType,
			Url:         contentFile.Url,
		}
	}
	return mediaContentFiles, nil
}
