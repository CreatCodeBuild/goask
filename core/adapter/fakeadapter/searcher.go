package fakeadapter

import (
	"goask/core/entity"
)

type Searcher struct {
	data *Data
}

func NewSearcher(data *Data) *Searcher {
	return &Searcher{data}
}

func (d *Searcher) Questions(search *string) ([]entity.Question, error) {
	if search == nil {
		return d.data.questions, nil
	}
	ret := make([]entity.Question, 0)
	for _, q := range d.data.questions {
		if match(q.Content, *search) {
			ret = append(ret, q)
		}
	}
	return ret, nil
}
