package adapter

import (
	"goask/core/entity"
)

type TagDAO interface {
	// Questions gets all questions which has this tag.
	Questions(entity.Tag) (entity.QuestionSet, error)
}
