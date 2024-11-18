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
			abi.Div(
				abi.Mul(
					abi.Sub(abi.Param("value"), "2"),
					"4",
				),
				"2",
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

	if method1.GetExecutions()[1] != "10" {
		t.Errorf("Method1 execution result error. Expected: 10, Actual: %v", method1.GetExecutions()[1])
	}

	if method1.GetExecutions()[2] != "6" {
		t.Errorf("Method1 execution result error. Expected: 6, Actual: %v", method1.GetExecutions()[2])
	}
}

func TestArithmeticOperators(t *testing.T) {
	interpreter := vm.NewInterpreter()
	interpreter.Init("transaction")

	post := &Method{}

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
			abi.PreciseSub(abi.Param("value1"), abi.Param("value2"), 2),
			// value1 = 30, value2 = 10
			// 1. value1 * 0.5 = 30 * 0.5 = 15.000 (scale=3)
			// 2. value2 / 2 = 10 / 2 = 5.000 (scale=3)
			// 3. 15.000 + 5.000 = 20.00 (scale=2)
			abi.PreciseAdd(
				abi.PreciseMul(abi.Param("value1"), "0.5", 3),
				abi.PreciseDiv(abi.Param("value2"), "2", 3),
				2,
			),
			// 복잡한 연산 추가
			// 1. value1 * 1.5 = 30 * 1.5 = 45.000 (scale=3)
			// 2. value2 * 0.8 = 10 * 0.8 = 8.000 (scale=3)
			// 3. if (45.000 > 40) then (45.000 - 8.000) else (45.000 + 8.000)
			// 4. 45.000 - 8.000 = 37.000 (scale=3)
			// 5. 37.000 * 2 = 74 (scale=0)
			abi.PreciseMul(
				abi.If(
					abi.Gt(
						abi.PreciseMul(abi.Param("value1"), "1.5", 3),
						"40",
					),
					abi.PreciseSub(
						abi.PreciseMul(abi.Param("value1"), "1.5", 3),
						abi.PreciseMul(abi.Param("value2"), "0.8", 3),
						3,
					),
					abi.PreciseAdd(
						abi.PreciseMul(abi.Param("value1"), "1.5", 3),
						abi.PreciseMul(abi.Param("value2"), "0.8", 3),
						3,
					),
				),
				"2",
				0,
			),
		},
	}

	signedData := NewSignedData()
	signedData.SetAttribute("value1", "30")
	signedData.SetAttribute("value2", "10")

	interpreter.Reset()
	interpreter.SetSignedData(signedData)
	interpreter.SetCode(methodSub)
	interpreter.SetPostProcess(post)
	_, err := interpreter.Execute()

	if !err {
		t.Errorf("Error occurred during execution: %v", err)
	}

	// 첫번째 실행 결과 검증 (30 - 10 = 20)
	if methodSub.GetExecutions()[0] != "20.00" {
		t.Errorf("First execution result error. Expected: 20, Actual: %v", methodSub.GetExecutions()[0])
	}

	// 두번째 실행 결과 검증 ((30 * 0.5) + (10 / 2) = 15 + 5 = 20)
	if methodSub.GetExecutions()[1] != "20.00" {
		t.Errorf("Second execution result error. Expected: 20.00, Actual: %v", methodSub.GetExecutions()[1])
	}

	// 세번째 실행 결과 검증 ((30 * 1.5) - (10 * 0.8) = 45 - 8 = 37)
	if methodSub.GetExecutions()[2] != "74" {
		t.Errorf("Third execution result error. Expected: 74, Actual: %v", methodSub.GetExecutions()[2])
	}
}

func TestLogicalOperators(t *testing.T) {
	interpreter := vm.NewInterpreter()
	interpreter.Init("transaction")

	post := &Method{}

	methodLogical := &Method{
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
			// AND 연산자 테스트
			abi.And(
				abi.Gt(abi.Param("value1"), "5"),
				abi.Lt(abi.Param("value2"), "30"),
			),
			// OR 연산자 테스트
			abi.Or(
				abi.Eq(abi.Param("value1"), "5"),
				abi.Gt(abi.Param("value2"), "15"),
			),
			// NOT 연산자 테스트
			// 복합 논리 연산 테스트
			abi.And(
				abi.Or(
					abi.Gt(abi.Param("value1"), "5"),
					abi.Lt(abi.Param("value2"), "10"),
				),
				abi.Eq(abi.Param("value1"), abi.Param("value2")),
			),
		},
	}

	signedData := NewSignedData()
	signedData.SetAttribute("value1", "10")
	signedData.SetAttribute("value2", "20")

	interpreter.Reset()
	interpreter.SetSignedData(signedData)
	interpreter.SetCode(methodLogical)
	interpreter.SetPostProcess(post)
	_, err := interpreter.Execute()

	if !err {
		t.Errorf("Error occurred during execution: %v", err)
	}

	// 첫번째 실행 결과 검증 (10 > 5 && 20 < 30 = true)
	if methodLogical.GetExecutions()[0] != true {
		t.Errorf("First execution result error. Expected: true, Actual: %v", methodLogical.GetExecutions()[0])
	}

	// 두번째 실행 결과 검증 (10 == 5 || 20 > 15 = true)
	if methodLogical.GetExecutions()[1] != true {
		t.Errorf("Second execution result error. Expected: true, Actual: %v", methodLogical.GetExecutions()[1])
	}

	// 세번째 실행 결과 검증 ((10 > 5 || 20 < 10) && (10 == 20) = false)
	if methodLogical.GetExecutions()[2] != false {
		t.Errorf("Third execution result error. Expected: false, Actual: %v", methodLogical.GetExecutions()[2])
	}
}

func TestConditionOperators(t *testing.T) {
	interpreter := vm.NewInterpreter()
	interpreter.Init("transaction")

	post := &Method{}

	methodCondition := &Method{
		Parameters: Parameters{
			"amount": NewParameter(map[string]interface{}{
				"name":         "amount",
				"requirements": true,
				"default":      "10",
			}),
			"rate": NewParameter(map[string]interface{}{
				"name":         "rate",
				"requirements": true,
				"default":      "0.1",
			}),
		},
		Executions: []Execution{
			abi.Condition(
				abi.Gt(abi.Param("amount"), "50"),
				"Condition:amount not gt 50",
			),
			abi.If(
				abi.Gt(abi.Param("amount"), "50"),
				abi.PreciseMul(abi.Param("amount"), "2", 2),
				abi.PreciseMul(abi.Param("amount"), "0.5", 2),
			),

			abi.If(
				abi.Gt(abi.Param("amount"), "200"),
				abi.PreciseMul(
					abi.PreciseMul(abi.Param("amount"), abi.Param("rate"), 3),
					"2",
					2,
				),
				abi.If(
					abi.Gt(abi.Param("amount"), "100"),
					abi.PreciseMul(abi.Param("amount"), abi.Param("rate"), 2),
					abi.PreciseMul(
						abi.Param("amount"),
						abi.PreciseMul(abi.Param("rate"), "0.5", 3),
						2,
					),
				),
			),
		},
	}

	signedData := NewSignedData()
	signedData.SetAttribute("amount", "10")
	signedData.SetAttribute("rate", "0.1")

	interpreter.Reset()
	interpreter.SetSignedData(signedData)
	interpreter.SetCode(methodCondition)
	interpreter.SetPostProcess(post)
	errMsg, _ := interpreter.Execute()

	if errMsg == "" {
		t.Errorf("Execution should be broken")
	}

	t.Logf("Error message: %v", errMsg)
}
