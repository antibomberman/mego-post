package repositories

import (
	"antibomberman/mego-post/internal/models"
	"github.com/jmoiron/sqlx"
)

type categoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) CategoryRepository {
	return &categoryRepository{
		db: db,
	}
}

func (r *categoryRepository) Find() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Select(&categories, "SELECT id,name,icon FROM categories")
	if err != nil {
		return []models.Category{}, err
	}
	return categories, nil
}

func (r *categoryRepository) Create(name, FileName string) (id string, err error) {
	query := `
        INSERT INTO categories (name, icon) VALUES ($1, $2) RETURNING id;
    `
	err = r.db.QueryRow(query, name, FileName).Scan(&id)
	return id, err
}
func (r *categoryRepository) Delete(id string) error {
	query := "DELETE FROM categories WHERE id = $1"
	_, err := r.db.Exec(query, id)
	return err
}

func (r *categoryRepository) ByPostId(postId string) ([]models.Category, error) {
	var categories []models.Category
	query := `
	SELECT categories.id, categories.name, categories.icon FROM categories
	JOIN post_categories ON post_categories.category_id = categories.id
	WHERE post_categories.post_id = $1
`
	err := r.db.Select(&categories, query, postId)
	if err != nil {
		return []models.Category{}, err
	}
	return categories, nil
}
func (r *categoryRepository) AddToPost(postId string, ids []string) error {
	query := "INSERT INTO post_categories (post_id, category_id) VALUES ($1, $2)"
	for _, id := range ids {
		_, err := r.db.Exec(query, postId, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *categoryRepository) RemoveFromPost(postId string) error {
	query := "DELETE FROM post_categories WHERE post_id = $1"
	_, err := r.db.Exec(query, postId)
	return err
}
