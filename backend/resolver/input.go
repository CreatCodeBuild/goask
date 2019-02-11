package resolver

import (
	"goask/core/entity"
)

type QuestionInput struct {
	entity.QuestionUpdate
	ID int32
}

type AnswerCreationInput struct {
	QuestionID int32
	Content    string
}
