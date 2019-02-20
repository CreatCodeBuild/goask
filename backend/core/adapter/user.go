package adapter

import (
	"fmt"
	"goask/core/entity"
)

type UserDAO interface {
	// CRUD
	UserByID(ID entity.ID) (entity.User, error)
	Users() ([]entity.User, error)
	CreateUser(name string) (entity.User, error)
	// Associated Questions
	QuestionsOfUser(UserID entity.ID) ([]entity.Question, error)
}

type ErrUserNotFound struct {
	ID entity.ID
}

func (e *ErrUserNotFound) Error() string {
	return fmt.Sprintf("user:%v not found", e.ID)
}

type ErrUserIsNotAuthorOfQuestion struct {
	UserID     entity.ID
	QuestionID entity.ID
}

func (e *ErrUserIsNotAuthorOfQuestion) Error() string {
	return fmt.Sprintf("user:%v is no the author of question:%v", e.UserID, e.QuestionID)
}
