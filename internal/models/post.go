package models

type Post struct {
	Id    int    `db:"id" json:"id"`
	Title string `db:"title" json:"title"`
}

type PostCreate struct {
	Title     string `json:"title" validate:"required" db:"title"`
	Content   string `json:"content" validate:"required" db:"content"`
	ImagePath string `json:"image_path" db:"image_path"`
	UserId    int    `json:"user_id" db:"user_id"`
}
type PostUpdate struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	ImagePath string `json:"image_path"`
}
