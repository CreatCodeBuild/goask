package fakeadapter

import "goask/core/entity"

type TagDAO struct {
	data *Data
}

func NewTagDAO(data *Data) *TagDAO {
	return &TagDAO{data}
}

func (t *TagDAO) Questions(tag entity.Tag) ([]entity.Question, error) {
	// todo
	return nil, nil
}
