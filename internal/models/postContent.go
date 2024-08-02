package models

type PostContent struct {
	Title            string            `json:"title"`
	Content          string            `json:"content"`
	PostContentFiles []PostContentFile `json:"files"`
}
