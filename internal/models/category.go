package models

type CategoryDetails struct {
	Id   string
	Name string
	Icon *File
}

type Category struct {
	Id   string
	Name string
	Icon string
}

type CategoryCreate struct {
	Name string
	Icon *FileCreate
}
type CategoryUpdate struct {
	Id   string
	Name string
	Icon *FileCreate
}
