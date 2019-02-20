package entity

import (
	"strconv"
)

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

type ID string

func NewIDInt(id int) ID {
	return ID(strconv.FormatInt(int64(id), 10))
}

func NewIDString(id string) ID {
	return ID(id)
}

func (id ID) EqualInt(i int32) bool {
	return id == ID(i)
}

func (id ID) IsEmpty() bool {
	return len(id) == 0
}

func (id ID) String() string {
	return string(id)
}
