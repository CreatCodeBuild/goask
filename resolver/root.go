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
	Searcher    adapter.Searcher
	log         logger
}

func NewStdResolver(QuestionDAO adapter.QuestionDAO,
	AnswerDAO adapter.AnswerDAO,
	UserDAO adapter.UserDAO,
	Searcher adapter.Searcher,
	logger logger,
) (stdResolver, error) {
	std := stdResolver{
		QuestionDAO: QuestionDAO,
		AnswerDAO:   AnswerDAO,
		UserDAO:     UserDAO,
		Searcher:    Searcher,
		log:         logger,
	}
	return std, std.check()
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
