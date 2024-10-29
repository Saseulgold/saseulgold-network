package main

import (
	"hello/pkg/core/util"
	"testing"
)

func MemoTest0(t *testing.T) {

	inc := func (var0) {
		return var0 + 1
	}

	f := util.MemoFn0(inc)
	r0 := f(0)

	if r0 != 1 {
		t.Fatalf("Inc result is not valid: %s", r0)
	}

}