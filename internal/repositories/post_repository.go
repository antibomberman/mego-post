package repositories

import (
	"antibomberman/mego-post/internal/models"
	"github.com/jmoiron/sqlx"
	"strconv"
	"time"
)

type postRepository struct {
	db *sqlx.DB
}

func NewPostRepository(db *sqlx.DB) PostRepository {
	return &postRepository{
		db: db,
	}
}

func (r *postRepository) Find(startIndex int, size int, sort string, search string, dateFrom, dateTo *time.Time) ([]models.Post, error) {
	var posts []models.Post

	query := `SELECT * FROM posts`
	if search != "" {
		query += ` WHERE title LIKE '%` + search + `%'`
	}
	if dateFrom != nil {
		query += ` AND created_at >= ` + dateFrom.Format("2006-01-02")
	}
	if dateTo != nil {
		query += ` AND created_at <= ` + dateTo.Format("2006-01-02")
	}
	switch sort {
	case "0":
		query += " ORDER BY created_at DESC"
	case "1":
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
func (r *postRepository) GetByAuthor(authorId string, startIndex int, size int, sort string) ([]models.Post, error) {
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

func (r *postRepository) Create(data models.PostCreate) (string, error) {
	res, err := r.db.Exec("INSERT INTO posts (title, author_id, type) values ($1, $2, $3)", data.Title, data.AuthorId, data.Type)
	if err != nil {
		return "", err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(id, 10), nil
}

func (r *postRepository) Delete(id, authorId string) error {
	_, err := r.db.Exec("DELETE FROM posts WHERE id = $1 AND author_id = $2;", id, authorId)
	if err != nil {
		return err
	}
	return nil
}

func (r *postRepository) Update(data models.PostUpdate) error {
	_, err := r.db.Exec("UPDATE posts SET title = $2, type = $3,updated_at = $4 WHERE id = $1;", data.Id, data.Title, data.Type, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (r *postRepository) Hide(id string) error {
	return nil
}

func (r *postRepository) CountByAuthor(authorId string) (int, error) {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM posts WHERE author_id = $1", authorId)
	if err != nil {
		return 0, err
	}
	return count, nil
}
