package adapter

import (
	"fmt"
	"goask/core/entity"
)

type AnswerDAO interface {
	AnswersOfQuestion(QuestionID entity.ID) []entity.Answer
	CreateAnswer(QuestionID entity.ID, Content string, AuthorID entity.ID) (entity.Answer, error)
	AcceptAnswer(AnswerID entity.ID, UserID entity.ID) (entity.Answer, error)
	DeleteAnswer(AnswerID entity.ID, UserID entity.ID) (entity.Answer, error)
}

// ErrQuestionOfAnswerNotFound is a data integrity error.
type ErrQuestionOfAnswerNotFound struct {
	QuestionID entity.ID
	AnswerID   entity.ID
}

func (e *ErrQuestionOfAnswerNotFound) Error() string {
	return fmt.Sprintf("question:%d of answer:%d not found", e.QuestionID, e.AnswerID)
}

// ErrAnswerNotFound
type ErrAnswerNotFound struct {
	ID entity.ID
}

func (e *ErrAnswerNotFound) Error() string {
	return fmt.Sprintf("answer:%d not found", e.ID)
}

// ErrUserIsNotAuthorOfAnswer
type ErrUserIsNotAuthorOfAnswer struct {
	AnswerID entity.ID
	UserID   entity.ID
}

func (e *ErrUserIsNotAuthorOfAnswer) Error() string {
	return fmt.Sprintf("user:%d is not the author of answer:%d found", e.UserID, e.AnswerID)
}
