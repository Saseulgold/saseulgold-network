package main

import (
	_ "fmt"
	. "hello/pkg/core"
	. "hello/pkg/core/vm"
	"testing"
)

func GetContract0() Contract {
	return NewContract()
}

func TestCase0(t *testing.T) {
	contract := NewContract()
	contract.SetName("test0")
	contract.SetVersion("0.0.1")
	contract.SetMachine("m")

	param0 := contract.AddParameter("amount0", IntegerFlag)

	const arg0 HInteger = 1
	const arg1 HInteger = 2

	addOp0 := OpAdd(param0, arg0)
	addOp1 := OpAdd(addOp0, arg1)

	contract.AddExecution(addOp1)

	var const0 HInteger = 1
	eq := OpEq(const0, addOp1)

	cond0 := OpCondition(eq)
	contract.AddExecution(cond0)

	instance := Instance()

	params := ParamValueMap{"amount0": HInteger(1)}
	instance.ExecuteContract(&contract, params)
}
