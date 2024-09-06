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
	case "NEWEST":
		query += " ORDER BY created_at DESC"
	case "OLDEST":
		query += " ORDER BY created_at ASC"
	case "MOST_LIKED":
		query += ""
	case "MOST_COMMENTED":
		query += ""
	default:
		query += " ORDER BY created_at DESC"
	}
	query += " OFFSET $1 LIMIT $2"
	err := r.db.Select(&posts, query, startIndex, size)
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

	query := `SELECT id, title, created_at FROM posts WHERE author_id = $1`

	switch sort {
	case "NEWEST":
		query += " ORDER BY created_at DESC"
	case "OLDEST":
		query += " ORDER BY created_at ASC"
	case "MOST_LIKED":
		query += ""
	case "MOST_COMMENTED":
		query += ""
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

func (r *postRepository) Create(AuthorId string, Type int, FileName, Title, Description string) (string, error) {
	var id int64

	err := r.db.QueryRow("INSERT INTO posts (author_id, type,image,title,description) VALUES ($1, $2, $3, $4, $5) RETURNING id", AuthorId, Type, FileName, Title, Description).Scan(&id)
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

func (r *postRepository) Update(id string, Type int, FileName, Title, Description string) error {
	_, err := r.db.Exec("UPDATE posts SET type = $2, image = $3,title = $4,description = $5,updated_at = $6 WHERE id = $1;", id, Type, FileName, Title, Description, time.Now())
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
