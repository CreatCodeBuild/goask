package fakeadapter

import (
	"goask/core/adapter/adaptertest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	data, err := NewData(BufferSerializer{})
	require.NoError(t, err)
	adaptertest.Data(t, data)
}
