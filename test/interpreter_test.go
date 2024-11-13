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
	signedData.SetAttribute("value", "5")
	interpreter.SetSignedData(signedData)

	post := &Method{Parameters: Parameters{},
		Executions: []Execution{},
	}

	// Create first test method - arithmetic operation test
	method1 := &Method{Parameters: Parameters{
		"value": NewParameter(map[string]interface{}{
			"name":         "value",
			"requirements": true,
			"default":      "3",
		}),
	},
		Executions: []Execution{
			abi.Add(abi.Add(abi.Param("value"), "2"), "3"),
			abi.If(
				abi.Lt(abi.Param("value"), "10"),
				abi.Add(abi.Param("value"), "5"),
				abi.Mul(abi.Param("value"), "3"),
			),
		},
	}
	t.Logf("Method1 executions: %v", method1.GetExecutions())

	// Execute method1
	interpreter.Reset()
	interpreter.SetSignedData(signedData)
	interpreter.SetCode(method1)
	interpreter.SetPostProcess(post)
	_, err := interpreter.Execute()

	if !err {
		t.Errorf("Error occurred during Method1 execution: %v", err)
	}
	if method1.GetExecutions()[0] != "10" {
		t.Errorf("Method1 execution result error. Expected: 10, Actual: %v", method1.GetExecutions()[0])
	}

	if method1.GetExecutions()[1] != "15" {
		t.Errorf("Method1 execution result error. Expected: 15, Actual: %v", method1.GetExecutions()[1])
	}
}

