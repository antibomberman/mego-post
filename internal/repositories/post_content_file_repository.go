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
	err := r.db.Select(&postContentFiles, "SELECT * FROM post_content_files WHERE post_content_id = $1", postContentId)
	if err != nil {
		return []models.PostContentFile{}, err
	}
	return postContentFiles, nil
}

func (r *postContentFileRepository) Create(create models.PostContentFileCreate) (string, error) {
	query := `
        INSERT INTO post_content_files (post_content_id, file_name, content_type)
        VALUES ($1, $2, $3, $4) RETURNING id;
    `
	id := ""
	err := r.db.QueryRowx(query, create.PostContentId, create.FileName, create.ContentType).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *postContentFileRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM post_content_files WHERE id = $1;", id)
	if err != nil {
		return err
	}
	return nil
}
