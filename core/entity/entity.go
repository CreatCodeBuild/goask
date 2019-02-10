package entity

type Answer struct {
	ID         ID
	QuestionID ID
	Content    string
	AuthorID   ID // The ID of the user who created it.
	Accepted   bool
}

type User struct {
	ID   ID
	Name string
}

type ID int64

func (id ID) Equal(i int32) bool {
	return id == ID(i)
}