func TestArithmeticOperators(t *testing.T) {
	interpreter := vm.NewInterpreter()
	interpreter.Init("transaction")

	post := &Method{}

	// OpAdd test
	methodAdd := &Method{
		Parameters: Parameters{
			"value1": NewParameter(map[string]interface{}{
				"name":         "value1",
				"requirements": true,
				"default":      "10",
			}),
			"value2": NewParameter(map[string]interface{}{
				"name":         "value2",
				"requirements": true,
				"default":      "20",
			}),
		},
		Executions: []Execution{
			abi.Add([]interface{}{abi.Param("value1"), abi.Param("value2")}),
		},
	}

	interpreter.Reset()
	interpreter.SetCode(methodAdd)
	interpreter.SetPostProcess(post)
	result, err := interpreter.Execute()
	if !err {
		t.Errorf("Error occurred during OpAdd execution: %v", err)
	}
	if result != "30" {
		t.Errorf("OpAdd result error. Expected: 30, Actual: %v", result)
	}

	// OpSub test
	methodSub := &Method{
		Parameters: Parameters{
			"value1": NewParameter(map[string]interface{}{
				"name":         "value1",
				"requirements": true,
				"default":      "30",
			}),
			"value2": NewParameter(map[string]interface{}{
				"name":         "value2",
				"requirements": true,
				"default":      "10",
			}),
		},
		Executions: []Execution{
			abi.Sub([]interface{}{abi.Param("value1"), abi.Param("value2")}),
		},
	}

	interpreter.Reset()
	interpreter.SetCode(methodSub)
	interpreter.SetPostProcess(post)
	result, err = interpreter.Execute()
	if !err {
		t.Errorf("Error occurred during OpSub execution: %v", err)
	}
	if result != "20" {
		t.Errorf("OpSub result error. Expected: 20, Actual: %v", result)
	}

	// OpMul test
	methodMul := &Method{
		Parameters: Parameters{
			"value1": NewParameter(map[string]interface{}{
				"name":         "value1",
				"requirements": true,
				"default":      "5",
			}),
			"value2": NewParameter(map[string]interface{}{
				"name":         "value2",
				"requirements": true,
				"default":      "4",
			}),
		},
		Executions: []Execution{
			abi.Mul([]interface{}{abi.Param("value1"), abi.Param("value2")}),
		},
	}

	interpreter.Reset()
	interpreter.SetCode(methodMul)
	interpreter.SetPostProcess(post)
	result, err = interpreter.Execute()
	if !err {
		t.Errorf("Error occurred during OpMul execution: %v", err)
	}
	if result != "20" {
		t.Errorf("OpMul result error. Expected: 20, Actual: %v", result)
	}

	// OpDiv test
	methodDiv := &Method{
		Parameters: Parameters{
			"value1": NewParameter(map[string]interface{}{
				"name":         "value1",
				"requirements": true,
				"default":      "20",
			}),
			"value2": NewParameter(map[string]interface{}{
				"name":         "value2",
				"requirements": true,
				"default":      "5",
			}),
		},
		Executions: []Execution{
			abi.Div([]interface{}{abi.Param("value1"), abi.Param("value2")}),
		},
	}

	interpreter.Reset()
	interpreter.SetCode(methodDiv)
	interpreter.SetPostProcess(post)
	result, err = interpreter.Execute()
	if !err {
		t.Errorf("Error occurred during OpDiv execution: %v", err)
	}
	if result != "4" {
		t.Errorf("OpDiv result error. Expected: 4, Actual: %v", result)
	}

	// OpPreciseAdd test
	methodPreciseAdd := &Method{
		Parameters: Parameters{
			"value1": NewParameter(map[string]interface{}{
				"name":         "value1",
				"requirements": true,
				"default":      "10.5",
			}),
			"value2": NewParameter(map[string]interface{}{
				"name":         "value2",
				"requirements": true,
				"default":      "20.3",
			}),
		},
		Executions: []Execution{
			abi.PreciseAdd(abi.Param("value1"), abi.Param("value2"), 2),
		},
	}

	interpreter.Reset()
	interpreter.SetCode(methodPreciseAdd)
	interpreter.SetPostProcess(post)
	result, err = interpreter.Execute()
	if !err {
		t.Errorf("Error occurred during OpPreciseAdd execution: %v", err)
	}
	if result != "30.80" {
		t.Errorf("OpPreciseAdd result error. Expected: 30.80, Actual: %v", result)
	}

	// OpPreciseSub test
	methodPreciseSub := &Method{
		Parameters: Parameters{
			"value1": NewParameter(map[string]interface{}{
				"name":         "value1",
				"requirements": true,
				"default":      "30.5",
			}),
			"value2": NewParameter(map[string]interface{}{
				"name":         "value2",
				"requirements": true,
				"default":      "10.2",
			}),
		},
		Executions: []Execution{
			abi.PreciseSub(abi.Param("value1"), abi.Param("value2"), 2),
		},
	}

	interpreter.Reset()
	interpreter.SetCode(methodPreciseSub)
	interpreter.SetPostProcess(post)
	result, err = interpreter.Execute()
	if !err {
		t.Errorf("Error occurred during OpPreciseSub execution: %v", err)
	}
	if result != "20.30" {
		t.Errorf("OpPreciseSub result error. Expected: 20.30, Actual: %v", result)
	}

	// OpPreciseMul test
	methodPreciseMul := &Method{
		Parameters: Parameters{
			"value1": NewParameter(map[string]interface{}{
				"name":         "value1",
				"requirements": true,
				"default":      "5.5",
			}),
			"value2": NewParameter(map[string]interface{}{
				"name":         "value2",
				"requirements": true,
				"default":      "2.0",
			}),
		},
		Executions: []Execution{
			abi.PreciseMul(abi.Param("value1"), abi.Param("value2"), 2),
		},
	}

	interpreter.Reset()
	interpreter.SetCode(methodPreciseMul)
	interpreter.SetPostProcess(post)
	result, err = interpreter.Execute()
	if !err {
		t.Errorf("Error occurred during OpPreciseMul execution: %v", err)
	}
	if result != "11.00" {
		t.Errorf("OpPreciseMul result error. Expected: 11.00, Actual: %v", result)
	}

	// OpPreciseDiv test
	methodPreciseDiv := &Method{
		Parameters: Parameters{
			"value1": NewParameter(map[string]interface{}{
				"name":         "value1",
				"requirements": true,
				"default":      "10.5",
			}),
			"value2": NewParameter(map[string]interface{}{
				"name":         "value2",
				"requirements": true,
				"default":      "2.0",
			}),
		},
		Executions: []Execution{
			abi.PreciseDiv(abi.Param("value1"), abi.Param("value2"), 2),
		},
	}

	interpreter.Reset()
	interpreter.SetCode(methodPreciseDiv)
	interpreter.SetPostProcess(post)
	result, err = interpreter.Execute()
	if !err {
		t.Errorf("Error occurred during OpPreciseDiv execution: %v", err)
	}
	if result != "5.25" {
		t.Errorf("OpPreciseDiv result error. Expected: 5.25, Actual: %v", result)
	}

	// OpScale test
	methodScale := &Method{
		Parameters: Parameters{
			"value": NewParameter(map[string]interface{}{
				"name":         "value",
				"requirements": true,
				"default":      "10.505",
			}),
		},
		Executions: []Execution{
			abi.Scale(abi.Param("value")),
		},
	}

	interpreter.Reset()
	interpreter.SetCode(methodScale)
	interpreter.SetPostProcess(post)
	result, err = interpreter.Execute()
	if !err {
		t.Errorf("Error occurred during OpScale execution: %v", err)
	}
	if result != 3 {
		t.Errorf("OpScale result error. Expected: 3, Actual: %v", result)
	}
}
