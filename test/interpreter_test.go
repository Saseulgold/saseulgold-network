package main

import (
	abi "hello/pkg/core/abi"
	. "hello/pkg/core/model"
	vm "hello/pkg/core/vm"
	"testing"
)

func TestInterpreterMethod(t *testing.T) {

	// Create and initialize interpreter
	interpreter := vm.NewInterpreter()
	interpreter.Init("transaction")

	// Create test data
	signedData := NewSignedData()
	signedData.SetAttribute("value", 5)

	post := &Method{Parameters: Parameters{},
		Executions: []Execution{},
	}

	// Create first test method - arithmetic operation test
	method1 := &Method{Parameters: Parameters{
		"value": NewParameter(map[string]interface{}{
			"name":         "value",
			"requirements": true,
			"default":      nil,
		}),
	},
		Executions: []Execution{
			abi.Mul([]interface{}{abi.Param("value"), 2}),
			abi.Add([]interface{}{abi.Param("value"), 10}),
			abi.Div([]interface{}{abi.Param("value"), 5}),
		},
	}

	t.Logf("Method1 executions: %v", method1.GetExecutions())

	// Execute method1
	interpreter.Reset()
	interpreter.SetCode(method1)
	interpreter.SetPostProcess(post)
	result1, err := interpreter.Execute()

	if !err {
		t.Errorf("Error occurred during Method1 execution: %v", err)
	}
	t.Logf("Method1 execution result: %v", result1)

	// Create second test method - conditional statement test
	method2 := &Method{
		Parameters: Parameters{
			"value": NewParameter(map[string]interface{}{
				"name":         "value",
				"requirements": true,
				"default":      nil,
			}),
		},
		Executions: []Execution{
			abi.Condition(abi.Gt(abi.Param("value"), 10), "Value must be greater than 10"),
			abi.If(abi.Lt(abi.Param("value"), 100),
				abi.Mul([]interface{}{abi.Param("value"), 2}),
				abi.Div([]interface{}{abi.Param("value"), 2})),
		},
	}

	// Execute method2
	interpreter.Reset()
	interpreter.SetCode(method2)
	interpreter.SetPostProcess(post)
	result2, err := interpreter.Execute()
	if !err {
		t.Logf("Expected error occurred during Method2 execution: %v", err)
	} else {
		t.Logf("Method2 execution result: %v", result2)
	}

	// Create third test method - complex operation test
	method3 := &Method{
		Parameters: Parameters{
			"value": NewParameter(map[string]interface{}{
				"name":         "value",
				"requirements": true,
				"default":      nil,
			}),
		},
		Executions: []Execution{
			abi.Condition(abi.IsNumeric(abi.Param("value")), "Must be numeric type"),
			abi.PreciseMul(abi.Param("value"), 1.5, 2),
			abi.If(abi.Gt(abi.Param("value"), 50),
				abi.Add([]interface{}{abi.Param("value"), 100}),
				abi.Sub([]interface{}{abi.Param("value"), 50})),
		},
	}

	// Execute method3
	interpreter.Reset()
	interpreter.SetCode(method3)
	interpreter.SetPostProcess(post)
	result3, err := interpreter.Execute()
	if !err {
		t.Errorf("Error occurred during Method3 execution: %v", err)
	}
	t.Logf("Method3 execution result: %v", result3)

}
