package models

type Author struct {
	Id         string
	FirstName  string
	MiddleName string
	LastName   string
	Email      string
	Phone      string
	Avatar     *Avatar
}
type Avatar struct {
	FileName string
	Url      string
}
