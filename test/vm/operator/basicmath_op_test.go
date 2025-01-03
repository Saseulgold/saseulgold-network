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
							abi.Get(abi.Param("user"), "point", nil),
						),
						"10000",
					),
					abi.Eq(abi.Get(abi.Param("user"), "vip", nil), "true"),
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
	interpreter.Reset(true)
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

func TestMinMaxOperators(t *testing.T) {
	// Initialize interpreter
	interpreter := vm.NewInterpreter()
	interpreter.Init("transaction")
	post := &Method{}

	// Test method definition
	methodTest := &Method{
		Parameters: Parameters{
			"hash1": NewParameter(map[string]interface{}{
				"name":         "hash1",
				"requirements": true,
			}),
			"hash2": NewParameter(map[string]interface{}{
				"name":         "hash2",
				"requirements": true,
			}),
			"num1": NewParameter(map[string]interface{}{
				"name":         "num1",
				"requirements": true,
			}),
			"num2": NewParameter(map[string]interface{}{
				"name":         "num2",
				"requirements": true,
			}),
		},
		Executions: []Execution{
			// Test Min with hashes
			abi.Min(abi.Param("hash1"), abi.Param("hash2")),

			// Test Max with hashes
			abi.Max(abi.Param("hash1"), abi.Param("hash2")),

			// Test Min with numbers
			abi.Min(abi.Param("num1"), abi.Param("num2")),

			// Test Max with numbers
			abi.Max(abi.Param("num1"), abi.Param("num2")),
		},
	}

	// First test with original order
	signedData := NewSignedData()
	signedData.SetAttribute("hash1", "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	signedData.SetAttribute("hash2", "fedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321")
	signedData.SetAttribute("num1", "100")
	signedData.SetAttribute("num2", "200")

	// Execute
	interpreter.Reset(true)
	interpreter.SetSignedData(signedData)
	interpreter.SetCode(methodTest)
	interpreter.SetPostProcess(post)
	_, err := interpreter.Execute()

	if err != nil {
		t.Errorf("Execution error: %v", err)
	}

	executions := methodTest.GetExecutions()

	// Verify Min hash result
	expectedMinHash := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
	if executions[0] != expectedMinHash {
		t.Errorf("Min hash test failed.\nExpected: %s\nGot: %v",
			expectedMinHash, executions[0])
	}

	// Verify Max hash result
	expectedMaxHash := "fedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321"
	if executions[1] != expectedMaxHash {
		t.Errorf("Max hash test failed.\nExpected: %s\nGot: %v",
			expectedMaxHash, executions[1])
	}

	// Verify Min number result
	expectedMinNum := "100"
	if executions[2] != expectedMinNum {
		t.Errorf("Min number test failed.\nExpected: %s\nGot: %v",
			expectedMinNum, executions[2])
	}

	// Verify Max number result
	expectedMaxNum := "200"
	if executions[3] != expectedMaxNum {
		t.Errorf("Max number test failed.\nExpected: %s\nGot: %v",
			expectedMaxNum, executions[3])
	}

	// Test with swapped hash values
	signedData2 := NewSignedData()
	signedData2.SetAttribute("hash1", "fedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321")
	signedData2.SetAttribute("hash2", "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	signedData2.SetAttribute("num1", "100")
	signedData2.SetAttribute("num2", "200")

	// Execute with swapped values
	interpreter.Reset(true)
	interpreter.SetSignedData(signedData2)
	interpreter.SetCode(methodTest)
	interpreter.SetPostProcess(post)
	_, err = interpreter.Execute()

	if err != nil {
		t.Errorf("Execution error with swapped values: %v", err)
	}

	executionsSwapped := methodTest.GetExecutions()

	// Verify Min hash result with swapped values
	expectedMinHash = "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
	if executionsSwapped[0] != expectedMinHash {
		t.Errorf("Min hash test with swapped values failed.\nExpected: %s\nGot: %v",
			expectedMinHash, executionsSwapped[0])
	}

	// Verify Max hash result with swapped values
	expectedMaxHash = "fedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321"
	if executionsSwapped[1] != expectedMaxHash {
		t.Errorf("Max hash test with swapped values failed.\nExpected: %s\nGot: %v",
			expectedMaxHash, executionsSwapped[1])
	}
}

func TestPowOperator(t *testing.T) {
	// Initialize interpreter
	interpreter := vm.NewInterpreter()
	interpreter.Init("transaction")
	post := &Method{}

	// Test method definition
	methodTest := &Method{
		Parameters: Parameters{
			"base": NewParameter(map[string]interface{}{
				"name":         "base",
				"requirements": true,
			}),
			"exponent": NewParameter(map[string]interface{}{
				"name":         "exponent",
				"requirements": true,
			}),
		},
		Executions: []Execution{
			// Test basic power operation
			abi.PrecisePow(abi.Param("base"), abi.Param("exponent"), "0"),

			// Test power with zero exponent
			abi.PrecisePow(abi.Param("base"), "0", "0"),

			// Test power with negative exponent
			abi.PrecisePow(abi.Param("base"), "-2", "0"),
		},
	}

	// Test data setup
	signedData := NewSignedData()
	signedData.SetAttribute("base", "2")
	signedData.SetAttribute("exponent", "3")

	// Execute
	interpreter.Reset(true)
	interpreter.SetSignedData(signedData)
	interpreter.SetCode(methodTest)
	interpreter.SetPostProcess(post)
	_, err := interpreter.Execute()

	if err != nil {
		t.Errorf("Execution error: %v", err)
	}

	executions := methodTest.GetExecutions()

	// Verify basic power operation (2^3 = 8)
	expectedPow := "8"
	if executions[0] != expectedPow {
		t.Errorf("Basic power test failed.\nExpected: %s\nGot: %v",
			expectedPow, executions[0])
	}

	// Verify power with zero exponent (2^0 = 1)
	expectedZeroPow := "1"
	if executions[1] != expectedZeroPow {
		t.Errorf("Power with zero exponent test failed.\nExpected: %s\nGot: %v",
			expectedZeroPow, executions[1])
	}

	// Verify power with negative exponent (2^-2 = 0.25)
	// because of scale is 0, so result is 0
	expectedNegPow := "0"
	if executions[2] != expectedNegPow {
		t.Errorf("Power with negative exponent test failed.\nExpected: %s\nGot: %v",
			expectedNegPow, executions[2])
	}
}
