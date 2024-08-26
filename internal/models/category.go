package models

type Category struct {
	Id   string
	Name string
	Icon File
}

type CategoryCreate struct {
	Id   string
	Name string
	Icon File
}
