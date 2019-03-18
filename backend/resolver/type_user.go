package resolver

import (
	"goask/core/entity"

	"github.com/graph-gophers/graphql-go"
)

type User struct {
	stdResolver
	entity entity.User
}

func (u User) ID() graphql.ID {
	return graphql.ID(u.entity.ID)
}

func (u User) Name() string {
	return u.entity.Name
}

func (u User) QuestionCount() (int32, error) {
	count, err := u.stdResolver.UserDAO.QuestionCount(u.entity.ID)
	return int32(count), err
}

func (u User) Questions() ([]Question, error) {
	questions, err := u.UserDAO.QuestionsOfUser(u.entity.ID)
	return QuestionAll(questions, u.stdResolver), err
}

func (u User) AnswerCount() (int32, error) {
	count, err := u.stdResolver.UserDAO.AnswerCount(u.entity.ID)
	return int32(count), err
}

func UserOne(user entity.User, stdResolver stdResolver) User {
	return User{entity: user, stdResolver: stdResolver}
}

func UserAll(users []entity.User, data stdResolver) []User {
	ret := make([]User, len(users))
	for i, user := range users {
		ret[i] = UserOne(user, data)
	}
	return ret
}

func ToEntityID(id graphql.ID) entity.ID {
	return entity.NewIDString(string(id))
}
