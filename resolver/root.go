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
	Searcher    adapter.Searcher
	QuestionDAO adapter.QuestionDAO
	AnswerDAO   adapter.AnswerDAO
	UserDAO     adapter.UserDAO
	TagDAO      adapter.TagDAO
	log         logger
}

func NewStdResolver(QuestionDAO adapter.QuestionDAO,
	AnswerDAO adapter.AnswerDAO,
	UserDAO adapter.UserDAO,
	Searcher adapter.Searcher,
	TagDAO adapter.TagDAO,
	logger logger,
) (stdResolver, error) {
	std := stdResolver{
		QuestionDAO: QuestionDAO,
		AnswerDAO:   AnswerDAO,
		UserDAO:     UserDAO,
		Searcher:    Searcher,
		TagDAO:      TagDAO,
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
	if r.Searcher == nil {
		return errors.New("stdResolver.Searcher is not initialized")
	}
	if r.TagDAO == nil {
		return errors.New("stdResolver.TagDAO is not initialized")
	}
	if r.log == nil {
		return errors.New("stdResolver.log is not initialized")
	}
	return nil
}
