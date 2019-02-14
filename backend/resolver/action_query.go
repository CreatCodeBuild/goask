package resolver

import (
	"goask/core/entity"
)

type QueryAction struct {
	stdResolver
	userSession UserSession
}

func (q QueryAction) Questions(args struct{ Search *string }) ([]Question, error) {
	questions, err := q.Searcher.Questions(args.Search)
	return QuestionAll(questions, q.stdResolver), err
}

func (q QueryAction) Question(args struct{ ID int32 }) (*Question, error) {
	question, err := q.QuestionDAO.QuestionByID(entity.ID(args.ID))
	questionResolver := QuestionOne(question, q.stdResolver)
	return &questionResolver, err
}

func (q QueryAction) GetUser(args struct{ ID int32 }) (*User, error) {
	user, err := q.UserDAO.UserByID(entity.ID(args.ID))
	if err != nil {
		return nil, err
	}
	userResolver := UserOne(user, q.stdResolver)
	return &userResolver, q.check()
}

func (q QueryAction) Users() ([]User, error) {
	users, err := q.UserDAO.Users()
	if err != nil {
		return nil, err
	}
	return UserAll(users, q.stdResolver), nil
}
