package resolver

import (
	"goask/core/entity"

	"github.com/graph-gophers/graphql-go"
)

// Question is the GraphQL resolver for Question type.
type Question struct {
	stdResolver
	entity entity.Question
}

func (q Question) ID() graphql.ID {
	return graphql.ID(q.entity.ID)
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

func (q Question) Tags() ([]Tag, error) {
	tags, err := q.QuestionDAO.Tags(q.entity.ID)
	return TagAll(tags.Slice(), q.stdResolver), err
}
