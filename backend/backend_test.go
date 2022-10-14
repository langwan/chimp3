package main

import (
	"testing"
)

func TestRangeConvert(t *testing.T) {
	t.Log(RangeConvert(-1+1, 0, 2, 0, 60))
	t.Log(RangeConvert(0.5+1, 0, 2, 0, 60))
	t.Log(RangeConvert(1+1, 0, 2, 0, 60))
}
