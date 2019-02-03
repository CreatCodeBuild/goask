package fakeadapter

import (
	"goask/core/adapter/adaptertest"
	"testing"
)

func Test(t *testing.T) {
	adaptertest.Data(t, &Data{})
}
