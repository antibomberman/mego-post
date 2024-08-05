package repositories

import (
	"antibomberman/mego-post/internal/models"
	"github.com/jmoiron/sqlx"
)

type postRepository struct {
	db *sqlx.DB
}

func NewPostRepository(db *sqlx.DB) PostRepository {
	return &postRepository{
		db: db,
	}
}

func (r *postRepository) Find(startIndex int, size int, search string, sort int) ([]models.Post, error) {
	var posts []models.Post

	query := `SELECT * FROM posts`
	if search != "" {
		query += ` WHERE title LIKE '%` + search + `%'`
	}
	switch sort {
	case 0:
		query += " ORDER BY created_at DESC"
	case 1:
		query += " ORDER BY created_at ASC"
	default:
		query += " ORDER BY created_at DESC"
	}

	err := r.db.Select(&posts, query+" OFFSET $1 LIMIT $2", startIndex, size)
	if err != nil {
		return nil, err
	}
	if len(posts) == 0 {
		return []models.Post{}, nil
	}

	return posts, nil
}
func (r *postRepository) GetByAuthor(authorId string, startIndex int, size int, sort int) ([]models.Post, error) {
	var posts []models.Post

	query := `SELECT id, title, created_at FROM posts WHERE user_id = $1`

	switch sort {
	case 0:
		query += " ORDER BY created_at DESC"
	case 1:
		query += " ORDER BY created_at ASC"
	default:
		query += " ORDER BY created_at DESC"
	}

	err := r.db.Select(&posts, query+" OFFSET $2 LIMIT $3", authorId, startIndex, size)
	if err != nil {
		return nil, err
	}
	if len(posts) == 0 {
		return []models.Post{}, nil
	}

	return posts, nil
}

func (r *postRepository) GetById(id string) (models.Post, error) {
	var post models.Post
	err := r.db.Get(&post, "SELECT * FROM posts WHERE id = $1", id)
	if err != nil {
		return models.Post{}, err
	}
	return post, nil
}

func (r *postRepository) Create(data models.PostCreate) (int, error) {
	return 0, nil
}

func (r *postRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = $1;", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *postRepository) Update(id string, data models.PostUpdate) error {

	return nil

}

func (r *postRepository) Hide(id string) error {
	return nil
}

func (r *postRepository) CountByAuthor(authorId string) (int, error) {
	return 0, nil
}

func (r *postRepository) GetContents(postId string) ([]models.PostContent, error) {
	var postContent []models.PostContent
	err := r.db.Select(&postContent, "SELECT id,title,content FROM post_contents WHERE post_id = $1", postId)
	if err != nil {
		return []models.PostContent{}, err
	}
	return postContent, nil
}
func (r *postRepository) GetContentFiles(postContentId string) ([]models.PostContentFile, error) {
	var postContentFiles []models.PostContentFile
	err := r.db.Select(&postContentFiles, "SELECT filename,url,size,type FROM post_content_files WHERE post_content_id = $1", postContentId)
	if err != nil {
		return []models.PostContentFile{}, err
	}
	return postContentFiles, nil
}
