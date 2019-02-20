package adapter

import (
	"fmt"
	"goask/core/entity"
)

type AnswerDAO interface {
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
	return fmt.Sprintf("question:%v of answer:%v not found", e.QuestionID, e.AnswerID)
}

// ErrAnswerNotFound
type ErrAnswerNotFound struct {
	ID entity.ID
}

func (e *ErrAnswerNotFound) Error() string {
	return fmt.Sprintf("answer:%v not found", e.ID)
}

// ErrUserIsNotAuthorOfAnswer
type ErrUserIsNotAuthorOfAnswer struct {
	AnswerID entity.ID
	UserID   entity.ID
}

func (e *ErrUserIsNotAuthorOfAnswer) Error() string {
	return fmt.Sprintf("user:%v is not the author of answer:%v found", e.UserID, e.AnswerID)
}
