package fakeadapter

import (
	"goask/core/adapter/adaptertest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	data, err := NewData(BufferSerializer{})
	require.NoError(t, err)

	userDAO := UserDAO{data: data}
	questionDAO := QuestionDAO{data: data, userDAO: &userDAO}
	answerDAO := AnswerDAO{data: data}
	adaptertest.All(t, &questionDAO, &answerDAO, &userDAO)
}
