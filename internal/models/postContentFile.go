package models

type PostContentFile struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	FileName string `json:"filename"`
	Size     int64  `json:"size"`
	Url      string `json:"url"`
	Type     int    `json:"type"`
}
