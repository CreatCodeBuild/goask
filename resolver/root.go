package resolver

import (
	"goask/core/adapter"

	"github.com/pkg/errors"
)

type Root struct {
	Query
	Mutation
}

type stdResolver struct {
	QuestionDAO adapter.QuestionDAO
	AnswerDAO   adapter.AnswerDAO
	UserDAO     adapter.UserDAO
	log         logger
}

func NewStdResolver(QuestionDAO adapter.QuestionDAO,
	AnswerDAO adapter.AnswerDAO,
	UserDAO adapter.UserDAO,
	logger logger,
) stdResolver {
	return stdResolver{
		QuestionDAO: QuestionDAO,
		AnswerDAO:   AnswerDAO,
		UserDAO:     UserDAO,
		log:         logger,
	}
}

func (r *stdResolver) check() error {
	if r.QuestionDAO == nil {
		return errors.New("stdResolver.QuestionDAO is not initialized")
	}
	if r.AnswerDAO == nil {
		return errors.New("stdResolver.AnswerDAO is not initialized")
	}
	if r.UserDAO == nil {
		return errors.New("stdResolver.UserDAO is not initialized")
	}
	if r.log == nil {
		return errors.New("stdResolver.log is not initialized")
	}
	return nil
}
