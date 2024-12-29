package main

import (
	abi "hello/pkg/core/abi"
	. "hello/pkg/core/model"
	vm "hello/pkg/core/vm"
	"testing"
)

func TestBasicOperators(t *testing.T) {
	// 인터프리터 초기화
	interpreter := vm.NewInterpreter()
	interpreter.Init("transaction")
	post := &Method{}

	// 테스트 메소드 정의
	methodTest := &Method{
		Parameters: Parameters{
			"amount": NewParameter(map[string]interface{}{
				"name":         "amount",
				"requirements": true,
			}),
			"user": NewParameter(map[string]interface{}{
				"name":         "user",
				"requirements": true,
			}),
			"items": NewParameter(map[string]interface{}{
				"name":         "items",
				"requirements": true,
			}),
		},
		Executions: []Execution{
			// 기본 비교 연산자 테스트
			abi.And(
				abi.Gt(abi.Param("amount"), "1000"),
				abi.Lt(abi.Param("amount"), "10000"),
			),

			// 중첩된 조건문 테스트
			abi.If(
				abi.Or(
					abi.And(
						abi.Gte(abi.Get(abi.Param("user"), "age", nil), "20"),
						abi.Lt(abi.Get(abi.Param("user"), "age", nil), "30"),
					),
					abi.Eq(abi.Get(abi.Param("user"), "vip", nil), "true"),
				),
				"할인 적용",
				"정상가",
			),

			// 중첩 연산자 테스트
			abi.If(
				abi.And(
					abi.Or(
						abi.Eq(abi.Get(abi.Param("user"), "membership", nil), "GOLD"),
						abi.Gte(abi.Param("amount"), "5000"),
					),
				),
				"사은품 증정",
				"일반 배송",
			),

			// Get 연산자 테스트
			abi.Get(abi.Param("user"), "email", nil),
		},
	}

	// 테스트 데이터 설정
	signedData := NewSignedData()
	signedData.SetAttribute("amount", "5500")
	signedData.SetAttribute("user", map[string]interface{}{
		"age":        "25",
		"vip":        "false",
		"membership": "GOLD",
		"email":      "test@company.com",
	})
	signedData.SetAttribute("items", map[string]interface{}{
		"list": []interface{}{"item1", "item2"},
	})

	// 실행
	interpreter.Reset(true)
	interpreter.SetSignedData(signedData)
	interpreter.SetCode(methodTest)
	interpreter.SetPostProcess(post)
	_, err := interpreter.Execute()

	if err != nil {
		t.Errorf("Error during execution: %v", err)
	}

	executions := methodTest.GetExecutions()

	// Validate basic comparison operator results
	expectedAmount := "5500"
	if executions[0] != true {
		t.Errorf("Basic operator test failed. Amount %s should be between 1000 and 10000", expectedAmount)
	}

	// Validate nested conditional results
	expectedDiscount := "할인 적용" // Regular price
	if executions[1] != expectedDiscount {
		t.Errorf("Nested conditional test failed.\nExpected: %s\nGot: %v\nUser age: 25, VIP: false",
			expectedDiscount, executions[1])
	}

	// Validate complex operator results
	expectedGift := "사은품 증정" // Gift included
	if executions[2] != expectedGift {
		t.Errorf("Complex operator test failed.\nExpected: %s\nGot: %v\nMembership: GOLD, Amount: %s",
			expectedGift, executions[2], expectedAmount)
	}

	// Validate Get operator results
	expectedEmail := "test@company.com"
	if executions[3] != expectedEmail {
		t.Errorf("Get operator test failed.\nExpected: %s\nGot: %v",
			expectedEmail, executions[3])
	}

	// Additional validation for data types
	if _, ok := executions[0].(bool); !ok {
		t.Errorf("First execution result should be boolean, got %T", executions[0])
	}

	if _, ok := executions[1].(string); !ok {
		t.Errorf("Second execution result should be string, got %T", executions[1])
	}

	if _, ok := executions[2].(string); !ok {
		t.Errorf("Third execution result should be string, got %T", executions[2])
	}

	if _, ok := executions[3].(string); !ok {
		t.Errorf("Fourth execution result should be string, got %T", executions[3])
	}
}
