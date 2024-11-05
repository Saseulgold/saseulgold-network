package main

import (
	"hello/pkg/util"
	"testing"
)

func MemoTest0(t *testing.T) {

	inc := func(var0 interface{}) interface{} {
		return var0.(int) + 1
	}

	f := util.MemoFn0("inc", inc)
	r0 := f(0)

	if r0 != 1 {
		t.Fatalf("Inc result is not valid: %s", r0)
	}

}
