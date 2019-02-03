package resolver

import (
	"goask/core/adapter"
	"goask/core/entity"
)

type Query struct {
	Data adapter.Data
}

func (q *Query) Questions(args struct{ Search *string }) ([]Question, error) {
	questions, err := q.Data.Questions(args.Search)
	return QuestionAll(questions, q.Data), err
}

func (q *Query) Question(args struct{ ID int32 }) (*Question, error) {
	question, err := q.Data.QuestionByID(entity.ID(args.ID))
	questionResolver := QuestionOne(question, q.Data)
	return &questionResolver, err
}

func (q *Query) GetUser(args struct{ ID int32 }) (*User, error) {
	user, err := q.Data.UserByID(entity.ID(args.ID))
	if err != nil {
		return nil, err
	}
	userResolver := UserOne(user, q.Data)
	return &userResolver, nil
}

func (q *Query) Users() ([]User, error) {
	users, err := q.Data.Users()
	if err != nil {
		return nil, err
	}
	return UserAll(users, q.Data), nil
}
