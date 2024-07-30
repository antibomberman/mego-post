package repositories

import (
	"antibomberman/mego-post/internal/models"
	"github.com/jmoiron/sqlx"
)

type PostRepository interface {
	Index(int, int) ([]models.Post, error)
	GetById(string) (models.Post, error)
	Create(models.PostCreate) (int, error)
	Delete(string) error
	Update(string, models.PostUpdate) error
}

type postRepository struct {
	db *sqlx.DB
}

func NewPostRepository(db *sqlx.DB) PostRepository {
	return &postRepository{
		db: db,
	}
}

func (r *postRepository) Index(startIndex int, size int) ([]models.Post, error) {
	var posts []models.Post

	err := r.db.Select(&posts, "SELECT id,title,created_at FROM posts OFFSET $1 LIMIT $2", startIndex, size)
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
	var postID int
	err := r.db.QueryRow(`
        INSERT INTO posts (user_id, title, content, image_path)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `, data.UserId, data.Title, data.Content, data.ImagePath).Scan(&postID)

	if err != nil {
		return 0, err
	}

	return postID, nil
}

func (r *postRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = $1;", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *postRepository) Update(id string, data models.PostUpdate) error {
	_, err := r.db.NamedExec("UPDATE posts SET title = :title, content = :content, image_path = :image_path WHERE id = :id", data)
	if err != nil {
		return err
	}
	return nil

}
