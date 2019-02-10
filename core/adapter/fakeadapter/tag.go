package fakeadapter

import (
	"goask/core/adapter"
	"goask/core/entity"

	"github.com/pkg/errors"
)

type TagDAO struct {
	data *Data
}

func NewTagDAO(data *Data) *TagDAO {
	return &TagDAO{data}
}

func (t *TagDAO) Questions(tag entity.Tag) (entity.QuestionSet, error) {

	questionIDs := t.data.tags.GetQuestionIDs(tag)

	questions := entity.QuestionSet{}
	for _, id := range questionIDs {
		q, ok := t.data.questions.Get(id)
		if !ok {
			return nil, errors.WithStack(&adapter.ErrQuestionNotFound{ID: id})
		}
		questions.Add(q)
	}
	return questions, nil
}
