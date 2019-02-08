package resolver

import (
	"goask/core/entity"
)

// Question is the GraphQL resolver for Question type.
type Question struct {
	stdResolver
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

func (q Question) Answers() ([]Answer, error) {
	answers, err := q.QuestionDAO.Answers(q.entity.ID)
	return AnswerAll(answers, q.stdResolver), err
}

func (q Question) Author() (User, error) {
	user, err := q.QuestionDAO.GetAuthor(q.entity.ID)
	return UserOne(user, q.stdResolver), err
}

func (q Question) VoteCount() (VoteCount, error) {
	up, down, err := q.QuestionDAO.VoteCount(q.entity.ID)
	return VoteCount{int32(up), int32(down)}, err
}

type VoteCount struct {
	up   int32
	down int32
}

func (v VoteCount) Up() int32 {
	return v.up
}

func (v VoteCount) Down() int32 {
	return v.down
}

// Answer is the GraphQL resolver for Answer type.
type Answer struct {
	stdResolver
	entity entity.Answer
}

func (a Answer) ID() int32 {
	return int32(a.entity.ID)
}

func (a Answer) Content() string {
	return a.entity.Content
}

func (a Answer) Question() (Question, error) {
	question, err := a.QuestionDAO.QuestionByID(a.entity.QuestionID)
	return QuestionOne(question, a.stdResolver), err
}

func (a Answer) Author() (User, error) {
	user, err := a.UserDAO.UserByID(a.entity.AuthorID)
	return UserOne(user, a.stdResolver), err
}

func (a Answer) Accepted() bool {
	return a.entity.Accepted
}

func QuestionOne(question entity.Question, data stdResolver) Question {
	return Question{
		entity:      question,
		stdResolver: data,
	}
}

func QuestionAll(questions []entity.Question, data stdResolver) []Question {
	ret := make([]Question, len(questions))
	for i, question := range questions {
		ret[i] = QuestionOne(question, data)
	}
	return ret
}

func AnswerOne(a entity.Answer, data stdResolver) Answer {
	return Answer{entity: a, stdResolver: data}
}

func AnswerAll(as []entity.Answer, data stdResolver) []Answer {
	answers := make([]Answer, len(as))
	for i, a := range as {
		answers[i] = AnswerOne(a, data)
	}
	return answers
}

type User struct {
	stdResolver
	entity entity.User
}

func (u User) ID() int32 {
	return int32(u.entity.ID)
}

func (u User) Name() string {
	return u.entity.Name
}

func (u User) Questions() ([]Question, error) {
	questions, err := u.UserDAO.QuestionsOfUser(u.entity.ID)
	return QuestionAll(questions, u.stdResolver), err
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
