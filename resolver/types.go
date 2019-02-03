package resolver

import (
	"goask/core/adapter"
	"goask/core/entity"
)

// Question is the GraphQL resolver for Question type.
type Question struct {
	data   adapter.Data // todo: let's think about the scope control here
	entity entity.Question
}

func (q Question) ID() int32 {
	return int32(q.entity.ID)
}

func (q Question) Title() string {
	return string(q.entity.Title)
}

func (q Question) Content() string {
	return string(q.entity.Content)
}

func (q Question) Answers() []Answer {
	answers := q.data.AnswersOfQuestion(q.entity.ID)
	return AnswerAll(answers, q.data)
}

func (q Question) Author() (User, error) {
	user, err := q.data.UserByID(q.entity.AuthorID)
	return UserOne(user, q.data), err
}

// Answer is the GraphQL resolver for Answer type.
type Answer struct {
	data   adapter.Data
	entity entity.Answer
}

func (a Answer) ID() int32 {
	return int32(a.entity.ID)
}

func (a Answer) Content() string {
	return a.entity.Content
}

func (a Answer) Question() (Question, error) {
	question, err := a.data.QuestionByID(a.entity.QuestionID)
	return QuestionOne(question, a.data), err
}

func (a Answer) Author() (User, error) {
	user, err := a.data.UserByID(a.entity.AuthorID)
	return UserOne(user, a.data), err
}

func (a Answer) Accepted() bool {
	return a.entity.Accepted
}

func QuestionOne(question entity.Question, data adapter.Data) Question {
	return Question{
		entity: question,
		data:   data,
	}
}

func QuestionAll(questions []entity.Question, data adapter.Data) []Question {
	ret := make([]Question, len(questions))
	for i, question := range questions {
		ret[i] = QuestionOne(question, data)
	}
	return ret
}

func AnswerOne(a entity.Answer, data adapter.Data) Answer {
	return Answer{entity: a, data: data}
}

func AnswerAll(as []entity.Answer, data adapter.Data) []Answer {
	answers := make([]Answer, len(as))
	for i, a := range as {
		answers[i] = AnswerOne(a, data)
	}
	return answers
}

type User struct {
	entity entity.User
	data   adapter.Data
}

func (u User) ID() int32 {
	return int32(u.entity.ID)
}

func (u User) Name() string {
	return u.entity.Name
}

func (u User) Questions() ([]Question, error) {
	questions, err := u.data.QuestionsByUserID(u.entity.ID)
	return QuestionAll(questions, u.data), err
}

func UserOne(user entity.User, data adapter.Data) User {
	return User{entity: user, data: data}
}

func UserAll(users []entity.User, data adapter.Data) []User {
	ret := make([]User, len(users))
	for i, user := range users {
		ret[i] = UserOne(user, data)
	}
	return ret
}
