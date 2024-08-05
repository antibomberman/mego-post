package models

type PostContentFile struct {
	FileName string `db:"file_name" json:"filename"`
	Size     int64  `db:"size" json:"size"`
	Url      string `db:"url" json:"url"`
	Type     int    `db:"type" json:"type"`
}
