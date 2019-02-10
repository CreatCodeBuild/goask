package resolver

import (
	"goask/core/entity"
)

type Tag struct {
	stdResolver
	entity entity.Tag
}

func (t Tag) Value() string {
	return string(t.entity)
}

func (t Tag) Questions() ([]Question, error) {
	questions, err := t.stdResolver.TagDAO.Questions(t.entity)
	if err != nil {
		t.stdResolver.log.Error(err)
		return nil, err
	}
	return QuestionAll(questions.Slice(), t.stdResolver), nil
}

func TagOne(tag entity.Tag, stdResolver stdResolver) Tag {
	return Tag{
		entity:      tag,
		stdResolver: stdResolver,
	}
}

func TagAll(tags []entity.Tag, stdResolver stdResolver) []Tag {
	ret := make([]Tag, len(tags))
	for i, tag := range tags {
		ret[i] = TagOne(tag, stdResolver)
	}
	return ret
}
