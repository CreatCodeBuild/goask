package resolver

import (
	"goask/core/entity"

	"github.com/graph-gophers/graphql-go"
)

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

func (a Answer) ID() graphql.ID {
	return graphql.ID(a.entity.ID)
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
