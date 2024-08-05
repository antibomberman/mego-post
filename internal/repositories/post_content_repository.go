package repositories

import (
	"antibomberman/mego-post/internal/models"
	"github.com/jmoiron/sqlx"
)

type postContentRepository struct {
	db *sqlx.DB
}

func NewPostContentRepository(db *sqlx.DB) PostContentRepository {
	return &postContentRepository{
		db: db,
	}
}

func (r *postContentRepository) Find(postId string) ([]models.PostContent, error) {
	var postContent []models.PostContent
	err := r.db.Select(&postContent, "SELECT id,title,content FROM post_contents WHERE post_id = $1", postId)
	if err != nil {
		return []models.PostContent{}, err
	}
	return postContent, nil
}

func (r *postContentRepository) Create(postContent models.PostContentCreate) (id string, err error) {
	query := `
        INSERT INTO post_contents (post_id, title, content)
            VALUES ($1, $2, $3) RETURNING id;
    `
	err = r.db.QueryRowx(query, postContent.PostId, postContent.Title, postContent.Content).Scan(&id)
	return id, err
}
