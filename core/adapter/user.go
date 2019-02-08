package adapter

import (
	"fmt"
	"goask/core/entity"
)

type UserDAO interface {
	UserByID(ID entity.ID) (entity.User, error)
	Users() ([]entity.User, error)
	CreateUser(name string) (entity.User, error)
	QuestionsOfUser(UserID entity.ID) ([]entity.Question, error)
}

type ErrUserNotFound struct {
	ID entity.ID
}

func (e *ErrUserNotFound) Error() string {
	return fmt.Sprintf("user:%d not found", e.ID)
}

type ErrUserIsNotAuthorOfQuestion struct {
	UserID     entity.ID
	QuestionID entity.ID
}

func (e *ErrUserIsNotAuthorOfQuestion) Error() string {
	return fmt.Sprintf("user:%d is no the author of question:%d", e.UserID, e.QuestionID)
}
