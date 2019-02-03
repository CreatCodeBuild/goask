package adapter

import (
	"fmt"
	"goask/core/entity"
)

type Data interface {
	QuestionDAO
	AnswerDAO
	UserDAO
}

type QuestionDAO interface {
	Questions(search *string) ([]entity.Question, error)
	QuestionByID(ID entity.ID) (entity.Question, error)
	CreateQuestion(post entity.Question) (entity.Question, error)
	UpdateQuestion(post entity.QuestionUpdate) (entity.Question, error)
	DeleteQuestion(userID entity.ID, questionID entity.ID) (entity.Question, error)
}

type ErrQuestionNotFound struct {
	ID entity.ID
}

func (e *ErrQuestionNotFound) Error() string {
	return fmt.Sprintf("question:%d not found", e.ID)
}

type ErrQuestionMutationDenied struct {
	UserID     entity.ID
	QuestionID entity.ID
}

func (e *ErrQuestionMutationDenied) Error() string {
	return fmt.Sprintf("user:%d is not authorized to delete question:%d", e.UserID, e.QuestionID)
}
