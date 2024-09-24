package main

import (
	"testing"
	"hello/pkg/core"
	"fmt"
)

func TestOps0(t *testing.T) {
	var arg0 core.HInteger = 3
	var param = core.NewParamValue("aa", 2)
	_instance := core.Instance()

	a := core.OpAdd(param, arg0)
	res := a.Eval(_instance)
	fmt.Println("a res: ", a.Res)

	if res != 5 {
		t.Errorf("Condition result is not valid")
	}
}
		
func TestOps1(t *testing.T) {
	var arg0 core.HInteger = 3
	var arg1 core.HInteger = 2
	var param = core.NewParamValue("aa", arg1)
	_instance := core.Instance()

	e := core.OpAdd(arg0, param)
	a := core.OpAdd(2, e)
	d := core.OpMul(a, 2)

	c := core.OpEq(14, d)
	e = core.OpCondition(c)

	res := e.Eval(_instance)

	if v, ok := res.(bool); !(ok && v) {
		t.Errorf("Condition result is not valid")
	}
}

func TestOps2(t *testing.T) {
	var arg0 core.HInteger = 3
	var param = core.NewParamValue("aa", 2)
	_instance := core.Instance()

	a := core.OpMul(param, arg0)
	b := core.OpEq(a, 6)
	res := b.Eval(_instance)

	if !(res.(bool)) {
		t.Errorf("Condition result is not valid")
	}
}
