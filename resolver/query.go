package resolver

import (
	"goask/core/entity"
)

type Query struct {
	stdResolver
}

func NewQuery(stdResolver stdResolver) Query {
	return Query{stdResolver}
}

func (q *Query) Questions(args struct{ Search *string }) ([]Question, error) {
	questions, err := q.QuestionDAO.Questions(args.Search)
	return QuestionAll(questions, q.stdResolver), err
}

func (q *Query) Question(args struct{ ID int32 }) (*Question, error) {
	question, err := q.QuestionDAO.QuestionByID(entity.ID(args.ID))
	questionResolver := QuestionOne(question, q.stdResolver)
	return &questionResolver, err
}

func (q *Query) GetUser(args struct{ ID int32 }) (*User, error) {
	user, err := q.UserDAO.UserByID(entity.ID(args.ID))
	if err != nil {
		return nil, err
	}
	userResolver := UserOne(user, q.stdResolver)
	return &userResolver, nil
}

func (q *Query) Users() ([]User, error) {
	users, err := q.UserDAO.Users()
	if err != nil {
		return nil, err
	}
	return UserAll(users, q.stdResolver), nil
}
