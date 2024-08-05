package models

type PostContentFile struct {
	FileName string `db:"file_name" json:"filename"`
	Size     int64  `db:"size" json:"size"`
	Path     string `db:"path" json:"path"`
	Type     int    `db:"type" json:"type"`
}
type PostContentFileCreate struct {
	PostContentId string `db:"post_content_id" json:"-"`
	FileName      string `db:"file_name" json:"filename"`
	Size          int64  `db:"size" json:"size"`
	Path          string `db:"path" json:"path"`
	Type          int    `db:"type" json:"type"`
}

type PostContentFileBinary struct {
	FileName string `db:"file_name" json:"filename"`
	//MimeType
	ContentType string `db:"content_type" json:"content_type"`
	Data        []byte `db:"data" json:"-"`
}
