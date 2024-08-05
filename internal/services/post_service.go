package services

import (
	"antibomberman/mego-post/internal/models"
	"antibomberman/mego-post/pkg/utils"
	postGrpc "github.com/antibomberman/mego-protos/gen/go/post"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
)
import "antibomberman/mego-post/internal/repositories"

type postService struct {
	postRepository repositories.PostRepository
}

func NewPostService(repo repositories.PostRepository, credentials grpc.DialOption) PostService {
	return &postService{postRepository: repo}
}

func (p *postService) Find(pageSize int, pageToken string, search string) ([]*postGrpc.PostDetail, string, error) {
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

func (p *postService) BuildPostDetails(posts []models.Post, pageSize, startIndex int) []*postGrpc.PostDetail {

	var postDetails []*postGrpc.PostDetail
	for _, post := range posts {
		postDetails = append(postDetails, p.BuildPostDetail(post))
	}

	return postDetails
}

func (p *postService) BuildPostDetail(post models.Post) *postGrpc.PostDetail {
	mediaContents, _ := p.GetMediaContents(post.Id)
	return &postGrpc.PostDetail{
		Id:       post.Id,
		Title:    post.Title,
		Contents: mediaContents,
		Author: &postGrpc.Author{
			Id: post.AuthorId,
		},
		CommentCount: 0,
		LikeCount:    0,
		RepostCount:  0,
		ViewCount:    0,
		CreatedAt:    timestamppb.New(post.CreatedAt.Time),
		UpdatedAt:    timestamppb.New(post.UpdatedAt.Time),
		DeletedAt:    timestamppb.New(post.DeletedAt.Time),
	}

}

func (p *postService) GetMediaContents(postId string) ([]*postGrpc.MediaContents, error) {
	contents, err := p.postRepository.GetContents(postId)
	if err != nil {
		return nil, err
	}

	var mediaContents []*postGrpc.MediaContents
	for _, content := range contents {
		mediaContentFiles, err := p.GetMediaContentFiles(content.Id)
		if err != nil {
			log.Printf("Error getting media content files for content %d: %v", content.Id, err)
			continue
		}
		mediaContents = append(mediaContents, &postGrpc.MediaContents{
			Files: mediaContentFiles,
		})
	}
	return mediaContents, nil
}

func (p *postService) GetMediaContentFiles(contentId string) ([]*postGrpc.MediaContentFiles, error) {
	contentFiles, err := p.postRepository.GetContentFiles(contentId)
	if err != nil {
		return nil, err
	}

	var mediaContentFiles []*postGrpc.MediaContentFiles
	for _, contentFile := range contentFiles {
		mediaContentFiles = append(mediaContentFiles, &postGrpc.MediaContentFiles{
			Url: contentFile.Url,
		})
	}
	return mediaContentFiles, nil
}
