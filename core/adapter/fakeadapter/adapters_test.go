package fakeadapter

import (
	"bytes"
	"goask/core/adapter/adaptertest"
	"testing"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	data, err := NewData(buf)
	require.NoError(t, err)
	adaptertest.Data(t, &data)
}
