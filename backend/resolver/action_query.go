package resolver

import (
	"github.com/graph-gophers/graphql-go"
)

type QueryAction struct {
	stdResolver
	userSession UserSession
}

func (q QueryAction) Questions(args struct{ Search *string }) ([]Question, error) {
	questions, err := q.Searcher.Questions(args.Search)
	return QuestionAll(questions, q.stdResolver), err
}

func (q QueryAction) Question(args struct{ ID graphql.ID }) (*Question, error) {
	question, err := q.QuestionDAO.QuestionByID(ToEntityID(args.ID))
	questionResolver := QuestionOne(question, q.stdResolver)
	return &questionResolver, err
}

func (q QueryAction) User(args struct{ ID graphql.ID }) (*User, error) {
	user, err := q.UserDAO.UserByID(ToEntityID(args.ID))
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
