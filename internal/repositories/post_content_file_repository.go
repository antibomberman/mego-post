package repositories

import (
	"antibomberman/mego-post/internal/models"
	"github.com/jmoiron/sqlx"
)

type postContentFileRepository struct {
	db *sqlx.DB
}

func NewPostContentFileRepository(db *sqlx.DB) PostContentFileRepository {
	return &postContentFileRepository{
		db: db,
	}
}

func (r *postContentFileRepository) Find(postContentId string) ([]models.PostContentFile, error) {
	var postContentFiles []models.PostContentFile
	err := r.db.Select(&postContentFiles, "SELECT filename,path,size,type FROM post_content_files WHERE post_content_id = $1", postContentId)
	if err != nil {
		return []models.PostContentFile{}, err
	}
	return postContentFiles, nil
}

func (r *postContentFileRepository) Create(create models.PostContentFileCreate) (string, error) {
	query := `
        INSERT INTO post_content_files (post_content_id, filename, path, size, type)
        VALUES ($1, $2, $3, $4, $5) RETURNING id;
    `
	id := ""
	err := r.db.QueryRowx(query, create.PostContentId, create.FileName, create.Path, create.Size, create.Type).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}
