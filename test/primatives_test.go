package main

import (
	"testing"
	"hello/pkg/core"
	_ "fmt"
)
		
func TestOps0(t *testing.T) {
	var arg0 float64 = 3
	var param = core.NewParam("aa", 2)

	e := core.Add(param, arg0)
	a := core.Add(2, e)
	d := core.Mul(a, 2)

	c := core.Eq(14, d)
	d = core.Condition(c, 3, 2)
	res := d.Eval()

	if res != 3 {
		t.Errorf("Condition result is not valid")
	}
}

func TestOps1(t *testing.T) {
	var arg0 core.HInteger = 3
	var param = core.NewParam("aa", 2)

	a := core.Mul(arg0, param)
	b := core.Eq(a, 6)
	res := b.Eval()

	if !(res.(bool)) {
		t.Errorf("Condition result is not valid")
	}
}