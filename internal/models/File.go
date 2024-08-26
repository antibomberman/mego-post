package models

type File struct {
	FileName    string `db:"file_name" json:"filename"`
	ContentType string `db:"content_type" json:"content_type"`
	Url         string `json:"url"`
}
type FileCreate struct {
	FileName    string `db:"file_name" json:"filename"`
	ContentType string `db:"content_type" json:"content_type"`
	Data        []byte `db:"data" json:"-"`
}
