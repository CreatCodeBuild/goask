package adapter

import (
	"fmt"
	"goask/core/entity"
)

// QuestionDAO is the data access object for questions.
type QuestionDAO interface {
	// CRUD
	QuestionByID(ID entity.ID) (entity.Question, error)
	CreateQuestion(post entity.Question, tags []entity.Tag) (entity.Question, error) // todo: change the signature to support tag
	UpdateQuestion(post entity.QuestionUpdate) (entity.Question, error)
	DeleteQuestion(userID entity.ID, questionID entity.ID) (entity.Question, error)
	// Associated Answers
	Answers(questionID entity.ID) ([]entity.Answer, error)
	// Votes
	VoteQuestion(userID entity.ID, questionID entity.ID, voteType entity.VoteType) (entity.Vote, error)
	VoteCount(questionID entity.ID) (up, down int, err error)
	// Associated Author
	GetAuthor(questionID entity.ID) (entity.User, error)
	// Tags
	Tags(questionID entity.ID) (entity.TagSet, error)
}

type Searcher interface {
	Questions(search *string) ([]entity.Question, error)
}

type ErrQuestionNotFound struct {
	ID entity.ID
}

func (e *ErrQuestionNotFound) Error() string {
	return fmt.Sprintf("question:%v not found", e.ID)
}

type ErrQuestionMutationDenied struct {
	UserID     entity.ID
	QuestionID entity.ID
}

func (e *ErrQuestionMutationDenied) Error() string {
	return fmt.Sprintf("user:%v is not authorized to delete question:%v", e.UserID, e.QuestionID)
}
