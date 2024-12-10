package main

import (
	abi "hello/pkg/core/abi"
	. "hello/pkg/core/model"
	vm "hello/pkg/core/vm"
	"testing"
)

func TestCombinedOperators(t *testing.T) {
	// 인터프리터 초기화
	interpreter := vm.NewInterpreter()
	interpreter.Init("transaction")
	post := &Method{}

	// 테스트 메소드 정의
	methodTest := &Method{
		Parameters: Parameters{
			"price": NewParameter(map[string]interface{}{
				"name":         "price",
				"requirements": true,
			}),
			"quantity": NewParameter(map[string]interface{}{
				"name":         "quantity",
				"requirements": true,
			}),
			"user": NewParameter(map[string]interface{}{
				"name":         "user",
				"requirements": true,
			}),
		},
		Executions: []Execution{
			// 산술 연산과 비교 연산자 조합 테스트
			abi.If(
				abi.Gt(
					abi.Mul(abi.Param("price"), abi.Param("quantity")),
					"50000",
				),
				"대량 주문 할인",
				"일반 주문",
			),

			// 산술 연산 결과를 이용한 조건부 로직
			abi.If(
				abi.And(
					abi.Gte(
						abi.Add(
							abi.Param("price"),
							abi.Get(abi.Param("user"), "point"),
						),
						"10000",
					),
					abi.Eq(abi.Get(abi.Param("user"), "vip"), "true"),
				),
				abi.Sub(abi.Param("price"), "1000"),
				abi.Param("price"),
			),

			// 복합 산술 연산
			abi.Div(
				abi.Mul(
					abi.Add(abi.Param("price"), "1000"),
					abi.Param("quantity"),
				),
				"2",
			),
		},
	}

	// 테스트 데이터 설정
	signedData := NewSignedData()
	signedData.SetAttribute("price", "12000")
	signedData.SetAttribute("quantity", "5")
	signedData.SetAttribute("user", map[string]interface{}{
		"point": "3000",
		"vip":   "true",
	})

	// 실행
	interpreter.Reset()
	interpreter.SetSignedData(signedData)
	interpreter.SetCode(methodTest)
	interpreter.SetPostProcess(post)
	_, err := interpreter.Execute()

	if err != nil {
		t.Errorf("실행 중 오류 발생: %v", err)
	}

	executions := methodTest.GetExecutions()

	// 대량 주문 할인 검증 (12000 * 5 = 60000 > 50000)
	expectedOrder := "대량 주문 할인"
	if executions[0] != expectedOrder {
		t.Errorf("대량 주문 테스트 실패.\n예상: %s\n결과: %v",
			expectedOrder, executions[0])
	}

	// VIP 할인 검증 (12000 + 3000 >= 10000 && vip == true)
	expectedPrice := "11000" // 12000 - 1000
	if executions[1] != expectedPrice {
		t.Errorf("VIP 할인 테스트 실패.\n예상: %s\n결과: %v",
			expectedPrice, executions[1])
	}

	// 복합 산술 연산 검증 ((12000 + 1000) * 5) / 2 = 32500
	expectedResult := "32500"
	if executions[2] != expectedResult {
		t.Errorf("복합 산술 연산 테스트 실패.\n예상: %s\n결과: %v",
			expectedResult, executions[2])
	}

	// 데이터 타입 검증
	if _, ok := executions[0].(string); !ok {
		t.Errorf("첫 번째 실행 결과는 문자열이어야 합니다. 현재 타입: %T", executions[0])
	}

	if _, ok := executions[1].(string); !ok {
		t.Errorf("두 번째 실행 결과는 문자열이어야 합니다. 현재 타입: %T", executions[1])
	}

	if _, ok := executions[2].(string); !ok {
		t.Errorf("세 번째 실행 결과는 문자열이어야 합니다. 현재 타입: %T", executions[2])
	}
}
