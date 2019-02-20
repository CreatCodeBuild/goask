package resolver

import (
	"goask/core/entity"

	"github.com/graph-gophers/graphql-go"
)

type QuestionInput struct {
	entity.QuestionUpdate
	ID graphql.ID
}

type AnswerCreationInput struct {
	QuestionID graphql.ID
	Content    string
}
