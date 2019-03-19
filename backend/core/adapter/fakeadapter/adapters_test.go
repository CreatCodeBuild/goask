package fakeadapter

import (
	"goask/core/adapter"
	"goask/core/adapter/adaptertest"
	"goask/id"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	adaptertest.All(t, NewDAO)
	adaptertest.Delete(t, NewDAO)
}

func NewDAO(t *testing.T) (questionDAO adapter.QuestionDAO, answerDAO adapter.AnswerDAO, userDAO adapter.UserDAO, tagDAO adapter.TagDAO) {
	data, err := NewData(BufferSerializer{})
	require.NoError(t, err)
	ud := &UserDAO{data: data}
	userDAO = ud
	questionDAO = NewQuestionDAO(data, ud, id.NewGenerator())
	answerDAO = &AnswerDAO{data: data}
	tagDAO = &TagDAO{data: data}
	return
}
