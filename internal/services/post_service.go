package services

import (
	"antibomberman/mego-post/internal/clients"
	"antibomberman/mego-post/internal/models"
	"antibomberman/mego-post/pkg/utils"
	"context"
	postPb "github.com/antibomberman/mego-protos/gen/go/post"
	userPb "github.com/antibomberman/mego-protos/gen/go/user"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
)
import "antibomberman/mego-post/internal/repositories"

type postService struct {
	postRepository repositories.PostRepository
	userClient     *clients.UserClient
}

func NewPostService(repo repositories.PostRepository, client *clients.UserClient) PostService {
	return &postService{postRepository: repo, userClient: client}
}

func (p *postService) Find(pageSize int, pageToken string, search string) ([]*postPb.PostDetail, string, error) {
	if pageSize < 1 {
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
		return nil, "", err
	}

	var nextPageToken string
	if len(posts) > pageSize {
		nextPageToken = utils.EncodePageToken(startIndex + pageSize)
		posts = posts[:pageSize]
	}
	postDetails := p.BuildPostDetails(posts, pageSize, startIndex)

	return postDetails, nextPageToken, nil

}

func (p *postService) BuildPostDetails(posts []models.Post, pageSize, startIndex int) []*postPb.PostDetail {

	var postDetails []*postPb.PostDetail
	for _, post := range posts {
		postDetails = append(postDetails, p.BuildPostDetail(post))
	}

	return postDetails
}

func (p *postService) BuildPostDetail(post models.Post) *postPb.PostDetail {
	mediaContents, _ := p.GetMediaContents(post.Id)
	return &postPb.PostDetail{
		Id:           post.Id,
		Title:        post.Title,
		Contents:     mediaContents,
		Author:       p.BuildPostAuthorDetail(post.AuthorId),
		CommentCount: 0,
		LikeCount:    0,
		RepostCount:  0,
		ViewCount:    0,
		CreatedAt:    timestamppb.New(post.CreatedAt.Time),
		UpdatedAt:    timestamppb.New(post.UpdatedAt.Time),
		DeletedAt:    timestamppb.New(post.DeletedAt.Time),
	}

}
func (p *postService) BuildPostAuthorDetail(authorId string) *postPb.Author {
	log.Printf("Getting user by id %s", authorId)
	user, err := p.userClient.GetById(context.Background(), &userPb.Id{Id: authorId})
	if err != nil {
		log.Printf("Error getting user by id %s: %v", authorId, err)
		return nil
	}
	return &postPb.Author{
		Id:         user.Id,
		FirstName:  user.FirstName,
		MiddleName: user.MiddleName,
		LastName:   user.LastName,
		Email:      user.Email,
		Phone:      user.Phone,
		Avatar:     user.Avatar,
	}

}

func (p *postService) GetMediaContents(postId string) ([]*postPb.MediaContents, error) {
	contents, err := p.postRepository.GetContents(postId)
	if err != nil {
		return nil, err
	}

	var mediaContents []*postPb.MediaContents
	for _, content := range contents {
		mediaContentFiles, err := p.GetMediaContentFiles(content.Id)
		if err != nil {
			log.Printf("Error getting media content files for content %d: %v", content.Id, err)
			continue
		}
		mediaContents = append(mediaContents, &postPb.MediaContents{
			Files: mediaContentFiles,
		})
	}
	return mediaContents, nil
}

func (p *postService) GetMediaContentFiles(contentId string) ([]*postPb.MediaContentFiles, error) {
	contentFiles, err := p.postRepository.GetContentFiles(contentId)
	if err != nil {
		return nil, err
	}

	var mediaContentFiles []*postPb.MediaContentFiles
	for _, contentFile := range contentFiles {
		mediaContentFiles = append(mediaContentFiles, &postPb.MediaContentFiles{
			Url: contentFile.Url,
		})
	}
	return mediaContentFiles, nil
}
