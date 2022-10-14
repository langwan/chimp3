package main

import (
	helperGrpc "github.com/langwan/langgo/helpers/grpc"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRangeConvert(t *testing.T) {
	t.Log(RangeConvert(-1+1, 0, 2, 0, 60))
	t.Log(RangeConvert(0.5+1, 0, 2, 0, 60))
	t.Log(RangeConvert(1+1, 0, 2, 0, 60))
}

func TestCall(t *testing.T) {
	response, code, err := helperGrpc.Call(&backend, "Next", "", nil)
	assert.NoError(t, err)
	assert.Equal(t, code, 0)
	t.Log(response)
}
func TestCall2(t *testing.T) {
	response, code, err := helperGrpc.Call(&backend, "Playing", "{is_play:1}", nil)
	assert.NoError(t, err)
	assert.Equal(t, code, 0)
	t.Log(response)
}
